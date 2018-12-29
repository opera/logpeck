#!/bin/bash

template_dir="logpeck.template"
if [ ! -d $template_dir ]; then
	echo "error: need $template_dir"
	exit 1
fi
if [ $GOPATH == "" ]; then 
	echo "error: none GOPATH"
	exit 1
fi

verfile="$GOPATH/src/github.com/opera/logpeck/version.go"
if [ ! -f $verfile ]; then 
	echo "error: version file not exist: $verfile"
	exit 1
fi
ver=`grep "const VersionString" $verfile | cut -d'"' -f 2`

pkg_dir="logpeck_$ver"
if [ -d $pkg_dir ] || [ -f $pkg_dir ]; then
	echo "error: package name exist, remove first: $pkg_dir"
	exit 1
fi

cp -r $template_dir $pkg_dir

sed -i "s/LOGPECK_VERSION/$ver/" $pkg_dir/DEBIAN/control

bin_dir="$pkg_dir/usr/bin"
mkdir -p $bin_dir

go build -v github.com/opera/logpeck/cmd/logpeckd
if [ $? -ne 0 ]; then 
	echo "error: build logpeckd error"
	exit 1
fi

mv logpeckd $bin_dir

dpkg -b $pkg_dir
