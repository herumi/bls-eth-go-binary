package main

import (
	"fmt"
	"github.com/herumi/bls-eth-go-binary/bls"
)

func sample1() {
	var sec bls.SecretKey
	sec.SetByCSPRNG()
	fmt.Printf("sec:%s\n", sec.SerializeToHexStr())
	pub := sec.GetPublicKey()
	fmt.Printf("1.pub:%s\n", pub.SerializeToHexStr())
	fmt.Printf("1.pub x=%x\n", pub)
	var P *bls.G1 = bls.CastFromPublicKey(pub)
	bls.G1Normalize(P, P)
	fmt.Printf("2.pub:%s\n", pub.SerializeToHexStr())
	fmt.Printf("2.pub x=%x\n", pub)
	fmt.Printf("P.X=%x\n", P.X.Serialize())
	fmt.Printf("P.Y=%x\n", P.Y.Serialize())
	fmt.Printf("P.Z=%x\n", P.Z.Serialize())
}

func sample2() {
	var sec bls.SecretKey
	sec.DeserializeHexStr("47b8192d77bf871b62e87859d653922725724a5c031afeabc60bcef5ff665138")

	msg := [40]byte{1, 2, 3}
	fmt.Printf("sec=%v\n", sec.Serialize())
	pub := sec.GetPublicKey()
	fmt.Printf("pub=%s\n", pub.SerializeToHexStr())
	sig := sec.SignHashWithDomain(msg[:])
	fmt.Printf("sig=%s\n", sig.SerializeToHexStr())
	fmt.Printf("verify=%v\n", sig.VerifyHashWithDomain(pub, msg[:]))
}

func sample3() {
	var sec bls.SecretKey
	b := make([]byte, 64)
	for i := 0; i < len(b); i++ {
		b[i] = 0xff
	}
	err := sec.SetLittleEndianMod(b)
	if err != nil {
		fmt.Printf("err")
		return
	}
	fmt.Printf("sec=%x\n", sec.Serialize())
}

func main() {
	bls.Init(bls.BLS12_381)
	sample1()
	sample2()
	sample3()
}
