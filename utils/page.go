package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	elastic "gopkg.in/olivere/elastic.v5"
)

type PageDataParam struct {
	Slide         uint32 // 活动方向 1: 向上滑动 2:向下滑动
	TopWholeValue string
	TopSameNum    uint32 // 页头的wholevalue相同个数
	BotWholeValue string // 页尾的wholevalue值
	BotSameNum    uint32 // 页尾的wholevalue相同个数
	Num           uint32 // 一次拉的条数
	RedisKey      string // 要拉的redis key列表
	SequenceType  uint32 // 1:正序 2:反序
}

type ZSetResult struct {
	Member string
	Score  string
}

const (
	POSITIVE_SEQUENCE uint32 = 1
	REVERSE_SEQUENCE  uint32 = 2
)
const (
	SLIDE_UP   uint32 = 1
	SLIDE_DWON uint32 = 2
)

func GetPageData(zoneDataParam *PageDataParam, friendDynamics ...bool) (mapScoreMember map[string][]string, scores []string, hasMore bool, err error) {
	posFlag := true //标志获取到的数据是否是正序的
	var strMin, strMax, strCmd string
	var count uint32 = 0
	switch zoneDataParam.Slide {
	case SLIDE_DWON:
		switch zoneDataParam.SequenceType {
		case POSITIVE_SEQUENCE:
			strCmd = "ZRANGEBYSCORE"
			strMin = "0"
			if zoneDataParam.TopWholeValue != "" {
				strMax = zoneDataParam.TopWholeValue
			} else {
				strMax = "+inf"
			}
		default:
			//从大到小
			// app客户端特殊处理
			if zoneDataParam.TopWholeValue == "" && zoneDataParam.BotWholeValue == "" {
				strCmd = "ZREVRANGEBYSCORE"
				strMin = "+inf"
				strMax = "0"
				posFlag = false
			} else {
				strCmd = "ZRANGEBYSCORE"
				strMax = "+inf"
				if zoneDataParam.TopWholeValue != "" {
					strMin = zoneDataParam.TopWholeValue
				} else {
					strMin = "0"
				}
				posFlag = true
			}
		}
		// 需要拉取的条数
		count = zoneDataParam.TopSameNum + zoneDataParam.Num
	default:
		switch zoneDataParam.SequenceType {
		case POSITIVE_SEQUENCE:
			strCmd = "ZRANGEBYSCORE"
			strMax = "+inf"
			if zoneDataParam.BotWholeValue != "" {
				strMin = zoneDataParam.BotWholeValue
			} else {
				strMin = "0"
			}
		default:
			//从大到小
			strMax = "0"
			strCmd = "ZREVRANGEBYSCORE"
			if zoneDataParam.BotWholeValue != "" {
				strMin = zoneDataParam.BotWholeValue
			} else {
				strMin = "+inf"
			}
			posFlag = false
		}
		count = zoneDataParam.BotSameNum + zoneDataParam.Num
	}

	mapScoreMember = make(map[string][]string)
	scores = make([]string, 0)

	redisKey := zoneDataParam.RedisKey
	dest := make([]string, 0)
	var tmp_err error
	// 多获取一条数据 以此来判断是否还有更多数据
	if dest, tmp_err = Cache.CommandReturnStringSlice(strCmd, redisKey, strMin, strMax, "WITHSCORES", "LIMIT", 0, count+1); tmp_err != nil {
		Info("get rediskey failed rediskey: %s error:%s", redisKey, tmp_err)
		return
	}

	// 这里需要除以2，因为拉出了score值
	if len(dest)/2 > int(count) {
		hasMore = true
		dest = dest[:count*2]
	}

	dest_len := len(dest)
	for i := 0; i < dest_len-1; i += 2 {
		if dest[i] == "" || dest[i+1] == "" {
			continue
		}
		if _, ok := mapScoreMember[dest[i+1]]; !ok {
			mapScoreMember[dest[i+1]] = make([]string, 0, 1)
		}
		mapScoreMember[dest[i+1]] = append(mapScoreMember[dest[i+1]], dest[i])
		scoreLen := len(scores)
		if scoreLen != 0 && scores[scoreLen-1] == dest[i+1] {
			// 去重
			continue
		}
		scores = append(scores, dest[i+1])
	}
	if !posFlag {
		//逆序存储 对其数据进行正序
		scoresLen := len(scores)
		for i := 0; i < scoresLen/2; i++ {
			scores[i], scores[scoresLen-1-i] = scores[scoresLen-1-i], scores[i]
		}
	}
	return
}

