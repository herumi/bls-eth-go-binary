//go:build ios && !iossimulator
// +build ios,!iossimulator

package bls

/*
#cgo LDFLAGS:-lbls384_256 -lstdc++ -lm
#cgo ios LDFLAGS:-L${SRCDIR}/lib/ios
*/
import "C"
