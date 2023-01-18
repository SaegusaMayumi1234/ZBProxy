package mcprotocol

import (
	"fmt"
	"io"
	"math"
	"net"

	"github.com/layou233/ZBProxy/common/buf"
)

type Conn struct {
	io.Reader
	io.Writer
	net.Conn
}

func StreamConn(conn net.Conn) Conn {
	return Conn{
		Reader: conn,
		Writer: conn,
		Conn:   conn,
	}
}

// ReadLimitedPacket likes ReadPacket, but limits the maximum number of packet content bytes to read to maxLen.
func (c Conn) ReadLimitedPacket(buffer *buf.Buffer, maxLen int) (err error) {
	length, _, err := ReadVarIntFrom(c.Reader)
	if err != nil {
		return
	}
	lengthInt := int(length)

	if lengthInt < 0 {
		return fmt.Errorf("incorrect packet length: %v", lengthInt)
	}
	if lengthInt > maxLen {
		return fmt.Errorf("packet max length exceeded: length=%v, max=%v", lengthInt, maxLen)
	}
	if buffer.FreeLen() < lengthInt {
		return fmt.Errorf("short buffer: free size=%v, need=%v", buffer.FreeLen(), lengthInt)
	}

	_, err = buffer.ReadFullFrom(c.Reader, lengthInt)
	return
}

// ReadPacket reads a full packet to buffer.
func (c Conn) ReadPacket(buffer *buf.Buffer) error {
	return c.ReadLimitedPacket(buffer, math.MaxInt)
}

// WritePacket appends packet length to packet head, and writes to Conn.
// Then reset the buffer to MaxVarIntLen.
// Note that the given buffer should have at least 5 bytes front headroom space.
func (c Conn) WritePacket(buffer *buf.Buffer) (err error) {
	l := buffer.Len()
	packetLengthBytes, packetLengthLen := EncodeVarInt(int32(l))
	header := buffer.ExtendHeader(MaxVarIntLen - packetLengthLen)
	copy(header, packetLengthBytes[:])

	_, err = c.Writer.Write(buffer.Bytes())
	buffer.Reset(MaxVarIntLen)
	return
}
