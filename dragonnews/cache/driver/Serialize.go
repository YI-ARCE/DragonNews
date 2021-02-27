package driver

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
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

//执行序列化
func Serialize(data interface{}) (string, error) {
	serialize, err := encodeToString(data)
	if err != nil {
		return "", err
	}
	return serialize, err
}

//反序列化
func UnSerizlize(serialize []byte, ps interface{}) error {
	return decodeToSource(string(serialize), ps)
}

//数据序列化输出字符串
func encodeToString(source interface{}) (string, error) {
	list := "dragon_s:"
	TSource := reflect.TypeOf(source)
	Stype := getTypes(int(TSource.Kind()))
	RSource := reflect.ValueOf(source)
	switch TSource.Kind() {
	case B:
		if RSource.Bool() {
			list += "B:1;"
		} else {
			list += "B:0;"
		}
		break
	//多条件判断均为输出类型一致,统一转换
	case I, I8, I16, I32, I64:
		list += Stype + ":" + strconv.FormatInt(RSource.Int(), 10) + ";"
		break
	case Ui, Ui8, Ui16, Ui32, Ui64:
		list += Stype + ":" + strconv.FormatUint(RSource.Uint(), 10) + ";"
		break
	//case Uintptr:
	//	list = list + Stype + ":"+strconv.FormatUint(uint64(RSource.UnsafeAddr()),10)+";"
	//	break
	case F32, F64:
		list += Stype + ":" + strconv.FormatFloat(RSource.Float(), 'E', -1, 64) + ";"
		break
	case Complex64, Complex128:
		list += Stype + ":" + strconv.FormatComplex(RSource.Complex(), 'E', -1, 128) + ";"
		break
	case A:
		len := RSource.Len()
		list += Stype + ":{"
		for i := 0; i < len; i++ {
			list += "K:" + strconv.Itoa(i) + ";" + refDataToString(RSource.Index(i))
		}
		list += "}"
		break
	case S:
		list += Stype + ":" + RSource.String() + ";"
		break
	case Struct:
		len := RSource.NumField()
		list += Stype + ":{"
		for i := 0; i < len; i++ {
			list += "K:" + TSource.Field(i).Name + ";" + refDataToString(RSource.Field(i))
		}
		list += "}"
		break
	case M:
		list += Stype + ":" + strconv.FormatInt(int64(RSource.Len()), 10) + ":{"
		Maps := RSource.MapRange()
		for Maps.Next() {
			list = list + refDataToString(Maps.Key()) + refDataToString(Maps.Value())
		}
		list += "}"
		break
	case Ptr:
		list += refDataToString(RSource.Elem())
		break
	default:
		cc := errors.New("不支持存储的类型:" + TSource.Kind().String())
		return "", cc
	}
	return list, nil
}

//获取反射数据
//与encodeToString区别在于传入源数据跟映射的反射数据
func refDataToString(source reflect.Value) string {
	list := ""
	TSource := source.Type()
	Stype := getTypes(int(TSource.Kind()))
	switch TSource.Kind() {
	case B:
		if source.Bool() {
			list += "B:1;"
		} else {
			list += "B:0;"
		}
		break
	//多条件判断均为输出类型一致,统一转换
	case I, I8, I16, I32, I64:
		list += Stype + ":" + strconv.FormatInt(source.Int(), 10) + ";"
		break
	case Ui, Ui8, Ui16, Ui32, Ui64:
		list += Stype + ":" + strconv.FormatUint(source.Uint(), 10) + ";"
		break
	//case Uintptr:
	//	list = list + Stype + ":"+strconv.FormatUint(uint64(source.UnsafeAddr()),10)+";"
	//	break
	case F32, F64:
		list += Stype + ":" + strconv.FormatFloat(source.Float(), 'E', -1, 64) + ";"
		break
	case Complex64, Complex128:
		list += Stype + ":" + strconv.FormatComplex(source.Complex(), 'E', -1, 128) + ";"
		break
	case A:
		len := source.Len()
		list += Stype + ":{"
		for i := 0; i < len; i++ {
			list = list + "K:" + strconv.Itoa(i) + ";" + refDataToString(source.Index(i))
		}
		list += "};"
		break
	case S:
		list += Stype + ":" + source.String() + ";"
		break
	case Struct:
		len := source.NumField()
		list += Stype + ":{"
		for i := 0; i < len; i++ {
			list = list + "K:" + TSource.Field(i).Name + ";" + refDataToString(source.Field(i))
		}
		list += "};"
		break
	case M:
		list = list + Stype + ":" + strconv.FormatInt(int64(source.Len()), 10) + ":{"
		Maps := source.MapRange()
		for Maps.Next() {
			list += refDataToString(Maps.Key()) + refDataToString(Maps.Value())
		}
		list += "};"
		break
	case Ptr:
		list += refDataToString(source.Elem())
		break
	}
	return list
}

