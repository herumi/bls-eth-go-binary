package bls

import (
	"unsafe"
)

func CastFromG1(in *G1) *PublicKey {
	return (*PublicKey)(unsafe.Pointer(in))
}

func CastToG1(in *PublicKey) *G1 {
	return (*G1)(unsafe.Pointer(in))
}

func CopyFromG1(in *G1) PublicKey {
	return *CastFromG1(in)
}

func CopyToG1(in *PublicKey) G1 {
	return *CastToG1(in)
}
