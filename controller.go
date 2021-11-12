package socket

type Controller interface {
	// TODO 这里后续可以做成一个interface。
	Handle(req []byte) interface{}
}
