package main

import (
	"fmt"
	"reflect"

	tagdebug "github.com/vishal2098govind/service/cmd/learn-reflect/tag-debug"
)

func debugTag(i interface{}) {
	val := reflect.ValueOf(i)
	typ := reflect.TypeOf(i)
	if val.Kind() == reflect.Struct {
		n := val.NumField()
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			fmt.Printf("field: %v\n", f.Name)
			fmt.Printf("field type: %v\n", f.Type.Name())
			fmt.Printf("field val: %v\n", val)
			fmt.Printf("tag: %v\n", f.Tag)
			conf := f.Tag.Get("conf")
			fmt.Printf("conf: %v\n", conf)
		}
	}
	fmt.Println("----")
}

func debugReflect(i interface{}) {
	t := reflect.TypeOf(i)
	fmt.Println("TypeOf(i): ", t)

	val := reflect.ValueOf(i)
	fmt.Println("ValueOf(i): ", val)

	if t.Kind() == reflect.Slice {
		for i := 0; i < val.Len(); i++ {
			fmt.Printf("val[%d]: %v\n", i, val.Index(i))
		}
	} else if t.Kind() == reflect.Map {
		iter := val.MapRange()
		for iter.Next() {
			k, v := iter.Key(), iter.Value()
			fmt.Printf("val[%v]: %v\n", k, v)
		}
	} else {
		fmt.Printf("val: %v\n", val)
	}
	fmt.Println("------")
}

func main() {
	var (
		num = 10
		str = "hello"
		slc = []int{1, 2, 3}
		mp  = map[string]int{"one": 1, "two": 2, "three": 3}
	)

	debugReflect(num)
	debugReflect(str)
	debugReflect(slc)
	debugReflect(mp)

	type conf struct {
		apiHost string `conf:"test=123"`
	}
	t := conf{}
	debugTag(t)
	debugTag(tagdebug.Conf{})
	debugTag(tagdebug.ConfPublic{
		ApiHost: "ApiHost",
	})
}
