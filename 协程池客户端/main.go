package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

// 设置console的Title方便区分用户名
func SetCmdTitle(title string) {
	kernel32, _ := syscall.LoadLibrary(`kernel32.dll`)
	sct, _ := syscall.GetProcAddress(kernel32, `SetConsoleTitleW`)
	syscall.Syscall(sct, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
	syscall.FreeLibrary(kernel32)
}

// 客户端
func main() {
	fmt.Printf("输入Console Title:")
	var name string
	fmt.Scanln(&name)
	SetCmdTitle(name)

	conn, err := net.Dial("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("err :", err)
		return
	}

	defer conn.Close() // 关闭连接
	inputReader := bufio.NewReader(os.Stdin)
	for {
		input, _ := inputReader.ReadString('\n') // 读取用户输入
		inputInfo := strings.Trim(input, "\r\n")
		_, err = conn.Write([]byte(fmt.Sprintf("[%s]", name) + inputInfo)) // 发送数据
		if err != nil {
			return
		}
		buf := [512]byte{}
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("recv failed, err:", err)
			return
		}
		fmt.Println(string(buf[:n]))
	}
}
