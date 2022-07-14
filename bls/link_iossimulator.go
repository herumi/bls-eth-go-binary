//go:build iossimulator
// +build iossimulator

package bls

/*
#cgo LDFLAGS:-lbls384_256 -lstdc++ -lm
#cgo ios LDFLAGS:-L${SRCDIR}/lib/iossimulator
*/
import "C"
