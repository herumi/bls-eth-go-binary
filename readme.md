# bls for eth with compiled static library


This repository contains compiled static library of https://github.com/herumi/bls with `BLS_ETH=1`

# How to run sample.go
```
go get github.com/herumi/bls-eth-go-binary/
go run sample.go
```

# How to build the static binary
The following steps are not necessary if you use compiled binary in this repository.

* Linux, Mac, Windows(mingw64)
clang++ is necessary
```
mkdir work
git clone https://github.com/herumi/mcl
git clone https://github.com/herumi/bls
git clone https://github.com/herumi/bls-eth-go-binary
cd bls-eth-go-binary
make CXX=clang++ # better performance than gcc
```

* Android
At first, setup Android SDK
```
make android
```

* iOS
```
make ios
```

