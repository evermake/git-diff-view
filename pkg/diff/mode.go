package diff

import (
	"bytes"
	"encoding/binary"
)

type Mode uint32

func parseMode(raw []byte) (Mode, error) {
	buff := bytes.NewReader(raw)

	var mode Mode
	if err := binary.Read(buff, binary.BigEndian, &mode); err != nil {
		return 0, err
	}

	return mode, nil
}
