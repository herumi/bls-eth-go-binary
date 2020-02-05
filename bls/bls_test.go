package bls

import (
	"encoding/csv"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func testUncompressedG1(t *testing.T, gen1 *G1) {
	buf, err := ioutil.ReadFile("tests/g1_uncompressed_valid_test_vectors.dat")
	if err != nil {
		t.Fatalf("ReadFile")
	}
	one := CastToPublicKey(gen1)
	var p1, p2 PublicKey
	for i := 0; i < 1000; i++ {
		if p1.DeserializeUncompressed(buf[i*96:(i+1)*96]) != nil {
			t.Fatalf("i=%d X.Deserialize", i)
		}
		if !p1.IsEqual(&p2) {
			t.Fatalf("i=%d p1=%x\np2=%x\n", i, p1.Serialize(), p2.Serialize())
		}
		p2.Add(one)
	}
}

func testCompressedG1(t *testing.T, gen1 *G1) {
	buf, err := ioutil.ReadFile("tests/g1_compressed_valid_test_vectors.dat")
	if err != nil {
		t.Fatalf("ReadFile")
	}
	one := CastToPublicKey(gen1)
	var p1, p2 PublicKey
	for i := 0; i < 1000; i++ {
		if p1.Deserialize(buf[i*48:(i+1)*48]) != nil {
			t.Fatalf("err i=%d\n", i)
		}
		if !p1.IsEqual(&p2) {
			t.Fatalf("p1=%x\np2=%x\n", p1.Serialize(), p2.Serialize())
		}
		p2.Add(one)
	}
}

func testUncompressedG2(t *testing.T, gen2 *G2) {
	buf, err := ioutil.ReadFile("tests/g2_uncompressed_valid_test_vectors.dat")
	if err != nil {
		t.Fatalf("ReadFile")
	}
	one := CastToSign(gen2)
	var p1, p2 Sign
	for i := 0; i < 1000; i++ {
		if p1.DeserializeUncompressed(buf[i*192:(i+1)*192]) != nil {
			t.Fatalf("i=%d X.Deserialize", i)
		}
		if !p1.IsEqual(&p2) {
			t.Fatalf("i=%d p1=%x\np2=%x\n", i, p1.Serialize(), p2.Serialize())
		}
		p2.Add(one)
	}
}

func testCompressedG2(t *testing.T, gen2 *G2) {
	buf, err := ioutil.ReadFile("tests/g2_compressed_valid_test_vectors.dat")
	if err != nil {
		t.Fatalf("ReadFile")
	}
	one := CastToSign(gen2)
	var p1, p2 Sign
	for i := 0; i < 1000; i++ {
		if p1.Deserialize(buf[i*96:(i+1)*96]) != nil {
			t.Fatalf("err i=%d\n", i)
		}
		if !p1.IsEqual(&p2) {
			t.Fatalf("p1=%x\np2=%x\n", p1.Serialize(), p2.Serialize())
		}
		p2.Add(one)
	}
}

func getSecPubHash() (*SecretKey, *PublicKey, []byte) {
	var sec SecretKey
	sec.SetByCSPRNG()
	pub := sec.GetPublicKey()
	var x Fp2
	x.D[0].SetByCSPRNG()
	x.D[1].SetByCSPRNG()
	hash := x.Serialize()
	return &sec, pub, hash
}

func testSignAndVerifyHash(t *testing.T) {
	sec, pub, hash := getSecPubHash()
	sig := sec.SignHash(hash)
	if sig == nil {
		t.Fatal("SignHash")
	}
	if !sig.VerifyHash(pub, hash) {
		t.Fatal("VerifyHash 1")
	}
	hash[0] = hash[0] + 1
	if sig.VerifyHash(pub, hash) {
		t.Fatal("VerifyHash 2")
	}
}

func getSecPubHashVec(n int) ([]PublicKey, [][]byte, []Sign) {
	pubVec := make([]PublicKey, n)
	hashVec := make([][]byte, n)
	sigVec := make([]Sign, n)
	var x Fp2
	var sec SecretKey
	for i := 0; i < n; i++ {
		sec.SetByCSPRNG()
		pubVec[i] = *sec.GetPublicKey()
		x.D[0].SetByCSPRNG()
		x.D[1].SetByCSPRNG()
		hashVec[i] = x.Serialize()
		sigVec[i] = *sec.SignHash(hashVec[i])
	}
	return pubVec, hashVec, sigVec
}

func testVerifyAggreageteHash(t *testing.T) {
	const N = 100
	pubVec, hashVec, sigVec := getSecPubHashVec(N)
	agg := sigVec[0]
	for i := 1; i < N; i++ {
		agg.Add(&sigVec[i])
	}
	if !agg.VerifyAggregateHashes(pubVec, hashVec) {
		t.Fatal("VerifyAggregateHashes 1")
	}
	hashVec[0][0] = hashVec[0][0] + 1
	if agg.VerifyAggregateHashes(pubVec, hashVec) {
		t.Fatal("VerifyAggregateHashes 2")
	}
}

func TestAreAllMsgDifferent(t *testing.T) {
	type V struct {
		s       string
		msgSize int
		result  bool
	}
	m := []V{V{"abcdabce", 4, true},
		V{"abcdabce", 2, false}, V{"abcdefgh", 2, true}, V{"xyzxyz", 2, true}, V{"xyzxyz", 3, false}}
	for _, v := range m {
		if AreAllMsgDifferent([]byte(v.s), v.msgSize) != v.result {
			t.Fatalf("err %v %v\n", v.s, v.msgSize)
		}
	}
}

func ethSignOneTest(t *testing.T, secHex string, msgHex string, sigHex string) {
	var sec SecretKey
	sec.DeserializeHexStr(secHex)
	pub := sec.GetPublicKey()
	msg, _ := hex.DecodeString(msgHex)
	sig := sec.SignByte(msg)
	if !sig.VerifyByte(pub, msg) {
		t.Fatalf("bad verify %v %v", secHex, msgHex)
	}
	if sig.SerializeToHexStr() != sigHex {
		t.Fatalf("bad sign %v %v", secHex, msgHex)
	}
}

func ethSignTest(t *testing.T) {
	secHex := "47b8192d77bf871b62e87859d653922725724a5c031afeabc60bcef5ff665138"
	msgHex := "0000000000000000000000000000000000000000000000000000000000000000"
	sigHex := "b2deb7c656c86cb18c43dae94b21b107595486438e0b906f3bdb29fa316d0fc3cab1fc04c6ec9879c773849f2564d39317bfa948b4a35fc8509beafd3a2575c25c077ba8bca4df06cb547fe7ca3b107d49794b7132ef3b5493a6ffb2aad2a441"

	ethSignOneTest(t, secHex, msgHex, sigHex)
	fileName := "tests/sign.txt"
	fp, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("can't open %v %v", fileName, err)
	}
	defer fp.Close()
	reader := csv.NewReader(fp)
	reader.Comma = ' '
	for {
		secHex, err := reader.Read()
		if err == io.EOF {
			break
		}
		msgHex, _ := reader.Read()
		sigHex, _ := reader.Read()
		ethSignOneTest(t, secHex[1], msgHex[1], sigHex[1])
	}
}

