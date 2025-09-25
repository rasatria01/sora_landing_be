package http

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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

	switch key {
	case "id":
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
			return res, fmt.Errorf("unsupported type for id")
		}

	case "slug":
		value = ctx.Param(key)
		if value == "" {
			return res, errors.New("slug is required")
		}
		// slug only makes sense as a string
		switch any(res).(type) {
		case string:
			res = any(value).(T)
		default:
			return res, fmt.Errorf("unsupported type for slug (must be string)")
		}

	case "ids":
		value = ctx.Query(key)
		if value == "" {
			return res, errors.New("ids are required")
		}
		res = any(strings.Split(value, ",")).(T)

	default:
		return res, fmt.Errorf("key %q is not supported", key)
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
