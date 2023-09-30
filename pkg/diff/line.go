package diff

import "github.com/bluekeyes/go-gitdiff/gitdiff"

type Line struct {
	Op                       gitdiff.LineOp
	Content                  string
	NumberInSrc, NumberInDst int64
}
