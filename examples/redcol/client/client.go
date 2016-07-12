package main

import (
	"bufio"
	"fmt"
	"github.com/weizetao/gotcp/examples/redcol"
	"log"
	"net"
	"time"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8989")
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	r := bufio.NewReaderSize(conn, 0)

	// ping <--> pong
	for i := 0; i < 3; i++ {
		// write
		redPkt := &redcol.RedPacket{}
		redPkt.SetCmd("hello", "wiky")
		redcol.SyncWritePacket(conn, redPkt)

		// read
		p, err := redcol.SyncReadPacket(r)
		if err == nil {
			fmt.Printf("Server reply:[%v]\n", p.Cmd())
		}

		time.Sleep(2 * time.Second)
	}

	conn.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