//解码反序列化
func decodeToSource(source string, ps interface{}) error {
	Len := len(source)
	if Len < 8 || source[0:8] != "dragon_s" {
		return errors.New("非序列数据")
	} else {
		source = source[9:]
		pos := strings.IndexAny(source, ":")
		Class := source[0:pos]
		err := unClassJudge(Class, source[pos+1:], reflect.ValueOf(ps).Elem())
		if err != nil {
			return err
		}
	}
	return nil
}

//类型判断
func unClassJudge(Class string, source string, ps reflect.Value) error {
	var error error
	error = nil
	switch Class {
	case "B":
		if source[0:1] == "1" {
			ps.SetBool(true)
		} else {
			ps.SetBool(false)
		}
		break
	case "I", "I8", "I16", "I32", "I64":
		var base int
		if Class == "I" || Class == "I32" {
			base = 32
		} else if Class == "I8" {
			base = 8
		} else if Class == "I16" {
			base = 16
		} else {
			base = 64
		}
		val, err := strconv.ParseInt(source[0:strings.IndexAny(source, ";")], 10, base)
		if err != nil {
			return err
		}
		ps.SetInt(val)
		break
	case "Ui", "Ui8", "Ui16", "Ui32", "Ui64":
		var base int
		if Class == "Ui" || Class == "Ui32" {
			base = 32
		} else if Class == "Ui8" {
			base = 8
		} else if Class == "Ui16" {
			base = 16
		} else {
			base = 64
		}
		val, err := strconv.ParseUint(source[0:strings.IndexAny(source, ";")], 10, base)
		if err != nil {
			return err
		}
		ps.SetUint(val)
		break
	case "F32", "F64":
		val, err := strconv.ParseFloat(source[0:strings.IndexAny(source, ";")], 10)
		if err != nil {
			return err
		}
		ps.SetFloat(val)
		break
	case "Complex64", "Complex128":
		val, err := strconv.ParseComplex(source[0:strings.IndexAny(source, ";")], 10)
		if err != nil {
			return err
		}
		ps.SetComplex(val)
		break
	case "A":
		source = source[1 : len(source)-1]
		psLen := ps.Len()
		for i := 0; i < psLen; i++ {
			keyPos := strings.IndexAny(source, ":")
			if source[0:keyPos] != "K" {
				error = errors.New("不正常的序列数据!")
				break
			}
			indexPos := strings.IndexAny(source, ";")
			index, err := strconv.ParseInt(source[keyPos+1:indexPos], 10, 64)
			if err != nil {
				error = errors.New("不正常的序列数据!")
				break
			}
			source = source[indexPos+1:]
			NewClass := source[0:strings.IndexAny(source, ":")]
			SLen, err := refUnClassJudge(NewClass, source, ps.Index(int(index)))
			if err != nil {
				return err
			}
			source = source[SLen:]
		}
		break
	case "S":
		ps.SetString(source[0:strings.IndexAny(source, ";")])
		break
	case "Struct":
		source = source[1 : len(source)-1]
		psLen := ps.NumField()
		for i := 0; i < psLen; i++ {
			keyPos := strings.IndexAny(source, ":")
			if source[0:keyPos] != "K" {
				error = errors.New("不正常的序列数据!")
				break
			}
			indexPos := strings.IndexAny(source, ";")
			name := source[keyPos+1 : indexPos]
			if ps.FieldByName(name).IsValid() == false {
				error = errors.New("不正常的序列数据!")
				break
			}
			source = source[indexPos+1:]
			NewClass := source[0:strings.IndexAny(source, ":")]
			if ps.FieldByName(name).CanSet() {
				SLen, err := refUnClassJudge(NewClass, source, ps.FieldByName(name))
				if err != nil {
					return err
				}
				source = source[SLen:]
			} else {
				SLen, err := unRefUnClassJudge(NewClass, source, ps.FieldByName(name))
				if err != nil {
					return err
				}
				source = source[SLen:]
			}
		}
		break
	case "M":
		MLen := strings.IndexAny(source, ":")
		MCount, err := strconv.ParseInt(source[0:MLen], 10, 64)
		if err != nil {
			return err
		}
		source = source[MLen+2 : len(source)-1]
		if ps.IsNil() {
			ps.Set(reflect.MakeMap(ps.Type()))
		}
		for i := int64(0); i < MCount; i++ {
			rk := reflect.New(ps.Type().Key()).Elem()
			rv := reflect.New(ps.Type().Elem()).Elem()
			NewClass := source[0:strings.IndexAny(source, ":")]
			KLen, err := refUnClassJudge(NewClass, source, rk)
			if err != nil {
				return err
			}
			source = source[KLen:]
			NewClass = source[0:strings.IndexAny(source, ":")]
			VLen, err := refUnClassJudge(NewClass, source, rv)
			if err != nil {
				return err
			}
			source = source[VLen:]
			ps.SetMapIndex(rk, rv)
		}
		break
	}
	return error
}

