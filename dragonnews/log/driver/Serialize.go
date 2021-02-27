package driver

import (
	"reflect"
	"strconv"
)

const (
	B = iota + 1
	I
	I8
	I16
	I32
	I64
	Ui
	Ui8
	Ui16
	Ui32
	Ui64
	Uintptr
	F32
	F64
	Complex64
	Complex128
	A
	C
	Func
	Interface
	M
	Ptr
	Slice
	S
	Struct
	UnsafePointer
)

func Serialize(source interface{}) (string, error) {

	serialize, err := encodeToString(source)
	if err != nil {
		return "", err
	}
	return serialize, err
}

//数据序列化输出字符串
func encodeToString(source interface{}) (string, error) {
	list := ""
	TSource := reflect.TypeOf(source)
	RSource := reflect.ValueOf(source)
	switch TSource.Kind() {
	case B:
		if RSource.Bool() {
			list += "true"
		} else {
			list += "false"
		}
		break
	//多条件判断均为输出类型一致,统一转换
	case I, I8, I16, I32, I64:
		list += strconv.FormatInt(RSource.Int(), 10)
		break
	case Ui, Ui8, Ui16, Ui32, Ui64:
		list += strconv.FormatUint(RSource.Uint(), 10)
		break
	//case Uintptr:
	//	list = list + Stype + ":"+strconv.FormatUint(uint64(RSource.UnsafeAddr()),10)+";"
	//	break
	case F32, F64:
		list += strconv.FormatFloat(RSource.Float(), 'E', -1, 64)
		break
	case Complex64, Complex128:
		list += strconv.FormatComplex(RSource.Complex(), 'E', -1, 128)
		break
	case A:
		len := RSource.Len()
		list += "{\n"
		for i := 0; i < len; i++ {
			list += "\t" + strconv.Itoa(i) + ": " + refDataToString(RSource.Index(i), 2) + "\n"
		}
		list += "}"
		break
	case S:
		list += RSource.String()
		break
	case Struct:
		len := RSource.NumField()
		list += "{\n"
		for i := 0; i < len; i++ {
			list += "\t" + TSource.Field(i).Name + ": " + refDataToString(RSource.Field(i), 2) + "\n"
		}
		list += "}"
		break
	case M:
		list += "{\n"
		Maps := RSource.MapRange()
		for Maps.Next() {
			list += "\t" + refDataToString(Maps.Key(), 1) + ": " + refDataToString(Maps.Value(), 2) + "\n"
		}
		list += "}"
		break
	case Ptr:
		list += refDataToString(RSource.Elem(), 1)
		break
	default:
		list += RSource.Type().Kind().String() + "\n"
	}
	return list, nil
}

//获取反射数据
//与encodeToString区别在于传入源数据跟映射的反射数据
func refDataToString(source reflect.Value, count int) string {
	list := ""
	TSource := source.Type()
	switch TSource.Kind() {
	case B:
		if source.Bool() {
			list = "true"
		} else {
			list = "false"
		}
		break
	//多条件判断均为输出类型一致,统一转换
	case I, I8, I16, I32, I64:
		list += strconv.FormatInt(source.Int(), 10)
		break
	case Ui, Ui8, Ui16, Ui32, Ui64:
		list += strconv.FormatUint(source.Uint(), 10)
		break
	//case Uintptr:
	//	list = list + Stype + ":"+strconv.FormatUint(uint64(source.UnsafeAddr()),10)+";"
	//	break
	case F32, F64:
		list += strconv.FormatFloat(source.Float(), 'E', -1, 64)
		break
	case Complex64, Complex128:
		list += strconv.FormatComplex(source.Complex(), 'E', -1, 128)
		break
	case A, Slice:
		len := source.Len()
		list += "{\n"
		for i := 0; i < len; i++ {
			list += insertT(count) + strconv.Itoa(i) + ": " + refDataToString(source.Index(i), count+1) + "\n"
		}
		list += insertT(count-1) + "}"
		break
	case S:
		if len(source.String()) < 1 {
			return list
		}
		list += source.String()
		break
	case Struct:
		len := source.NumField()
		list += "{\n"
		for i := 0; i < len; i++ {
			list += insertT(count) + TSource.Field(i).Name + ": " + refDataToString(source.Field(i), count+1) + "\n"
		}
		list += insertT(count-1) + "}"
		break
	case M:
		list += "{\n"
		Maps := source.MapRange()
		for Maps.Next() {
			list += insertT(count) + refDataToString(Maps.Key(), count+1) + ": " + refDataToString(Maps.Value(), count+1) + "\n"
		}
		list += insertT(count-1) + "}"
		break
	case Ptr, Interface:
		if source.IsZero() {
			break
		}
		list += refDataToString(source.Elem(), count)
		break
	default:
		list += source.Type().Kind().String()
		break
	}
	return list
}

func insertT(count int) string {
	str := ""
	for i := 0; i < count; i++ {
		str += "\t"
	}
	return str
}
