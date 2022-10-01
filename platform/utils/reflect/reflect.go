package reflect

import (
	"reflect"
	"strings"
)

func GetJSONFieldNames(model interface{}) []string {
	var names []string
	v := reflect.TypeOf(model)
	if v.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < v.NumField(); i++ {
		{
			// FIXME:
			// added this to not include populated arrays. This is not enough though.
			// trying to filter values that are from other related tables would not work
			// some kind of stable way should be implemented
			typeOfV := reflect.TypeOf(v.Field(i)).Kind()
			if typeOfV == reflect.Slice {
				continue
			}
		}
		jsonName := strings.Split(v.Field(i).Tag.Get("json"), ",")[0]
		if jsonName == "" {
			// construct construct fieldName
			fieldNameRunes := []rune(v.Field(i).Name)
			result := []rune{[]rune(strings.ToLower(string(fieldNameRunes[0])))[0]}
			for _, v := range fieldNameRunes[1:] {
				if v == []rune(strings.ToUpper(string(v)))[0] {
					result = append(result, '_', []rune(strings.ToLower(string(v)))[0])
				} else {
					result = append(result, v)
				}
			}
			jsonName = string(result)
		}
		names = append(names, jsonName)
	}
	return names
}
