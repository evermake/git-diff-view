package diff

type LineOperation rune

const (
	LineOperationModify = LineOperation('M')
	LineOperationAdd    = LineOperation('A')
	LineOperationDelete = LineOperation('D')
)

type Line struct {
	Operation LineOperation
	Src       LineState
	Dst       LineState
}

type LineState struct {
	Number  int64
	Content string
}