//导出数据的循环复原
func refUnClassJudge(Class string, source string, ps reflect.Value) (int, error) {
	var error error
	var Len int
	switch Class {
	case "B":
		Len = strings.IndexAny(source, ";")
		if source[2:Len] == "1" {
			ps.SetBool(true)
		} else {
			ps.SetBool(false)
		}
		Len = 1
		break
	case "I", "I8", "I16", "I32", "I64":
		var base int
		if Class == "I" || Class == "I32" {
			base = 32
		} else if Class == "I8" {
			base = 8
		} else if Class == "I16" {
			base = 16
		} else {
			base = 64
		}
		Len = strings.IndexAny(source, ";")
		val, err := strconv.ParseInt(source[strings.IndexAny(source, ":"):Len], 10, base)
		if err != nil {
			return 0, err
		}
		ps.SetInt(val)
		break
	case "Ui", "Ui8", "Ui16", "Ui32", "Ui64":
		var base int
		if Class == "Ui" || Class == "Ui32" {
			base = 32
		} else if Class == "Ui8" {
			base = 8
		} else if Class == "Ui16" {
			base = 16
		} else {
			base = 64
		}
		Len = strings.IndexAny(source, ";")
		val, err := strconv.ParseUint(source[strings.IndexAny(source, ":"):Len], 10, base)
		if err != nil {
			return 0, err
		}
		ps.SetUint(val)
		break
	case "F32", "F64":
		Len = strings.IndexAny(source, ";")
		val, err := strconv.ParseFloat(source[strings.IndexAny(source, ":"):Len], 10)
		if err != nil {
			return 0, err
		}
		ps.SetFloat(val)
		break
	case "Complex64", "Complex128":
		Len = strings.IndexAny(source, ";")
		val, err := strconv.ParseComplex(source[strings.IndexAny(source, ":"):Len], 10)
		if err != nil {
			return 0, err
		}
		ps.SetComplex(val)
		break
	case "A":
		source = source[2 : len(source)-1]
		psLen := ps.Len()
		for i := 0; i < psLen; i++ {
			keyPos := strings.IndexAny(source, ":")
			if source[0:keyPos] != "K" {
				error = errors.New("不正常的序列数据!")
				break
			}
			indexPos := strings.IndexAny(source, ";")
			index, err := strconv.ParseInt(source[keyPos+1:indexPos], 10, 64)
			if err != nil {
				error = errors.New("不正常的序列数据!")
				break
			}
			source = source[indexPos+1:]
			NewClass := source[0:strings.IndexAny(source, ":")]
			SLen, err := refUnClassJudge(NewClass, source, ps.Index(int(index)))
			if err != nil {
				return 0, err
			}
			SLen = SLen + 1
			source = source[SLen:]
			Len = Len + SLen
		}
		break
	case "S":
		Len = strings.IndexAny(source, ";")
		ps.SetString(source[2:Len])
		break
	case "Struct":
		source = source[8:]
		psLen := ps.NumField()
		for i := 0; i < psLen; i++ {
			keyPos := strings.IndexAny(source, ":")
			if source[0:keyPos] != "K" {
				error = errors.New("不正常的序列数据!")
				break
			}
			indexPos := strings.IndexAny(source, ";")
			name := source[keyPos+1 : indexPos]
			if ps.FieldByName(name).IsValid() == false {
				error = errors.New("不正常的序列数据!")
				break
			}
			Len = Len + indexPos
			source = source[indexPos+1:]
			NewClass := source[0:strings.IndexAny(source, ":")]
			if ps.FieldByName(name).CanSet() {
				FLen, err := refUnClassJudge(NewClass, source, ps.FieldByName(name))
				if err != nil {
					return 0, err
				}
				source = source[FLen:]
				Len = Len + FLen
			} else {
				FLen, err := unRefUnClassJudge(NewClass, source, ps.FieldByName(name))
				if err != nil {
					return 0, err
				}
				source = source[FLen:]
				Len = Len + FLen
			}
		}
		Len = Len + 8 + 2
		break
	case "M":
		source = source[2:]
		MLen := strings.IndexAny(source, ":")
		MCount, err := strconv.ParseInt(source[0:MLen], 10, 64)
		if err != nil {
			error = err
		}
		source = source[MLen+2:]
		if ps.IsNil() {
			ps.Set(reflect.MakeMap(ps.Type()))
		}
		for i := int64(0); i < MCount; i++ {
			rk := reflect.New(ps.Type().Key()).Elem()
			rv := reflect.New(ps.Type().Elem()).Elem()
			NewClass := source[0:strings.IndexAny(source, ":")]
			FLen, err := refUnClassJudge(NewClass, source, rk)
			if err != nil {
				return 0, err
			}
			source = source[FLen:]
			Len = Len + FLen
			NewClass = source[0:strings.IndexAny(source, ":")]
			FLen, err = refUnClassJudge(NewClass, source, rv)
			if err != nil {
				return 0, err
			}
			Len = Len + FLen
			source = source[FLen:]
			ps.SetMapIndex(rk, rv)
		}
		Len = Len + 2 + MLen + 2 + 1
		break
	default:
		error = errors.New("暂未支持的序列类型:" + Class)
	}
	return Len + 1, error
}

