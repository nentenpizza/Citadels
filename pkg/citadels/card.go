package citadels

type Card interface {
	Type() string
	MakeAction() *Action
}
