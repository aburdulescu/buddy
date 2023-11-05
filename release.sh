#!/bin/sh

set -e

VERSION=$1

if [ -z $VERSION ]
then
    echo "VERSION not set"
    exit 1
fi

set -x

git tag $VERSION
git push origin $VERSION
GOPROXY=proxy.golang.org go list -m bandr.me/p/buddy@$VERSION
