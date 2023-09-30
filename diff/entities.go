package diff

type Status struct {
	Type       StatusType
	Percentage *int
}

type Mode uint32

type State struct {
	Mode Mode
	SHA1 []byte
	Path string
}

type Diff struct {
	Status Status
	Src    State
	Dst    State
}
