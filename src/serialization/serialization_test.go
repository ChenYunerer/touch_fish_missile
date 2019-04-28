package serialization

import (
	"bytes"
	"testing"
)

func TestSerialization(t *testing.T) {
	b := make([]byte, 0)
	b = append(b, []byte("123456789098756123456789098756123456789098756123456789098756123456789098756")...)
	buf := &bytes.Buffer{}
	buf.Write(b)

	c := make([]byte, 2)
	buf.Read(c)
	println(c)
	buf.Next(2)
	buf.Read(c)
	println(c)
}
