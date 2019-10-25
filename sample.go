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

	msg := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 87, 33, 13, 72, 155, 73, 4, 185, 87, 46, 230, 247, 159, 191, 7, 148, 85, 120, 129, 175, 102, 169, 241, 139, 189, 44, 244, 68, 119, 60, 28, 101, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 225, 95, 237, 38, 188, 142, 181, 147, 233, 183, 232, 13, 219, 92, 94, 79, 19, 174, 172, 105, 133, 207, 4, 113, 115, 242, 140, 138, 44, 215, 244, 77}
	fmt.Printf("sec=%v\n", sec.Serialize())
	pub := sec.GetPublicKey()
	fmt.Printf("pub=%s\n", pub.SerializeToHexStr())
	sig := sec.SignHash(msg)
	fmt.Printf("sig=%s\n", sig.SerializeToHexStr())
	var x bls.Fp2
	if x.Deserialize(msg) != nil {
		fmt.Printf("ERR D")
		return
	}
	fmt.Printf("D[0]=%s\n", x.D[0].GetString(16))
	fmt.Printf("D[1]=%s\n", x.D[1].GetString(16))
	var Q bls.G2
	bls.MapToG2(&Q, &x)
	bls.G2Mul(&Q, &Q, bls.CastFromSecretKey(&sec))
	fmt.Printf("Q =%x\n", Q.Serialize())
	fmt.Printf("ok=%s\n", "b9d1bf921b3dd048bdce38c2ceac2a2a8093c864881f2415f22b198de935ffa791707855c1656dc21a7af2d502bb46590151d645f062634c3b2cb79c4ed1c4a4b8b3f19f0f5c76965c651553e83d153ff95353735156eff77692f7a62ae653fb")
}

func main() {
	bls.Init(bls.BLS12_381)
	sample1()
	sample2()
}
