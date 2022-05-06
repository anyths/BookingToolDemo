package modules

import (
	"bookcourt/utils"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func SelectMode(l *utils.Logs) string {
	var input string
	for {
		fmt.Println("请选择模式: 0 - 自动模式 / 1 - 立即执行  / 2 - 锁单模式")
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("请入0/1/2: ")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "0" || input == "1" || input == "2" || input == "3" {
			break
		} else {
			l.Err("无法识别输入: " + input)
			continue
		}
	}
	return input
}
