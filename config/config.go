package config

import (
	"common/cache"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/coreos/etcd/client"
)

var Conf *Configure

type MysqlConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Addr     string `json:"addr"`
	DB       string `json:"db"`
	Charset  string `json:"charset"`
	MaxIdle  int    `json:"max_idle"`
	MaxConn  int    `json:"max_conn"`
}

type RPCConfig struct {
	Scheme         string `json:"scheme"`
	Name           string `json:"name"`
	UpdateInterval uint32 `json:"update_interval"`
}

type CDNConfig struct {
	VideoUrl string `json:"video_url"`
}

type OtherConfig struct {
	ShareUrl     string `json:"share_url"`
	WikiUrl      string `json:"wiki_url"`
	Timeout      uint64 `json:"timeout"`
	AppKey       string `json:"app_key"`
	AppSecret    string `json:"app_secret"`
	CrowFoundH5  string `json:"crowfound_h5"`
	CrowOrderUrl string `json:"croworder_url"`
}

type WxpayConfig struct {
	AppID        string `json:"app_id"`
	AppKey       string `json:"app_key"`
	AppSecret    string `json:"app_secret"`
	MchID        string `json:"mch_id"`
	Body         string `json:"body"`
	TradeType    string `json:"trade_type"`
	Package      string `json:"package"`
	RequestUrl   string `json:"request_url"`
	NotifyUrl    string `json:"notify_url"`
	ActNotifyUrl string `json:"act_notify_url"`
	RefundUrl    string `json:"refund_url"`
}

type AlipayConfig struct {
	NotifyUrl    string `json:"notify_url"`
	ActNotifyUrl string `json:"act_notify_url"`
	Appid        string `json:"appid"`
	Charset      string `json:"charset"`
	RefundUrl    string `json:"refund_url"`
}

type PayConfig struct {
	Wxpay  WxpayConfig
	Alipay AlipayConfig
}

type ElasticConfig struct {
	Addr string `json:"addr"`
}

type SpiderConfig struct {
	ProxyAddrs []string `json:"proxy_addrs"`
}

type Configure struct {
	DBConfig       MysqlConfig
	CacheConfig    cache.RedisConfig
	RPCBasePath    string
	CrewRPC        RPCConfig
	LiveRPC        RPCConfig
	UserRPC        RPCConfig
	AuthRPC        RPCConfig
	ActivityRPC    RPCConfig
	CommunityRPC   RPCConfig
	DispatcherRPC  RPCConfig
	HashtagRPC     RPCConfig
	InformationRPC RPCConfig
	LabelRPC       RPCConfig
	ReplayRPC      RPCConfig
	SearchRPC      RPCConfig
	ShortVideoRPC  RPCConfig
	CommentRPC     RPCConfig
	Elastic        ElasticConfig
	Cdn            CDNConfig
	Other          OtherConfig
	Pay            PayConfig
	Spider         SpiderConfig
}

