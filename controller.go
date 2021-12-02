package socket

import "reflect"

type Controller interface {
	Handle(req []byte) interface{}
}

type Filter interface {
	MyselfRequest(requestData []byte) bool
}

// FilterFunc of custom defined function type.
type FilterFunc func(requestData []byte) bool

func (ff FilterFunc) MyselfRequest(requestData []byte) bool {
	return ff(requestData)
}

// https://stackoverflow.com/questions/27803654/explanation-of-checking-if-value-implements-interface

var FILTERITERN = reflect.TypeOf((*Filter)(nil)).Elem()
