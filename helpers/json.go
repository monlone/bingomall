package helper

import "encoding/json"

func Json(v interface{}) string {
	temp, _ := json.Marshal(v)
	return string(temp)
}
