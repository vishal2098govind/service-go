package parser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// expands each field if root struct has nested struct
// example:
/*
	type T struct {
		A string
		B struct {
			C string
		}
	}

	Field will be computed for each of A and C,
	Field for A: Field{Name:"A", Key: []string{"A"}}
	Field for C: Field{Name:"C", Key: []string{"B", "C"}} -> can be referenced as B_C
*/
func ExtractFields(inp interface{}) ([]Field, error) {
	val := reflect.ValueOf(inp)
	if val.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("t must be of type struct")
	}
	val = val.Elem()
	typ := val.Type()

	fields := []Field{}
	for i := range val.NumField() {
		v := val.Field(i)
		t := typ.Field(i)
		field := Field{
			Name: t.Name,
			Key:  []string{},
		}

		field.Key = append(field.Key, t.Name)

		if t.Type.Kind() == reflect.Struct {
			nestedVal := v.Addr().Interface()
			nestedFields, err := ExtractFields(nestedVal)
			if err != nil {
				return nil, err
			}
			field.Name = nestedFields[len(nestedFields)-1].Name
			for _, v := range nestedFields {
				field.Key = append(field.Key, v.Key...)
			}
		}

		if t.Type.Kind() != reflect.Struct {
			confTag := t.Tag.Get("conf")
			if strings.Trim(confTag, " ") != "" {
				confs := strings.SplitN(confTag, ":", 2)
				if len(confs) == 2 {
					switch confs[0] {
					case "default":
						fmt.Printf("default for %v:%v\n", field.Name, confs[1])
						switch t.Type.Kind() {
						case reflect.String:
							v.SetString(confs[1])
						case reflect.Int:
							in, err := strconv.ParseInt(confs[1], 0, t.Type.Bits())
							if err != nil {
								return nil, fmt.Errorf("invalid default value: %w", err)
							}
							v.SetInt(in)
						}
					}
				}
			}
		}

		fields = append(fields, field)
	}

	return fields, nil
}

// represents each field of root struct
type Field struct {
	Name string
	Key  []string
	Tag  string
}
