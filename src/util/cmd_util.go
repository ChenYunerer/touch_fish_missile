package util

import "fmt"

//print normal msg to cmd
func PrintMsgToCmd(a ...interface{}) {
	fmt.Println(a...)
}

//print system info or system msg to cmd
func PrintSysNotifyToCmd(a ...interface{}) {
	fmt.Println("⬇⬇⬇⬇⬇⬇⬇⬇⬇⬇⬇⬇⬇⬇⬇⬇⬇⬇⬇⬇⬇⬇⬇⬇⬇⬇⬇")
	fmt.Println("SYS NOTIFY：", a[:])
	fmt.Println("⬆⬆⬆⬆⬆⬆⬆⬆⬆⬆⬆⬆⬆⬆⬆⬆⬆⬆⬆⬆⬆⬆⬆⬆⬆⬆⬆")
}
