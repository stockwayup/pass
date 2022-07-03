package event

//go:generate msgp

type Validate struct {
	Input    []byte `msgp:"input"`
	Password []byte `msgp:"password"`
	Salt     []byte `msgp:"salt"`
}
