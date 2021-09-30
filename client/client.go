package client

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"os"
	"os/signal"
	"reflect"
	"robot/GameMsg"
	"sync"
)

type Client struct {
	C          net.Conn
	sendQ      chan proto.Message
	msgHandler func(GameMsg.MsgId, proto.Message)
}

func (r *Client) Init() {
	r.sendQ = make(chan proto.Message, 10)
	conn, err := net.Dial("tcp", ServerAddr)
	if err != nil {
		Log.Fatal("client init fail", zap.Error(err))
	}
	r.C = conn

	r.msgHandler(NetworkConnected, nil)
}

func (r *Client) Run() {
	wg := &sync.WaitGroup{}
	//closeSign := make(chan struct{})

	wg.Add(1)
	go func() {
		for {
			if err := r.ReadMsg(); err != nil {
				Log.Error("ReadMsg error", zap.Error(err))
				Log.Info("ReadMsg goroutine exited")
				wg.Done()
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		r.writeMsgLoop()
		Log.Info("SendMsg goroutine exited")
		wg.Done()
	}()

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt)

	//<-time.After(20 * time.Second)
	s := <-signalCh
	Log.Info("received signal", zap.Any("signal", s))
	if err := r.C.Close(); err != nil {
		log.Fatalf("conn.Close error: %v", err)
	}
	//close(closeSign)
	close(r.sendQ)
	close(signalCh)

	wg.Wait()
}

func (r *Client) SendMsg(msg proto.Message) {
	r.sendQ <- msg
}

func (r *Client) writeMsgLoop() {

	for msg := range r.sendQ {
		reflect.TypeOf(msg)
		id, ok := messageId[reflect.TypeOf(msg)]
		if !ok {
			Log.Error("Not registered message", zap.Any("msg_type", reflect.TypeOf(msg)))
			continue
		}
		_msg, _ := proto.Marshal(msg)
		msgLen := 4 + len(_msg)
		m := make([]byte, 4+msgLen)

		// 默认使用大端序
		ByteOrder.PutUint32(m, uint32(msgLen))
		ByteOrder.PutUint32(m[4:], uint32(id))
		copy(m[8:], _msg)

		// 发送消息
		//todo: handle error
		if _, err := r.C.Write(m); err != nil {
			log.Fatalf("conn.Write error: %v", err)
		}

		fmt.Printf("> %-30s| %s\n", id, JsonString(msg))

		//Log.Debugw("send message", zap.Any("MsgID", id.String()), zap.String("message", JsonString(msg)))
	}
}

func (r *Client) ReadMsg() error {
	data := make([]byte, 4)
	_, err := r.C.Read(data)
	if err != nil {
		return err
	}

	msgLen := ByteOrder.Uint32(data)

	data = make([]byte, msgLen)
	_, err = r.C.Read(data)
	if err != nil {
		return err
	}

	id := ByteOrder.Uint32(data)
	msgId := GameMsg.MsgId(id)
	typ, ok := messageType[msgId]
	if !ok {
		fmt.Printf("unknown msg id %v(%v)\n", msgId, id)
		return nil
	}
	msg := reflect.New(typ.Elem()).Interface().(proto.Message)

	if err := proto.Unmarshal(data[4:msgLen], msg); err != nil {
		//log.Fatalf("proto.Unmarshal %+v msgLen:%v error:%v", msgId, msgLen, err)
		Log.Fatal("proto.Unmarshal", zap.Int32("msg_id", int32(msgId)),
			zap.Uint32("msg_len", msgLen), zap.Error(err))
	}

	//fmt.Printf("MsgLen: %d\n", msgLen)

	fmt.Printf("< %-30s| %v\n", msgId, JsonString(msg))
	if code := fetchReturnCode(msg); code != GameMsg.ReturnCode_OK {
		Log.Debug("Request failed", zap.Any("ReturnCode", code))
	}

	r.msgHandler(msgId, msg)
	return nil
}

func fetchReturnCode(msg interface{}) GameMsg.ReturnCode {
	typ := reflect.TypeOf(msg).Elem()
	val := reflect.ValueOf(msg).Elem()
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Type == reflect.TypeOf(new(GameMsg.ReturnCode)) {
			return val.Field(i).Elem().Interface().(GameMsg.ReturnCode)
		}
	}
	return GameMsg.ReturnCode_OK
}