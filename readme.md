# bls for eth with compiled static library

# This is under construction!!!

This repository contains compiled static library of https://github.com/herumi/bls with `BLS_ETH=1`

# How to run sample.go
```
go get github.com/herumi/bls-eth-go-binary/
go run sample.go
```

# How to build the static binary
The following steps are not necessary if you use compiled binary in this repository.

* Linux, Mac, Windows(mingw64)
```
mkdir work
git clone htpps://github.com/herumi/mcl
git clone htpps://github.com/herumi/bls
cd mcl
make src/base64.ll
make BIT=32 src/base32.ll
cd ../bls
make minimized_static BLS_ETH=1 MIN_WITH_XBYAK=1 LIB_DIR=${GOPATH}/src/github.com/herumi/bls-eth-go-binary/bls/lib/${GOOS}/${GOARCH}/
```

* Android
```
cd android
ndk-build
```

* iOS
```
make
```

Copy each static library `libbls384_256.a` to `src/bls/lib/<os>/<arch>/`.