func GetConfig(cli client.Client) *Configure {
	keysAPI := client.NewKeysAPI(cli)
	configure := &Configure{}
	mysqlBasePath := "/config/paas10000/mysql"
	if dbcfg := getMysqlConfig(keysAPI, mysqlBasePath); dbcfg != nil {
		configure.DBConfig = *dbcfg
	}
	redisBasePath := "/config/paas10000/redis"
	if redisConfig := getCacheConfig(keysAPI, redisBasePath); redisConfig != nil {
		configure.CacheConfig = *redisConfig
	}
	rpcBasePath := "/config/paas10000/microserver/basepath"
	if resp, err := keysAPI.Get(context.Background(), rpcBasePath, &client.GetOptions{
		Recursive: false,
	}); err == nil && resp != nil && resp.Node != nil {
		configure.RPCBasePath = resp.Node.Value
	}
	// cdn
	cdnBasePath := "/config/paas10000/microserver/live/cdn"
	if conf := getCDNConfig(keysAPI, cdnBasePath); conf != nil {
		configure.Cdn = *conf
	}
	// 其他
	otherBasePath := "/config/paas10000/other"
	if conf := getOtherConfig(keysAPI, otherBasePath); conf != nil {
		configure.Other = *conf
	}

	// 支付
	payBasePath := "/config/paas10000/pay"
	if conf := getPayConfig(keysAPI, payBasePath); conf != nil {
		configure.Pay = *conf
	}

	// 剧组
	crewBasePath := "/config/paas10000/microserver/crew"
	if conf := getRPCConfig(keysAPI, crewBasePath); conf != nil {
		configure.CrewRPC = *conf
	}
	// 直播
	liveBasePath := "/config/paas10000/microserver/live"
	if conf := getRPCConfig(keysAPI, liveBasePath); conf != nil {
		configure.LiveRPC = *conf
	}
	// 用户
	userBasePath := "/config/paas10000/microserver/user"
	if conf := getRPCConfig(keysAPI, userBasePath); conf != nil {
		configure.UserRPC = *conf
	}
	// 授权
	authBasePath := "/config/paas10000/microserver/auth"
	if conf := getRPCConfig(keysAPI, authBasePath); conf != nil {
		configure.AuthRPC = *conf
	}
	// 活动
	activityBasePath := "/config/paas10000/microserver/activity"
	if conf := getRPCConfig(keysAPI, activityBasePath); conf != nil {
		configure.ActivityRPC = *conf
	}
	// 社区
	communityBasePath := "/config/paas10000/microserver/community"
	if conf := getRPCConfig(keysAPI, communityBasePath); conf != nil {
		configure.CommunityRPC = *conf
	}
	// dispatcher
	dispatcherBasePath := "/config/paas10000/microserver/dispatcher"
	if conf := getRPCConfig(keysAPI, dispatcherBasePath); conf != nil {
		configure.DispatcherRPC = *conf
	}
	// hashtag
	hashtagBasePath := "/config/paas10000/microserver/hashtag"
	if conf := getRPCConfig(keysAPI, hashtagBasePath); conf != nil {
		configure.HashtagRPC = *conf
	}
	// information
	informationBasePath := "/config/paas10000/microserver/information"
	if conf := getRPCConfig(keysAPI, informationBasePath); conf != nil {
		configure.InformationRPC = *conf
	}
	// 标签
	labelBasePath := "/config/paas10000/microserver/label"
	if conf := getRPCConfig(keysAPI, labelBasePath); conf != nil {
		configure.LabelRPC = *conf
	}
	// replay
	replayBasePath := "/config/paas10000/microserver/replay"
	if conf := getRPCConfig(keysAPI, replayBasePath); conf != nil {
		configure.ReplayRPC = *conf
	}
	// search
	searchBasePath := "/config/paas10000/microserver/search"
	if conf := getRPCConfig(keysAPI, searchBasePath); conf != nil {
		configure.SearchRPC = *conf
	}
	// shortvideo
	shortvideoBasePath := "/config/paas10000/microserver/shortvideo"
	if conf := getRPCConfig(keysAPI, shortvideoBasePath); conf != nil {
		configure.ShortVideoRPC = *conf
	}
	// comment
	commentBasePath := "/config/paas10000/microserver/comment"
	if conf := getRPCConfig(keysAPI, commentBasePath); conf != nil {
		configure.CommentRPC = *conf
	}
	// elastic
	elasticBasePath := "/config/paas10000/microserver/elastic"
	if conf := getElasticConfig(keysAPI, elasticBasePath); conf != nil {
		configure.Elastic = *conf
	}
	// spider config
	spiderBasePath := "/config/paas10000/innerserver/spider"
	if conf := getSpiderConfig(keysAPI, spiderBasePath); conf != nil {
		configure.Spider = *conf
	}

	Conf = configure
	return configure
}

func getCacheConfig(api client.KeysAPI, basePath string) *cache.RedisConfig {
	if resp, err := api.Get(context.Background(), basePath, &client.GetOptions{
		Recursive: true,
	}); err == nil && resp != nil && resp.Node != nil {
		cacheCfg := &cache.RedisConfig{}
		networkKey := basePath + "/network"
		addrKey := basePath + "/addr"
		pwdKey := basePath + "/password"
		ctoKey := basePath + "/connecttimeout"
		rtoKey := basePath + "/readtimeout"
		wtoKey := basePath + "/writetimeout"
		maxActiveKey := basePath + "/maxactivite"
		maxIdleKey := basePath + "/maxidel"
		idletoKey := basePath + "/ideltimeout"
		waitKey := basePath + "/wait"
		for _, node := range resp.Node.Nodes {
			switch node.Key {
			case networkKey:
				cacheCfg.Network = node.Value
			case addrKey:
				cacheCfg.Address = node.Value
			case pwdKey:
				cacheCfg.Password = node.Value
			case ctoKey:
				tmp, _ := strconv.ParseUint(node.Value, 10, 64)
				cacheCfg.ConnectTimeout = int(tmp)
			case rtoKey:
				tmp, _ := strconv.ParseUint(node.Value, 10, 64)
				cacheCfg.ReadTimeout = int(tmp)
			case wtoKey:
				tmp, _ := strconv.ParseUint(node.Value, 10, 64)
				cacheCfg.WriteTimeout = int(tmp)
			case maxIdleKey:
				tmp, _ := strconv.ParseUint(node.Value, 10, 64)
				cacheCfg.MaxIdle = int(tmp)
			case idletoKey:
				tmp, _ := strconv.ParseUint(node.Value, 10, 64)
				cacheCfg.IdleTimeout = int(tmp)
			case maxActiveKey:
				tmp, _ := strconv.ParseUint(node.Value, 10, 64)
				cacheCfg.MaxActive = int(tmp)
			case waitKey:
				cacheCfg.Wait = node.Value == "true"
			}
		}
		return cacheCfg
	} else {
		fmt.Printf("get redis config failed. error:%s", err)
		os.Exit(1)
	}
	return nil
}

