# ios
XCODEPATH=$(shell xcode-select -p)
IOS_CLANG=$(XCODEPATH)/Toolchains/XcodeDefault.xctoolchain/usr/bin/clang
IOS_AR=${XCODEPATH}/Toolchains/XcodeDefault.xctoolchain/usr/bin/ar
PLATFORM?="iPhoneOS"
IOS_MIN_VERSION?=7.0
IOS_CFLAGS=-fembed-bitcode -fno-common -DPIC -fPIC -Dmcl_EXPORTS
IOS_CFLAGS+=-DMCL_USE_VINT -DMCL_VINT_FIXED_BUFFER -DMCL_DONT_USE_OPENSSL -DMCL_DONT_USE_XBYAK -DMCL_LLVM_BMI2=0 -DMCL_USE_LLVM=1 -DMCL_SIZEOF_UNIT=8 -std=c++11 -Wall -Wextra -Wformat=2 -Wcast-qual -Wcast-align -Wwrite-strings -Wfloat-equal -Wpointer-arith -O3 -DNDEBUG
IOS_CFLAGS+=-I../mcl/include -I../bls/include
IOS_CFLAGS+=-DBLS_ETH -DBLS_SWAP_G
IOS_LDFLAGS=-dynamiclib -Wl,-flat_namespace -Wl,-undefined -Wl,suppress
CURVE_BIT?=384_256
IOS_OBJS=$(IOS_OUTDIR)/fp.o $(IOS_OUTDIR)/base64.o $(IOS_OUTDIR)/bls_c$(CURVE_BIT).o
IOS_LIB=libbls$(CURVE_BIT)

GOMOBILE_ARCHS=armv7 arm64

all:
	@for target in $(GOMOBILE_ARCHS); do \
		$(MAKE) ios ARCH=$$target PLATFORM="iPhoneOS"; \
	done

../mcl/src/base64.ll:
	$(MAKE) -C ../mcl src/base64.ll

ios: ../mcl/src/base64.ll
	@echo "Building iOS $(ARCH)..."
	$(eval IOS_OUTDIR=ios/$(ARCH))
	$(eval IOS_SDK_PATH=$(XCODEPATH)/Platforms/$(PLATFORM).platform/Developer/SDKs/$(PLATFORM).sdk)
	$(eval IOS_COMMON=-arch $(ARCH) -isysroot $(IOS_SDK_PATH) -mios-version-min=$(IOS_MIN_VERSION))
	@mkdir -p $(IOS_OUTDIR)
	$(IOS_CLANG) $(IOS_COMMON) $(IOS_CFLAGS) -c ../mcl/src/fp.cpp -o $(IOS_OUTDIR)/fp.o
	$(IOS_CLANG) $(IOS_COMMON) $(IOS_CFLAGS) -c ../mcl/src/base64.ll -o $(IOS_OUTDIR)/base64.o
	$(IOS_CLANG) $(IOS_COMMON) $(IOS_CFLAGS) -c ../bls/src/bls_c$(CURVE_BIT).cpp -o $(IOS_OUTDIR)/bls_c$(CURVE_BIT).o
	ar cru $(IOS_OUTDIR)/$(IOS_LIB).a $(IOS_OBJS)
	ranlib $(IOS_OUTDIR)/$(IOS_LIB).a

update:
	cp ../bls/include/bls/bls.h bls/include/bls/.
	cp ../bls/include/bls/bls384_256.h bls/include/bls/.
	cp ../mcl/include/mcl/bn.h bls/include/mcl/.
	cp ../mcl/include/mcl/bn_c384_256.h bls/include/mcl/.
	cp ../mcl/include/mcl/curve_type.h bls/include/mcl/.
	patch -o - -p0 ../bls/ffi/go/bls/mcl.go <patch/mcl.patch > bls/mcl.go
	patch -o - -p0 ../bls/ffi/go/bls/bls.go <patch/bls.patch > bls/bls.go

.PHONY: ios
