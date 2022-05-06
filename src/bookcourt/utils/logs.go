package utils

import (
	"fmt"
	"os"
)

type Logs []byte

func (l *Logs) Add(s string) {
	fmt.Printf("[logs]: %s\n", s)
	str := GetNowTime() + " >> [logs]: " + s + "\n"
	*l = append(*l, []byte(str)...)
}
func (l *Logs) Err(s string) {
	fmt.Printf("[!Err]: %s\n", s)
	str := GetNowTime() + " >> [!Err]: " + s + "\n"
	*l = append(*l, []byte(str)...)
}
func (l *Logs) Write() {
	path := GetLogFile()
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("log写入失败...")
		return
	}
	defer file.Close()

	_, err = file.Write(*l)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	*l = []byte{}
}
