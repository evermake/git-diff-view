package diff

import (
	"bytes"
	"encoding/binary"
	"unicode"
	"unicode/utf8"
)

type Status struct {
	Type       StatusType
	Percentage *int
}

type StatusType uint8

const (
	StatusUnknown StatusType = iota
	StatusAdd
	StatusCopy
	StatusDelete
	StatusModify
	StatusRename
	StatusChangeType
	StatusUnmerged
)

func parseStatusType(letter rune) (StatusType, bool) {
	t, ok := map[rune]StatusType{
		'A': StatusAdd,
		'C': StatusCopy,
		'D': StatusDelete,
		'M': StatusModify,
		'R': StatusRename,
		'T': StatusChangeType,
		'U': StatusUnmerged,
		'X': StatusUnknown,
	}[letter]

	return t, ok
}

func parseStatusPercentage(raw []byte) (int, error) {
	buff := bytes.NewReader(raw)

	var percentage int
	if err := binary.Read(buff, binary.BigEndian, &percentage); err != nil {
		return 0, err
	}

	return percentage, nil
}

func parseStatus(raw []byte) (status Status, err error) {
	fields := bytes.FieldsFunc(raw, func(r rune) bool {
		return unicode.IsDigit(r)
	})

	statusLetter, _ := utf8.DecodeRune(fields[0])

	var ok bool
	status.Type, ok = parseStatusType(statusLetter)
	if !ok {
		err = ErrUnknownStatusLetter
		return
	}

	switch status.Type {
	case StatusCopy, StatusRename:
		percentage, err := parseStatusPercentage(fields[1])
		if err != nil {
			return
		}

		status.Percentage = &percentage
	case StatusModify:
		if len(fields) == 2 {
			percentage, err := parseStatusPercentage(fields[1])
			if err != nil {
				return
			}

			status.Percentage = &percentage
		}
	}

	return
}
