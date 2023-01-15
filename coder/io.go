package coder

import (
	"io"
)

func runOnStream(r io.Reader, w io.Writer, base64HandlerFunc func([]byte, int, byte, int, int) ([]byte, byte, int, int)) error {
	b := byte(0x00)
	o := 0
	l := 0
	var processedBuff []byte
	buff := make([]byte, 1024)
	for {
		size, err := r.Read(buff)
		if err == io.EOF || size == 0 {
			break
		} else if err != nil {
			return err
		}

		processedBuff, b, o, l = base64HandlerFunc(buff, size, b, o, l)
		w.Write(processedBuff)
	}

	if o != 0 {
		buff, _, _, _ = base64HandlerFunc([]byte{}, 0, b, o, l)
		w.Write(buff)
	}

	w.Write([]byte{'\n'})

	return nil
}
