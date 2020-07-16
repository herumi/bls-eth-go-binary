#!/bin/bash

# Cross compile for 32bit MIPS architecture

# define the usage
usage () {
	echo "Arguments:"
	echo "  -buildroot <path to buildroot>"
	echo "  -lib \"<list of additional libraries to link in compile>\""
	echo ""
}

# define the path to the buildroot
BUILDROOT_PATH=""
USER_LIBS=""
INCLUDE_LIBS=""
bDebug=0


###################
# parse arguments #
while [ "$1" != "" ]
do
	case "$1" in
		# options
		buildroot|-buildroot|--buildroot)
			shift
			BUILDROOT_PATH="$1"
			shift
		;;
		include|-include|--include)
			shift
			INCLUDE_LIBS="$1"
			shift
		;;
		lib|-lib|--lib)
			shift
			USER_LIBS="$1"
			shift
		;;
		-d|--d|debug|-debug|--debug)
			bDebug=1
			shift
		;;
		-h|--h|help|-help|--help)
			usage
			exit
		;;
	    *)
			echo "ERROR: Invalid Argument: $1"
			usage
			exit
		;;
	esac
done

# check to ensure correct arguments
if [ "$BUILDROOT_PATH" = "" ]
then
	echo "ERROR: missing path to buildroot"
	echo ""
	usage
	exit
fi


# define the toolchain and target names
TOOLCHAIN_NAME="toolchain-mipsel_24kc_gcc-7.3.0_musl"
TARGET_NAME="target-mipsel_24kc_musl"

# define the relative paths
STAGING_DIR_RELATIVE="staging_dir"
TOOLCHAIN_RELATIVE="$STAGING_DIR_RELATIVE/$TOOLCHAIN_NAME"
TARGET_RELATIVE="$STAGING_DIR_RELATIVE/$TARGET_NAME"

# define the toolchain paths
TOOLCHAIN="$BUILDROOT_PATH/$TOOLCHAIN_RELATIVE"
TOOLCHAIN_BIN="$BUILDROOT_PATH/$TOOLCHAIN_RELATIVE/bin"

TOOLCHAIN_INCLUDE="$BUILDROOT_PATH/$TOOLCHAIN_RELATIVE/include"
TOOLCHAIN_LIB="$BUILDROOT_PATH/$TOOLCHAIN_RELATIVE/lib"
TOOLCHAIN_USR_INCLUDE="$BUILDROOT_PATH/$TOOLCHAIN_RELATIVE/usr/include"
TOOLCHAIN_USR_LIB="$BUILDROOT_PATH/$TOOLCHAIN_RELATIVE/usr/lib"

# define the target paths
TARGET="$BUILDROOT_PATH/$TARGET_RELATIVE"

TARGET_INCLUDE="$BUILDROOT_PATH/$TARGET_RELATIVE/include"
TARGET_LIB="$BUILDROOT_PATH/$TARGET_RELATIVE/lib"
TARGET_USR_INCLUDE="$BUILDROOT_PATH/$TARGET_RELATIVE/usr/include"
TARGET_USR_LIB="$BUILDROOT_PATH/$TARGET_RELATIVE/usr/lib"

export STAGING_DIR="BUILDROOT_PATH/$STAGING_DIR_RELATIVE"

# define the compilers and such
TOOLCHAIN_CC="$TOOLCHAIN_BIN/mipsel-openwrt-linux-gcc"
TOOLCHAIN_CXX="$TOOLCHAIN_BIN/mipsel-openwrt-linux-g++"
TOOLCHAIN_LD="$TOOLCHAIN_BIN/mipsel-openwrt-linux-ld"

TOOLCHAIN_AR="$TOOLCHAIN_BIN/mipsel-openwrt-linux-ar"
TOOLCHAIN_RANLIB="$TOOLCHAIN_BIN/mipsel-openwrt-linux-ranlib"




# define the FLAGS
INCLUDE_LINES="-I $TOOLCHAIN_USR_INCLUDE -I $TOOLCHAIN_INCLUDE -I $TARGET_USR_INCLUDE -I $TARGET_INCLUDE -I $INCLUDE_LIBS"
TOOLCHAIN_CFLAGS="-Os -pipe -mno-branch-likely -mips32r2 -mtune=24kc -fno-caller-saves -fno-plt -fhonour-copts -Wno-error=unused-but-set-variable -Wno-error=unused-result -msoft-float -mips16 -minterlink-mips16 -Wformat -Werror=format-security -fstack-protector -D_FORTIFY_SOURCE=1 -Wl,-z,now -Wl,-z,relro"
#TOOLCHAIN_CFLAGS="-Os -pipe -mno-branch-likely -mips32r2 -mtune=34kc -fno-caller-saves -fhonour-copts -Wno-error=unused-but-set-variable -Wno-error=unused-result -msoft-float -mips16 -minterlink-mips16 -fpic"
TOOLCHAIN_CFLAGS="$TOOLCHAIN_CFLAGS $INCLUDE_LINES"

TOOLCHAIN_CXXFLAGS="$TOOLCHAIN_CFLAGS"
#TOOLCHAIN_CXXFLAGS="-Os -pipe -mno-branch-likely -mips32r2 -mtune=34kc -fno-caller-saves -fhonour-copts -Wno-error=unused-but-set-variable -Wno-error=unused-result -msoft-float -mips16 -minterlink-mips16 -fpic"
TOOLCHAIN_CXXFLAGS="$TOOLCHAIN_CXXFLAGS $INCLUDE_LINES"

TOOLCHAIN_LDFLAGS="-L$TOOLCHAIN_USR_LIB -L$TOOLCHAIN_LIB -L$TARGET_USR_LIB -L$TARGET_LIB"

# debug
if [ $bDebug -eq 1 ]; then
	echo "CC=$TOOLCHAIN_CC"
	echo "CXX=$TOOLCHAIN_CXX"
	echo "LD=$TOOLCHAIN_LD"
	echo "CFLAGS=$TOOLCHAIN_CFLAGS"
	echo "LDFLAGS=$TOOLCHAIN_LDFLAGS"
	echo "USER_LIBS=$USER_LIBS"
	echo ""
fi

# first run make clean
make clean
# run the make command
make \
	CC="$TOOLCHAIN_CC" \
	CXX="$TOOLCHAIN_CXX" \
	LD="$TOOLCHAIN_LD" \
	CFLAGS="$TOOLCHAIN_CFLAGS" \
	LDFLAGS="$TOOLCHAIN_LDFLAGS" \
	LIB="$USER_LIBS" \
	AR="$TOOLCHAIN_AR"
