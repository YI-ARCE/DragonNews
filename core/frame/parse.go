package frame

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

const (
	b = iota + 1
	i
	i8
	i16
	i32
	i64
	ui
	ui8
	ui16
	ui32
	ui64
	uiptr
	f32
	f64
	complexity64
	complexity128
	a
	c
	function
	interfaces
	m
	ptr
	slice
	s
	structs
	unsafePointer
)

func serialize(source interface{}) string {
	serialize, err := encodeToString(source)
	if err != nil {
		fmt.Println(`framePrintError:`, err)
		return ""
	}
	return serialize
}

// 数据序列化输出字符串
func encodeToString(source interface{}) (string, error) {
	if source == nil {
		return `<nil>`, nil
	}
	if v, ok := source.(int); ok {
		if v == 0 {
			return `0`, nil
		}
	}
	list := ""
	TSource := reflect.TypeOf(source)
	RSource := reflect.ValueOf(source)
	if RSource.IsZero() {
		return `<nil>`, nil
	}
	switch TSource.Kind() {
	case b:
		if RSource.Bool() {
			list = "true"
		} else {
			list = "false"
		}
		break
	//多条件判断均为输出类型一致,统一转换
	case i, i8, i16, i32, i64:
		list = strconv.FormatInt(RSource.Int(), 10)
		break
	case ui, ui8, ui16, ui32, ui64:
		list = strconv.FormatUint(RSource.Uint(), 10)
		break
	//case uiptr:
	//	list = list + Stype + ":"+strconv.FormatUint(uint64(RSource.UnsafeAddr()),10)+";"
	//	break
	case f32, f64:
		list = strconv.FormatFloat(RSource.Float(), 'f', -1, 64)
		break
	case complexity64, complexity128:
		list = strconv.FormatComplex(RSource.Complex(), 'f', -1, 128)
		break
	case a, slice:
		if RSource.IsZero() {
			return `<nil>`, nil
		}
		if RSource.Len() < 1 {
			return `<nil>`, nil
		}
		l := RSource.Len()
		list = "[ "
		strings := ""
		empty := 0
		for i := 0; i < l; i++ {
			empty++
			strings = strJoin(strings, strconv.Itoa(i), ":", refDataToString(RSource.Index(i), 2), ` , `)
		}
		strings = strings[:len(strings)-3]
		if empty > 0 {
			list = strJoin(list, strings)
		}
		list = strJoin(list, " ]")
		break
	case s:
		list += RSource.String()
		break
	case structs:
		l := RSource.NumField()
		list = "{ "
		strings := ""
		empty := 0
		for i := 0; i < l; i++ {
			empty++
			strings = strJoin(strings, TSource.Field(i).Name, ":", refDataToString(RSource.Field(i), 2), " , ")
		}
		strings = strings[:len(strings)-3]
		if empty > 0 {
			list = strJoin(list, strings)
		}
		list = strJoin(list, " }")
		break
	case m:
		if RSource.Len() < 1 {
			return `<nil>`, nil
		}
		list = "[ "
		strings := ""
		empty := 0
		Maps := RSource.MapRange()
		for Maps.Next() {
			empty++
			strings = strJoin(strings, refDataToString(Maps.Key(), 1), ":", refDataToString(Maps.Value(), 2), " , ")
		}
		strings = strings[:len(strings)-3]
		if empty > 0 {
			list = strJoin(list, strings)
		}
		list = strJoin(list, " ]")
		break
	case ptr:
		list = strJoin(list, refDataToString(RSource.Elem(), 1))
		break
	case uiptr:
		list = strJoin(list, `(uintptr) `, fmt.Sprintf("%v", RSource.Interface()))
	case interfaces:
		list = strJoin(list, refDataToString(RSource.Elem(), 1))
	default:
		list = strJoin(RSource.Type().Kind().String(), "\n")
	}
	return list, nil
}

// 获取反射数据
// 与encodeToString区别在于传入源数据跟映射的反射数据
func refDataToString(source reflect.Value, count int) string {
	list := ""
	TSource := source.Type()
	if source.IsZero() {
		return `<nil>`
	}
	switch TSource.Kind() {
	case b:
		if source.Bool() {
			list = "true"
		} else {
			list = "false"
		}
		break
	//多条件判断均为输出类型一致,统一转换
	case i, i8, i16, i32, i64:
		list = strconv.FormatInt(source.Int(), 10)
		break
	case ui, ui8, ui16, ui32, ui64:
		list = strconv.FormatUint(source.Uint(), 10)
		break
	//case uiptr:
	//	list = list + Stype + ":"+strconv.FormatUint(uint64(source.UnsafeAddr()),10)+";"
	//	break
	case f32, f64:
		list = strconv.FormatFloat(source.Float(), 'f', -1, 64)
		break
	case complexity64, complexity128:
		list = strconv.FormatComplex(source.Complex(), 'f', -1, 128)
		break
	case a, slice:
		if source.Len() < 1 {
			return `<nil>`
		}
		l := source.Len()
		list = "[ "
		str := ""
		empty := 0
		for i := 0; i < l; i++ {
			empty++
			str = strJoin(str, strconv.Itoa(i), ":", refDataToString(source.Index(i), count+1), " , ")
		}
		str = str[:len(str)-3]
		if empty > 0 {
			list = strJoin(list, str)
		}
		list = strJoin(list, " ]")
		break
	case s:
		list = strJoin(list, source.String())
		break
	case structs:
		l := source.NumField()
		list = "{ "
		str := ""
		empty := 0
		for i := 0; i < l; i++ {
			empty++
			str = strJoin(str, TSource.Field(i).Name, ":", refDataToString(source.Field(i), count+1), " , ")
		}
		str = str[:len(str)-3]
		if empty > 0 {
			list = strJoin(list, str)
		}
		list = strJoin(list, " }")
		break
	case m:
		if source.Len() < 1 {
			return `<nil>`
		}
		list += "[ "
		str := ""
		Maps := source.MapRange()
		empty := 0
		for Maps.Next() {
			empty++
			str = strJoin(str, refDataToString(Maps.Key(), count+1), ":", refDataToString(Maps.Value(), count+1), " , ")
		}
		str = str[:len(str)-3]
		if empty > 0 {
			list = strJoin(list, str)
		}
		list = strJoin(list, " ]")
		break
	case ptr:
		list = strJoin(list, refDataToString(source.Elem(), count))
		break
	case interfaces:
		if source.IsZero() {
			break
		}
		is := 1
		for ; is < 5; is++ {
			source = source.Elem()
			if source.Type().Kind().String() != "interface" || source.Type().Kind().String() != "ptr" {
				break
			}
		}
		if is == 4 {
			list = strJoin(list, source.Type().Kind().String())
			break
		}
		list = strJoin(list, refDataToString(source, count))
		break
	default:
		list = strJoin(list, source.Type().Kind().String())
		break
	}
	return list
}

func insertT(count int) string {
	str := ""
	for i := 0; i < count; i++ {
		str = strBuild(str, "\t")
	}
	return str
}

// 多次少于32字符时用此方法
//
//	正常用此方法即可
func strBuild(str string, s ...string) string {
	b := strings.Builder{}
	b.Write([]byte(str))
	for _, v := range s {
		b.Write(*(*[]byte)(unsafe.Pointer(&v)))
	}
	str = b.String()
	return str
}

// 多次大于32字符拼接时用此方法
func strJoin(str string, s ...string) string {
	bs := []string{str}
	for _, v := range s {
		bs = append(bs, v)
	}
	return strings.Join(bs, "")
}
