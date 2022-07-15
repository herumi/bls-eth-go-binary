-include ETH.cfg
SRC_DIR?=src/
BLS_DIR=$(SRC_DIR)/bls
MCL_DIR=$(BLS_DIR)/mcl
all:
	$(MAKE) -f $(BLS_DIR)/Makefile.onelib BLS_DIR=$(BLS_DIR) MCL_DIR=$(MCL_DIR) OUT_DIR=$(shell pwd) ETH_CFLAGS=$(ETH_CFLAGS) all
ios:
	$(MAKE) -f $(BLS_DIR)/Makefile.onelib BLS_DIR=$(BLS_DIR) MCL_DIR=$(MCL_DIR) OUT_DIR=$(shell pwd) ETH_CFLAGS=$(ETH_CFLAGS) ios
ios_simulator:
	$(MAKE) -f $(BLS_DIR)/Makefile.onelib BLS_DIR=$(BLS_DIR) MCL_DIR=$(MCL_DIR) OUT_DIR=$(shell pwd) ETH_CFLAGS=$(ETH_CFLAGS) ios_simulator

NDK_BUILD?=ndk-build
android:
	$(MAKE) -f $(BLS_DIR)/Makefile.onelib BLS_DIR=$(BLS_DIR) MCL_DIR=$(MCL_DIR) OUT_DIR=$(shell pwd) ETH_CFLAGS=$(ETH_CFLAGS) NDK_BUILD=$(NDK_BUILD) BLS_LIB_SHARED=$(BLS_LIB_SHARED) android

update:
	cp $(BLS_DIR)/include/bls/bls.h bls/include/bls/.
	cp $(BLS_DIR)/include/bls/bls384_256.h bls/include/bls/.
	cp $(MCL_DIR)/include/mcl/bn.h bls/include/mcl/.
	cp $(MCL_DIR)/include/mcl/bn_c384_256.h bls/include/mcl/.
	cp $(MCL_DIR)/include/mcl/curve_type.h bls/include/mcl/.
	patch -o - -p0 ../bls/ffi/go/bls/mcl.go <patch/mcl.patch > bls/mcl.go
	patch -o - -p0 ../bls/ffi/go/bls/bls.go <patch/bls.patch > bls/bls.go

update_patch:
	-diff -up $(BLS_DIR)/ffi/go/bls/mcl.go bls/mcl.go > patch/mcl.patch
	-diff -up $(BLS_DIR)/ffi/go/bls/bls.go bls/bls.go > patch/bls.patch

clean:
	$(MAKE) -C $(MCL_DIR) clean
	$(MAKE) -C $(BLS_DIR) clean
	$(RM) -rf obj/*.o android/obj/* bls/lib/android/*

.PHONY: android ios each_ios clean
