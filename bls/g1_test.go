package bls

import (
	"io/ioutil"
	"testing"
)

func testUncompressed(t *testing.T, gen *G1) {
	buf, err := ioutil.ReadFile("tests/g1_uncompressed_valid_test_vectors.dat")
	if err != nil {
		t.Fatalf("ReadFile")
	}
	var p1, p2 G1
	p1.Z.SetString("1", 10) // affine
	for i := 0; i < 1000; i++ {
		if p1.DeserializeUncompressed(buf[i*96 : (i+1)*96]) != nil {
			t.Fatalf("i=%d X.Deserialize", i)
		}
		if !p1.IsEqual(&p2) {
			t.Fatalf("p1=%x\np2=%x\n", p1.Serialize(), p2.Serialize())
		}
		G1Add(&p2, &p2, gen)
	}
}

func testCompressed(t *testing.T, gen *G1) {
	buf, err := ioutil.ReadFile("tests/g1_compressed_valid_test_vectors.dat")
	if err != nil {
		t.Fatalf("ReadFile")
	}
	var p1, p2 G1
	for i := 0; i < 1000; i++ {
		if p1.Deserialize(buf[i*48 : (i+1)*48]) != nil {
			t.Fatalf("err i=%d\n", i)
		}
		if !p1.IsEqual(&p2) {
			t.Fatalf("p1=%x\np2=%x\n", p1.Serialize(), p2.Serialize())
		}
		G1Add(&p2, &p2, gen)
	}
}

func Test(t *testing.T) {
	if Init(BLS12_381) != nil {
		t.Fatalf("Init")
	}
	var gen G1
	if gen.SetString("1 0x17f1d3a73197d7942695638c4fa9ac0fc3688c4f9774b905a14e3a3f171bac586c55e83ff97a1aeffb3af00adb22c6bb 0x08b3f481e3aaa0f1a09e30ed741d8ae4fcf5e095d5d00af600db18cb2c04b3edd03cc744a2888ae40caa232946c5e7e1", 16) != nil {
		t.Fatalf("SetString")
	}
	testCompressed(t, &gen)
}
