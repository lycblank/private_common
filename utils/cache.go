/*
 Package xxx

 Copyright(c) 2000-2016 Quqian Technology Company Limit.
 All rights reserved.
 创建人：  longyongchun
 创建日期：2016-09-30 10:19:05
 修改记录：
 ----------------------------------------------------------------
 修改人        |  修改日期    				|    备注
 ----------------------------------------------------------------
 longyongchun  | 2016-09-30 10:19:05		|    创建文件
 ----------------------------------------------------------------
*/

package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/astaxie/beego/config"

	"common/cache"
)

var Cache *cache.Cache

func InitCache(conf *cache.RedisConfig) {
	Cache = new(cache.Cache)
	if Cache.RedisPool(*conf) == nil {
		val, _ := json.Marshal(*conf)
		Error("", "init cache error. redisConfig: %s", string(val))
		fmt.Printf("init cache error. redisConfig: %s\n\n", string(val))
		os.Exit(1)
	}
}

func Init(iniConf config.Configer, runmode string) {
	redisConfig := cache.RedisConfig{
		Network:  iniConf.String(runmode + "::network"),
		Address:  iniConf.String(runmode + "::address"),
		Password: iniConf.String(runmode + "::password"),
		DB:       iniConf.String(runmode + "::db"),
	}
	redisConfig.ConnectTimeout, _ = iniConf.Int(runmode + "::connecttimeout")
	redisConfig.ReadTimeout, _ = iniConf.Int(runmode + "::readtimeout")
	redisConfig.WriteTimeout, _ = iniConf.Int(runmode + "::writetimeout")
	redisConfig.MaxActive, _ = iniConf.Int(runmode + "::maxactive")
	redisConfig.MaxIdle, _ = iniConf.Int(runmode + "::maxidle")
	redisConfig.IdleTimeout, _ = iniConf.Int(runmode + "::idletimeout")
	redisConfig.Wait, _ = iniConf.Bool(runmode + "::wait")
	Cache = new(cache.Cache)
	if Cache.RedisPool(redisConfig) == nil {
		val, _ := json.Marshal(redisConfig)
		Error("", "init cache error. redisConfig: %s", string(val))
		fmt.Printf("init cache error. redisConfig: %s\n\n", string(val))
		os.Exit(1)
	}
}
