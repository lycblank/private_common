package httpin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func FillRequireParam(ctx *context.Context, param interface{}, logID ...interface{}) {
	v := reflect.ValueOf(param)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		SystemError(500, "query is not struct")
	}

	numField := v.NumField()
	for i := 0; i < numField; i++ {
		fillParam(ctx, v.Type().Field(i), v.Field(i))
	}

	if beego.BConfig.RunMode == "dev" {
		datas, _ := json.Marshal(param)
		var logid interface{}
		if logID != nil && len(logID) > 0 {
			logid = logID[0]
		}
		Debug(logid, "required:%s", string(datas))
	}
}

func fillParam(ctx *context.Context, field reflect.StructField, val reflect.Value) {
	//fmt.Println(ctx.Input.Context.Request.Header)
	switch field.Tag.Get("in") {
	case "header":
		fillHeaderParam(ctx, field, val)
	case "body":
		fillBodyParam(ctx, field, val)
	case "path":
		fillPathParam(ctx, field, val)
	case "query":
		fillQueryParam(ctx, field, val)
	default:
		if !val.IsValid() || !val.CanAddr() || !val.CanSet() {
			SystemError(500, "value can not set")
		}
		FillRequireParam(ctx, val.Addr().Interface())
	}
}

func fillBodyParam(ctx *context.Context, field reflect.StructField, val reflect.Value) {
	if !val.IsValid() || !val.CanAddr() || !val.CanSet() {
		SystemError(500, "value can not set")
	}

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	body := ctx.Input.RequestBody
	switch field.Tag.Get("decode") {
	case "base64":
		body, _ = base64.StdEncoding.DecodeString(string(body))
	}

	if err := json.Unmarshal(body, val.Addr().Interface()); err != nil {
		BadParamer(500, fmt.Sprintf("%s is not json error:%s", "body", err))
	}

	CheckBodyFiled(val.Addr().Interface())
}

func CheckBodyFiled(param interface{}, logID ...interface{}) {
	v := reflect.ValueOf(param)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		SystemError(500, "body is not struct")
	}
	numField := v.NumField()
	for i := 0; i < numField; i++ {
		if v.Field(i).Kind() == reflect.Struct { //递归检查
			CheckBodyFiled(v.Field(i).Addr().Interface())
		}
		field := v.Type().Field(i)
		val := v.Field(i)
		if reflect.DeepEqual(val.Interface(), reflect.Zero(val.Type()).Interface()) {
			if field.Tag.Get("required") == "true" {
				jsonTag := field.Tag.Get("json")
				if jsonTag == "" {
					jsonTag = field.Name
				}
				def := field.Tag.Get("default")
				fillValue(v.Field(i), def, true, jsonTag)
			}
		}
	}
}

func fillHeaderParam(ctx *context.Context, field reflect.StructField, val reflect.Value) {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		// 没有json标记，取变量的名称
		jsonTag = field.Name
	}
	var v string
	if v = ctx.Input.Header(jsonTag); v == "" {
		// 尝试将第一个字母转换为大写 测试环境使用 beego Swagger的bug
		firstUpperJsonTag := strings.ToUpper(jsonTag[0:1]) + string(jsonTag[1:])
		if v = ctx.Input.Header(firstUpperJsonTag); v == "" {
			v = field.Tag.Get("default")
		}
	}
	required := field.Tag.Get("required") == "true"
	if v == "" && required {
		BadParamer(400, fmt.Sprintf("%s is not exists in header", jsonTag))
	}
	fillValue(val, v, required, jsonTag)
}

func fillPathParam(ctx *context.Context, field reflect.StructField, val reflect.Value) {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		// 没有json标记，取变量的名称
		jsonTag = field.Name
	}
	var v string
	if v = ctx.Input.Param(fmt.Sprintf(":%s", jsonTag)); v == "" {
		v = field.Tag.Get("default")
	}
	required := field.Tag.Get("required") == "true"
	if v == "" && required {
		BadParamer(400, fmt.Sprintf("%s is not exists in path", jsonTag))
	}
	fillValue(val, v, required, jsonTag)
}

func fillQueryParam(ctx *context.Context, field reflect.StructField, val reflect.Value) {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		// 没有json标记，取变量的名称
		jsonTag = field.Name
	}
	var v string
	if v = ctx.Input.Query(jsonTag); v == "" {
		v = field.Tag.Get("default")
	}
	required := field.Tag.Get("required") == "true"
	if v == "" && required {
		BadParamer(400, fmt.Sprintf("%s is not exists in query", jsonTag))
	}
	fillValue(val, v, required, jsonTag)
}

func fillValue(val reflect.Value, requestVal string, required bool, jsonTag string) {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	kind := val.Kind()
	switch {
	case kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64:
		if requestVal == "" && !required {
			requestVal = "0"
		}
		v, err := strconv.ParseUint(requestVal, 10, 64)
		if err != nil {
			BadParamer(400, fmt.Sprintf("%s is not uint", jsonTag))
		}
		if !val.IsValid() || !val.CanSet() {
			SystemError(500, "value can not set")
		}
		val.SetUint(v)
	case kind == reflect.Int:
		if requestVal == "" && !required {
			requestVal = "0"
		}
		v, err := strconv.ParseInt(requestVal, 10, 64)
		if err != nil {
			BadParamer(400, fmt.Sprintf("%s is not uint", jsonTag))
		}
		if !val.IsValid() || !val.CanSet() {
			SystemError(500, "value can not set")
		}
		val.SetInt(v)
	case kind == reflect.String:
		val.SetString(requestVal)
	default:
		BadParamer(400, "type is not support")
	}
}
