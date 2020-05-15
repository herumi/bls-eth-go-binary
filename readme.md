[![Build Status](https://travis-ci.org/herumi/bls-eth-go-binary.png)](https://travis-ci.org/herumi/bls-eth-go-binary)
# bls for eth with compiled static library

This repository contains compiled static library of https://github.com/herumi/bls with `BLS_ETH=1`.

* SecretKey; Fr
* PublicKey; G1
* Sign; G2

# News
- 2020/May/15 `EthModeDraft07` is added for [draft-07](https://www.ietf.org/id/draft-irtf-cfrg-hash-to-curve-07.txt).
- 2020/Apr/20 `EthModeDraft06` is default. Call `SetETHmode(EthModeDraft05)` to use older evrsion.
- 2020/Mar/26 The signature value in `SetETHmode(2)` has changed because of changing DST in hash-to-curve function.
- 2020/Mar/17 This library supports [eth2.0 functions](https://github.com/ethereum/eth2.0-specs/blob/dev/specs/phase0/beacon-chain.md#bls-signatures). But the spec of hash-to-curve function may be changed.

* `SetETHmode(EthModeDraft05)` supports draft-irtf-cfrg-hash-to-curve-05, which will be removed in the near future.
* `SetETHmode(EthModeFraft06)` supports [draft-irtf-cfrg-hash-to-curve-06](https://cfrg.github.io/draft-irtf-cfrg-hash-to-curve/draft-irtf-cfrg-hash-to-curve.txt) at March 2020. But it has not yet fully tested.

Init as the followings:

```
Init(BLS12_381)
SetETHmode(1) // or SetETHmode(2)
```

then, you can use the following functions.

bls-eth-go-binary | eth2.0 spec name|
------|-----------------|
SecretKey::SignByte|Sign|
PublicKey::VerifyByte|Verify|
Sign::Aggregate|Aggregate|
Sign::FastAggregateVerify|FastAggregateVerify|
Sign::AggregateVerifyNoCheck|AggregateVerify|

The size of message must be 32 byte.

Check functions:
- VerifySignatureOrder ; make `deserialize` check the correctness of the order
- Sign::IsValidOrder ; check the correctness of the order
- VerifyPublicKeyOrder ; make `deserialize` check the correctness of the order
- PublicKey::IsValidOrder ; check the correctness of the order
- AreAllMsgDifferent ; check that all messages are different each other

# How to run `examples/sample.go`

```
go get github.com/herumi/bls-eth-go-binary/
go run examples/sample.go
```

# How to build the static library
The following steps are not necessary if you use compiled binary in this repository.

## Linux, Mac, Windows(mingw64)
```
mkdir work
cd work
git clone https://github.com/herumi/mcl
git clone https://github.com/herumi/bls
git clone https://github.com/herumi/bls-eth-go-binary
cd bls-eth-go-binary
make CXX=clang++
```

clang generates better binary than gcc.

## Android
```
make android
```

## iOS
```
make ios
```

# How to use the static library from C
```
#define BLS_ETH
#include <mcl/bn_c384_256.h>
#include <bls/bls.h>
```

## Author
MITSUNARI Shigeo(herumi@nifty.com)

## Sponsors welcome
[GitHub Sponsor](https://github.com/sponsors/herumi)
