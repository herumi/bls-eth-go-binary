include ../mcl/common.mk
ETH_CFLAGS=-DBLS_ETH -DBLS_SWAP_G

UNIT?=8

MIN_CFLAGS=-std=c++03 -O3 -DNDEBUG -DMCL_DONT_USE_OPENSSL -DMCL_LLVM_BMI2=0 -DMCL_USE_LLVM=1 -DMCL_USE_VINT -DMCL_SIZEOF_UNIT=$(UNIT) -DMCL_VINT_FIXED_BUFFER -DMCL_MAX_BIT_SIZE=384 -DCYBOZU_DONT_USE_EXCEPTION -DCYBOZU_DONT_USE_STRING -D_FORTIFY_SOURCE=0 -I../bls/include -I../mcl/include $(ETH_CFLAGS) $(CFLAGS_USER)
OBJ_DIR=obj

all: ../mcl/src/base$(BIT).ll
ifeq ($(CPU),x86-64)
	$(eval _ARCH=amd64)
ifeq ($(OS),mingw64)
	$(eval _OS=windows)
endif
ifeq ($(OS),Linux)
	$(eval _OS=linux)
	$(eval MIN_CFLAGS=$(MIN_CFLAGS) -fPIC)
endif
ifeq ($(OS),mac)
	$(eval _OS=darwin)
	$(eval MIN_CFLAGS=$(MIN_CFLAGS) -fPIC)
endif
endif
ifeq ($(CPU),aarch64)
	$(eval _ARCH=arm64)
ifeq ($(OS),Linux)
	$(eval _OS=linux)
	$(eval MIN_CFLAGS=$(MIN_CFLAGS) -fPIC)
endif
endif
	$(eval LIB_DIR=bls/lib/$(_OS)/$(_ARCH))
	-mkdir -p $(LIB_DIR)
	$(CXX) -c -o $(OBJ_DIR)/fp.o ../mcl/src/fp.cpp $(MIN_CFLAGS)
	$(CXX) -c -o $(OBJ_DIR)/base$(BIT).o ../mcl/src/base$(BIT).ll $(MIN_CFLAGS)
	$(CXX) -c -o $(OBJ_DIR)/bls_c384_256.o ../bls/src/bls_c384_256.cpp $(MIN_CFLAGS)
	$(AR) $(LIB_DIR)/libbls384_256.a $(OBJ_DIR)/bls_c384_256.o $(OBJ_DIR)/fp.o $(OBJ_DIR)/base$(BIT).o

BASE_LL=../mcl/src/base64.ll ../mcl/src/base32.ll

../mcl/src/base64.ll:
	$(MAKE) -C ../mcl src/base64.ll

../mcl/src/base32.ll:
	$(MAKE) -C ../mcl src/base32.ll BIT=32

ANDROID_TARGET=armeabi-v7a arm64-v8a x86_64
android: $(BASE_LL)
	@ndk-build -C android/jni NDK_DEBUG=0
	@for target in $(ANDROID_TARGET); do \
		mkdir -p bls/lib/android/$$target; \
		cp android/obj/local/$$target/libbls384_256.a bls/lib/android/$$target/; \
	done

# ios
XCODEPATH=$(shell xcode-select -p)
IOS_CLANG=$(XCODEPATH)/Toolchains/XcodeDefault.xctoolchain/usr/bin/clang
IOS_AR=${XCODEPATH}/Toolchains/XcodeDefault.xctoolchain/usr/bin/ar
PLATFORM?="iPhoneOS"
IOS_MIN_VERSION?=7.0
IOS_CFLAGS=-fembed-bitcode -fno-common -DPIC -fPIC -Dmcl_EXPORTS
IOS_CFLAGS+=-DMCL_USE_VINT -DMCL_VINT_FIXED_BUFFER -DMCL_DONT_USE_OPENSSL -DMCL_DONT_USE_XBYAK -DMCL_LLVM_BMI2=0 -DMCL_USE_LLVM=1 -std=c++11 -Wall -Wextra -Wformat=2 -Wcast-qual -Wcast-align -Wwrite-strings -Wfloat-equal -Wpointer-arith -O3 -DNDEBUG $(ETH_CFLAGS)
IOS_CFLAGS+=-I../mcl/include -I../bls/include
IOS_LDFLAGS=-dynamiclib -Wl,-flat_namespace -Wl,-undefined -Wl,suppress
CURVE_BIT?=384_256
IOS_LIB=libbls$(CURVE_BIT).a
IOS_LIBS=ios/armv7/$(IOS_LIB) ios/arm64/$(IOS_LIB) ios/x86_64/$(IOS_LIB) ios/i386/$(IOS_LIB)

