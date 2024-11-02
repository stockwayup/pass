package event

//go:generate msgp

type Generate struct {
	Password []byte `msg:"password"`
}
