package main

import (
	"goland-IM-System/server/log"
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan map[string]string //发送消息的用户add：消息
	conn   net.Conn
	server *Server
}

//创建一个用户
func NewUser(conn net.Conn, server *Server) *User {

	userAddr := conn.RemoteAddr().String()

	user := &User{
		userAddr,
		userAddr,
		make(chan map[string]string),
		conn,
		server,
	}
	//当启动一个用户时，就会监听消息。
	go user.ListenMessage()
	return user
}

//用户上线业务
func (t *User) Online() {
	t.server.mapLock.Lock()
	t.server.OnlineMap[t.Name] = t
	t.server.mapLock.Unlock()

	t.server.BroadCast(t, "【在线】")
}

//用户的下线业务

func (t *User) Offline() {
	t.server.mapLock.Lock()
	delete(t.server.OnlineMap, t.Name)
	t.server.mapLock.Unlock()

	t.server.BroadCast(t, "【离线】")
}

//用户发送消息
func (t *User) DoMessage(msg string) {

	//fmt.Println(msg)
	//查询多少用户在线
	//if msg=="who"
	if strings.HasPrefix(msg, "who") {
		t.server.mapLock.Lock()
		for _, user := range t.server.OnlineMap {
			onlineMsg := "（其他人）ip:【" + user.Addr + "】" + "用户名:【" + user.Name + "】:" + "【在线】; \n"
			if user.Addr == t.Addr {
				onlineMsg = "（你自己）ip:【" + user.Addr + "】" + "用户名:【" + user.Name + "】:" + "【在线】; \n"
			}

			t.SendMsg(onlineMsg)
		}
		t.server.mapLock.Unlock()
		//修改名称
	} else if len(msg) > 7 && msg[:7] == "rename|" { //这里得输入rename|名称| 因为后面会带不知道空格还是啥的东西
		newName := strings.Split(msg, "|")[1]
		_, ok := t.server.OnlineMap[newName]
		if ok {
			t.SendMsg("此用户名已被使用！！！\n")
		} else {
			t.server.mapLock.Lock()
			delete(t.server.OnlineMap, t.Name)

			t.server.OnlineMap[newName] = t
			t.server.mapLock.Unlock()
			t.Name = newName
			t.SendMsg("您的用户名已改变成:【" + newName + "】\n")
		}
		//私聊协议
	} else if len(msg) > 4 && msg[:3] == "to|" {
		//消息格式：to|张三|消息内容

		//1 获取对方用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			t.SendMsg("您的消息发送的格式不正确\n")
			return
		}

		//2 根据用户名 得到对方User对象
		remoteUser, ok := t.server.OnlineMap[remoteName]
		if !ok {
			t.SendMsg("此用户 【" + remoteName + "】 不存在")
			return
		}
		//3 获取消息内容，通过对方的User 对象将消息发送
		content := strings.Split(msg, "|")[2]

		if content == "" {
			t.SendMsg("消息不可为空\n")
			return
		}
		if remoteUser.Name == t.Name {
			remoteUser.SendMsg("不允许给自己发送消息！！！")
		} else {
			remoteUser.SendMsg("【" + t.Name + "】" + ":" + content)
		}

		//发送消息 广播除了自己的
	} else {
		t.server.BroadCast(t, msg)
	}
}

func (t *User) SendMsg(msg string) {
	log.Logger.Info.Println("客户端 " + t.Addr + "  【" + t.Name + "】   " + "  收到消息： " + msg)
	t.conn.Write([]byte(msg + "\n"))
}

//监听当前User channel 的方法，一旦有消息，就直接发送给客户端
func (t *User) ListenMessage() {
	for {
		msgModel := <-t.C
		var sendUser = ""
		var sendMsg = ""
		for k, v := range msgModel { //这里只会执行一次 如果msgModel
			sendUser = k
			sendMsg = v

			//只对不是自己发送消息
			if sendUser != t.Addr {
				t.SendMsg(sendMsg + "\n")
			}
		}

	}
}