func ethAggregateTest(t *testing.T) {
	msgHexTbl := []string{ "b2a0bd8e837fc2a1b28ee5bcf2cddea05f0f341b375e51de9d9ee6d977c2813a5c5583c19d4e7db8d245eebd4e502163076330c988c91493a61b97504d1af85fdc167277a1664d2a43af239f76f176b215e0ee81dc42f1c011dc02d8b0a31e32", "b2deb7c656c86cb18c43dae94b21b107595486438e0b906f3bdb29fa316d0fc3cab1fc04c6ec9879c773849f2564d39317bfa948b4a35fc8509beafd3a2575c25c077ba8bca4df06cb547fe7ca3b107d49794b7132ef3b5493a6ffb2aad2a441", "a1db7274d8981999fee975159998ad1cc6d92cd8f4b559a8d29190dad41dc6c7d17f3be2056046a8bcbf4ff6f66f2a360860fdfaefa91b8eca875d54aca2b74ed7148f9e89e2913210a0d4107f68dbc9e034acfc386039ff99524faf2782de0e", }
	sigHex := "973ab0d765b734b1cbb2557bcf52392c9c7be3cd21d5bd28572d99f618c65e921f0dd82560cc103feb9f000c23c00e660e1364ed094f137e1045e73116cd75903af446df3c357540a4970ec367a7f7fa7493a5db27ca322c48d57740908585e8"
	n := len(msgHexTbl)
	sigVec := make([]Sign, n)
	for i, sigHex := range msgHexTbl {
		var t Sign
		t.DeserializeHexStr(sigHex)
		sigVec[i] = t
	}
	var aggSig Sign
	aggSig.Aggregate(sigVec)
	s := aggSig.SerializeToHexStr()
	if s != sigHex {
		t.Fatalf("bad aggregate %v %v\n", s, sigHex)
	}
}

