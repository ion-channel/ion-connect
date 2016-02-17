#!/bin/bash
## dependency-scan-job.sh
## usage: dependency-scan-job.sh name hash url timeout(optional, default: 120s)
## desc: This script puts together a basic workflow of scan, then transfer
## and artifact of interest.  The scan runs, and if 'finished' with no negative
## scan results, it will proceed and transfer (pushing the artifact in/up).
##
## Copyright (C) 2015 Selection Pressure LLC
##
## This software may be modified and distributed under the terms
## of the MIT license.  See the LICENSE file for details.

function error {
  echo "Error: dependency-scan-job.sh path type timeout(optional, default: 120s)"
}

if ! [ -x "$(command -v ion-connect)" ]; then
  echo 'dependency-scan-job.sh: ion-connect is not installed.' >&2
  exit
fi

if ! [ -x "$(command -v jq)" ]; then
  echo 'dependency-scan-job.sh: jq is not installed.' >&2
  exit
fi

if [ -z "$1" ]; then
  error
  exit
fi

if [ -z "$2" ]; then
  error
  exit
fi


if [ -z $4 ]; then
  TIMEOUT=120
else
  TIMEOUT=$4
fi

SCANRESULT=`ion-connect scanner scan-dependencies --type $2 $1`
if [ "$?" != "0" ]; then
  echo $SCANRESULT
  exit 1
fi

SCANSTATUS=`echo $SCANRESULT | jq -r '.status'`

if [ "$SCANSTATUS" = "accepted" ]; then
  SCANID=`echo $SCANRESULT | jq -r '.id'`
else
  echo "ERROR: Failed to scan in Ion"
  echo $SCANRESULT
  exit 1
fi

while [[ $SCANSTATUS = "accepted" ]]; do
  COUNTER=1
  if [[ $COUNTER -lt $TIMEOUT ]]; then
    sleep 1
    GETSCANRESULT=`ion-connect scanner get-dependencies $SCANID`
    SCANSTATUS=`echo $GETSCANRESULT | jq -r '.status'`
  else
    echo "ERROR: ion-connect has timed out"
    exit 1
  fi
done

printf "$GETSCANRESULT"
