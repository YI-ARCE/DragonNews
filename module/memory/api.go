package memory

import (
	"encoding/binary"
	"fmt"
	"golang.org/x/sys/windows"
	"math"
	"sync"
	"syscall"
	"time"
	"unsafe"
	"yiarce/core/frame"
	"yiarce/core/timing"
)

var (
	kernel32               = windows.NewLazySystemDLL("kernel32.dll")
	procReadProcessMemory  = kernel32.NewProc("ReadProcessMemory")
	procWriteProcessMemory = kernel32.NewProc("WriteProcessMemory")
	procVirtualProtectEx   = kernel32.NewProc("VirtualProtectEx")
)

func OpenProcessFromName(name string) []uint32 {

	// 获取资源管理器进程ID
	handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		fmt.Printf("CreateToolhelp32Snapshot failed: %s", err)
		return []uint32{}
	}
	defer windows.CloseHandle(handle)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	err = windows.Process32First(handle, &entry)
	var ids []uint32
	for err == nil {
		if windows.UTF16ToString(entry.ExeFile[:]) == name {
			ids = append(ids, entry.ProcessID)
		}
		err = windows.Process32Next(handle, &entry)
	}
	return ids
}

func OpenProcessFromPid(id uint32) Ctx {
	handle, err := windows.OpenProcess(0x1F0FFF, false, id)
	if err != nil {
		fmt.Printf("OpenProcess failed: %s\n", err)
		return Ctx{}
	}
	return Ctx{
		handle:  handle,
		uintPtr: uintptr(handle),
		maps:    locksMap{locks: make(map[uintptr]*bool), lock: &sync.RWMutex{}},
	}
}

// SetValFromString
// 设置字符类型值
func (c Ctx) SetValFromString(address uintptr, newVal string) {
	c.SetValFromBytes(address, *(*[]byte)(unsafe.Pointer(&newVal)))
}

// SetValFromBytes
// 设置字节类型值 该类型为全类型通用
func (c Ctx) SetValFromBytes(address uintptr, newVal []byte) {
	procWriteProcessMemory.Call(c.uintPtr, address, uintptr(unsafe.Pointer(&newVal[0])), uintptr(len(newVal)), 0)
}

// SetValFromInt
// 设置int类型值
func (c Ctx) SetValFromInt(address uintptr, newVal int, isLock ...bool) {
	var flag *bool
	if c.maps.locks[address] != nil {
		flag = c.maps.locks[address]
	} else {
		ff := false
		flag = &ff
		c.maps.locks[address] = flag
	}
	if len(isLock) > 0 {
		*flag = true
		timing.Anonymous(func() bool {
			if *flag {
				procWriteProcessMemory.Call(c.uintPtr, address, uintptr(unsafe.Pointer(&newVal)), uintptr(4), 0)
				return true
			} else {
				return false
			}
		}, time.Second/5/2).Start()
	} else {
		procWriteProcessMemory.Call(c.uintPtr, address, uintptr(unsafe.Pointer(&newVal)), uintptr(4), 0)
	}
}

// SetValFromFloat
// 设置float类型值
func (c Ctx) SetValFromFloat(address uintptr, newVal float64, isLock ...bool) {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, math.Float32bits(float32(newVal)))
	var flag *bool
	if c.maps.locks[address] != nil {
		flag = c.maps.locks[address]
	} else {
		ff := false
		flag = &ff
		c.maps.locks[address] = flag
	}
	if len(isLock) > 0 {
		*flag = true
		c.SetValFromBytes(address, bytes)
		timing.Anonymous(func() bool {
			c.SetValFromBytes(address, bytes)
			if *flag {
				return true
			} else {
				return false
			}
		}, time.Second/5/2).Start()
	} else {
		*c.maps.locks[address] = false
		c.SetValFromBytes(address, bytes)
	}
}

// GetValFromInt
// 读内存int类型值
func (c Ctx) GetValFromInt(address uintptr) (val int) {
	_, _, _ = procReadProcessMemory.Call(c.uintPtr, address, uintptr(unsafe.Pointer(&val)), uintptr(4), 0, 0)
	return
}

func (c Ctx) BytesToFloat(bytes []byte) float64 {
	return float64(math.Float32frombits(binary.LittleEndian.Uint32(bytes)))
}

