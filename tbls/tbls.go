package tbls

import (
	bls "github.com/herumi/bls-eth-go-binary/bls"
)

// Function to generate random IDs for N parties of DKG
func GenerateRandomIds(numberOfParties uint) (idsAsHex []string) {

	bls.Init(bls.BLS12_381)

	for i := 0; i < int(numberOfParties); i++ {

		secretkey := bls.SecretKey{}

		secretkey.SetByCSPRNG()

		idsAsHex = append(idsAsHex, secretkey.SerializeToHexStr())

	}

	return

}

// Function to generate verification vector and secret shares to distribute them among N parties
func GenerateTbls(threshold uint, partiesIDs []string) (verificationVector, sharesForOtherParties []string) {

	bls.Init(bls.BLS12_381)

	var secretVector []bls.SecretKey

	// Generate a secretVector and verificationVector

	for i := 0; i < int(threshold); i++ {

		secretkey := bls.SecretKey{}

		secretkey.SetByCSPRNG()

		secretVector = append(secretVector, secretkey)

		publicKey := secretkey.GetPublicKey()

		verificationVector = append(verificationVector, publicKey.SerializeToHexStr())

	}

	// Generate key shares for other parties

	for _, hexID := range partiesIDs {

		sk := bls.SecretKey{}

		sk.SetByCSPRNG()

		emptyID := &bls.ID{}

		emptyID.DeserializeHexStr(hexID)

		sk.Set(secretVector, emptyID)

		sharesForOtherParties = append(sharesForOtherParties, sk.SerializeToHexStr())

	}

	return

}

// Function to verify share which you'll get from other parties based on their verification vector + your ID
func VerifyShare(yourIdAsHex, secretShareAsHex string, verificationVectorAsHex []string) bool {

	bls.Init(bls.BLS12_381)

	pubKeyToVerifyShare := bls.PublicKey{}

	yourID := &bls.ID{}

	yourID.DeserializeHexStr(yourIdAsHex)

	// Deserialize verification vector ( []Hex => []bls.PublicKey )

	vvec := make([]bls.PublicKey, len(verificationVectorAsHex))

	for index, hexValueInVerificationVector := range verificationVectorAsHex {

		vvec[index].DeserializeHexStr(hexValueInVerificationVector)

	}

	pubKeyToVerifyShare.Set(vvec, yourID)

	// Deserialize secret share that you get from other signer

	secretShare := bls.SecretKey{}

	secretShare.DeserializeHexStr(secretShareAsHex)

	pubKeyFromShare := *secretShare.GetPublicKey()

	return pubKeyFromShare.IsEqual(&pubKeyToVerifyShare)

}

// Function to get the rootPubKey from verification vectors of group members
func DeriveRootPubKey(verificationVectors ...[]string) string {

	bls.Init(bls.BLS12_381)

	/*

		<verificationVectors> is a slice of verification vectors that you should get

		from all the group members (N/N)

	*/

	// First of all - take the first slice and create pubkeys based on it

	futureRootPubKey := []bls.PublicKey{}

	for _, hexPubKeyAsPartOfVerificationVector := range verificationVectors[0] {

		pubKeyTemplate := bls.PublicKey{}

		pubKeyTemplate.DeserializeHexStr(hexPubKeyAsPartOfVerificationVector)

		futureRootPubKey = append(futureRootPubKey, pubKeyTemplate)

	}

	// Now iterate over verification vectors with indexes 1,2,... and aggregtate with pubkeys with appropriate indexes (futureRootPubKey[1],futureRootPubKey[2],...)

	for indexOfVector, verificationVectorBySomeSigner := range verificationVectors {

		// Omit the first vector because we already add it to <futureRootPubKey>

		if indexOfVector > 0 {

			for indexOfPubKey, hexPubKeyAsPartOfVerificationVector := range verificationVectorBySomeSigner {

				pubKeyTemplate := bls.PublicKey{}

				pubKeyTemplate.DeserializeHexStr(hexPubKeyAsPartOfVerificationVector)

				// Aggregate

				futureRootPubKey[indexOfPubKey].Add(&pubKeyTemplate)

			}

		}

	}

	return futureRootPubKey[0].SerializeToHexStr()

}

// Function to generate partial signature as one of the T/N signers
// Theese signatures by T parties will be aggregated later and can be verified with rootPubKey(pubkey of group)
func GeneratePartialSignature(yourIdAsHex, message string, secretSharesAsHex []string) (partialSignatureAsHex string) {

	bls.Init(bls.BLS12_381)

	// Recover ID

	emptyID := bls.ID{}

	emptyID.DeserializeHexStr(yourIdAsHex)

	// Recover the secret based on secret shares received by you from other T signers

	groupSecret := bls.SecretKey{}

	for _, secretShareAsHex := range secretSharesAsHex {

		emptyShare := bls.SecretKey{}

		emptyShare.DeserializeHexStr(secretShareAsHex)

		// Aggregate

		groupSecret.Add(&emptyShare)

	}

	// Finally - generate signature

	partialSignatureAsHex = groupSecret.Sign(message).SerializeToHexStr()

	return

}

// Now based on partial signatures + IDs of T/N signers - recover the root signature
func BuildRootSignature(partialSignaturesAsHex, idsOfSignersAsHex []string) (rootSignatureAsHex string) {

	bls.Init(bls.BLS12_381)

	// Deserialize partial signatures

	partialSignatures := []bls.Sign{}

	for _, serializedPartialSignature := range partialSignaturesAsHex {

		templateForSigna := bls.Sign{}

		templateForSigna.DeserializeHexStr(serializedPartialSignature)

		partialSignatures = append(partialSignatures, templateForSigna)

	}

	// Deserialize IDs

	ids := []bls.ID{}

	for _, serializedID := range idsOfSignersAsHex {

		templateForID := bls.ID{}

		templateForID.DeserializeHexStr(serializedID)

		ids = append(ids, templateForID)

	}

	//___________ Now recover and return in serialized form___________

	rootSignature := bls.Sign{}

	rootSignature.Recover(partialSignatures, ids)

	rootSignatureAsHex = rootSignature.SerializeToHexStr()

	return

}

func VerifyRootSignature(rootPubKeyAsHex, rootSignaAsHex, message string) bool {

	bls.Init(bls.BLS12_381)

	rootPubKey := bls.PublicKey{}

	rootPubKey.DeserializeHexStr(rootPubKeyAsHex)

	rootSignature := bls.Sign{}

	rootSignature.DeserializeHexStr(rootSignaAsHex)

	return rootSignature.Verify(&rootPubKey, message)

}
