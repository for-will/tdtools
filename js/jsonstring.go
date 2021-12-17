package js

import (
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func IdentJson(d interface{}) string {
	out, _ := jsoniter.MarshalIndent(d, "", "  ")
	return string(out)
}

func MinifyJson(d interface{}) string {
	s, _ := jsoniter.MarshalToString(d)
	return s
}

func PbMinifyJson(m proto.Message) string {
	jsb, _ := protojson.MarshalOptions{}.Marshal(m)
	return string(jsb)
}
