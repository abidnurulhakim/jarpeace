package helper

import (
	"math/rand"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func GetDecoder(result interface{}) (*mapstructure.Decoder, error) {
	return mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:          "json",
		Result:           result,
		WeaklyTypedInput: false})
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ConvertStringToSlice(value interface{}) []string {
	var newValue []string
	if reflect.TypeOf(value).Kind() != reflect.Slice {
		newValue = []string{value.(string)}
	} else {
		for _, val := range value.([]interface{}) {
			newValue = append(newValue, val.(string))
		}
	}
	return newValue
}

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomString(length int) string {
	return StringWithCharset(length, charset)
}

func RandInt(min int, max int) int {
	for true {
		i := seededRand.Intn(max)
		if i >= min {
			return i
		}
	}
	return 0
}
