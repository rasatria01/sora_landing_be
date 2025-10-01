package utils

import (
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"sora_landing_be/pkg/storage"
	"strings"
	"time"
	"unicode"

	"github.com/segmentio/ksuid"
)

func Fallback[T any, U bool | func() bool](option1, option2 T, terms U) T {
	var termsType = reflect.TypeOf(terms)
	var result = option2

	switch termsType.Kind() {
	case reflect.Bool:
		if reflect.ValueOf(terms).Bool() {
			result = option1
		}
	case reflect.Func:
		if val := reflect.ValueOf(terms).Call([]reflect.Value{}); val[0].Bool() {
			result = option1
		}
	default:
		panic("unhandled default case")
	}
	return result
}

func SafelyDereference[T any](input *T) (output T) {
	var res T
	if input == nil {
		return res
	}
	return *input
}

func SafelyReference[T any](input T) (output *T) {
	return &input
}

func ToSnakeCase(str string) string {
	var charArr []string

	for i, char := range str {
		if i > 0 && unicode.IsUpper(rune(str[i-1])) && unicode.IsUpper(rune(str[i])) {
			charArr = append(charArr, string(char))
			continue
		}

		if i > 0 && unicode.IsUpper(rune(str[i])) {
			charArr = append(charArr, "_")
		}

		if i > 1 && unicode.IsUpper(rune(str[i-1])) && unicode.IsUpper(rune(str[i-2])) && unicode.IsLower(rune(str[i])) {
			charArr, _ = InsertAtIndex(charArr, i-1, "_")
		}

		charArr = append(charArr, string(char))
	}

	return strings.ToLower(strings.Join(charArr, ""))
}

func InsertAtIndex[T any](slice []T, index int, value T) ([]T, error) {
	if index < 0 || index > len(slice) {
		return nil, fmt.Errorf("index out of range")
	}

	newSlice := make([]T, len(slice)+1)
	copy(newSlice, slice[:index])
	newSlice[index] = value
	copy(newSlice[index+1:], slice[index:])

	return newSlice, nil
}

func MaskingString(str string, limitWordsShow ...int) string {
	var (
		wordLen   = len(str)
		limitWord = limitWordsShow[0]
	)
	if wordLen == 0 {
		return ""
	}

	if limitWord == 0 {
		limitWord = wordLen
	}

	var result string
	showLen := int(math.Ceil(0.1 * float64(wordLen)))

	for i := 0; i < limitWord; i++ {
		if i > showLen {
			result += "*"
			continue
		}
		result += string(str[i])
	}
	return result
}

type DomainMapping[T any] interface {
	ToDomain(parentID string) T
}

type DomainMappingWithoutParent[T any] interface {
	ToDomain() T
}

func ToDomainArray[T, U any](parentID *string, items []T) []U {
	domainItems := make([]U, 0, len(items))
	for i, val := range items {
		if impl, ok := reflect.Indirect(reflect.ValueOf(val)).Interface().(DomainMapping[U]); ok && parentID != nil {
			domainItems = slices.Insert(domainItems, i, impl.ToDomain(*parentID))
		} else if impl, ok := reflect.Indirect(reflect.ValueOf(val)).Interface().(DomainMappingWithoutParent[U]); ok {
			domainItems = slices.Insert(domainItems, i, impl.ToDomain())
		} else {
			return nil
		}
	}
	return domainItems
}

func ToDeletedObject[T any](requests, record []T) []string {
	var (
		deletedIDs []string
		reqMap     = make(map[string]bool)
	)

	for _, t := range requests {
		id := reflect.Indirect(reflect.ValueOf(t)).FieldByName("ID").String()
		reqMap[id] = true
	}

	for _, request := range record {
		id := reflect.Indirect(reflect.ValueOf(request)).FieldByName("ID").String()
		if id != "" && !reqMap[id] {
			deletedIDs = append(deletedIDs, id)
		}
	}

	return deletedIDs
}

func IsFileSizeExceedMb(fileSize int64, maxSizeMb int64) bool {
	return fileSize > maxSizeMb*1000*1000
}

func IsDocumentFile(fileName string) (bool, []string, string) {
	extension := filepath.Ext(fileName)
	return slices.Contains(storage.DocumentFileType, extension), storage.DocumentFileType, extension
}

func GenerateKeyFile(fileName string) string {
	fileName, ext := strings.TrimSuffix(fileName, filepath.Ext(fileName)), filepath.Ext(fileName)
	return fmt.Sprintf("%s_%s%s", fileName, ksuid.New().String(), ext)
}

func GetCurrentTimeBasedOnLocation(timezone string) (time.Time, error) {
	currentTime, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}

	return time.Now().In(currentTime), nil
}

func GetDiffDataByField[T comparable](data1, data2 any, fields ...string) ([]T, error) {
	var (
		data1Type     = reflect.Indirect(reflect.ValueOf(data1))
		data2Type     = reflect.Indirect(reflect.ValueOf(data2))
		data1ValueMap = make(map[string]T)
		res           = make([]T, 0, len(fields))
	)

	if data1Type.Kind() != data2Type.Kind() {
		return nil, errors.New("type not same")
	}

	for _, field := range fields {
		data1ValueMap[field] = data1Type.FieldByName(field).Interface().(T)
	}

	for _, field := range fields {
		data2Value := data2Type.FieldByName(field).Interface()
		if data1ValueMap[field] != data2Value {
			res = append(res, data1ValueMap[field])
		}
	}

	return res, nil
}

func GetFileName(fullName string) string {
	return regexp.MustCompile(`_[^_.]*\.`).ReplaceAllString(fullName, ".")
}

func Contains[T comparable](arr []T, str T) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}
