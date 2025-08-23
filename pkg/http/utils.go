package http

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
	"strings"
)

func BindData(ctx *gin.Context, obj any) error {
	if ctx.Request.ContentLength > 0 && ctx.GetHeader("Content-Type") == "application/json" {
		if err := ctx.ShouldBindJSON(obj); err != nil {
			return err
		}

		SanitizeStruct(obj)
		return nil
	}

	if err := ctx.ShouldBind(obj); err != nil {
		return err
	}

	return nil
}

func BindParams[T any](ctx *gin.Context, key string) (T, error) {
	var res T
	var value string

	if key == "id" {
		value = ctx.Param(key)
		if value == "" {
			return res, errors.New("id is required")
		}

		switch any(res).(type) {
		case string:
			res = any(value).(T)
		case int64:
			id, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return res, err
			}
			res = any(id).(T)
		default:
			return res, fmt.Errorf("unsupported type")
		}
	} else if key == "ids" {
		value = ctx.Query(key)
		if value == "" {
			return res, errors.New("ids are required")
		}
		res = any(strings.Split(value, ",")).(T)
	} else {
		return res, fmt.Errorf("key is not supported")
	}

	return res, nil
}

func SanitizeStruct(payload interface{}) {
	v := reflect.ValueOf(payload).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String {
			sanitizedValue := strings.TrimSpace(field.String())
			field.SetString(sanitizedValue)
		}
	}
}
