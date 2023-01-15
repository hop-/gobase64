package coder

import (
	"io"
	"math"
)

func DecodeStream(r io.Reader, w io.Writer) error {
	return runOnStream(r, w, decodeChunk)
}

func decodeChunk(buff []byte, buffLen int, b byte, o int, lineLen int) ([]byte, byte, int, int) {
	if buffLen == 0 {
		return []byte{}, 0x00, 0, lineLen
	}

	decodedBuffLen := int(math.Ceil(float64(buffLen-(lineLen+buffLen)/76) * 3 / 4))
	if o != 0 {
		decodedBuffLen++
	}

	decodedBuff := make([]byte, decodedBuffLen)
	decodedBuff[0] = b

	j := 0
	offset := o
	for i := 0; i < buffLen; i++ {
		if lineLen > 0 && lineLen%76 == 0 {
			lineLen = 0
			continue
		}
		lineLen++

		cb := charToByte[buff[i]]
		switch offset {
		case 0:
			offset = 6
			decodedBuff[j] = cb << 2
		case 6:
			offset = 4
			decodedBuff[j] |= cb >> 4
			j++
			if j >= decodedBuffLen {
				break
			}
			decodedBuff[j] = cb << 4
		case 4:
			offset = 2
			decodedBuff[j] |= cb >> 2
			j++
			if j >= decodedBuffLen {
				break
			}
			decodedBuff[j] = (cb & 0x03) << 6
		case 2:
			offset = 0
			decodedBuff[j] |= cb
			j++
		}
	}
	if offset != 0 {
		return decodedBuff[:decodedBuffLen-1], decodedBuff[len(decodedBuff)-1], offset, lineLen
	}
	return decodedBuff, 0x00, offset, lineLen
}
