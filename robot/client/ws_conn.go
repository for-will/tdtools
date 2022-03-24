package client

import (
	"golang.org/x/net/websocket"
	"log"
)

type Connect interface {
	WriteMsg(b []byte)
	ReadMsg() ([]byte, error)
	Close()
}

type WsConn struct {
	conn *websocket.Conn
}

func (w *WsConn) WriteMsg(b []byte) {
	//log.Print(b)
	if _, err := w.conn.Write(b); err != nil {
		log.Fatalf("WsConn WriteMsg error: %v", err)
	}
}

func (w *WsConn) ReadMsg() ([]byte, error) {
	var b = make([]byte, 16*1024)
	n, err := w.conn.Read(b)
	if err != nil {
		//log.Fatalf("WsConn ReadMsg err: %v", err)
		return nil, err
	}
	return b[:n], nil
}

func (w *WsConn) Close() {
	w.conn.Close()
}

//func Connect() {
//	con, err := websocket.Dial("ws://172.16.1.218:3653/", "", "http://localhost/")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer con.Close()
//}

func NewWsConn(url string) *WsConn {
	con, err := websocket.Dial(url, "", "http://localhost/")
	if err != nil {
		log.Fatal(err)
	}
	return &WsConn{conn: con}
}
