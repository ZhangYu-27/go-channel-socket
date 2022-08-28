package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type AddrList struct {
	Type string `json:"type"`
	Info string `json:"info"`
}

type InputData struct {
	Type      string `json:"type"`
	InputAddr string `json:"input_addr"`
	Message   string `json:"message"`
}

var sendAddr string
var addrSlice []string

func printAddr(str string, myAddr string) {
	addrSlice = []string{}
	json.Unmarshal([]byte(str), &addrSlice)
	fmt.Println("当前在线的服务列表，输入序号开始发送消息")
	for i := 0; i < len(addrSlice); i++ {
		if addrSlice[i] != myAddr {
			fmt.Printf("%v.  %v\n", i, addrSlice[i])
		}
	}
	fmt.Println("请输入序号选择")
}

func main() {
	fmt.Println("开始链接 长时间无响应就是")
	conn, err := net.Dial("tcp", "101.33.233.165:47777")
	if err != nil {
		fmt.Println("", err)
		return
	}
	fmt.Println("------------------------------")
	fmt.Println("欢迎使用能用级别的聊天软件,键入exit退出")
	fmt.Println("本机地址与端口", conn.LocalAddr())
	fmt.Println("------------------------------")
	go func() {
		for {
			buf := make([]byte, 1024)
			num, err := conn.Read(buf)
			if err != nil {
				fmt.Println("", err)
				return
			}
			serverInfo := AddrList{}
			json.Unmarshal(buf[:num], &serverInfo)
			fmt.Println("接受数据中")
			if serverInfo.Type == "serverList" {
				printAddr(serverInfo.Info, conn.LocalAddr().String())
			} else {
				fmt.Println("------------------------------")
				fmt.Println("他说", serverInfo.Info)
				fmt.Println("------------------------------")
			}
		}

	}()

	for {
		reader := bufio.NewReader(os.Stdin)
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("", err)
			return
		}
		str = strings.Trim(str, "\n")
		if str == "exit" {
			fmt.Println("期待您下次使用")
			break
		}
		sendInfo := &InputData{}
		if sendAddr == "" {
			inputInt, err := strconv.Atoi(str)
			if err != nil {
				fmt.Println(err)
				str = "serverlist"
				sendInfo = &InputData{
					Type: "serverlist",
				}
			} else {
				sendAddr = addrSlice[inputInt]
				fmt.Println("设置成功，当前发送对象", sendAddr, "输入信息回车发送")
				continue
			}
		} else {
			sendInfo = &InputData{
				Type:      "message",
				InputAddr: sendAddr,
				Message:   str,
			}
		}
		jso, err := json.Marshal(sendInfo)
		_, err = conn.Write(jso)
		if err != nil {
			fmt.Println("发送失败", err)
		}
		fmt.Println("+++++++++++++++++++")
		fmt.Println("你说", strings.Trim(str, "\n"))
		fmt.Println("+++++++++++++++++++")
	}
}
