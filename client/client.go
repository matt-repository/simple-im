package main

import (
	"flag"
	"fmt"
	io "io"
	"net"
	"os"
	"strconv"
	"strings"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Conn       net.Conn
	Flag       int
}

//建立客户端
func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		Flag:       99,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.Conn = conn
	return client
}

var serverIp = "127.0.0.1"
var serverPort = 6999

// ./client -ip 127.0.0.1 -port 6999
func init() {

	fmt.Println("准备连接冯骎的IM服务器...")
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址（默认是127.0.0.1）")
	flag.IntVar(&serverPort, "port", 6999, "设置服务器端口是（默认是6999）")
}

//处理server回应的消息，直接显示到标准输出即可
func (client *Client) DealResponse() {
	//将收到的消息打印到控制台，永久阻塞监听
	io.Copy(os.Stdout, client.Conn)
}

//菜单
func (client *Client) menu() bool {
	var command string
	fmt.Printf("  \n")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>菜单Start<<<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Println("1.发送群消息")
	fmt.Println("2.私发消息")
	fmt.Println("3.更新用户名")
	fmt.Println("4.查询在线用户")
	fmt.Println("0.退出")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>菜单End<<<<<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Printf("  \n")
	fmt.Scanln(&command)

	flag, err := strconv.Atoi(command)
	if err != nil {
		fmt.Println("你输入的字符串必须为0~4之间的数字！！！")
		return true
	}
	if flag >= 0 && flag <= 4 {
		client.Flag = flag
		return true
	} else {
		fmt.Println("你输入的字符串必须为0~4之间的数字！！！")
		return false
	}
}

//公聊逻辑
func (client *Client) PublicChat() {
	var chatMsg string
	for {
		fmt.Println(">>>>请输入聊天内容")
		fmt.Scanln(&chatMsg)
		if strings.Trim(chatMsg, " ") == "exit" {
			break
		} else if len(chatMsg) > 0 {
			sendMsg := chatMsg + "\r\n<----------------------------------------------------------来自群消息" + "\n"
			_, err := client.Conn.Write([]byte(sendMsg))
			fmt.Println("你对大家说:", chatMsg)
			if err != nil {
				fmt.Println("conn write err:", err)
				break
			}
			break
		} else {
			fmt.Println("输入的消息不可为空")
		}
		chatMsg = ""
	}

}

//查询在线用户
func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.Conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return
	}
}

//私聊模式
func (client *Client) PrivateChat() {
	var chatMsg string
	var remoteName string
	var inputValue string
	for {
		fmt.Println(">>>>请输入聊天内容  如你想对小明（用户名）发送你好  你需输入： 小明:你好")
		fmt.Println(">>>>现在用户在线情况为")
		client.SelectUsers()
		fmt.Scanln(&inputValue)
		if strings.Trim(inputValue, " ") == "exit" {
			break
		} else if len(strings.Split(inputValue, "：")) == 2 || len(strings.Split(inputValue, ":")) == 2 {
			data := strings.Split(inputValue, ":")
			if len(data) <= 1 {
				data = strings.Split(inputValue, "：")
			}
			remoteName = data[0]
			chatMsg = data[1] + "\r\n<----------------------------------------------------------来自私人消息"
			sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
			fmt.Println("你对【" + remoteName + "】说:" + data[1])
			_, err := client.Conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err:", err)
				break
			}
			break
		} else {
			fmt.Println("您输入的消息不符合规范,请重新输入!!!")
		}
		remoteName = ""
		chatMsg = ""
	}

}

//修改用户名
func (client *Client) UpdateName() bool {
	fmt.Println(">>>请输入您的新用户名：")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "|"
	_, err := client.Conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}

	return true
}

//客户端启动主程序
func (client *Client) Run() {
	for client.Flag != 0 { //==0会退出
		for client.menu() != true {

		}
		switch client.Flag {
		case 1:
			client.PublicChat()
			//fmt.Println("公聊模式")

		case 2:
			client.PrivateChat()
			//fmt.Println("私聊模式")

		case 3:
			client.UpdateName()
			//fmt.Println("更新用户名")

		case 4:
			client.SelectUsers()
		}
		client.Flag = 99
	}
}

func main() {
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("ip:", serverIp, "port:", serverPort)
		fmt.Println(">>>>连接服务器失败。。。")
		select {}
		return
	}

	fmt.Println(">>>>>连接服务器成功")

	go client.DealResponse()
	//这里阻塞一下 。防止直接结束
	//select {}
	client.Run()
}