func GetPostitivePageData(zoneDataParam *PageDataParam, scores []string, mapScoreMember map[string][]string) (res []ZSetResult, err error) {
	count := uint32(0)
	sameScoreDelCount := uint32(0)
	res = make([]ZSetResult, 0)
	// 参数非法时，直接取前面的数据
	if zoneDataParam.TopWholeValue == "" || zoneDataParam.BotWholeValue == "" {
		for _, score := range scores {
			if members, ok := mapScoreMember[score]; ok {
				for _, member := range members {
					res = append(res, ZSetResult{
						Member: member,
						Score:  score,
					})
					count += 1
					if count == zoneDataParam.Num {
						break
					}
				}
			}
			// 退出第二层循环
			if count == zoneDataParam.Num {
				break
			}
		}
		return
	}
	//拉老数据
	if zoneDataParam.Slide == SLIDE_DWON {
		for _, score := range scores {
			if score > zoneDataParam.TopWholeValue {
				// 新数据不进行处理
				break
			}
			if members, ok := mapScoreMember[score]; ok {
				for _, member := range members {
					if score == zoneDataParam.TopWholeValue && sameScoreDelCount < zoneDataParam.TopSameNum {
						//去重处理
						sameScoreDelCount += 1
						continue
					}
					res = append(res, ZSetResult{
						Member: member,
						Score:  score,
					})
					count += 1
					if count == zoneDataParam.Num {
						break
					}
				}
			}
			// 退出第二层循环
			if count == zoneDataParam.Num {
				break
			}
		}
	} else { //新数据
		for _, score := range scores {
			if score < zoneDataParam.BotWholeValue {
				//比底端老数据还新的数据不处理
				continue
			}
			if members, ok := mapScoreMember[score]; ok {
				for _, member := range members {
					if score == zoneDataParam.BotWholeValue && sameScoreDelCount < zoneDataParam.BotSameNum {
						//去重处理
						sameScoreDelCount += 1
						continue
					}
					res = append(res, ZSetResult{
						Member: member,
						Score:  score,
					})
					count += 1
					if count == zoneDataParam.Num {
						break
					}
				}
			}
			if count == zoneDataParam.Num {
				break
			}
		}
	}
	return
}

func GetReversePageData(zoneDataParam *PageDataParam, scores []string, mapScoreMember map[string][]string) (res []ZSetResult, err error) {
	count := uint32(0)
	sameScoreDelCount := uint32(0)
	res = make([]ZSetResult, 0)
	scoreLen := len(scores)
	// 参数非法时，直接取前面的数据
	if zoneDataParam.TopWholeValue == "" || zoneDataParam.BotWholeValue == "" {
		for i := scoreLen - 1; i >= 0; i -= 1 {
			score := scores[i]
			if members, ok := mapScoreMember[score]; ok {
				for _, member := range members {
					res = append(res, ZSetResult{
						Member: member,
						Score:  score,
					})
					count += 1
					if count == zoneDataParam.Num {
						break
					}
				}
			}
			// 退出第二层循环
			if count == zoneDataParam.Num {
				break
			}
		}
		return
	}
	//拉新数据
	if zoneDataParam.Slide == SLIDE_DWON {
		topScore, _ := strconv.ParseFloat(zoneDataParam.TopWholeValue, 32)
		for i := scoreLen - 1; i >= 0; i -= 1 {
			score := scores[i]
			fScore, _ := strconv.ParseFloat(score, 32)
			if fScore < topScore {
				// 老数据不进行处理
				break
			}
			if members, ok := mapScoreMember[score]; ok {
				for _, member := range members {
					if score == zoneDataParam.TopWholeValue && sameScoreDelCount < zoneDataParam.TopSameNum {
						//去重处理
						sameScoreDelCount += 1
						continue
					}
					res = append(res, ZSetResult{
						Member: member,
						Score:  score,
					})
					count += 1
					if count == zoneDataParam.Num {
						break
					}
				}
			}
			// 退出第二层循环
			if count == zoneDataParam.Num {
				break
			}
		}
	} else { //老数据
		bottomScore, _ := strconv.ParseFloat(zoneDataParam.BotWholeValue, 32)
		for i := scoreLen - 1; i >= 0; i -= 1 {
			score := scores[i]
			fScore, _ := strconv.ParseFloat(score, 32)
			if fScore > bottomScore {
				//比底端老数据还新的数据不处理
				continue
			}
			if members, ok := mapScoreMember[score]; ok {
				for _, member := range members {
					if score == zoneDataParam.BotWholeValue && sameScoreDelCount < zoneDataParam.BotSameNum {
						//去重处理
						sameScoreDelCount += 1
						continue
					}
					res = append(res, ZSetResult{
						Member: member,
						Score:  score,
					})
					count += 1
					if count == zoneDataParam.Num {
						break
					}
				}
			}
			if count == zoneDataParam.Num {
				break
			}
		}

	}
	return
}

