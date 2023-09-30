package diff

type State struct {
	Mode Mode
	SHA1 []byte
	Path string
}
