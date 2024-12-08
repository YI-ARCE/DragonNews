package log

import (
	"reflect"
	"strconv"
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

func serialize(source interface{}) (string, error) {
	serialize, err := encodeToString(source)
	if err != nil {
		return "", err
	}
	return serialize, err
}

// 数据序列化输出字符串
func encodeToString(source interface{}) (string, error) {
	list := ""
	TSource := reflect.TypeOf(source)
	RSource := reflect.ValueOf(source)
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
		l := RSource.Len()
		list = "["
		strings := ""
		empty := 0
		for i := 0; i < l; i++ {
			empty++
			strings = strJoin(strings, "\t", strconv.Itoa(i), ": ", refDataToString(RSource.Index(i), 2), "\n")
		}
		if empty > 0 {
			list = strJoin(list, "\n", strings)
		}
		list = strJoin(list, "]")
		break
	case s:
		list += RSource.String()
		break
	case structs:
		l := RSource.NumField()
		list = "{"
		strings := ""
		empty := 0
		for i := 0; i < l; i++ {
			empty++
			strings = strJoin(strings, "\t", TSource.Field(i).Name, ": ", refDataToString(RSource.Field(i), 2), "\n")
		}
		if empty > 0 {
			list = strJoin(list, "\n", strings)
		}
		list = strJoin(list, "}")
		break
	case m:
		list = "["
		strings := ""
		empty := 0
		Maps := RSource.MapRange()
		for Maps.Next() {
			empty++
			strings = strJoin(strings, "\t", refDataToString(Maps.Key(), 1), ": ", refDataToString(Maps.Value(), 2), "\n")
		}
		if empty > 0 {
			list = strJoin(list, "\n", strings)
		}
		list = strJoin(list, "]")
		break
	case ptr:
		list = strJoin(list, refDataToString(RSource.Elem(), 1))
		break
	case interfaces:
		if RSource.IsZero() {
			break
		}
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
		l := source.Len()
		list = "["
		str := ""
		empty := 0
		for i := 0; i < l; i++ {
			empty++
			str = strJoin(str, insertT(count), strconv.Itoa(i), ": ", refDataToString(source.Index(i), count+1), "\n")
		}
		if empty > 0 {
			list = strJoin(list, "\n", str, insertT(count-1))
		}
		list = strJoin(list, "]")
		break
	case s:
		list = strJoin(list, source.String())
		break
	case structs:
		l := source.NumField()
		list = "{"
		str := ""
		empty := 0
		for i := 0; i < l; i++ {
			empty++
			str = strJoin(str, insertT(count), TSource.Field(i).Name, ": ", refDataToString(source.Field(i), count+1), "\n")
		}
		if empty > 0 {
			list = strJoin(list, "\n", str, insertT(count-1))
		}
		list = strJoin(list, "}")
		break
	case m:
		list += "["
		str := ""
		Maps := source.MapRange()
		empty := 0
		for Maps.Next() {
			empty++
			str = strJoin(str, insertT(count), refDataToString(Maps.Key(), count+1), ": ", refDataToString(Maps.Value(), count+1), "\n")
		}
		if empty > 0 {
			list = strJoin(list, "\n", str, insertT(count-1))
		}
		list = strJoin(list, "]")
		break
	case ptr:
		if source.IsZero() {
			break
		}
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