func GetPagingParam(scores []string) (top string, bottom string) {
	if len(scores) == 0 {
		return
	}
	topValue := scores[0]
	scoreLen := len(scores)
	bottomValue := scores[scoreLen-1]

	topNum := uint32(0)
	bottomNum := uint32(0)

	for _, score := range scores {
		if score == topValue {
			topNum += 1
		} else {
			break
		}
	}

	for i := scoreLen - 1; i >= 0; i -= 1 {
		if scores[i] == bottomValue {
			bottomNum += 1
		} else {
			break
		}
	}

	top = JoinPage(topValue, topNum)
	bottom = JoinPage(bottomValue, bottomNum)
	return top, bottom
}

func JoinPage(wholeValue string, sameNum uint32) (page string) {
	return fmt.Sprintf("%s+%d", wholeValue, sameNum)
}
func JoinPageWithScore(score float64, wholeValue float64, sameNum uint32) (page string) {
	return fmt.Sprintf("%v+%f+%d", score, wholeValue, sameNum)
}

func GetTailSameNum(score string, results ...ZSetResult) uint32 {
	sameNum := uint32(0)
	if results != nil {
		resLen := len(results)
		for i := resLen - 1; i >= 0; i-- {
			if results[i].Score == score {
				sameNum += 1
			} else {
				break
			}
		}
	}
	return sameNum
}

func SplitPage(page string) (wholeValue string, sameNum uint32, err error) {
	if page == "" {
		return "", 0, nil
	}

	mems := strings.Split(page, "+")
	if mems == nil || len(mems) < 2 {
		return "", 0, errors.New("member is not enough")
	}
	num, err := strconv.ParseUint(mems[1], 10, 64)
	if err != nil {
		return "", 0, err
	}
	return mems[0], uint32(num), nil
}

func SplitPageWithScore(page string) (score string, wholeValue string, sameNum uint32, err error) {
	if page == "" {
		return "", "", 0, nil
	}

	mems := strings.Split(page, "+")
	if mems == nil || len(mems) < 2 {
		return "", "", 0, errors.New("member is not enough")
	} else if len(mems) == 2 {
		num, err := strconv.ParseUint(mems[1], 10, 64)
		if err != nil {
			return "", "", 0, err
		}
		return "", mems[0], uint32(num), nil
	} else {
		num, err := strconv.ParseUint(mems[2], 10, 64)
		if err != nil {
			return "", "", 0, err
		}
		return mems[0], mems[1], uint32(num), nil
	}
	return "", "", 0, nil
}

func IsEmpty(mems ...string) bool {
	if len(mems) == 0 {
		return true
	}
	for _, mem := range mems {
		if mem != "" {
			return false
		}
	}
	return true
}

type ElasticSort struct {
	Field string
	Asc   bool
}

type ElasticPageDataParam struct {
	Slide            uint32 // 活动方向 1: 向上滑动 2:向下滑动
	Top              string
	Bottom           string // 页尾的wholevalue值
	Num              uint32 // 一次拉的条数
	SequenceType     uint32 // 1:正序 2:反序
	Indexs           []string
	Types            []string
	Addr             string
	Match            map[string]string
	Query            elastic.Query
	Sorts            []ElasticSort
	PageFilter       string //页过滤字段 默认为time
	ExtraFilterField string //扩展过滤域 主要用于二级排序 比如短视频的置顶使用tab_video_score
}

type ObjectItem struct {
	Type  string
	ID    uint64
	Time  float64
	Score float64
}

type PageDataResult struct {
	Items      []ObjectItem
	Top        string
	Bottom     string
	HasMore    bool
	TotalCount uint64
}

type Logger struct {
}

