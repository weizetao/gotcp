package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/weizetao/gotcp"
	"github.com/weizetao/gotcp/examples/redcol"
)

type Callback struct{}

func (this *Callback) OnConnect(c *gotcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	fmt.Println("OnConnect:", addr)
	return true
}
func (this *Callback) OnDial() (*net.TCPConn, error) {
	return nil, nil
}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	redcolPacket := p.(*redcol.RedPacket)
	fmt.Printf("OnMessage:[%s]\n", redcolPacket.Cmd())

	redPk := &redcol.RedPacket{}
	redPk.SetCmd("hello", "redcol", "wrold")

	c.AsyncWritePacket(redPk, time.Second)
	return true
}

func (this *Callback) OnClose(c *gotcp.Conn) {
	fmt.Println("OnClose:", c.GetExtraData())
}

func main() {
	runtime.GOMAXPROCS(1)

	// creates a tcp listener
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":8989")
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// creates a server
	config := &gotcp.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}
	srv := gotcp.NewServer(config, &Callback{}, &redcol.RedProtocol{})

	// starts service
	go srv.Start(listener, time.Second)
	fmt.Println("listening:", listener.Addr())

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	srv.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
