package sanitizer

import (
	"github.com/microcosm-cc/bluemonday"
	"reflect"
)

func SanitizeData(v interface{}) interface{} {
	sanitizer := bluemonday.UGCPolicy()
	vv := reflect.ValueOf(v)

	for i := 0; i < vv.Elem().NumField(); i++  {
		if vv.Elem().Field(i).Kind() == reflect.String {
			vv.Elem().Field(i).SetString(sanitizer.Sanitize(vv.Elem().Field(i).String()))
		}
	}

	return vv.Elem().Interface()
}
