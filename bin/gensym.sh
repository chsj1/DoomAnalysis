#!/bin/sh

VERSION=$1

cd "./dump/$VERSION"
cp "../../lib/${VERSION}"_libgame.so ./libgame.so
../../tools/dump_syms libgame.so > libgame.so.sym
rm -f libgame.so

cd "../../lib"
touch "$VERSION".txt

