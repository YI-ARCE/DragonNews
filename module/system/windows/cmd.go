package windows

import (
	"fmt"
	"os"
	"os/exec"
)

// OpenDirOrFileWindow 打开文件/夹窗口
func OpenDirOrFileWindow(path string) {
	cmd := exec.Command("cmd", "/c", "explorer", path)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// WaitToArbitraryKey 等待任意键输入后继续执行
func WaitToArbitraryKey(str string) {
	fmt.Println(str + `,按任意键返回...`)
	fmt.Scanf("%d\n")
}

func WaitToArbitraryKey2(str string) {
	fmt.Println(str)
	fmt.Scanf("%d\n")
}

func Clear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
