#!/bin/sh

dir=$(dirname "$0")
cd "${dir}" || exit 1
rm symlink.md
ln -s fileA.md symlink.md
