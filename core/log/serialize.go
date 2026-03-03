package log

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Serializer 序列化器接口
type Serializer interface {
	Serialize(source interface{}) (string, error)
}

// DefaultSerializer 默认序列化器
type DefaultSerializer struct{}

// Serialize 序列化任意类型数据为字符串
func (ds *DefaultSerializer) Serialize(source interface{}) (string, error) {
	return ds.encodeToString(source, 0)
}

// serialize 全局序列化函数，保持向后兼容
func serialize(source interface{}) (string, error) {
	serializer := &DefaultSerializer{}
	return serializer.Serialize(source)
}

// encodeToString 核心序列化实现
func (ds *DefaultSerializer) encodeToString(source interface{}, indent int) (string, error) {
	// 处理nil值
	if source == nil {
		return "<nil>", nil
	}

	// 获取反射信息
	value := reflect.ValueOf(source)
	typ := value.Type()
	kind := typ.Kind()

	// 根据类型进行序列化
	switch kind {
	case reflect.Bool:
		if value.Bool() {
			return "true", nil
		}
		return "false", nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10), nil

	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'f', -1, 64), nil

	case reflect.Complex64, reflect.Complex128:
		return strconv.FormatComplex(value.Complex(), 'f', -1, 128), nil

	case reflect.String:
		return value.String(), nil

	case reflect.Array, reflect.Slice:
		return ds.serializeSlice(value, indent)

	case reflect.Struct:
		return ds.serializeStruct(value, indent)

	case reflect.Map:
		return ds.serializeMap(value, indent)

	case reflect.Ptr:
		if value.IsNil() {
			return "<nil>", nil
		}
		return ds.encodeToString(value.Elem().Interface(), indent)

	case reflect.Interface:
		if value.IsNil() {
			return "<nil>", nil
		}
		// 解引用接口值
		elem := value.Elem()
		for i := 0; i < 5 && (elem.Kind() == reflect.Interface || elem.Kind() == reflect.Ptr); i++ {
			if elem.IsNil() {
				return "<nil>", nil
			}
			elem = elem.Elem()
		}
		return ds.encodeToString(elem.Interface(), indent)

	default:
		return fmt.Sprintf("<%s>", kind.String()), nil
	}
}

// serializeSlice 序列化数组和切片
func (ds *DefaultSerializer) serializeSlice(value reflect.Value, indent int) (string, error) {
	length := value.Len()
	if length == 0 {
		return "[]", nil
	}

	var builder strings.Builder
	builder.WriteString("[\n")

	indentStr := strings.Repeat("\t", indent+1)

	for i := 0; i < length; i++ {
		builder.WriteString(indentStr)
		builder.WriteString(strconv.Itoa(i))
		builder.WriteString(": ")

		itemValue := value.Index(i)
		itemStr, err := ds.encodeToString(itemValue.Interface(), indent+1)
		if err != nil {
			return "", err
		}
		builder.WriteString(itemStr)
		builder.WriteString("\n")
	}

	builder.WriteString(strings.Repeat("\t", indent))
	builder.WriteString("]")

	return builder.String(), nil
}

// serializeStruct 序列化结构体
func (ds *DefaultSerializer) serializeStruct(value reflect.Value, indent int) (string, error) {
	typ := value.Type()
	fieldCount := value.NumField()
	if fieldCount == 0 {
		return "{}", nil
	}

	var builder strings.Builder
	builder.WriteString("{\n")

	indentStr := strings.Repeat("\t", indent+1)

	for i := 0; i < fieldCount; i++ {
		field := typ.Field(i)
		fieldValue := value.Field(i)

		// 跳过未导出字段
		if !field.IsExported() {
			continue
		}

		builder.WriteString(indentStr)
		builder.WriteString(field.Name)
		builder.WriteString(": ")

		fieldStr, err := ds.encodeToString(fieldValue.Interface(), indent+1)
		if err != nil {
			return "", err
		}
		builder.WriteString(fieldStr)
		builder.WriteString("\n")
	}

	builder.WriteString(strings.Repeat("\t", indent))
	builder.WriteString("}")

	return builder.String(), nil
}

// serializeMap 序列化映射
func (ds *DefaultSerializer) serializeMap(value reflect.Value, indent int) (string, error) {
	if value.Len() == 0 {
		return "[]", nil
	}

	var builder strings.Builder
	builder.WriteString("[\n")

	indentStr := strings.Repeat("\t", indent+1)

	// 遍历map
	iter := value.MapRange()
	for iter.Next() {
		builder.WriteString(indentStr)

		// 序列化key
		keyStr, err := ds.encodeToString(iter.Key().Interface(), indent+1)
		if err != nil {
			return "", err
		}
		builder.WriteString(keyStr)
		builder.WriteString(": ")

		// 序列化value
		valueStr, err := ds.encodeToString(iter.Value().Interface(), indent+1)
		if err != nil {
			return "", err
		}
		builder.WriteString(valueStr)
		builder.WriteString("\n")
	}

	builder.WriteString(strings.Repeat("\t", indent))
	builder.WriteString("]")

	return builder.String(), nil
}
