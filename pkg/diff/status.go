package diff

import (
	"bytes"
	"encoding/binary"
	"unicode"
	"unicode/utf8"
)

type Status struct {
	Type  StatusType
	Score *int
}

type StatusType rune

const (
	StatusAdd        = StatusType('A')
	StatusCopy       = StatusType('C')
	StatusDelete     = StatusType('D')
	StatusModify     = StatusType('M')
	StatusRename     = StatusType('R')
	StatusChangeType = StatusType('T')
	StatusUnmerged   = StatusType('U')
	StatusUnknown    = StatusType('X')
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

func parseStatusScore(raw []byte) (int, error) {
	buff := bytes.NewReader(raw)

	var score int
	if err := binary.Read(buff, binary.BigEndian, &score); err != nil {
		return 0, err
	}

	return score, nil
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
		score, err := parseStatusScore(fields[1])
		if err != nil {
			return Status{}, err
		}

		status.Score = &score
	case StatusModify:
		if len(fields) == 2 {
			score, err := parseStatusScore(fields[1])
			if err != nil {
				return Status{}, err
			}

			status.Score = &score
		}
	}

	return
}
