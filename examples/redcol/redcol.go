package redcol

import (
	"bufio"
	"fmt"
	"github.com/weizetao/gotcp"
	"io"
)

type RedPacket struct {
	resp *Resp
	cmd  string
	keys [][]byte
}

func (this *RedPacket) Cmd() string {
	return this.cmd
}

func (this *RedPacket) Serialize() []byte {
	return this.resp.Bytes()
}
func (this *RedPacket) SetCmd(cmd string, keys ...interface{}) error {
	if cmd == "" {
		return fmt.Errorf("ERR: cmd is not allowed empty!")
	}
	this.cmd = cmd
	for _, key := range keys {
		this.keys = append(this.keys, formatCommandArg(key))
	}
	return nil
}

func (this *RedPacket) AppendKeys(keys ...interface{}) error {
	if this.cmd == "" {
		return fmt.Errorf("ERR: cmd is not allowed empty!")
	}
	for _, key := range keys {
		this.keys = append(this.keys, formatCommandArg(key))
	}
	return nil
}

func SyncReadPacket(r *bufio.Reader) (*RedPacket, error) {
	p := &RedPacket{}

	resp, err := Parse(r)
	if err != nil {
		return nil, err
	}

	cmd, keys, err := resp.GetOpKeys()
	if err != nil {
		return nil, err
	}

	p.resp = resp
	p.cmd = string(cmd)
	p.keys = keys

	return p, nil
}
func SyncWritePacket(w io.Writer, p *RedPacket) error {
	return WriteCmdKeys(w, p.cmd, p.keys)
}

type RedProtocol struct {
}

func (this *RedProtocol) ReadPacket(r *bufio.Reader) (gotcp.Packet, error) {
	return SyncReadPacket(r)
}

func (this *RedProtocol) WritePacket(w *bufio.Writer, p gotcp.Packet) error {
	redsyncPk := p.(*RedPacket)
	return SyncWritePacket(w, redsyncPk)
}
