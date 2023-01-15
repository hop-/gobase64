package coder

import (
	"io"
)

func EncodeStream(r io.Reader, w io.Writer) error {
	return runOnStream(r, w, encodeChunk)
}

func encodeChunk(baseBuff []byte, baseBuffLen int, b byte, o int, lineLen int) ([]byte, byte, int, int) {
	var buff []byte
	if o != 0 {
		if baseBuffLen == 0 {
			buff = []byte{b, 0x00}
		} else {
			buff = append([]byte{b}, baseBuff[:baseBuffLen]...)
		}
	} else {
		if baseBuffLen == 0 {
			return []byte{}, 0x00, 0, lineLen
		}

		buff = baseBuff[:baseBuffLen]
	}

	encodedBuffLen := baseBuffLen * 4 / 3
	encodedBuffLen += (lineLen + encodedBuffLen) / 76

	encodedBuff := make([]byte, encodedBuffLen)

	j := 0
	offset := o
	for i := 0; i < encodedBuffLen; i++ {
		if lineLen > 0 && lineLen%76 == 0 {
			encodedBuff[i] = '\n'
			lineLen = 0
			i++
			if i >= encodedBuffLen {
				break
			}
		}
		lineLen++

		var c byte
		switch offset {
		case 0:
			c = buff[j] >> 2
			offset = 6
		case 6:
			c = (buff[j]&0x03)<<4 | buff[j+1]>>4
			offset = 4
			j++
		case 4:
			c = buff[j]&0x0F<<2 | buff[j+1]>>6
			offset = 2
			j++
		case 2:
			c = buff[j] & 0x3F
			offset = 0
			j++
		}
		encodedBuff[i] = byteToChar[c]
	}

	return encodedBuff, buff[len(buff)-1], offset, lineLen
}
