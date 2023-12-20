package webssh

import (
	"fmt"
	"fqhWeb/internal/consts"
	"fqhWeb/pkg/api"
	"time"

	"github.com/gorilla/websocket"
)

// 发送命令给机器(读取前端的命令发送机器)
func (s *WebSshService) Recv(websshConn *websocket.Conn, sshConn *api.SSHConnect, quit chan int) {
	defer Quit(quit)
	var (
		bytes []byte
		err   error
	)
	websshConn.SetReadDeadline(time.Now().Add(consts.PongPeriod))
	websshConn.SetPongHandler(func(appData string) error {
		websshConn.SetReadDeadline(time.Now().Add(consts.PongPeriod))
		return nil
	})

	for {
		if _, bytes, err = websshConn.ReadMessage(); err != nil {
			return
		}
		if len(bytes) > 0 {
			if _, e := sshConn.StdinPipe.Write(bytes); e != nil {
				return
			} else {
				fmt.Println("===========Recv==========" + string(bytes))
			}
		}
	}
}

// 读取机器给我的返回信息(读取机器的信息发送给前端)
func (s *WebSshService) Output(websshConn *websocket.Conn, sshConn *api.SSHConnect, quit chan int) {
	defer Quit(quit)
	var (
		read int
		err  error
	)
	tickPing := time.NewTicker(consts.PingPeriod)
	defer tickPing.Stop()
	tick := time.NewTicker(60 * time.Millisecond)
	defer tick.Stop()
Loop:
	// 无限循环直到退出
	for {
		select {
		case <-tick.C:
			i := make([]byte, 1024)
			if read, err = sshConn.StdoutPipe.Read(i); err != nil {
				fmt.Println(err)
				break Loop
			}
			if err = s.WebSshSendText(websshConn, i[:read]); err != nil {
				fmt.Println(err)
				break Loop
			} else {
				fmt.Println("===========Output==========" + string(i[:read]))
			}
		case <-tickPing.C:
			if err := websshConn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				fmt.Println(err)
				break Loop
			}
		}
	}
}

func Quit(quit chan<- int) {
	quit <- 1
}
