package mysql

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

func parseResponseData(rows *sql.Rows, i interface{}) {
	//ct,_ := rows.ColumnTypes()
	st := reflect.TypeOf(i)
	sv := reflect.ValueOf(i)
	if st.Kind() == reflect.Ptr {
		st = st.Elem()
		sv = sv.Elem()
	}
	sl := st.NumField()
	mapping := map[string]int{}
	for si := 0; si < sl; si++ {
		tags := strings.Split(st.Field(si).Tag.Get(`json`), `,`)
		mapping[tags[0]] = si
	}
	t, _ := rows.ColumnTypes()
	data := make([]interface{}, len(t))
	for di, columnType := range t {
		index, flag := mapping[columnType.Name()]
		if !flag {
			s := ``
			data[di] = &s
			continue
		}
		ptr := sv.Field(index).Addr().Interface()
		//fmt.Println(`结构体映射字段:`, columnType.Name(), `对应的导出字段:`, st.Field(index).Name, `对应的真实字段指针是:`, ptr)
		switch st.Field(index).Type.Kind() {
		case reflect.String:
			data[di] = ptr.(*string)
		case reflect.Int:
			ptrs, ok := ptr.(*int)
			if ok {
				data[di] = ptrs
			} else {
				data[di] = (*int)(reflect.ValueOf(ptr).UnsafePointer())
			}
		case reflect.Int8:
			data[di] = ptr.(*int8)
		case reflect.Int16:
			data[di] = ptr.(*int16)
		case reflect.Int32:
			data[di] = ptr.(*int32)
		case reflect.Int64:
			data[di] = ptr.(*int64)
		case reflect.Uint:
			data[di] = ptr.(*uint)
		case reflect.Uint8:
			data[di] = ptr.(*uint8)
		case reflect.Uint16:
			v := uint16(0)
			data[di] = &v
		case reflect.Uint32:
			v := uint32(0)
			data[di] = &v
		case reflect.Uint64:
			v := uint64(0)
			data[di] = &v
		case reflect.Float32:
			v := float32(0)
			data[di] = &v
		case reflect.Float64:
			v := float64(0)
			data[di] = &v
		case reflect.Bool:
			v := false
			data[di] = &v
		default:
			panic(`sql结果转化出现暂不支持的类型:` + columnType.ScanType().Kind().String())
		}
	}
	rows.Next()
	err := rows.Scan(data...)
	if err != nil {
		fmt.Println(`sql异常`, err)
	}
	rows.Close()
}
