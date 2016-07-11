package redsync

import (
	"bufio"
	"fmt"
	"github.com/weizetao/gotcp"
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

type RedProtocol struct {
}

func (this *RedProtocol) ReadPacket(r *bufio.Reader) (gotcp.Packet, error) {
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

func (this *RedProtocol) WritePacket(w *bufio.Writer, p gotcp.Packet) error {
	redsyncPk := p.(*RedPacket)
	// return redsyncPk.resp.WriteTo(w)
	return WriteCmdKeys(w, redsyncPk.cmd, redsyncPk.keys)
}
