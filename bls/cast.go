package bls

import (
	"unsafe"
)

// SecretKey

func CastFromSecretKey(in *SecretKey) *Fr {
	return (*Fr)(unsafe.Pointer(in))
}

func CastToSecretKey(in *Fr) *SecretKey {
	return (*SecretKey)(unsafe.Pointer(in))
}

func CopyFromSecretKey(in *SecretKey) Fr {
	return *CastFromSecretKey(in)
}

func CopyToSecretKey(in *Fr) SecretKey {
	return *CastToSecretKey(in)
}

// PublicKey

func CastFromPublicKey(in *PublicKey) *G1 {
	return (*G1)(unsafe.Pointer(in))
}

func CastToPublicKey(in *G1) *PublicKey {
	return (*PublicKey)(unsafe.Pointer(in))
}

func CopyFromPublicKey(in *PublicKey) G1 {
	return *CastFromPublicKey(in)
}

func CopyToPublicKey(in *G1) PublicKey {
	return *CastToPublicKey(in)
}

// Sign

func CastFromSign(in *Sign) *G2 {
	return (*G2)(unsafe.Pointer(in))
}

func CastToSign(in *G2) *Sign {
	return (*Sign)(unsafe.Pointer(in))
}

func CopyFromSign(in *Sign) G2 {
	return *CastFromSign(in)
}

func CopyToSign(in *G2) Sign {
	return *CastToSign(in)
}
