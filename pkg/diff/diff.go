package diff

import (
	"bytes"
)

type Diff struct {
	Status   Status
	Lines    []*Line
	Src      State
	Dst      State
	IsBinary bool
}

// parseDiff parses raw diff line (:100644 100644 bcd1234 0123456 M file0) according to
// https://git-scm.com/docs/git-diff
func parseDiff(raw []byte) (diff Diff, err error) {
	// first trailing colon is removed
	raw = bytes.TrimPrefix(raw, []byte(":"))

	fields := bytes.Fields(raw)

	if len(fields) < 5 || len(fields) > 6 {
		err = ErrMalformedDiff
		return
	}

	nextField := func() func() []byte {
		n := -1

		return func() []byte {
			n++
			return fields[n]
		}
	}()

	diff.Src.Mode, err = parseMode(nextField())
	if err != nil {
		return
	}

	diff.Dst.Mode, err = parseMode(nextField())
	if err != nil {
		return
	}

	diff.Src.SHA1 = nextField()
	diff.Dst.SHA1 = nextField()

	diff.Status, err = parseStatus(nextField())
	if err != nil {
		return
	}

	diff.Src.Path = string(nextField())

	if diff.Status.Type == StatusCopy || diff.Status.Type == StatusRename {
		if len(fields) != 6 {
			err = ErrMalformedDiff
			return
		}

		diff.Dst.Path = string(nextField())
	} else {
		diff.Dst.Path = diff.Src.Path
	}

	return
}