ios:
	$(MAKE) each_ios PLATFORM="iPhoneOS" ARCH=armv7 BIT=32 UNIT=4
	$(MAKE) each_ios PLATFORM="iPhoneOS" ARCH=arm64 BIT=64 UNIT=8
	$(MAKE) each_ios PLATFORM="iPhoneSimulator" ARCH=x86_64 BIT=64 UNIT=8
	$(MAKE) each_ios PLATFORM="iPhoneSimulator" ARCH=i386 BIT=32 UNIT=4
	@echo $(IOS_LIBS)
	@mkdir -p bls/lib/ios
	lipo $(IOS_LIBS) -create -output bls/lib/ios/$(IOS_LIB)

each_ios: $(BASE_LL)
	@echo "Building iOS $(ARCH) BIT=$(BIT) UNIT=$(UNIT)"
	$(eval IOS_CFLAGS=$(IOS_CFLAGS) -DMCL_SIZEOF_UNIT=$(UNIT))
	@echo IOS_CFLAGS=$(IOS_CFLAGS)
	$(eval IOS_OUTDIR=ios/$(ARCH))
	$(eval IOS_SDK_PATH=$(XCODEPATH)/Platforms/$(PLATFORM).platform/Developer/SDKs/$(PLATFORM).sdk)
	$(eval IOS_COMMON=-arch $(ARCH) -isysroot $(IOS_SDK_PATH) -mios-version-min=$(IOS_MIN_VERSION))
	@mkdir -p $(IOS_OUTDIR)
	$(IOS_CLANG) $(IOS_COMMON) $(IOS_CFLAGS) -c ../mcl/src/fp.cpp -o $(IOS_OUTDIR)/fp.o
	$(IOS_CLANG) $(IOS_COMMON) $(IOS_CFLAGS) -c ../mcl/src/base$(BIT).ll -o $(IOS_OUTDIR)/base$(BIT).o
	$(IOS_CLANG) $(IOS_COMMON) $(IOS_CFLAGS) -c ../bls/src/bls_c$(CURVE_BIT).cpp -o $(IOS_OUTDIR)/bls_c$(CURVE_BIT).o
	ar cru $(IOS_OUTDIR)/$(IOS_LIB) $(IOS_OUTDIR)/fp.o $(IOS_OUTDIR)/base$(BIT).o $(IOS_OUTDIR)/bls_c$(CURVE_BIT).o
	ranlib $(IOS_OUTDIR)/$(IOS_LIB)

update:
	cp ../bls/include/bls/bls.h bls/include/bls/.
	cp ../bls/include/bls/bls384_256.h bls/include/bls/.
	cp ../mcl/include/mcl/bn.h bls/include/mcl/.
	cp ../mcl/include/mcl/bn_c384_256.h bls/include/mcl/.
	cp ../mcl/include/mcl/curve_type.h bls/include/mcl/.
	patch -o - -p0 ../bls/ffi/go/bls/mcl.go <patch/mcl.patch > bls/mcl.go
	patch -o - -p0 ../bls/ffi/go/bls/bls.go <patch/bls.patch > bls/bls.go

update_patch:
	-diff -up ../bls/ffi/go/bls/mcl.go bls/mcl.go > patch/mcl.patch
	-diff -up ../bls/ffi/go/bls/bls.go bls/bls.go > patch/bls.patch

.PHONY: android ios each_ios
