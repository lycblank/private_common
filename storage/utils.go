/*
 Package xxx

 Copyright(c) 2000-2016 Quqian Technology Company Limit.
 All rights reserved.
 创建人：  longyongchun
 创建日期：2016-11-15 12:59:47
 修改记录：
 ----------------------------------------------------------------
 修改人        |  修改日期    				|    备注
 ----------------------------------------------------------------
 longyongchun  | 2016-11-15 12:59:47		|    创建文件
 ----------------------------------------------------------------
*/

package storage

import (
	"common/httpin"

	"github.com/astaxie/beego/orm"
)

type DbaTasker interface {
	Exec(db orm.Ormer) error
	Rollback(db orm.Ormer) error
}

type Task func(db orm.Ormer) error

func (t *Task) Exec(db orm.Ormer) error {
	return (*t)(db)
}

func (t *Task) Rollback(db orm.Ormer) error {
	return nil
}

func GenerateDbaTask(task Task) DbaTasker {
	return &task
}

func ExecTrans(tasks ...DbaTasker) error {
	mysqlORM := orm.NewOrm()
	successTask := make([]DbaTasker, 0)
	isCommit := false
	defer func() {
		// 捕获异常，并执行回滚操作
		if v := recover(); v != nil || !isCommit {
			for _, task := range successTask {
				task.Rollback(mysqlORM)
			}
			mysqlORM.Rollback()
			//向上层抛出异常
			if v != nil {
				panic(v)
			}
		}
	}()

	if err := mysqlORM.Begin(); err != nil {
		httpin.SystemError(500, "db error")
	}

	for _, task := range tasks {
		if task != nil {
			if err := task.Exec(mysqlORM); err != nil {
				if err == EXIT_QUERY {
					return err
				}
				panic(err)
			}
		}
		successTask = append(successTask, task)
	}

	mysqlORM.Commit()
	isCommit = true
	return nil
}