func getMysqlConfig(api client.KeysAPI, basePath string) *MysqlConfig {
	if resp, err := api.Get(context.Background(), basePath, &client.GetOptions{
		Recursive: true,
	}); err == nil && resp != nil && resp.Node != nil {
		dbcfg := &MysqlConfig{}
		userKey := basePath + "/user"
		pwdKey := basePath + "/password"
		addrKey := basePath + "/addr"
		dbKey := basePath + "/db"
		charsetKey := basePath + "/charset"
		maxIdleKey := basePath + "/max_idle"
		maxConnKey := basePath + "/max_conn"
		for _, node := range resp.Node.Nodes {
			switch node.Key {
			case userKey:
				dbcfg.User = node.Value
			case pwdKey:
				dbcfg.Password = node.Value
			case addrKey:
				dbcfg.Addr = node.Value
			case dbKey:
				dbcfg.DB = node.Value
			case charsetKey:
				dbcfg.Charset = node.Value
			case maxIdleKey:
				idle, _ := strconv.ParseUint(node.Value, 10, 64)
				dbcfg.MaxIdle = int(idle)
			case maxConnKey:
				conn, _ := strconv.ParseUint(node.Value, 10, 64)
				dbcfg.MaxConn = int(conn)
			}
		}
		return dbcfg
	} else {
		fmt.Printf("get mysql config failed. error:%s", err)
		os.Exit(1)
	}
	return nil
}

func getOtherConfig(api client.KeysAPI, basePath string) *OtherConfig {
	if resp, err := api.Get(context.Background(), basePath, &client.GetOptions{
		Recursive: true,
	}); err == nil && resp != nil && resp.Node != nil {
		othercfg := &OtherConfig{}
		shareurlKey := basePath + "/shareurl"
		timeoutKey := basePath + "/timeout"
		wikiurlKey := basePath + "/wikiurl"
		appkeyKey := basePath + "/appkey"
		appsecretKey := basePath + "/appsecret"
		crowfound := basePath + "/crowfunding_h5"
		croworderurl := basePath + "/orderurl"
		for _, node := range resp.Node.Nodes {
			switch node.Key {
			case shareurlKey:
				othercfg.ShareUrl = node.Value
			case wikiurlKey:
				othercfg.WikiUrl = node.Value
			case appkeyKey:
				othercfg.AppKey = node.Value
			case appsecretKey:
				othercfg.AppSecret = node.Value
			case timeoutKey:
				othercfg.Timeout, _ = strconv.ParseUint(node.Value, 10, 64)
			case crowfound:
				othercfg.CrowFoundH5 = node.Value
			case croworderurl:
				othercfg.CrowOrderUrl = node.Value
			}
		}
		return othercfg
	} else {
		fmt.Printf("get %s other config failed. error:%s", basePath, err)
		os.Exit(1)
	}
	return nil
}

