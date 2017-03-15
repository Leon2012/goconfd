package reflect

import (
	"reflect"
)

func GetTag(s interface{}, fieldName, tagName string) string {
	t := reflect.TypeOf(s)
	field, b := t.Elem().FieldByName(fieldName)
	if !b {
		return ""
	}
	tag := field.Tag.Get(tagName)
	return tag
}