func testEth(t *testing.T) {
	SetETHmode(1)
	ethAggregateTest(t)
	ethSignTest(t)
}

func Test(t *testing.T) {
	if Init(BLS12_381) != nil {
		t.Fatalf("Init")
	}
	var gen1 G1
	if gen1.SetString("1 0x17f1d3a73197d7942695638c4fa9ac0fc3688c4f9774b905a14e3a3f171bac586c55e83ff97a1aeffb3af00adb22c6bb 0x08b3f481e3aaa0f1a09e30ed741d8ae4fcf5e095d5d00af600db18cb2c04b3edd03cc744a2888ae40caa232946c5e7e1", 16) != nil {
		t.Fatalf("gen1.SetString")
	}
	var gen2 G2
	if gen2.SetString("1 0x024aa2b2f08f0a91260805272dc51051c6e47ad4fa403b02b4510b647ae3d1770bac0326a805bbefd48056c8c121bdb8 0x13e02b6052719f607dacd3a088274f65596bd0d09920b61ab5da61bbdc7f5049334cf11213945d57e5ac7d055d042b7e 0x0ce5d527727d6e118cc9cdc6da2e351aadfd9baa8cbdd3a76d429a695160d12c923ac9cc3baca289e193548608b82801 0x0606c4a02ea734cc32acd2b02bc28b99cb3e287e85a763af267492ab572e99ab3f370d275cec1da1aaa9075ff05f79be", 16) != nil {
		t.Fatalf("gen2.SetString")
	}

	testUncompressedG1(t, &gen1)
	testCompressedG1(t, &gen1)
	testUncompressedG2(t, &gen2)
	testCompressedG2(t, &gen2)
	testSignAndVerifyHash(t)
	testVerifyAggreageteHash(t)
	testEth(t)
}

func BenchmarkPairing(b *testing.B) {
	b.StopTimer()
	err := Init(BLS12_381)
	if err != nil {
		b.Fatal(err)
	}
	var P G1
	var Q G2
	var e GT
	P.HashAndMapTo([]byte("abc"))
	Q.HashAndMapTo([]byte("abc"))
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		Pairing(&e, &P, &Q)
	}
	b.StopTimer()
}

func BenchmarkSignHash(b *testing.B) {
	b.StopTimer()
	err := Init(BLS12_381)
	if err != nil {
		b.Fatal(err)
	}
	sec, _, hash := getSecPubHash()
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		sec.SignHash(hash)
	}
	b.StopTimer()
}

func BenchmarkVerifyHash(b *testing.B) {
	b.StopTimer()
	err := Init(BLS12_381)
	if err != nil {
		b.Fatal(err)
	}
	sec, pub, hash := getSecPubHash()
	sig := sec.SignHash(hash)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		sig.VerifyHash(pub, hash)
	}
	b.StopTimer()
}

func BenchmarkVerifyAggreageteHash(b *testing.B) {
	b.StopTimer()
	err := Init(BLS12_381)
	if err != nil {
		b.Fatal(err)
	}
	const N = 50
	pubVec, hashVec, sigVec := getSecPubHashVec(N)
	agg := sigVec[0]
	for i := 1; i < N; i++ {
		agg.Add(&sigVec[i])
	}
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		agg.VerifyAggregateHashes(pubVec, hashVec)
	}
}

func BenchmarkDeserialization1(b *testing.B) {
	b.StopTimer()
	err := Init(BLS12_381)
	if err != nil {
		b.Fatal(err)
	}
	VerifyOrderG1(false)
	VerifyOrderG2(false)
	const N = 50
	sec, _, hash := getSecPubHash()
	sig := sec.SignHash(hash)
	buf := sig.Serialize()
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		sig.Deserialize(buf)
	}
}

func BenchmarkDeserialization2(b *testing.B) {
	b.StopTimer()
	err := Init(BLS12_381)
	if err != nil {
		b.Fatal(err)
	}
	VerifyOrderG1(true)
	VerifyOrderG2(true)
	const N = 50
	sec, _, hash := getSecPubHash()
	sig := sec.SignHash(hash)
	buf := sig.Serialize()
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		sig.Deserialize(buf)
	}
}
