/*
 Package xxx

 Copyright(c) 2000-2016 Quqian Technology Company Limit.
 All rights reserved.
 创建人：  longyongchun
 创建日期：2016-10-25 15:46:44
 修改记录：
 ----------------------------------------------------------------
 修改人        |  修改日期    				|    备注
 ----------------------------------------------------------------
 longyongchun  | 2016-10-25 15:46:44		|    创建文件
 ----------------------------------------------------------------
*/

package utils

import (
	"common/config"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	gocontext "context"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"

	"github.com/astaxie/beego/orm"

	"github.com/golang/protobuf/proto"

	"github.com/coreos/etcd/client"
)

var PAGE_SIZE = uint32(10)

func init() {
	orm.RegisterModel(new(UnqiueID))
}

type MicroTask func(conf *config.Configure, etcdaddrs []string)
type WatchTask func(action string, val string)

func Run(svrName string, tasks ...MicroTask) {
	// 初始化环境
	runmode := GetEnv("runmode")
	if runmode != "dev" && runmode != "prod" {
		fmt.Printf("runmode is invalidate maybe value: dev or prod\n")
		os.Exit(1)
	}
	addrs := GetEnv("etcd_addrs")
	etcdaddrs := strings.Split(addrs, ",")
	if len(etcdaddrs) == 0 {
		fmt.Printf(`etcd addr is invalidate maybe value: http://127.0.0.1:2379\n`)
		os.Exit(1)
	}
	initEtcd(etcdaddrs)

	configure := config.GetConfig(GetEtcdClient())
	if configure == nil {
		fmt.Printf("get configure failed")
		os.Exit(1)
	}
	InitDB(&configure.DBConfig)
	InitCache(&configure.CacheConfig)
	InitLog(svrName, runmode)

	go WatchConfig(map[string]WatchTask{
		fmt.Sprintf("/config/paas10000/microserver/%s/log-mode", svrName): WatchLogMode,
	})

	monitorAddr := GetEnv(fmt.Sprintf("monitor_%s_addr", svrName))
	go func() {
		ln, err := net.Listen("tcp", monitorAddr)
		if err != nil {
			Error("", "listen monitor addr failed. addr:%s error:%s", monitorAddr, err)
			return
		}
		for {
			data := make([]byte, 200)
			reqData := fmt.Sprintf("%s:%s", svrName, "health")
			if conn, err := ln.Accept(); err == nil {
				if n, err := conn.Read(data); err == nil {
					if string(data[:n]) == reqData {
						// 健康检查
						conn.Write(data[:n])
					}
				}
				conn.Close()
			}
		}
	}()
	for _, task := range tasks {
		task(configure, etcdaddrs)
	}
}

func WatchLogMode(action string, val string) {
	Info("", "watch the log mode change action:%s log mode:%s", action, val)
	if action != "set" && action != "update" && action != "create" && action != "compareAndSwap" {
		return
	}
	logLevel := -1
	switch val {
	case "debug":
		logLevel = logs.LevelDebug
	case "info":
		logLevel = logs.LevelInformational
	case "notice":
		logLevel = logs.LevelNotice
	case "warning":
		logLevel = logs.LevelWarning
	case "error":
		logLevel = logs.LevelError
	case "critical":
		logLevel = logs.LevelCritical
	case "alert":
		logLevel = logs.LevelAlert
	case "emergency":
		logLevel = logs.LevelEmergency
	default:
		Info("", "the log mode[%s] is not support. level(debug, info, notice, warning, error, critical, alert, emergency)", val)
		return
	}
	logModel = val
	if logLevel == logs.LevelDebug {
		orm.Debug = true
		redis.Debug = true
	} else {
		orm.Debug = false
		redis.Debug = false
	}
	log.SetLevel(logLevel)
	accessLog.SetLevel(logLevel)
	return
}

func WatchConfig(keyFunc map[string]WatchTask) {
	cli := GetEtcdClient()
	keysAPI := client.NewKeysAPI(cli)
	for key, f := range keyFunc {
		watcher := keysAPI.Watcher(key, &client.WatcherOptions{Recursive: false})
		go func(watcher client.Watcher, f WatchTask) {
			for {
				if r, err := watcher.Next(gocontext.Background()); err == nil {
					fmt.Println("test")
					f(r.Action, r.Node.Value)
				}
			}
		}(watcher, f)
	}
}

func GetEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		fmt.Printf(`hasn't %s env var\n`, name)
		os.Exit(1)
	}
	return value
}

type UnqiueID struct {
	Id   uint64 `orm:"auto;column(id)"`
	Time uint64 `orm:"column(time)"`
}

func (ud UnqiueID) TableName() string {
	return "t_unique_id"
}

func (ud *UnqiueID) Create(db orm.Ormer) error {
	if _, err := db.Insert(ud); err != nil {
		return err
	}
	return nil
}

