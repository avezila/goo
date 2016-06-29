package proxy

import (
	"encoding/binary"
	"fmt"
	"net/http"
	"os"

	"./grpc/saver"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
	"google.golang.org/grpc"
)

type Proxy struct {
	client saver.SaverClient
}

func New(server *grpc.Server) (*Proxy, error) {
	port := os.Getenv("WS_PORT")
	if port == "" {
		port = "5354"
	}
	proxy := &Proxy{}
	RegisterProxyServer(server, proxy)

	conn, err := grpc.Dial("saver:5353", grpc.WithInsecure())
	if err != nil {
		if err != nil {
			return nil, err
		}
	}
	proxy.client = saver.NewSaverClient(conn)

	http.Handle("/grpc0", websocket.Handler(proxy.GrpcProxyPushRow))
	http.Handle("/grpc1", websocket.Handler(proxy.GrpcProxyFilterRows))
	go func() {
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}()
	return proxy, nil
}

// Echo the data received on the WebSocket.
func (p *Proxy) GrpcProxyPushRow(ws *websocket.Conn) {
	for {
		buf := make([]byte, 4)
		_, err := ws.Read(buf)
		if err != nil {
			ws.Close()
			return
		}
		len := binary.LittleEndian.Uint32(buf)

		buf = make([]byte, len)
		ws.Read(buf)
		if err != nil {
			ws.Close()
			fmt.Println("err", err)
			return
		}
		req := &saver.PushRowReq{}
		if err := proto.Unmarshal(buf, req); err != nil {
			ws.Close()
			fmt.Println("err", err)
			return
		}

		res, err := p.client.PushRow(context.Background(), req)
		if err != nil {
			ws.Close()
			fmt.Println("err", err)
			return
		}

		buf, err = proto.Marshal(res)
		if err != nil {
			ws.Close()
			fmt.Println("err", err)
			return
		}

		websocket.Message.Send(ws, buf)
	}
	ws.Close()
}
func (p *Proxy) GrpcProxyFilterRows(ws *websocket.Conn) {
	for {
		buf := make([]byte, 4)
		_, err := ws.Read(buf)
		if err != nil {
			ws.Close()
			return
		}
		len := binary.LittleEndian.Uint32(buf)
		buf = make([]byte, len)

		ws.Read(buf)
		if err != nil {
			ws.Close()
			fmt.Println("err", err)
			return
		}
		req := &saver.FilterRowsReq{}
		if err := proto.Unmarshal(buf, req); err != nil {
			ws.Close()
			fmt.Println("err", err)
			return
		}

		res, err := p.client.FilterRows(context.Background(), req)
		if err != nil {
			ws.Close()
			fmt.Println("err", err)
			return
		}

		buf, err = proto.Marshal(res)
		if err != nil {
			ws.Close()
			fmt.Println("err", err)
			return
		}
		res2 := &saver.FilterRowsRes{}
		proto.Unmarshal(buf, res2)
		websocket.Message.Send(ws, buf)
	}
	ws.Close()
}
