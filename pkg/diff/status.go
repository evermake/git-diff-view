package diff

import (
	"strconv"
	"unicode"
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
	n, err := strconv.ParseInt(string(raw), 10, 64)
	if err != nil {
		return 0, err
	}

	return int(n), nil
}

func parseStatus(raw []byte) (status Status, err error) {
	var i int
	for i = 0; i < len(raw); i++ {
		if !unicode.IsLetter(rune(raw[i])) {
			break
		}
	}

	statusLetter := rune(raw[:i][0])

	var ok bool
	status.Type, ok = parseStatusType(statusLetter)
	if !ok {
		err = ErrUnknownStatusLetter
		return
	}

	switch status.Type {
	case StatusCopy, StatusRename:
		score, err := parseStatusScore(raw[i:])
		if err != nil {
			return Status{}, err
		}

		status.Score = &score
	case StatusModify:
		if i < len(raw) {
			score, err := parseStatusScore(raw[i:])
			if err != nil {
				return Status{}, err
			}

			status.Score = &score
		}
	}

	return
}
