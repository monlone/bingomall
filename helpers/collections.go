package helper

import (
	"bingomall/helpers/convention"
	"reflect"
)

func ModelObjectToSlice(inArr interface{}, key string) (ret []uint64) {
	v := reflect.ValueOf(inArr)
	var array = make(map[string]string)
	if v.Kind() == reflect.Slice {
		l := v.Len()
		for i := 0; i < l; i++ {
			one := v.Index(i).Interface()
			vo := reflect.ValueOf(one)
			var value string
			if vo.Kind() == reflect.Ptr {
				value = vo.Elem().FieldByName(key).String()
			} else {
				value = vo.FieldByName(key).String()
			}

			if len(value) > 0 {
				if _, ok := array[value]; !ok {
					array[value] = value
				}
			}
		}
		for _, val := range array {
			ret = append(ret, convention.StringToUint64(val))
		}
	}
	return ret
}
