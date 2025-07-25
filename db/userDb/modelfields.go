package userDb

import (
	"github.com/helays/utils/v2/map/syncMapWrapper"
	"github.com/helays/utils/v2/tools"
	"reflect"
)

var modelFieldsCache = &syncMapWrapper.SyncMap[string, *modelFieldTypes]{}

type fieldTypes struct {
	jsonTagName string
	defaultVal  string
	fieldName   string
	dblike      string
	kind        reflect.Kind
}

type modelFieldTypes struct {
	tableName string
	fieldsMap map[string]fieldTypes //
	fields    []string              // 字段顺序
}

func getModelFields(model any, alias string) *modelFieldTypes {
	if model == nil {
		return nil
	}
	structName, isStruct, isAnonymous, mt := tools.IsStruct(model)
	if !isStruct || isAnonymous {
		return nil
	}

	modelFields, _ := modelFieldsCache.LoadOrStoreFunc(structName, func() (*modelFieldTypes, bool) {
		v := reflect.ValueOf(model)
		_modelFields := autoGetStructFieldJsonTag(mt, v)
		_modelFields.tableName = alias
		if _modelFields.tableName == "" {
			tbName := v.MethodByName("TableName")
			if tbName.IsValid() {
				_modelFields.tableName = tbName.Call([]reflect.Value{})[0].String()
			} else {
				_modelFields.tableName = tools.SnakeString(structName)
			}
		}
		return &_modelFields, true
	})
	return modelFields
}

func autoGetStructFieldJsonTag(mt reflect.Type, mv reflect.Value) modelFieldTypes {
	type retypes struct {
		mt reflect.Type
		mv reflect.Value
	}

	fieldCache := modelFieldTypes{
		fieldsMap: make(map[string]fieldTypes),
		fields:    nil,
	}

	stack := []retypes{{mt: mt, mv: mv}}
	for len(stack) > 0 {
		item := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		t := item.mt
		if t.Kind() != reflect.Struct {
			continue
		}
		v := item.mv
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tagName := field.Tag.Get("json")
			kind := field.Type.Kind()
			// 处理嵌套结构体
			if kind == reflect.Struct && field.Tag.Get("gorm") == "" && tagName == "" {
				nextElement := v.Field(i).Interface()
				stack = append(stack, retypes{
					mt: reflect.TypeOf(nextElement),
					mv: reflect.ValueOf(nextElement),
				})
				continue
			}
			if tagName == "" {
				continue
			}
			fieldName := tools.SnakeString(field.Name)

			fieldCache.fields = append(fieldCache.fields, fieldName)
			fieldCache.fieldsMap[fieldName] = fieldTypes{
				jsonTagName: tagName,
				defaultVal:  field.Tag.Get("default"),
				fieldName:   fieldName,
				dblike:      field.Tag.Get("dblike"),
				kind:        kind,
			}
		}
	}
	return fieldCache
}
