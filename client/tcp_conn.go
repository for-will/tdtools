package client

import (
	"encoding/binary"
	"go.uber.org/zap"
	"log"
	"net"
)

type TcpConn struct {
	conn net.Conn
}

func (tc *TcpConn) WriteMsg(b []byte) {
	m := make([]byte, 4+len(b))

	binary.LittleEndian.PutUint32(m, uint32(len(b)))

	copy(m[4:], b)
	if _, err := tc.conn.Write(m); err != nil {
		log.Fatalf("TcpConn WriteMsg error: %v", err)
	}
}

func (tc *TcpConn) ReadMsg() ([]byte, error) {
	data := make([]byte, 4)
	_, err := tc.conn.Read(data)
	if err != nil {
		return nil, err
	}

	msgLen := ByteOrder.Uint32(data)

	data = make([]byte, msgLen)
	_, err = tc.conn.Read(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (tc *TcpConn) Close() {
	if err := tc.conn.Close(); err != nil {
		log.Printf("TcpConn Close err: %v", err)
	}
}

func NewTcpConn(ServerAddr string) *TcpConn {

	conn, err := net.Dial("tcp", ServerAddr)
	if err != nil {
		Log.Fatal("NewTcpConn", zap.Error(err))
	}
	return &TcpConn{conn: conn}
}
