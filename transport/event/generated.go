package event

//go:generate msgp

type Generated struct {
	Hash []byte `msg:"hash"`
	Salt []byte `msg:"salt"`
}
