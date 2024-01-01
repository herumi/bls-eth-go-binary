# TBLS implementation with DKG procedure

> Credits to [ameba23](https://github.com/ameba23)

## Intro

This is the improvement for default BLS abilities based on this lib. With TBLS & DKG you and other parts can generate a single group pubkey and sets of secret shares to control this pubkey in decentralized way using `T/N` agreements (`T<N`). This is so-called threshold signatures. To use them - follow the tutorial


## 1. Generate random IDs

Saying, you and 2 your friends want to generate a multisig address in some cryptocurrency network and spend money only in case `2 of 3` friends agree with some solution.

Let's start with IDs generation

### Use package in your Golang project

```shell
go get -u github.com/herumi/bls-eth-go-binary/tbls
```

Then, generate random IDs

```go

totalNumberOfSigners := 3

randomIDs := tbls.GenerateRandomIds(totalNumberOfSigners)

fmt.Println("IDs are => ", randomIDs)
```

Output

```
IDs are =>  [   
    "70b631f57c3636805d5f8b93d6db7184977a9a6e5de427915724309f60a7fd32", // this is ID for user 1 "4049877b89f7b1e4cecebb613cea8944c23e41baabfe5628e7ae8aa86dc349af", // ...for user 2 "73dbd67216c2c9d5d6fd5a20b6c549466680e8ee471f8362ee39157aa9674f9b"  // ...for user 3
]
```

> This procedure can be done by one of the signers, no matter who. These IDs can be public and shareable

Now, based on these values - generate secret shares and verification vectors


## 2. Generate secret shares and verification vectors

After each of signers know own ID and IDs of other parties - generate secret shares and verification vectors. To do it in your Golang code

```go

threshold := 2 // because 2/3

randomIDs := []string{

    "70b631f57c3636805d5f8b93d6db7184977a9a6e5de427915724309f60a7fd32", "4049877b89f7b1e4cecebb613cea8944c23e41baabfe5628e7ae8aa86dc349af", "73dbd67216c2c9d5d6fd5a20b6c549466680e8ee471f8362ee39157aa9674f9b"}


// Each group member do it individually
vvec1, secretShares1 := tbls.GenerateTbls(threshold, randomIDs)
vvec2, secretShares2 := tbls.GenerateTbls(threshold, randomIDs)
vvec3, secretShares3 := tbls.GenerateTbls(threshold, randomIDs)
```

Now, distribute these values among all the signers in appropriate order

```go

//______User 1(ID=70b631f57c3636805d5f8b93d6db7184977a9a6e5de427915724309f60a7fd32)_________

// To itself
vvec1, secretShares1[0]

// Send to user2 (ID=4049877b89f7b1e4cecebb613cea8944c23e41baabfe5628e7ae8aa86dc349af)
vvec1, secretShares1[1]

// Send to user3 (ID=73dbd67216c2c9d5d6fd5a20b6c549466680e8ee471f8362ee39157aa9674f9b)
vvec1, secretShares1[2]




//______User 2(ID=4049877b89f7b1e4cecebb613cea8944c23e41baabfe5628e7ae8aa86dc349af)_________

// Send to user1 (ID=70b631f57c3636805d5f8b93d6db7184977a9a6e5de427915724309f60a7fd32)
vvec2, secretShares2[0]

// To itself
vvec2, secretShares2[1]

// Send to user3 (ID=73dbd67216c2c9d5d6fd5a20b6c549466680e8ee471f8362ee39157aa9674f9b)
vvec2, secretShares2[2]




//_____User 3(ID=73dbd67216c2c9d5d6fd5a20b6c549466680e8ee471f8362ee39157aa9674f9b)__________

// Send to user1 (ID=70b631f57c3636805d5f8b93d6db7184977a9a6e5de427915724309f60a7fd32)
vvec3, secretShares3[0]

// Send to user2 (ID=4049877b89f7b1e4cecebb613cea8944c23e41baabfe5628e7ae8aa86dc349af)
vvec3, secretShares3[1]

// To itself
vvec3, secretShares3[2]

```

After this step group members should have:

```shell

User 1 - vvec1,vvec2,vvec3,secretShares1[0],secretShares2[0],secretShares3[0]
User 2 - vvec1,vvec2,vvec3,secretShares1[1],secretShares2[1],secretShares3[1]
User 3 - vvec1,vvec2,vvec3,secretShares1[2],secretShares2[2],secretShares3[2]

```

## 3. After you get verification vector & secret share from other group member - you can verify it

For this - use verification vector and secret share that you'll get from another member + your ID from [step1](#1-generate-random-ids)

In case smth failed - `don't send money to group public key` because you'll lose control over it and group can ignore your choices.

```go
// Group member 1 should verify shares from members 2 and 3

isShareFromMember2Ok := tbls.VerifyShare("70b631f57c3636805d5f8b93d6db7184977a9a6e5de427915724309f60a7fd32",secretShares2[0])
isShareFromMember3Ok := tbls.VerifyShare("70b631f57c3636805d5f8b93d6db7184977a9a6e5de427915724309f60a7fd32",secretShares3[0])

shouldMember1Trust := isShareFromMember2Ok && isShareFromMember3Ok

// Group member 2 should verify shares from members 1 and 3

isShareFromMember1Ok := tbls.VerifyShare("4049877b89f7b1e4cecebb613cea8944c23e41baabfe5628e7ae8aa86dc349af",secretShares1[1])
isShareFromMember3Ok := tbls.VerifyShare("4049877b89f7b1e4cecebb613cea8944c23e41baabfe5628e7ae8aa86dc349af",secretShares3[1])

shouldMember2Trust := isShareFromMember1Ok && isShareFromMember3Ok

// Group member 3 should verify shares from members 1 and 2

isShareFromMember1Ok := tbls.VerifyShare("73dbd67216c2c9d5d6fd5a20b6c549466680e8ee471f8362ee39157aa9674f9b",secretShares1[2])
isShareFromMember2Ok := tbls.VerifyShare("73dbd67216c2c9d5d6fd5a20b6c549466680e8ee471f8362ee39157aa9674f9b",secretShares2[2])

shouldMember3Trust := isShareFromMember1Ok && isShareFromMember2Ok
```

## 4. Derive group public key

To identify your group - get the master public key based on verification vectors from all the group members. To do it in Golang code:

```go
rootPubKey := tbls.DeriveRootPubKey(vvec1, vvec2, vvec3) // a29ac612857e8c0d86c0e2e1d1794421c55a5c2ae13585dc7e9e22075cb2d6dc6b1312a36d188f421ffb83647e515f18
```

> Now you can securely send money to `a29ac612857e8c0d86c0e2e1d1794421c55a5c2ae13585dc7e9e22075cb2d6dc6b1312a36d188f421ffb83647e515f18`


## 5. Generate partial signatures by `T/N` members

Later, `2/3` of group members(in our case) decide to spend cryptocurrency from this group wallet. To do this, they need to generate partial signatures for further aggregation into a master signature on behalf of the entire group, which can be verified using the group’s master public key(see previous [step](#4-derive-group-public-key)).

Saying, members 1 and 2 generate partial signatures while member 3 is AFK(died, disagree, etc.)


```go
secretSharesFor1 := []string{secretShares1[0], secretShares2[0], secretShares3[0]}
secretSharesFor2 := []string{secretShares1[1], secretShares2[1], secretShares3[1]}

msg := "Buy BTC ₿"

partialSignature1 := tbls.GeneratePartialSignature("70b631f57c3636805d5f8b93d6db7184977a9a6e5de427915724309f60a7fd32", msg, secretSharesFor1)
partialSignature2 := tbls.GeneratePartialSignature("4049877b89f7b1e4cecebb613cea8944c23e41baabfe5628e7ae8aa86dc349af", msg, secretSharesFor2)

fmt.Println("Partial signature 1 is => ", partialSignature1)
fmt.Println("Partial signature 2 is => ", partialSignature2)
```

Output:

```
Partial signature 1 is =>  90af35830e8ab7059c6258cea0ac31e409a84f22449a75004c510e0c334446f36389b3106c45f3f5dcaeee9f97ddfb37023014db7630fcb2ffed03e9b0f0e66e83fc7e6d7b2f9fbd5b50e4fdf0f81a33cf830f752326af1d1bf8e5ad5c06d1fd
Partial signature 2 is =>  8c4975e80d24e0d453039630e83e6f8aeb8b59ece958e1740d258bc3742547e167fd59f63b2e226704c4f3ca7213a84a06d658693a140578d6c08d566447c41f620ab5c7da700f8498455bb47131817f01202afcd84c8f1d01f3eabd9c7beb6a
```

## 6. Aggregate partial signatures

Now, anyone can take these 2 partial signatures and IDs of signers and aggregate them to get the master signature

```go

hexIDOfUser1 := "70b631f57c3636805d5f8b93d6db7184977a9a6e5de427915724309f60a7fd32"
hexIDOfUser2 := "4049877b89f7b1e4cecebb613cea8944c23e41baabfe5628e7ae8aa86dc349af"

rootSignature := tbls.BuildRootSignature([]string{partialSignature1, partialSignature2}, []string{hexIDOfUser1, hexIDOfUser2})

fmt.Println("Root signature is => ", rootSignature)
```

Output:

```
Root signature is =>  ae7225c9a28e504ba5ab64cf083fef4ea907a2592283ff8c3ba5cee10b04be6dccae95fea9cac7c211ba8b4851928a9007bf80c466384ddef1d25e367e3cc31acdaf9acea296e0f98c0d6ed3dc45e99b0a2ce453daf7bc9b16d808f4300f6a94
```


## 7. Verify master signature with group pubkey

```go

msg := "Buy BTC ₿"

rootPubKey := "a29ac612857e8c0d86c0e2e1d1794421c55a5c2ae13585dc7e9e22075cb2d6dc6b1312a36d188f421ffb83647e515f18"

rootSignature := "ae7225c9a28e504ba5ab64cf083fef4ea907a2592283ff8c3ba5cee10b04be6dccae95fea9cac7c211ba8b4851928a9007bf80c466384ddef1d25e367e3cc31acdaf9acea296e0f98c0d6ed3dc45e99b0a2ce453daf7bc9b16d808f4300f6a94"

fmt.Println("Is root signature ok ? => ", tbls.VerifyRootSignature(rootPubKey, rootSignature, msg))
```****