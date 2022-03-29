package client

import (
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"log"
	"market/GameMsg"
	"os"
	"os/signal"
	"reflect"
	"sync"
)

type Client struct {
	ServerTCP  string
	ServerWS   string
	C          Connect
	sendQ      chan proto.Message
	msgHandler func(GameMsg.MsgId, proto.Message)
}

func (r *Client) Init() {
	r.sendQ = make(chan proto.Message, 10)

	if r.ServerTCP != "" {

		r.C = NewTcpConn(r.ServerTCP)
	} else if r.ServerWS != "" {
		r.C = NewWsConn(r.ServerWS)
	} else {
		log.Fatal("Server addr invalid")
	}

	r.msgHandler(NetworkConnected, nil)
	//log.SetFlags(0)
}

func (r *Client) Run() {
	wg := &sync.WaitGroup{}

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
	r.C.Close()
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
		m := make([]byte, 4+len(_msg))

		ByteOrder.PutUint32(m, uint32(id))
		copy(m[4:], _msg)

		//log.Printf("\x1b[40m> %-30s| %s\x1b[0m\n", id, Msg2Json(msg))
		LogSndMsg(id, msg)
		// 发送消息
		//todo: handle error
		r.C.WriteMsg(m)
	}
}

func (r *Client) ReadMsg() error {

	data, err := r.C.ReadMsg()
	if err != nil {
		return err
	}

	id := ByteOrder.Uint32(data)
	msgId := GameMsg.MsgId(id)
	typ, ok := messageType[msgId]
	if !ok {
		Log.Error("unknown msg id", zap.String("id", msgId.String()), zap.Uint32("value", id))
		return nil
	}
	msg := reflect.New(typ.Elem()).Interface().(proto.Message)

	if err := proto.Unmarshal(data[4:], msg); err != nil {
		//log.Fatalf("proto.Unmarshal %+v msgLen:%v error:%v", msgId, msgLen, err)
		Log.Fatal("proto.Unmarshal", zap.Int32("msg_id", int32(msgId)),
			zap.Int("msg_len", len(data)), zap.Error(err))
	}
	if code := fetchReturnCode(msg); code != GameMsg.ReturnCode_OK {
		LogErrMsg(msgId, msg)
	} else {
		LogRcvMsg(msgId, msg)
	}

	r.msgHandler(msgId, msg)
	return nil
}

type ResponseMessage interface {
	GetReturnCode() GameMsg.ReturnCode
}

func fetchReturnCode(msg interface{}) GameMsg.ReturnCode {
	//typ := reflect.TypeOf(msg).Elem()
	//val := reflect.ValueOf(msg).Elem()
	//ReturnCodeType := reflect.TypeOf(GameMsg.ReturnCode(0))
	//for i := 0; i < typ.NumField(); i++ {
	//	if typ.Field(i).Type == ReturnCodeType {
	//		return val.Field(i).Interface().(GameMsg.ReturnCode)
	//	}
	//}

	if rsp, ok := msg.(ResponseMessage); ok {
		return rsp.GetReturnCode()
	}
	return GameMsg.ReturnCode_OK
}
