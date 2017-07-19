package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"msplugin"
	"net/url"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	metrics "github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx"
	"github.com/smallnest/rpcx/codec"
	"github.com/smallnest/rpcx/core"
	"github.com/smallnest/rpcx/plugin"
)

var accessRPCLog *logs.BeeLogger

func SetRPCAccessLogger(loger *logs.BeeLogger) {
	accessRPCLog = loger
}

type Task func() error

type ClientTask func(client *rpcx.Client) error

func ExecRPCServerQuery(ctx context.Context, req interface{}, resp interface{}, tasks ...Task) (returnerr error) {
	startTime := time.Now()
	defer func() {
		status := "SUCCESS"
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				returnerr = e
			} else {
				returnerr = fmt.Errorf("%+v", r)
			}
			Error("", "unknown error:%s", returnerr)
		}

		if accessRPCLog != nil {
			if returnerr != nil {
				status = fmt.Sprintf("FAIL[%s]", returnerr)
				Error("", "server error:%s", returnerr)
			}
			if vs, ok := core.FromMapContext(ctx); ok {
				traceid := ""
				spanid := ""
				method := ""
				if s, ok := vs[":header"]; ok {
					if values, ok := s.(core.Header); ok {
						if ids, ok := values["trace_id"]; ok && len(ids) > 0 {
							traceid = ids[0]
						}
						if ids, ok := values["span_id"]; ok && len(ids) > 0 {
							spanid = ids[0]
						}
					}
				}
				if s, ok := vs[":method"]; ok {
					if ss, ok := s.(string); ok && ss != "" {
						method = ss
					}
				}
				reqContext, _ := json.Marshal(req)
				accessRPCLog.Info("|%s|%s|%s|%f|%s|%s|",
					traceid, spanid, method,
					float64(time.Now().UnixNano()-startTime.UnixNano())/1000000, string(reqContext), status)
			}
		}

		respContext, _ := json.Marshal(resp)
		Debug(ctx, "resp: %s", string(respContext))

	}()

	for _, task := range tasks {
		if err := task(); err != nil {
			returnerr = err
			return
		}
	}
	return
}

func ExecRPCClientQuery(s rpcx.ClientSelector, ctx context.Context, tasks ...ClientTask) (returnerr error) {
	client := rpcx.NewClient(s)
	defer func() {
		client.Close()
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				returnerr = e
			} else {
				returnerr = fmt.Errorf("%+v", r)
			}
		}
	}()
	client.ClientCodecFunc = codec.NewJSONRPC2ClientCodec
	// 失败时最多重试4次
	client.FailMode = rpcx.Failover
	client.Retries = 4

	for _, task := range tasks {
		if err := task(client); err != nil {
			return err
		}
	}
	return nil
}

func WrapRPCClientExec(client *rpcx.Client, method string, arg interface{}, reply interface{}, headers map[string][]string) error {
	rpcHeaders := map[string][]string{}
	target := func() error {
		cid := uint64(0)
		sid := ""
		tid := ""
		if cids, ok := headers["child_id"]; ok && len(cids) > 0 {
			cid, _ = strconv.ParseUint(cids[0], 10, 64)
		}
		if sids, ok := headers["span_id"]; ok && len(sids) > 0 {
			sid = sids[0]
		}
		if tids, ok := headers["trace_id"]; ok && len(tids) > 0 {
			tid = tids[0]
		}
		if len(headers["child_id"]) == 0 {
			headers["child_id"] = []string{}
		}

		headers["child_id"][0] = fmt.Sprintf("%d", cid+1)

		rpcHeaders["trace_id"] = []string{tid}
		rpcHeaders["span_id"] = []string{fmt.Sprintf("%s.%d", sid, cid)}
		rpcHeaders["child_id"] = []string{"1"}

		defer func() {
			if logModel == "debug" {
				datas, _ := json.Marshal(reply)
				Debug(core.NewContext(context.Background(), core.Header(rpcHeaders)), "reply: %s", string(datas))
			}
		}()

		argDatas, err := json.Marshal(arg)
		if err != nil {
			return err
		}
		argMap := map[string]interface{}{}
		d := json.NewDecoder(bytes.NewReader(argDatas))
		d.UseNumber()
		err = d.Decode(&argMap)
		if err != nil {
			return err
		}

		// 追加headers内容
		headerMap := (url.Values)(rpcHeaders)
		argMap["paas10000_rpc_header"] = headerMap.Encode()

		ctx := core.NewContext(context.Background(), core.Header(rpcHeaders))
		return client.Call(ctx, method, argMap, reply)
	}
	var decoFunc func() error
	WithRunTime(core.NewContext(context.Background(), core.Header(rpcHeaders)), method, &decoFunc, target)
	return decoFunc()
}

func GetRpcHeadersFromReq(ctx context.Context) map[string][]string {
	if vs, ok := core.FromMapContext(ctx); ok {
		traceid := ""
		spanid := ""
		childid := ""
		if s, ok := vs[":header"]; ok {
			if values, ok := s.(core.Header); ok {
				if ids, ok := values["trace_id"]; ok && len(ids) > 0 {
					traceid = ids[0]
				}
				if ids, ok := values["span_id"]; ok && len(ids) > 0 {
					spanid = ids[0]
				}
				if ids, ok := values["child_id"]; ok && len(ids) > 0 {
					childid = ids[0]
				}
			}
		}

		return map[string][]string{
			"trace_id": []string{traceid},
			"span_id":  []string{spanid},
			"child_id": []string{childid},
		}
	}
	return map[string][]string{}
}

type MicroServerConfig struct {
	Scheme         string
	Addr           string
	EtcdAddrs      []string
	BasePath       string
	Name           string
	UpdateInterval uint32
}

func CreateAndStartMicroServer(conf MicroServerConfig, svr interface{}) {
	server := rpcx.NewServer()
	server.ServerCodecFunc = codec.NewJSONRPC2ServerCodec
	// 注册rpcx.NewServer
	rplugin := &plugin.EtcdRegisterPlugin{
		ServiceAddress: fmt.Sprintf("%s@%s", conf.Scheme, conf.Addr),
		EtcdServers:    conf.EtcdAddrs,
		BasePath:       conf.BasePath,
		Metrics:        metrics.NewRegistry(),
		Services:       []string{conf.Name},
		UpdateInterval: time.Duration(conf.UpdateInterval) * time.Second,
	}
	rplugin.Start()
	server.PluginContainer.Add(rplugin)
	server.PluginContainer.Add(plugin.NewMetricsPlugin())
	server.PluginContainer.Add(msplugin.NewServerPostReadRequestPlugin())
	server.RegisterName(conf.Name, svr, "weight=1&m=devops")
	server.Serve(conf.Scheme, conf.Addr)
}
