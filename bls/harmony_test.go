package bls

import "testing"
import "fmt"
import "encoding/hex"

func TestHarmony(t *testing.T) {
	if Init(BLS12_381) != nil {
		t.Fatalf("Init")
	}
	if SetETHmode(0) != nil {
		t.Fatalf("SetMapToMode")
	}
	SetETHserialization(false)
	SetMapToMode(0)
	var gen PublicKey
	gen.SetHexString("1 4f58f3d9ee829f9a853f80b0e32c2981be883a537f0c21ad4af17be22e6e9959915ec21b7f9d8cc4c7315f31f3600e5 1212110eb10dbc575bccc44dcd77400f38282c4728b5efac69c0b4c9011bd27b8ed608acd81f027039216a291ac636a8")
	SetGeneratorOfPublicKey(&gen)

	var sec SecretKey
	bytes := []byte{185, 10, 0, 237, 196, 231, 235, 253, 95, 221, 36, 224, 109, 252, 10, 222, 14, 197, 82, 107, 220, 29, 208, 123, 169, 190, 98, 181, 240, 198, 35, 98}
	sec.Deserialize(bytes)
	sig := sec.SignHash(bytes)
	fmt.Printf("sig=%v\n", sig.SerializeToHexStr())
	//	assertBool("check sig", byteToHexStr(sig.serialize()).equals("2f4ff940216b2f13d75a231b988cd16ef22b45a4709df3461d9baeebfeaafeb54fad86ea7465212f35ceb0af6fe86b1828cf6de9099cefe233d97e0523ba6c0f5eecf4db71f7b1ae08cd098547946abbd0329fdac14d27102f2a1891e9188a19"));A

	sec.DeserializeHexStr("d243f3f029c188a5b1c4c098f5719cbc967184ef962b5c5d6c72693c92c1f725")
	pub := sec.GetPublicKey()
	fmt.Printf("pub=%v\n", pub.SerializeToHexStr())
	if pub.SerializeToHexStr() != "15ad529698be1f6164fd50416d1991a04b977d9014b43ae9d014ca50ae634829182632d54c7188b3a53e0b77ae4c9e87" {
		t.Errorf("bad pub")
	}
	msg, _ := hex.DecodeString("1100000000000000000000000000000000000000000000000000000000000031")
	sig = sec.SignByte(msg)
	if sig.SerializeToHexStr() != "2f2071370054ded61081a069dc657d7111ae909fc08a38ae1f3a967c4e4c1b97b7d2ee8adc512ded7213dc0ab2c9db158451eca54a4db26ec013e108d561d8e0eede883b73074d49143ea8e9291903aabce2a2b4f153f99822e4512dac86ae18" {
		t.Errorf("bad sig")
	}
	sig = sec.SignHash(msg)
	if sig.SerializeToHexStr() != "807f6fb074cfdd0318501ef894b127a3f71754745dcff2163dd128091efd065c4ac419e3d7d5428a0b94ada19de657176fb4ccfdb3105f869d6351e503ba20e9d6ef55f179510941db4d131c3678e9d9316090c99bd7d62107d83f6fe5f6f115" {
		t.Errorf("bad sig2")
	}
}
