package fltk_bridge

/*
#include "secret_input.h"
*/
import "C"
import "unsafe"

type SecretInput struct {
	Input
}

func NewSecretInput(x, y, w, h int, text ...string) *SecretInput {
	i := &SecretInput{}
	initWidget(i, unsafe.Pointer(C.go_fltk_new_Secret_Input(C.int(x), C.int(y), C.int(w), C.int(h), cStringOpt(text))))
	return i
}
