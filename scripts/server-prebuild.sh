#!/usr/bin/env bash
pushd client/ || exit
yarn build
popd || exit
rm -rf server/config.sample.nt
cp config/config.sample.nt server/
rm -rf server/client-dist
cp -r client/dist server/client-dist
