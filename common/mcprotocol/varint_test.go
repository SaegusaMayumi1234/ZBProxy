package mcprotocol

import (
	"bytes"
	"io"
	"testing"

	"github.com/Tnze/go-mc/net/packet"
	"github.com/layou233/ZBProxy/common/buf"
)

// samples from https://wiki.vg/Protocol#VarInt_and_VarLong

func checkWrite(t *testing.T, n int32, result [MaxVarIntLen]byte) {
	buffer := buf.NewSize(MaxVarIntLen + 1)
	buffer.WriteZeroN(MaxVarIntLen)
	buffer.FullReset()
	defer buffer.Release()
	_, err := VarInt(n).WriteTo(buffer)
	if err != nil {
		return
	}
	t.Log("VarInt", n, "WriteTo", buffer.Bytes())
	buffer.Truncate(MaxVarIntLen)
	if !bytes.Equal(buffer.Bytes(), result[:]) {
		t.Fatalf("VarInt WriteTo error: got %v, expect %v", buffer.Bytes(), result)
	}
	return
}

func checkRead(t *testing.T, n int32, result [MaxVarIntLen]byte) {
	buffer := buf.As(result[:])
	defer buffer.Release()
	vi, _, err := ReadVarIntFrom(buffer)
	if err != nil {
		return
	}
	t.Log("VarInt", vi, "ReadFrom", result)
	if n != vi {
		t.Fatalf("VarInt ReadFrom error: got %v, expect %v", vi, n)
	}
	return
}

func TestVarInt_WriteTo(t *testing.T) {
	packet.VarInt(-1).WriteTo(io.Discard)
	checkWrite(t, 0, [MaxVarIntLen]byte{0})
	checkWrite(t, 1, [MaxVarIntLen]byte{1})
	checkWrite(t, 2, [MaxVarIntLen]byte{2})
	checkWrite(t, 127, [MaxVarIntLen]byte{127})
	checkWrite(t, 128, [MaxVarIntLen]byte{128, 1})
	checkWrite(t, 255, [MaxVarIntLen]byte{255, 1})
	checkWrite(t, 25565, [MaxVarIntLen]byte{221, 199, 1})
	checkWrite(t, 2097151, [MaxVarIntLen]byte{255, 255, 127})
	checkWrite(t, 2147483647, [MaxVarIntLen]byte{255, 255, 255, 255, 7})
	checkWrite(t, -1, [MaxVarIntLen]byte{255, 255, 255, 255, 15})
	checkWrite(t, -2147483648, [MaxVarIntLen]byte{128, 128, 128, 128, 8})
}

func TestReadFrom(t *testing.T) {
	checkRead(t, 0, [MaxVarIntLen]byte{0})
	checkRead(t, 1, [MaxVarIntLen]byte{1})
	checkRead(t, 2, [MaxVarIntLen]byte{2})
	checkRead(t, 127, [MaxVarIntLen]byte{127})
	checkRead(t, 128, [MaxVarIntLen]byte{128, 1})
	checkRead(t, 255, [MaxVarIntLen]byte{255, 1})
	checkRead(t, 25565, [MaxVarIntLen]byte{221, 199, 1})
	checkRead(t, 2097151, [MaxVarIntLen]byte{255, 255, 127})
	checkRead(t, 2147483647, [MaxVarIntLen]byte{255, 255, 255, 255, 7})
	checkRead(t, -1, [MaxVarIntLen]byte{255, 255, 255, 255, 15})
	checkRead(t, -2147483648, [MaxVarIntLen]byte{128, 128, 128, 128, 8})
}
