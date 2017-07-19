package utils

import (
	"common/config"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/go-sql-driver/mysql"
)

var (
	DB_QUERY_WHERE_EMPTY = errors.New("where is empty")
	DB_RECORD_NOT_EXISTS = errors.New("record is not exists")
	DB_DUPLICATE         = errors.New("database duplicate")
)

func RecordExists(db orm.Ormer, ptrStructOrTableName interface{}, where map[string]interface{}) bool {
	query := db.QueryTable(ptrStructOrTableName)
	for k, v := range where {
		query = query.Filter(k, v)
	}
	return query.Exist()
}

func RecordUpdate(db orm.Ormer, ptrStructOrTableName interface{}, where map[string]interface{}, updates map[string]interface{}) error {
	if where == nil || len(where) == 0 {
		return DB_QUERY_WHERE_EMPTY
	}
	// 获取表名
	tableName := ""
	if table, ok := ptrStructOrTableName.(string); ok {
		tableName = snakeString(table)
	} else {
		tableName = getTableName(reflect.ValueOf(ptrStructOrTableName))
	}
	update := db.QueryTable(tableName)
	// 组装where条件
	for k, v := range where {
		update = update.Filter(k, v)
	}
	// 执行数据库更新
	num, err := update.Update(orm.Params(updates))
	if err != nil {
		Error("", "update %s table failed. where:%+v, updates:%+v, error: %s", tableName, where, updates, err)
		if IsDBDuplicate(err) {
			return DB_DUPLICATE
		}
		return err
	}
	if num == 0 {
		if !RecordExists(db, tableName, where) {
			Error("", "update %s table failed. where %+v, updates:%+v, error:%s", tableName, where, updates, DB_RECORD_NOT_EXISTS)
			return DB_RECORD_NOT_EXISTS
		}
	}
	return nil
}

// =============== 从beego orm中获取的函数 ================
// get table name. method, or field name. auto snaked.
func getTableName(val reflect.Value) string {
	ind := reflect.Indirect(val)
	fun := val.MethodByName("TableName")
	if fun.IsValid() {
		vals := fun.Call([]reflect.Value{})
		if len(vals) > 0 {
			val := vals[0]
			if val.Kind() == reflect.String {
				return val.String()
			}
		}
	}
	return snakeString(ind.Type().Name())
}

// snake string, XxYy to xx_yy
func snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// =============== 从beego orm中获取的函数 ================

//MYSQL错误码
const (
	MYSQL_CODE_DUPLICATE uint16 = 1062
)

// 判断是否是MySQL冲突
func IsDBDuplicate(err error) bool {
	if dberr, ok := err.(*mysql.MySQLError); ok {
		return dberr.Number == MYSQL_CODE_DUPLICATE
	}
	return false
}

func InitDB(conf *config.MysqlConfig) {
	dataSource := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=%s`, conf.User, conf.Password, conf.Addr, conf.DB, conf.Charset)
	if err := orm.RegisterDriver("mysql", orm.DRMySQL); err != nil {
		fmt.Printf("register mysql driver failed. error:%s\n", err)
		os.Exit(1)
	}
	if err := orm.RegisterDataBase("default", "mysql", dataSource, conf.MaxIdle, conf.MaxConn); err != nil {
		fmt.Printf("register mysql data base failed. error:%s", err)
		os.Exit(1)
	}
}
