#!/bin/bash

RELEASE=0.4.0-187

#TODO: REMOVE BEFORE PROVIDING TO PUBLIC
TOKEN=6678b9316aff31d12efb00ff57d0633e5e87e092

echo "Downloading ion-connect-$RELEASE.tar.gz"
curl -sS -L -H "Accept:application/octet-stream" https://api.github.com/repos/ion-channel/ion-connect/releases/assets/1413222\?access_token\=$TOKEN -o $TMPDIR/ion-connect-$RELEASE.tar.gz

echo "Untaring $TMPDIR/ion-connect-$RELEASE.tar.gz"
tar xfvz $TMPDIR/ion-connect-$RELEASE.tar.gz -C $TMPDIR

echo "Copying linux binaries"
cp $TMPDIR/ion-connect/linux/bin/ion-connect /usr/local/bin
