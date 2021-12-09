package js

import (
	jsoniter "github.com/json-iterator/go"
)

func IdentJson(d interface{}) string {
	out, _ := jsoniter.MarshalIndent(d, "", "  ")
	return string(out)
}

func MinifyJson(d interface{}) string {
	s, _ := jsoniter.MarshalToString(d)
	return s
}
