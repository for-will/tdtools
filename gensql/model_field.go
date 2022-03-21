package main

import (
	"github.com/ahmetb/go-linq/v3"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type ModelField struct {
	Name string
	Type string
	Tag  reflect.StructTag
}

func (f ModelField) SqlName() string {
	return snakeCase(f.Name)
}

func (f ModelField) SqlType() string {

	switch f.Type {
	case "bool":
		return "bool"
	case "int32":
		return "int"
	case "time.Time":
		return "timestamp"
	default:
		log.Fatalf("unknown SqlType: %#v", f)
	}
	return "bad"
}

func (f ModelField) SqlDefault() string {
	switch f.Type {
	case "time.Time":
		return "current_timestamp()"
	default:
		return ""
	}
}

type IndexDesc struct {
	Index  string
	Keys   []*IndexKey
	Unique bool
}

type IndexKey struct {
	KeyName string
	Order   int
}

func modelIndies(m *Model) []*IndexDesc {

	var indiesMap = map[string]*IndexDesc{}
	for _, field := range m.Fields {

		// 索引&唯一索引
		indies := haveIndex(m.DbTableName(), field)
		for _, v := range indies {
			if idx, ok := indiesMap[v.Index]; ok {
				idx.Keys = append(idx.Keys, v.Keys...)
			} else {
				indiesMap[v.Index] = v
			}
		}
	}

	var indexList []*IndexDesc
	linq.From(indiesMap).SelectT(func(value linq.KeyValue) interface{} {
		return value.Value
	}).OrderByT(func(i *IndexDesc) interface{} {
		return i.Index
	}).ToSlice(&indexList)

	for _, idx := range indexList {
		linq.From(idx.Keys).OrderByT(func(k *IndexKey) int {
			return k.Order
		}).ToSlice(&idx.Keys)
	}
	return indexList
}

func haveIndex(tbn string, f *ModelField) (indies []*IndexDesc) {
	tag := f.Tag.Get("db")

	words := strings.Split(tag, ",")
	for _, word := range words {
		kv := strings.Split(word, ":")
		idx := &IndexDesc{
			Unique: false,
		}
		if kv[0] == "index" {
			if len(kv) == 1 {
				idx.Index = "idx_" + tbn + "_" + f.SqlName()
			} else {
				idx.Index = kv[1]
			}
		} else if kv[0] == "unique" {
			idx.Unique = true
			if len(kv) == 1 {
				idx.Index = "uni_" + tbn + "_" + f.SqlName()
			} else {
				idx.Index = kv[1]
			}
		} else {
			continue
		}

		key := &IndexKey{
			KeyName: f.SqlName(),
			Order:   0,
		}
		if len(kv) == 3 {
			key.Order, _ = strconv.Atoi(kv[2])
		}
		idx.Keys = append(idx.Keys, key)
		indies = append(indies, idx)
	}
	return indies
}