func unRefUnClassJudge(Class string, source string, ps reflect.Value) (int, error) {
	var error error
	error = nil
	var Len int
	Ups := unsafe.Pointer(ps.UnsafeAddr())
	switch Class {
	case "B":
		Len = strings.IndexAny(source, ";")
		if source[2:Len] == "1" {
			*((*bool)(Ups)) = true
		} else {
			*((*bool)(Ups)) = false
		}
		Len = 1
		break
	case "I", "I8", "I16", "I32", "I64":
		var base int
		if Class == "I" || Class == "I32" {
			base = 32
		} else if Class == "I8" {
			base = 8
		} else if Class == "I16" {
			base = 16
		} else {
			base = 64
		}
		Len = strings.IndexAny(source, ";")
		val, err := strconv.ParseInt(source[strings.IndexAny(source, ":")+1:Len], 10, base)
		if err != nil {
			error = err
			break
		}
		*((*int64)(Ups)) = val
		break
	case "Ui", "Ui8", "Ui16", "Ui32", "Ui64":
		var base int
		if Class == "Ui" || Class == "Ui32" {
			base = 32
		} else if Class == "Ui8" {
			base = 8
		} else if Class == "Ui16" {
			base = 16
		} else {
			base = 64
		}
		Len = strings.IndexAny(source, ";")
		val, err := strconv.ParseUint(source[strings.IndexAny(source, ":"):Len], 10, base)
		if err != nil {
			error = err
			break
		}
		*((*uint64)(Ups)) = val
		break
	case "F32", "F64":
		Len = strings.IndexAny(source, ";")
		val, err := strconv.ParseFloat(source[strings.IndexAny(source, ":"):Len], 10)
		if err != nil {
			error = err
			break
		}
		*((*float64)(Ups)) = val
		break
	case "Complex64", "Complex128":
		Len = strings.IndexAny(source, ";")
		val, err := strconv.ParseComplex(source[strings.IndexAny(source, ":"):Len], 10)
		if err != nil {
			error = err
			break
		}
		*((*complex128)(Ups)) = val
		break
	case "A":
		source = source[2 : len(source)-1]
		psLen := ps.Len()
		for i := 0; i < psLen; i++ {
			keyPos := strings.IndexAny(source, ":")
			if source[0:keyPos] != "K" {
				error = errors.New("不正常的序列数据!")
				break
			}
			indexPos := strings.IndexAny(source, ";")
			index, err := strconv.ParseInt(source[keyPos+1:indexPos], 10, 64)
			if err != nil {
				error = errors.New("不正常的序列数据!")
				break
			}
			source = source[indexPos+1:]
			NewClass := source[0:strings.IndexAny(source, ":")]
			SLen, err := unRefUnClassJudge(NewClass, source, ps.Index(int(index)))
			if err != nil {
				return 0, err
			}
			SLen = SLen + 1
			source = source[SLen:]
			Len = Len + SLen
		}
		break
	case "S":
		Len = strings.IndexAny(source, ";")
		*((*string)(Ups)) = source[2:Len]
		break
	case "Struct":
		source = source[8:]
		psLen := ps.NumField()
		for i := 0; i < psLen; i++ {
			keyPos := strings.IndexAny(source, ":")
			if source[0:keyPos] != "K" {
				error = errors.New("不正常的序列数据!")
				break
			}
			indexPos := strings.IndexAny(source, ";")
			name := source[keyPos+1 : indexPos]
			if ps.FieldByName(name).IsValid() == false {
				error = errors.New("不正常的序列数据!")
				break
			}
			Len = Len + indexPos
			source = source[indexPos+1:]
			NewClass := source[0:strings.IndexAny(source, ":")]
			if ps.FieldByName(name).CanSet() {
				FLen, err := refUnClassJudge(NewClass, source, ps.FieldByName(name))
				if err != nil {
					return 0, err
				}
				source = source[FLen:]
				Len = Len + FLen
			} else {
				FLen, err := unRefUnClassJudge(NewClass, source, ps.FieldByName(name))
				if err != nil {
					return 0, err
				}
				source = source[FLen:]
				Len = Len + FLen
			}
		}
		Len = Len + 8 + 2
		break

	case "M":
		source = source[2:]
		MLen := strings.IndexAny(source, ":")
		MCount, err := strconv.ParseInt(source[0:MLen], 10, 64)
		if err != nil {
			error = err
		}
		source = source[MLen+2:]
		if ps.IsNil() {
			ps.Set(reflect.MakeMap(ps.Type()))
		}
		for i := int64(0); i < MCount; i++ {
			rk := reflect.New(ps.Type().Key()).Elem()
			rv := reflect.New(ps.Type().Elem()).Elem()
			NewClass := source[0:strings.IndexAny(source, ":")]
			FLen, err := refUnClassJudge(NewClass, source, rk)
			if err != nil {
				return 0, err
			}
			source = source[FLen:]
			Len = Len + FLen
			NewClass = source[0:strings.IndexAny(source, ":")]
			FLen, err = refUnClassJudge(NewClass, source, rv)
			if err != nil {
				return 0, err
			}
			Len = Len + FLen
			source = source[FLen:]
			ps.SetMapIndex(rk, rv)
		}
		Len = Len + 2 + MLen + 2 + 1
		break
	default:
		error = errors.New("暂未支持的序列类型:" + Class)
		break
	}
	return Len + 1, error
}

