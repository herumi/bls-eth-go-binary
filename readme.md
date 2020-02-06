# bls for eth with compiled static library

This repository contains compiled static library of https://github.com/herumi/bls with `BLS_ETH=1`.

* SecretKey; Fr
* PublicKey; G1
* Sign; G2

# News
The new [eth2.0 functions](https://github.com/ethereum/eth2.0-specs/blob/dev/specs/phase0/beacon-chain.md#bls-signatures) are supported.

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

# How to run sample.go
```
go get github.com/herumi/bls-eth-go-binary/
go run sample.go
```

# How to build the static library
The following steps are not necessary if you use compiled binary in this repository.

```
mkdir work
git clone https://github.com/herumi/mcl
git clone https://github.com/herumi/bls
```

* Linux, Mac, Windows(mingw64)
clang generates better binary than gcc.
```
make CXX=clang++
```

* Android
```
make android
```

* iOS
```
make ios
```

# How to use the static library from C
```
#define BLS_ETH
#include <mcl/bn_c384_256.h>
#include <bls/bls.h>
```
