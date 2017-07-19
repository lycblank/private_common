/*
 Package xxx

 Copyright(c) 2000-2016 Quqian Technology Company Limit.
 All rights reserved.
 创建人：  longyongchun
 创建日期：2016-11-15 13:48:46
 修改记录：
 ----------------------------------------------------------------
 修改人        |  修改日期    				|    备注
 ----------------------------------------------------------------
 longyongchun  | 2016-11-15 13:48:46		|    创建文件
 ----------------------------------------------------------------
*/

package storage

import "github.com/astaxie/beego/orm"

type cacheTask struct {
	execTask     Task
	rollbackTask Task
}

func (ct *cacheTask) Exec(db orm.Ormer) error {
	if ct.execTask != nil {
		return ct.execTask(db)
	}
	return nil
}

func (ct *cacheTask) Rollback(db orm.Ormer) error {
	if ct.rollbackTask != nil {
		return ct.rollbackTask(db)
	}
	return nil
}

func GenerateCacheTask(execTask Task, rollbackTask Task) DbaTasker {
	return &cacheTask{
		execTask:     execTask,
		rollbackTask: rollbackTask,
	}
}
