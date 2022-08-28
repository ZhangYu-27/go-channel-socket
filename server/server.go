package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
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

var userMap map[string]chan string
var addrMap map[net.Addr]bool
var addrSlice []string

func init() {
	userMap = make(map[string]chan string)
	addrMap = make(map[net.Addr]bool)
}

func makeJson(typ string, info string) []byte {
	addrSlice = make([]string, 0)
	json, err := json.Marshal(AddrList{typ, info})
	if err != nil {
		fmt.Println("转换失败", err)
		return []byte{}
	}
	return json
}

func sendServerList(conn net.Conn) {

	for key, value := range addrMap {
		if value == true {
			addrSlice = append(addrSlice, key.String())
		}
	}
	str, _ := json.Marshal(addrSlice)
	conn.Write(makeJson("serverList", string(str)))
}

func process(conn net.Conn) {
	addr := conn.RemoteAddr()
	userMap[addr.String()] = make(chan string)
	addrMap[addr] = true
	defer conn.Close()
	//向客户端发送初始化信息
	go sendServerList(conn)
	//下面这个匿名函数是将其他人放入管道的消息发送给客户端
	go func() {
		fmt.Println("读取管道的信息")
		for {
			fmt.Println(userMap)

			select {
			case str := <-userMap[addr.String()]:
				fmt.Println("开始发送")
				jso := makeJson("message", str)
				fmt.Println(jso)
				num, err := conn.Write(jso)
				if err != nil {
					fmt.Println("-------", err)
					return
				}
				fmt.Printf("向客户端%v发送了%v个数据", conn.RemoteAddr(), num)
			}
		}

	}()
	for {
		buf := make([]byte, 1024)
		fmt.Println("服务器在等", addr)
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				delete(addrMap, addr)
				delete(userMap, addr.String())
				fmt.Println("客户端退出")
				return
			}
			fmt.Println("错误信息：=", err)
			return
		}
		str := string(buf[:n])
		fmt.Printf(str)
		info := &InputData{}
		json.Unmarshal([]byte(str), info)
		fmt.Println(info)
		if info.Type == "serverlist" {
			go sendServerList(conn)
			fmt.Println(str)
			continue
		} else {
			fmt.Println("放入管道数据")
			userMap[info.InputAddr] <- info.Message
		}
	}
}

func main() {
	fmt.Println("服务器开始监听")
	listen, err := net.Listen("tcp", "0.0.0.0:47777")
	if err != nil {
		fmt.Println("监听错误", err)
		return
	}
	defer listen.Close()
	for {
		fmt.Println("等待客户端链接")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("链接失败")
		} else {
			fmt.Println(conn.RemoteAddr())
			go process(conn)
		}
	}
}