func getPayConfig(api client.KeysAPI, basePath string) *PayConfig {
	paycfg := &PayConfig{}
	// wxpay
	if resp, err := api.Get(context.Background(), basePath+"/wxpay", &client.GetOptions{
		Recursive: true,
	}); err == nil && resp != nil && resp.Node != nil {
		wxpaycfg := WxpayConfig{}
		wxBasePath := basePath + "/wxpay"
		for _, node := range resp.Node.Nodes {
			switch node.Key {
			case wxBasePath + "/app_id":
				wxpaycfg.AppID = node.Value
			case wxBasePath + "/app_key":
				wxpaycfg.AppKey = node.Value
			case wxBasePath + "/app_secret":
				wxpaycfg.AppSecret = node.Value
			case wxBasePath + "/mch_id":
				wxpaycfg.MchID = node.Value
			case wxBasePath + "/body":
				wxpaycfg.Body = node.Value
			case wxBasePath + "/trade_type":
				wxpaycfg.TradeType = node.Value
			case wxBasePath + "/package":
				wxpaycfg.Package = node.Value
			case wxBasePath + "/request_url":
				wxpaycfg.RequestUrl = node.Value
			case wxBasePath + "/notify_url":
				wxpaycfg.NotifyUrl = node.Value
			case wxBasePath + "/act_notify_url":
				wxpaycfg.ActNotifyUrl = node.Value
			case wxBasePath + "/refundurl":
				wxpaycfg.RefundUrl = node.Value
			}
		}
		paycfg.Wxpay = wxpaycfg
	} else {
		fmt.Printf("get %s wxpay config failed. error:%s", basePath, err)
		return nil
	}

	// alipay
	if resp, err := api.Get(context.Background(), basePath+"/alipay", &client.GetOptions{
		Recursive: true,
	}); err == nil && resp != nil && resp.Node != nil {
		alipaycfg := AlipayConfig{}
		aliBasePath := basePath + "/alipay"
		for _, node := range resp.Node.Nodes {
			switch node.Key {
			case aliBasePath + "/notify_url":
				alipaycfg.NotifyUrl = node.Value
			case aliBasePath + "/act_notify_url":
				alipaycfg.ActNotifyUrl = node.Value
			case aliBasePath + "/appid":
				alipaycfg.Appid = node.Value
			case aliBasePath + "/refundurl":
				alipaycfg.RefundUrl = node.Value
			case aliBasePath + "/charset":
				alipaycfg.Charset = node.Value
			}
		}
		paycfg.Alipay = alipaycfg
	} else {
		fmt.Printf("get %s wxpay config failed. error:%s", basePath, err)
		return nil
	}

	return paycfg
}

func getCDNConfig(api client.KeysAPI, basePath string) *CDNConfig {
	if resp, err := api.Get(context.Background(), basePath, &client.GetOptions{
		Recursive: true,
	}); err == nil && resp != nil && resp.Node != nil {
		cdncfg := &CDNConfig{}
		videourlKey := basePath + "/videourl"
		for _, node := range resp.Node.Nodes {
			switch node.Key {
			case videourlKey:
				cdncfg.VideoUrl = node.Value
			}
		}
		return cdncfg
	} else {
		fmt.Printf("get %s cdn config failed. error:%s", basePath, err)
		os.Exit(1)
	}
	return nil
}

func getRPCConfig(api client.KeysAPI, basePath string) *RPCConfig {
	if resp, err := api.Get(context.Background(), basePath, &client.GetOptions{
		Recursive: true,
	}); err == nil && resp != nil && resp.Node != nil {
		rpccfg := &RPCConfig{}
		schemeKey := basePath + "/scheme"
		nameKey := basePath + "/name"
		uiKey := basePath + "/updateinterval"
		for _, node := range resp.Node.Nodes {
			switch node.Key {
			case schemeKey:
				rpccfg.Scheme = node.Value
			case nameKey:
				rpccfg.Name = node.Value
			case uiKey:
				tmp, _ := strconv.ParseUint(node.Value, 10, 64)
				rpccfg.UpdateInterval = uint32(tmp)
			}
		}
		return rpccfg
	} else {
		fmt.Printf("get %s rpc config failed. error:%s", basePath, err)
		os.Exit(1)
	}
	return nil
}

func getSpiderConfig(api client.KeysAPI, basePath string) *SpiderConfig {
	if resp, err := api.Get(context.Background(), basePath, &client.GetOptions{
		Recursive: true,
	}); err == nil && resp != nil && resp.Node != nil {
		spidercfg := &SpiderConfig{}
		addrsKey := basePath + "/addrs"
		for _, node := range resp.Node.Nodes {
			switch node.Key {
			case addrsKey:
				spidercfg.ProxyAddrs = strings.Split(node.Value, "|")
			}
		}
		return spidercfg
	} else {
		fmt.Printf("get %s spider config failed. error:%s", basePath, err)
		os.Exit(1)
	}
	return nil
}
func getElasticConfig(api client.KeysAPI, basePath string) *ElasticConfig {
	if resp, err := api.Get(context.Background(), basePath, &client.GetOptions{
		Recursive: true,
	}); err == nil && resp != nil && resp.Node != nil {
		elasticcfg := &ElasticConfig{}
		addrKey := basePath + "/addr"
		for _, node := range resp.Node.Nodes {
			switch node.Key {
			case addrKey:
				elasticcfg.Addr = node.Value
			}
		}
		return elasticcfg
	} else {
		fmt.Printf("get %s rpc config failed. error:%s", basePath, err)
		os.Exit(1)
	}
	return nil
}
