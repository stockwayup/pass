package event

//go:generate msgp

type Validate struct {
	Input    []byte `msg:"input"`
	Password []byte `msg:"password"`
	Salt     []byte `msg:"salt"`
}