func ParseRequiredFromBytes(required []byte, cs proto.Message, args ...interface{}) error {
	var v interface{} = ""
	if len(args) > 0 {
		v = args[0]
	}

	if err := proto.Unmarshal(required, cs); err != nil {
		Error(v, "unmarshaling error: %s", err)
		return err
	}
	// 打印真实的请求数据
	val, _ := json.Marshal(cs)
	Debug(v, "require: %s", string(val))

	return nil
}

func MarshalToResponseString(response proto.Message, args ...interface{}) (string, error) {
	var v interface{} = ""
	if len(args) > 0 {
		v = args[0]
	}

	// 打印应答数据
	val, _ := json.Marshal(response)
	Debug(v, "response: %s", string(val))

	r, err := proto.Marshal(response)
	if err != nil {
		Error(v, "marshaling error: %s", err)
		return base64.StdEncoding.EncodeToString(r), err
	}

	return base64.StdEncoding.EncodeToString(r), nil
}

func MarshalToProtoString(response proto.Message) (string, error) {
	r, err := proto.Marshal(response)
	if err != nil {
		Error("", "marshaling error: %s", err)
		return "", err
	}

	return string(r), nil
}

/*func SendToClient(ctx *context.Context, status int, resp string) {
	//ctx.Output.SetStatus(status)
	ctx.WriteString(resp)
}*/

// 获取特定的日志ID
func GetUniqueLogID(cmd string) string {
	nowTime := time.Now()
	return fmt.Sprintf("%d_%s", nowTime.Unix(), cmd)
}

func GetStructNotZeroField(arg interface{}) []string {
	v := reflect.ValueOf(arg)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	result := make([]string, 0)
	if v.Kind() == reflect.Struct {
		numField := v.NumField()
		for i := 0; i < numField; i++ {
			if !reflect.DeepEqual(v.Field(i).Interface(), reflect.Zero(v.Field(i).Type()).Interface()) {
				result = append(result, v.Type().Field(i).Name)
			}
		}
	}
	return result
}

func GetHeaderUint32Value(header http.Header, key string) (uint32, error) {
	val := header.Get(key)
	v, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint32(v), err
}

func GetHeaderStringValue(header http.Header, key string) (string, error) {
	return header.Get(key), nil
}

func GetDataFromQueryParam(ctx *context.Context, key string) string {
	//Debug("", "%+v", ctx.Input.Context.Request.Form)
	return ctx.Input.Query(key)
}

func GetDataFromBody(ctx *context.Context, key string) string {
	if ctx.Request.Form == nil {
		ctx.Request.ParseForm()
	}
	//Debug("", "%+v", ctx.Request.Form)
	mapParams := ctx.Request.Form
	data := ""
	if params, ok := mapParams["data"]; ok {
		if params != nil && len(params) > 0 {
			data = params[0]
		}
	}
	return data
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

//// 裁剪出无html标签和乱码的内容，可以指定字节长度
func CutUtf8ContentWithoutHtml(src string, limitCount uint32) (result string, err error) {
	chars := []rune(src)
	res := []rune{}
	closed := true
	for _, char := range chars {
		if char == rune('<') {
			closed = false
		}
		if char == rune('>') {
			closed = true
			continue
		}
		if closed && char != rune('\n') && char != rune('\r') && char != rune(' ') {
			if limitCount > uint32(len(res)) {
				res = append(res, char)
			}

			if limitCount < uint32(len(res)) {
				break
			}
		}
	}
	if uint32(len(res)) >= limitCount {
		res = append(res, []rune("...")...)
	}
	return string(res), nil
}

func GetUint64MaxAndMin(nums ...uint64) (max uint64, min uint64) {
	if len(nums) == 0 {
		return 0, 0
	}
	max = nums[0]
	min = nums[0]
	start := 0
	if len(nums)%2 == 0 {
		min = nums[0]
		max = nums[1]
		if min > max {
			min, max = max, min
		}
		start = 2
	} else {
		max = nums[0]
		min = nums[1]
		start = 1
	}
	for i := start; i < len(nums); i += 2 {
		if nums[i] > nums[i+1] {
			if max < nums[i] {
				max = nums[i]
			}
			if min > nums[i+1] {
				min = nums[i+1]
			}
		} else {
			if max < nums[i+1] {
				max = nums[i+1]
			}
			if min > nums[i] {
				min = nums[i]
			}
		}
	}
	return max, min
}

func CheckAllEmptyString(src ...string) bool {
	result := true
	for _, s := range src {
		if s != "" {
			result = false
			break
		}
	}
	return result
}
func GetMysqlUniqueID() uint64 {
	ud := &UnqiueID{
		Time: uint64(time.Now().Unix()),
	}
	if err := ud.Create(orm.NewOrm()); err != nil {
		return 0
	}
	return ud.Id
}

func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}