func (c Ctx) BytesToInt(bytes []byte) int {
	if len(bytes) == 4 {
		return int(binary.LittleEndian.Uint32(bytes))
	} else {
		return int(binary.LittleEndian.Uint16(bytes))
	}
}

// GetValFromFloat
// 读内存int类型值
func (c Ctx) GetValFromFloat(address uintptr) (val float64) {
	bytes := c.GetValFromByte(address, 4)
	val = float64(math.Float32frombits(binary.LittleEndian.Uint32(bytes)))
	return
}

// GetValFromByte
// 读内存byte类型值
func (c Ctx) GetValFromByte(address uintptr, length int) []byte {
	var bt []byte
	bt = make([]byte, length)
	_, _, _ = procReadProcessMemory.Call(c.uintPtr, address, uintptr(unsafe.Pointer(&bt[0])), uintptr(length), 0, 0)
	return bt
}

// SearchValIntFrom4
// int类型搜索(4字节)
func (c Ctx) SearchValIntFrom4(val int) {
	// 待定是否使用
}

func (c Ctx) EnumProcess(moduleName string) uintptr {
	var model []uintptr
	model = append(model, 0)
	//model = make([]uintptr, 111)
	var lens uint32
	err := windows.EnumProcessModulesEx(c.handle, (*windows.Handle)(&model[0]), uint32(0), &lens, windows.LIST_MODULES_ALL)
	if err != nil {
		frame.Println(err.Error())
	}
	moduleCount := int(lens) / int(unsafe.Sizeof(model[0]))
	model = make([]uintptr, moduleCount)
	err = windows.EnumProcessModulesEx(c.handle, (*windows.Handle)(&model[0]), lens, &lens, windows.LIST_MODULES_ALL)
	if err != nil {
		frame.Println(err.Error())
		return 0
	}
	for i := 0; i < moduleCount; i++ {
		addr := model[i]
		moduleNameBuf := make([]uint16, 256)
		err := windows.GetModuleBaseName(c.handle, windows.Handle(model[i]), &moduleNameBuf[0], 256)
		if err != nil {
			fmt.Println(`跳过打印`)
			continue
		}
		if syscall.UTF16ToString(moduleNameBuf[:]) == moduleName {
			return addr
		}
	}
	return 0
}

func (c Ctx) GetModuleName(address uintptr) {
	moduleNameBuf := make([]uint16, 256)
	err := windows.GetModuleBaseName(c.handle, windows.Handle(address), &moduleNameBuf[0], 256)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(syscall.UTF16ToString(moduleNameBuf[:]))
}

func (c Ctx) Skew(address uintptr, val ...int64) uintptr {
	addr := uintptr(0)
	for _, i2 := range val {
		addr = address + uintptr(i2)
		values := 0
		_, _, _ = procReadProcessMemory.Call(c.uintPtr, addr, uintptr(unsafe.Pointer(&values)), uintptr(8), 0, 0)
		//frame.Println(frame.PrintDisAbleDebugInfo, `偏移前地址:`, fmt.Sprintf(`0X%X`, address), `+ 偏移量:`, fmt.Sprintf(`%X`, i2), `->`, `偏移后地址:`, fmt.Sprintf(`0X%X`, addr))
		address = uintptr(values)
	}
	return addr
}

func (c Ctx) Lock(address uintptr) {
	o := uint32(0)
	err := windows.VirtualProtectEx(c.handle, address, uintptr(4), windows.PAGE_READONLY, &o)
	fmt.Println(err)
}

func (c Ctx) UnLock(address uintptr) {
	o := uint32(0)
	err := windows.VirtualProtectEx(c.handle, address, uintptr(4), windows.PAGE_READWRITE, &o)
	fmt.Println(err)
}

func Int2BytesCode(num ...int) {
	for _, i2 := range num {
		arr := []byte{0, 0, 0, 0}
		binary.LittleEndian.PutUint32(arr, uint32(i2))
		str := fmt.Sprintf("%x", arr)
		count := len(str) / 2
		for i := 0; i < count; i++ {
			fmt.Print(str[i*2:(i+1)*2] + ` `)
		}
	}
}
