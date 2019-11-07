# bls for eth with compiled static library

This repository contains compiled static library of https://github.com/herumi/bls with `BLS_ETH=1`.

* SecretKey; Fr
* PublicKey; G1
* Signature; G2

# How to run sample.go
```
go get github.com/herumi/bls-eth-go-binary/
go run sample.go
```

# How to build the static binary
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
