package main

import (
	"fmt"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pkg/pool/goroutine"
	"goland-IM-System/server/comm/frame"
	"goland-IM-System/server/log"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Server struct {
	*gnet.EventServer
	Ip        string
	Port      int
	multicore  bool
	async      bool
	codec      gnet.ICodec
	workerPool *goroutine.Pool

	OnlineMap map[string]*User  //z
	//mapLock   sync.RWMutex           //用户和消息的读写锁
	//Message   chan map[string]string //发送消息的用户add：消息

}

//创建服务端
func CustomCodecServe(ip string, port int,multicore,async bool , codec gnet.ICodec)  {
	codec = frame.Frame{}
	server := &Server{
		_,
		ip,
		port,
		multicore,
		async,
		codec,
		goroutine.Default(),
		make(map[string]*User),
		//sync.RWMutex{},
		//make(chan map[string]string),
	}
	addr:="tcp://"+ip+":"+strconv.Itoa( port)
	err:= gnet.Serve(server, addr, gnet.WithMulticore(multicore), gnet.WithTCPKeepAlive(time.Minute*5), gnet.WithCodec(codec))
	if err!=nil{
		log.Logger.Error.Println(err)
	}
}

//监听消息
func (t *Server) ListenMessage() {
	for {
		msgModel := <-t.Message

		t.mapLock.Lock()
		for _, cli := range t.OnlineMap {
			cli.C <- msgModel
		}
		t.mapLock.Unlock()
	}
}

//广播消息
func (t *Server) BroadCast(user *User, msg string) {
	sendMsg := "【" + user.Name + "】:" + msg

	log.Logger.Info.Println("广播消息:" + sendMsg)

	log.Logger.Info.Println("此时在线用户有：")
	for _, v := range t.OnlineMap {
		log.Logger.Info.Println(strings.Replace(v.Name, "\n", "", -1))

	}

	t.Message <- map[string]string{
		user.Addr: sendMsg,
	}
}

//接收连接后的操作
func (s *Server) Hanlder(conn net.Conn) {
	user := NewUser(conn, t)
	user.Online()

	//处理客户端发送的消息，进行广播

	//判断此用户是否活跃
	isLive := make(chan bool)

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				log.Logger.Error.Println("Conn Read err:" + err.Error())
				return
			}
			//msg:=string(buf[:n-1])  //这里去除最后一个换行符，其实可以不用去掉
			msg := string(buf[:n])
			//客户端发送消息进行处理
			user.DoMessage(msg)
			isLive <- true
		}
	}()

	//阻塞当前handler //  因为这里涉及到 channel 要阻塞那边能够拿到才能结束这个goroutine
	for {
		//要么获取isLive，要么开始计时是否有10秒，当isLive时。它就重新计时了
		select {
		case <-isLive:
			//不做任何操作
		case <-time.After(time.Minute * 10):
			user.SendMsg("您已10分钟不发送消息，超时了已强制将您退出")
			close(user.C)
			conn.Close()
			return
		}
	}

}

//接受消息
func (cs *Server) React(framePayload []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	user := NewUser(c, cs)
	user.Online()
	// packet decode
	fmt.Println(string(framePayload))
	msg := string(framePayload)
	//客户端发送消息进行处理
	user.DoMessage(msg)
	out=framePayload
	return
}

//服务端启动
func (t *Server) Start() {
	//socket listen
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", t.Ip, t.Port))
	if err != nil {
		log.Logger.Error.Println("net.Listen err:" + err.Error())
		return
	}

	//close listen socket
	defer listen.Close()
	//监听 客户端上线的消息
	go t.ListenMessage()
	for {
		//accept
		conn, err1 := listen.Accept()
		if err1 != nil {
			log.Logger.Error.Println("listener accept err:" + err.Error())
			continue
		}
		//do handler
		go t.Hanlder(conn)
	}

}