func CopyStructSameField(dst interface{}, src interface{}) error {
	d := reflect.ValueOf(dst)
	if d.Kind() != reflect.Ptr {
		return fmt.Errorf("dst is not ptr")
	}
	d = d.Elem()
	if d.Kind() != reflect.Struct {
		return fmt.Errorf("dst is not struct")
	}

	s := reflect.ValueOf(src)
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}
	if s.Kind() != reflect.Struct {
		return fmt.Errorf("src is not struct")
	}
	srcmap := map[string]reflect.Value{}
	numField := s.NumField()
	for i := 0; i < numField; i++ {
		srcmap[s.Type().Field(i).Name] = s.Field(i)
	}

	numField = d.NumField()
	for i := 0; i < numField; i++ {
		name := d.Type().Field(i).Name
		if v, ok := srcmap[name]; ok {
			if v.Type().Kind() == d.Type().Field(i).Type.Kind() {
				if v.Type().Kind() == reflect.Struct {
					CopyStructSameField(d.Field(i).Addr().Interface(), v.Interface())
				} else if v.Type().Kind() == reflect.Slice {
					d.Field(i).Set(reflect.MakeSlice(d.Field(i).Type(), 0, v.Len()))
					if t, err := CopySlice(d.Field(i).Interface(), v.Interface()); err == nil {
						tmp := reflect.ValueOf(t)
						d.Field(i).Set(reflect.MakeSlice(d.Field(i).Type(), tmp.Len(), tmp.Len()))
						reflect.Copy(d.Field(i), reflect.ValueOf(t))
					}
				} else {
					d.Field(i).Set(v)
				}
			}
		}
	}
	return nil
}

func CopySlice(dst interface{}, src interface{}) (ddd interface{}, err error) {
	d := reflect.ValueOf(dst)
	if d.Kind() == reflect.Ptr {
		d = d.Elem()
	}
	if d.Kind() != reflect.Slice {
		return dst, fmt.Errorf("dst is not slice")
	}

	s := reflect.ValueOf(src)
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}
	if s.Kind() != reflect.Slice {
		return dst, fmt.Errorf("src is not slice")
	}

	count := s.Len()
	for i := 0; i < count; i++ {
		t := reflect.New(d.Type().Elem())
		if t.Type().Elem().Kind() == reflect.Struct && s.Index(i).Type().Kind() == reflect.Struct {
			CopyStructSameField(t.Interface(), s.Index(i).Interface())
		} else if t.Kind() == reflect.Slice && s.Index(i).Kind() == reflect.Slice {
			CopySlice(t, s.Field(i))
		} else if t.Kind() != reflect.Slice && t.Kind() != reflect.Struct && s.Index(i).Kind() != reflect.Slice && s.Index(i).Kind() != reflect.Struct {
			t.Set(s.Field(i))
		}
		d = reflect.Append(d, t.Elem())
	}
	return d.Interface(), nil
}

func GetCommonScore() string {
	nowTime := time.Now()
	score := fmt.Sprintf("%d.%d", nowTime.Unix(), nowTime.Nanosecond())
	return score
}

func GetRandomMember(mems []string, num int) (result []string, err error) {
	total := len(mems)
	if total <= num {
		num = total
	}
	result = make([]string, 0, num)
	for i := 1; i <= num; i++ {
		idx := rand.New(rand.NewSource(time.Now().Unix())).Intn(total)
		result = append(result, mems[idx])
		mems[idx], mems[total-1] = mems[total-1], mems[idx]
		total -= 1
	}
	return result, nil
}

func GetRpcHeaders(ctx *context.Context, addChild ...bool) map[string][]string {
	result := map[string][]string{
		"trace_id": []string{ctx.Input.Header("trace_id")},
		"span_id":  []string{ctx.Input.Header("span_id")},
		"child_id": []string{ctx.Input.Header("child_id")},
	}
	add := true
	if len(addChild) > 0 {
		add = addChild[0]
	}
	if add {
		cidstr := ctx.Input.Header("child_id")
		cid, _ := strconv.ParseUint(cidstr, 10, 64)
		ctx.Input.Context.Request.Header.Set("child_id", fmt.Sprintf("%d", cid+1))
	}
	return result
}

func WithRunTime(ctx gocontext.Context, method string, decoPtr, fn interface{}) {
	decoratedFunc := reflect.ValueOf(decoPtr).Elem()
	targetFunc := reflect.ValueOf(fn)

	v := reflect.MakeFunc(decoratedFunc.Type(), func(in []reflect.Value) (out []reflect.Value) {
		defer func(t time.Time) {
			Info(ctx, "function:%s exec time:%v", method, time.Since(t))
		}(time.Now())
		out = targetFunc.Call(in)
		return
	})

	decoratedFunc.Set(v)
	return
}
