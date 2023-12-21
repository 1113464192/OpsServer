package webssh

import (
	"bytes"
	"encoding/json"
	"fqhWeb/internal/consts"
	"fqhWeb/internal/model"
	"fqhWeb/pkg/api"
	"time"

	"github.com/gorilla/websocket"
)

func flushCombOutput(w *WebsshBufferWriter, wsConn *websocket.Conn) error {
	if w.Buffer.Len() != 0 {
		err := wsConn.WriteMessage(websocket.BinaryMessage, w.Buffer.Bytes())
		if err != nil {
			return err
		}
		w.Buffer.Reset()
	}
	return nil
}

func (s *SSHConnect) wsQuit(ch chan struct{}) {
	s.Once.Do(func() {
		close(ch)
	})
}

// 向websocket发送服务器返回的信息
func (s *SSHConnect) WsSend(wsConn *websocket.Conn, quitCh chan struct{}) {
	defer s.wsQuit(quitCh)
	pingTick := time.NewTicker(consts.WebsshPingPeriod)
	defer pingTick.Stop()

	tick := time.NewTicker(consts.WebsshReadMessageTickerDuration)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			//write combine output bytes into websocket response
			if err := flushCombOutput(s.CombineOutput, wsConn); err != nil {
				if e := WebSsh().WebSshSendErr(wsConn, "发送服务器返回信息到websocket失败: "+err.Error()); e != nil {
					s.Logger.Error("Webssh", "发送错误信息至websocket失败", err)
				}
				s.Logger.Error("Webssh", "发送服务器返回信息到websocket失败", err)
				return
			}
			// 发送ping至websocket
		case <-pingTick.C:
			if err := wsConn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				if e := WebSsh().WebSshSendErr(wsConn, "发送Ping至websocket失败: "+err.Error()); e != nil {
					s.Logger.Error("Webssh", "发送错误信息至websocket失败", err)
				}
				s.Logger.Error("Webssh", "发送Ping至websocket失败", err)
				return
			}

		case <-quitCh:
			return
		}
	}
}

func (s *SSHConnect) WsRec(wsConn *websocket.Conn, quitCh chan struct{}) {
	//tells other go routine quit
	defer s.wsQuit(quitCh)
	// 处理pong消息
	wsConn.SetReadDeadline(time.Now().Add(consts.WebsshPongWait))
	wsConn.SetPongHandler(func(appData string) error {
		// 重置读取存活时间
		wsConn.SetReadDeadline(time.Now().Add(consts.WebsshPongWait))
		return nil
	})

	for {
		select {
		case <-quitCh:
			return
		default:
			// read websocket msg
			_, wsData, err := wsConn.ReadMessage()
			if err != nil {
				if e := WebSsh().WebSshSendErr(wsConn, "接收websocket发送的信息失败: "+err.Error()); e != nil {
					s.Logger.Error("Webssh", "发送错误信息至websocket失败", err)
				}
				s.Logger.Error("Webssh", "接收websocket发送的信息失败", err)
				return
			}

			// 每次传输一个或多个char
			if len(wsData) > 0 {
				// resize 或者 粘贴
				resize := api.WindowSize{}
				err = json.Unmarshal(wsData, &resize)
				if err != nil {
					goto SEND
				}
				if resize.Hight > 0 && resize.Weight > 0 {
					if err = s.Session.WindowChange(resize.Hight, resize.Weight); err != nil {
						if e := WebSsh().WebSshSendErr(wsConn, "变更WindowSize失败: "+err.Error()); e != nil {
							s.Logger.Error("Webssh", "发送错误信息至websocket失败", err)
						}
						s.Logger.Error("Webssh", "变更WindowSize失败", err)
						break
					}
				} else {
					goto SEND
				}
				break
			}
			// 服务器的返回发送给websocket
		SEND:
			if _, err = s.StdinPipe.Write(wsData); err != nil {
				if e := WebSsh().WebSshSendErr(wsConn, "发送服务器信息到前端失败: "+err.Error()); e != nil {
					s.Logger.Error("Webssh", "发送错误信息至websocket失败", err)
				}
				s.Logger.Error("Webssh", "发送服务器信息到前端失败", err)
				break
			}
			// 入库
			decodeBytes := bytes.TrimSpace(wsData)
			if len(decodeBytes) > consts.WebsshMaxRecordLength {
				decodeBytes = decodeBytes[:consts.WebsshMaxRecordLength]
			}
			websshRecord := &model.WebsshRecord{
				UserId: s.User.ID,
				HostId: s.Host.ID,
				Ipv4:   s.Host.Ipv4,
				Ipv6:   s.Host.Ipv6,
				Actlog: decodeBytes,
			}
			if err = WebSsh().RecordCreate(websshRecord); err != nil {
				if e := WebSsh().WebSshSendErr(wsConn, "用户操作记录入库失败: "+err.Error()); e != nil {
					s.Logger.Error("Webssh", "发送错误信息至websocket失败", err)
				}
				s.Logger.Error("Webssh", "用户操作记录入库失败", err)
				return
			}
		}
	}
}

func (s *SSHConnect) SessionWait(quitChan chan struct{}) {
	if err := s.Session.Wait(); err != nil {
		s.Logger.Error("Webssh", "Session.Wait报错", err)
		s.wsQuit(quitChan)
	}
}

// // 发送命令给机器(读取前端的命令发送机器)
// func (s *WebSshService) Recv(websshConn *websocket.Conn, sshConn *api.SSHConnect, quit chan int) {
// 	defer Quit(quit)
// 	var (
// 		bytes []byte
// 		err   error
// 	)
// 	websshConn.SetReadDeadline(time.Now().Add(consts.PongPeriod))
// 	websshConn.SetPongHandler(func(appData string) error {
// 		websshConn.SetReadDeadline(time.Now().Add(consts.PongPeriod))
// 		return nil
// 	})

// 	for {
// 		if _, bytes, err = websshConn.ReadMessage(); err != nil {
// 			return
// 		}
// 		if len(bytes) > 0 {
// 			if _, e := sshConn.StdinPipe.Write(bytes); e != nil {
// 				return
// 			} else {
// 				fmt.Println("===========Recv==========" + string(bytes))
// 			}
// 		}
// 	}
// }

// // 读取机器给我的返回信息(读取机器的信息发送给前端)
// func (s *WebSshService) Output(websshConn *websocket.Conn, sshConn *api.SSHConnect, quit chan int) {
// 	defer Quit(quit)
// 	var (
// 		read int
// 		err  error
// 	)
// 	tickPing := time.NewTicker(consts.PingPeriod)
// 	defer tickPing.Stop()
// 	tick := time.NewTicker(60 * time.Millisecond)
// 	defer tick.Stop()
// Loop:
// 	// 无限循环直到退出
// 	for {
// 		select {
// 		case <-tick.C:
// 			i := make([]byte, 1024)
// 			if read, err = sshConn.StdoutPipe.Read(i); err != nil {
// 				fmt.Println(err)
// 				break Loop
// 			}
// 			if err = s.WebSshSendText(websshConn, i[:read]); err != nil {
// 				fmt.Println(err)
// 				break Loop
// 			} else {
// 				fmt.Println("===========Output==========" + string(i[:read]))
// 			}
// 		case <-tickPing.C:
// 			if err := websshConn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
// 				fmt.Println(err)
// 				break Loop
// 			}
// 		}
// 	}
// }

func Quit(quit chan<- int) {
	quit <- 1
}
