package gotcp

import (
	"bufio"
)

type Packet interface {
	Serialize() []byte
}

type Protocol interface {
	ReadPacket(r *bufio.Reader) (Packet, error)
	WritePacket(w *bufio.Writer, p Packet) error
}