//获取类型
func getTypes(types int) string {
	rType := ""
	switch types {
	case I:
		rType = "I"
		break
	case I8:
		rType = "I8"
		break
	case I16:
		rType = "I16"
		break
	case I32:
		rType = "I32"
		break
	case I64:
		rType = "I64"
		break
	case Ui:
		rType = "Ui"
		break
	case Ui8:
		rType = "Ui8"
		break
	case Ui16:
		rType = "Ui16"
		break
	case Ui32:
		rType = "Ui32"
		break
	case Ui64:
		rType = "Ui64"
		break
	case Uintptr:
		rType = "Uintptr"
		break
	case F32:
		rType = "F32"
		break
	case F64:
		rType = "F64"
		break
	case Complex64:
		rType = "Complex64"
		break
	case Complex128:
		rType = "Complex128"
		break
	case A:
		rType = "A"
		break
	case C:
		rType = "C"
		break
	case Func:
		rType = "Func"
		break
	case Interface:
		rType = "Interface"
		break
	case M:
		rType = "M"
		break
	case Ptr:
		rType = "Ptr"
		break
	case Slice:
		rType = "Slice"
		break
	case S:
		rType = "S"
		break
	case Struct:
		rType = "Struct"
		break
	case UnsafePointer:
		rType = "UnsafePointer"
		break
	default:
		rType = "nil"
		break
	}
	return rType
}