func (l *Logger) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func GetPageDataFromElastic(zoneDataParam *ElasticPageDataParam) (result PageDataResult, err error) {
	// 拆分top bottom
	topScore, topWholeValue, topSameNum, err := SplitPageWithScore(zoneDataParam.Top)
	if err != nil {
		return result, err
	}
	botScore, botWholeValue, botSameNum, err := SplitPageWithScore(zoneDataParam.Bottom)
	if err != nil {
		return result, err
	}
	pageFilter := zoneDataParam.PageFilter
	if pageFilter == "" {
		pageFilter = "time"
	}
	client := GetElasticClient()
	searchServer := client.Search()
	if len(zoneDataParam.Indexs) > 0 {
		searchServer = searchServer.Index(zoneDataParam.Indexs...)
	}
	if len(zoneDataParam.Types) > 0 {
		searchServer = searchServer.Type(zoneDataParam.Types...)
	}

	// 组装match条件
	boolQuery := elastic.NewBoolQuery()
	if len(zoneDataParam.Match) == 0 {
		if zoneDataParam.Query != nil {
			boolQuery = boolQuery.Must(zoneDataParam.Query)
		} else {
			boolQuery = boolQuery.Must(elastic.NewMatchAllQuery())
		}
	} else {
		for key, value := range zoneDataParam.Match {
			boolQuery = boolQuery.Must(elastic.NewMatchQuery(key, value))
		}
	}

	if len(zoneDataParam.Sorts) > 0 {
		for _, s := range zoneDataParam.Sorts {
			searchServer = searchServer.Sort(s.Field, s.Asc)
		}
	}

	topRemove := true
	removeCount := uint32(0)
	count := uint32(0)
	postFlag := true
	switch zoneDataParam.Slide {
	case SLIDE_DWON:
		switch zoneDataParam.SequenceType {
		case POSITIVE_SEQUENCE:
			// 正序向下滑 time >=0
			if topWholeValue == "" {
				// time >= 0   正序排列
				if zoneDataParam.ExtraFilterField != "" {
					searchServer = searchServer.Sort(zoneDataParam.ExtraFilterField, true)
					boolQuery = boolQuery.Filter(elastic.NewRangeQuery(zoneDataParam.ExtraFilterField).Gte("0"))
				}
				boolQuery = boolQuery.Filter(elastic.NewRangeQuery(pageFilter).Gte("0"))
				searchServer = searchServer.Sort(pageFilter, true)
			} else {
				// time <= top 倒序排列
				if zoneDataParam.ExtraFilterField != "" {
					searchServer = searchServer.Sort(zoneDataParam.ExtraFilterField, false)
					boolQuery = boolQuery.Filter(elastic.NewRangeQuery(zoneDataParam.ExtraFilterField).Lte(topScore))
				}
				boolQuery = boolQuery.Filter(elastic.NewRangeQuery(pageFilter).Lte(topWholeValue))
				searchServer = searchServer.Sort(pageFilter, false)
				postFlag = false
			}
		default:
			// 逆序向下滑
			if topWholeValue == "" {
				// time >= 0 倒序排列
				if zoneDataParam.ExtraFilterField != "" {
					searchServer = searchServer.Sort(zoneDataParam.ExtraFilterField, false)
					boolQuery = boolQuery.Filter(elastic.NewRangeQuery(zoneDataParam.ExtraFilterField).Gte("0"))
				}
				boolQuery = boolQuery.Filter(elastic.NewRangeQuery(pageFilter).Gte("0"))
				searchServer = searchServer.Sort(pageFilter, false)
			} else {
				// time >= top 倒序排列
				if zoneDataParam.ExtraFilterField != "" {
					searchServer = searchServer.Sort(zoneDataParam.ExtraFilterField, false)
					boolQuery = boolQuery.Filter(elastic.NewRangeQuery(zoneDataParam.ExtraFilterField).Gte(topScore))
				}
				boolQuery = boolQuery.Filter(elastic.NewRangeQuery(pageFilter).Gte(topWholeValue))
				searchServer = searchServer.Sort(pageFilter, false)
			}
		}
		removeCount = topSameNum
		topRemove = false
		// 需要拉取的条数
		count = topSameNum + zoneDataParam.Num
	default:
		switch zoneDataParam.SequenceType {
		case POSITIVE_SEQUENCE:
			// 正序向上滑
			if botWholeValue == "" {
				// time >= 0 正序排序
				if zoneDataParam.ExtraFilterField != "" {
					boolQuery = boolQuery.Filter(elastic.NewRangeQuery(zoneDataParam.ExtraFilterField).Gte("0"))
					searchServer = searchServer.Sort(zoneDataParam.ExtraFilterField, true)
				}
				boolQuery = boolQuery.Filter(elastic.NewRangeQuery(pageFilter).Gte("0"))
				searchServer = searchServer.Sort(pageFilter, true)
			} else {
				// time >= bot 正序排序
				if zoneDataParam.ExtraFilterField != "" {
					boolQuery = boolQuery.Filter(elastic.NewRangeQuery(zoneDataParam.ExtraFilterField).Gte(botScore))
					searchServer = searchServer.Sort(zoneDataParam.ExtraFilterField, true)
				}
				boolQuery = boolQuery.Filter(elastic.NewRangeQuery(pageFilter).Gte(botWholeValue))
				searchServer = searchServer.Sort(pageFilter, true)
			}
		default:
			// 逆序向上滑
			if botWholeValue == "" {
				// time >= 0 逆序排序
				if zoneDataParam.ExtraFilterField != "" {
					boolQuery = boolQuery.Filter(elastic.NewRangeQuery(zoneDataParam.ExtraFilterField).Gte("0"))
					searchServer = searchServer.Sort(zoneDataParam.ExtraFilterField, false)
				}
				boolQuery = boolQuery.Filter(elastic.NewRangeQuery(pageFilter).Gte("0"))
				searchServer = searchServer.Sort(pageFilter, false)
			} else {
				// time <= bot 逆序排序
				if zoneDataParam.ExtraFilterField != "" {
					boolQuery = boolQuery.Filter(elastic.NewRangeQuery(zoneDataParam.ExtraFilterField).Lte(botScore))
					searchServer = searchServer.Sort(zoneDataParam.ExtraFilterField, false)
				}
				boolQuery = boolQuery.Filter(elastic.NewRangeQuery(pageFilter).Lte(botWholeValue))
				searchServer = searchServer.Sort(pageFilter, false)
			}
		}
		removeCount = botSameNum
		// 需要拉取的条数
		count = botSameNum + zoneDataParam.Num
	}
	searchServer = searchServer.Sort("id", false).From(0).Size(int(count))
	res, err := searchServer.Query(boolQuery).Do(context.Background())
	if err != nil {
		return result, err
	}

	items := []ObjectItem{}
	if res != nil && res.Hits != nil && len(res.Hits.Hits) > 0 {
		result.HasMore = res.Hits.TotalHits > int64(count)
		if res.Hits.TotalHits > int64(removeCount) {
			result.TotalCount = uint64(res.Hits.TotalHits - int64(removeCount))
		}
		for _, hit := range res.Hits.Hits {
			if hit.Source != nil {
				tmp := map[string]interface{}{}
				if err := json.Unmarshal([]byte(*hit.Source), &tmp); err == nil {
					id := uint64(0)
					t, _ := tmp["id"].(float64)
					id = uint64(t + 0.5)
					tm, _ := tmp["time"].(float64)
					score := float64(0.0)
					if s, ok := tmp[zoneDataParam.ExtraFilterField]; ok {
						score, _ = s.(float64)
					}
					items = append(items, ObjectItem{
						ID:    id,
						Type:  hit.Type,
						Time:  tm,
						Score: score,
					})
				} else {
					fmt.Println(err)
				}
			}
		}
	}

	if !postFlag {
		//需要逆序结果 对其数据进行正序
		itemsLen := len(items)
		for i := 0; i < itemsLen/2; i++ {
			items[i], items[itemsLen-1-i] = items[itemsLen-1-i], items[i]
		}
	}

	// 去掉上一次访问过的数据
	if topRemove {
		items = items[removeCount:]
	} else {
		items = items[:len(items)-int(removeCount)]
	}
	// 计算top与bottom
	if len(items) > 0 {
		topValue := items[0].Time
		topScore := items[0].Score
		bottomValue := items[len(items)-1].Time
		botScore := items[len(items)-1].Score

		topNum := uint32(0)
		bottomNum := uint32(0)

		for _, item := range items {
			if math.Abs(topScore-item.Score) <= 0.0000001 && math.Abs(item.Time-topValue) <= 0.0000001 {
				topNum += 1
			} else {
				break
			}
		}

		for i := len(items) - 1; i >= 0; i -= 1 {
			if math.Abs(botScore-items[i].Score) <= 0.0000001 && math.Abs(items[i].Time-bottomValue) <= 0.0000001 {
				bottomNum += 1
			} else {
				break
			}
		}

		result.Top = JoinPageWithScore(topScore, topValue, topNum)
		result.Bottom = JoinPageWithScore(botScore, bottomValue, bottomNum)
		result.Items = items
	}
	return result, nil
}
