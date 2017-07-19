/*
 Package xxx

 Copyright(c) 2000-2016 Quqian Technology Company Limit.
 All rights reserved.
 创建人：  longyongchun
 创建日期：2016-09-13 12:51:37
 修改记录：
 ----------------------------------------------------------------
 修改人        |  修改日期    				|    备注
 ----------------------------------------------------------------
 longyongchun  | 2016-09-13 12:51:37		|    创建文件
 ----------------------------------------------------------------
*/

package utils

import (
	"common/httpin"
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"github.com/smallnest/rpcx/core"

	bcontext "github.com/astaxie/beego/context"
)

var log = logs.NewLogger(10000)
var accessLog = logs.NewLogger(10000)

var logModel string

func InitLog(svrName string, runmode string) *logs.BeeLogger {
	if err := log.SetLogger("file", fmt.Sprintf(`{"filename":"logs/%s.log"}`, svrName)); err != nil {
		fmt.Printf("set logger file failed. error:%s\n\n", err)
		os.Exit(1)
	}
	log.EnableFuncCallDepth(true)
	// 设置等级为3 log信息需要在封一层
	log.SetLogFuncCallDepth(4)
	log.Async()

	if err := accessLog.SetLogger("file", fmt.Sprintf(`{"filename":"logs/%s-access.log"}`, svrName)); err != nil {
		fmt.Printf("set logger file failed. error:%s\n\n", err)
		os.Exit(1)
	}
	accessLog.Async()

	SetRPCAccessLogger(accessLog)
	httpin.SetLogger(log)

	orm.DebugLog = log
	redis.DebugLog = log

	if runmode == "prod" {
		logModel = "info"
		log.SetLevel(logs.LevelInformational)
		accessLog.SetLevel(logs.LevelInformational)
	}
	if runmode == "dev" {
		logModel = "debug"
		orm.Debug = true
		redis.Debug = true
	}
	return log
}

func SetLogger(loger *logs.BeeLogger) {
	log = loger
}

func Info(v interface{}, format string, args ...interface{}) {
	logString := GetLogIDString(v)
	if log != nil {
		log.Info(logString+format, args...)
	}
}

func Error(v interface{}, format string, args ...interface{}) {
	logString := GetLogIDString(v)
	if log != nil {
		// 打印堆栈信息
		buf := make([]byte, 1024)
		size := runtime.Stack(buf, true)
		log.Error(logString+"%s", string(buf[:size]))

		log.Error(logString+format, args...)
	}
}

func Debug(v interface{}, format string, args ...interface{}) {
	logString := GetLogIDString(v)
	if log != nil {
		log.Debug(logString+format, args...)
	}
}

func Warning(v interface{}, format string, args ...interface{}) {
	logString := GetLogIDString(v)
	if log != nil {
		log.Warning(logString+format, args...)
	}
}

func GetLogIDString(arg interface{}) string {
	var values map[string][]string
	if ctx, ok := arg.(context.Context); ok {
		if values = GetRpcHeadersFromReq(ctx); values == nil {
			return ""
		} else {
			if v, ok := core.FromContext(ctx); ok {
				values = map[string][]string(v)
			}
		}
	}
	if ctx, ok := arg.(*bcontext.Context); ok {
		if values = GetRpcHeaders(ctx, false); values == nil {
			return ""
		}
	}
	if values != nil {
		traceId := ""
		if ids := values["trace_id"]; ids != nil && len(ids) > 0 {
			traceId = ids[0]
		}
		spanId := ""
		if ids := values["span_id"]; ids != nil && len(ids) > 0 {
			spanId = ids[0]
		}
		return fmt.Sprintf("|%s|%s|", traceId, spanId)
	}
	return ""
}

type ElasticErrorLog struct {
}

func (errLog *ElasticErrorLog) Printf(format string, v ...interface{}) {
	Error("", format, v...)
}

type ElasticInfoLog struct {
}

func (infoLog *ElasticInfoLog) Printf(format string, v ...interface{}) {
	Info("", format, v...)
}

type ElasticDebugLog struct {
}

func (debugLog *ElasticDebugLog) Printf(format string, v ...interface{}) {
	Debug("", format, v...)
}
