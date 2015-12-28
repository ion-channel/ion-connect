#!/bin/bash
## process-ion-job.sh
## usage: process-ion-job.sh name hash url timeout(optional, default: 120s)
## desc: This script puts together a basic workflow of scan, then transfer
## and artifact of interest.  The scan runs, and if 'finished' with no negative
## scan results, it will proceed and transfer (pushing the artifact in/up).
##
## Copyright (C) 2015 Selection Pressure LLC
##
## This software may be modified and distributed under the terms
## of the MIT license.  See the LICENSE file for details.

function error {
  echo "Error: process-ion-job.sh name hash url timeout(optional, default: 120s)"
}

if ! [ -x "$(command -v ion-connect)" ]; then
  echo 'process-ion-job.sh: ion-connect is not installed.' >&2
  exit
fi

if ! [ -x "$(command -v jq)" ]; then
  echo 'process-ion-job.sh: jq is not installed.' >&2
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

if [ -z "$3" ]; then
  error
  exit
fi

if [ -z $4 ]; then
  TIMEOUT=120
else
  TIMEOUT=$4
fi

SCANRESULT=`ion-connect scanner scan-artifact-url --checksum $2 --project $1 --url $3`
SCANSTATUS=`echo $SCANRESULT | jq -r '.status'`

if [ "$SCANSTATUS" = "accepted" ]; then
  SCANID=`echo $SCANRESULT | jq -r '.id'`
else
  echo "ERROR: Failed to scan in Ion"
  echo $SCANRESULT
  exit 1
fi

#echo "SS: $SCANSTATUS - SR: $SCANRESULT"
## TODO: What is the status if it finds something, or what is the indicator
## that we need to fail this loop?

while [[ $SCANSTATUS != "finished" ]]; do
  COUNTER=1
  if [[ $COUNTER -lt $TIMEOUT ]]; then
    sleep 1
    GETSCANRESULT=`ion-connect scanner get-scan --id $SCANID`
    SCANSTATUS=`echo $GETSCANRESULT | jq -r '.status'`
  else
    echo "ERROR: ion-connect has timed out"
    exit 1
  fi
done

## We have completed the scan portion, now push the artifact

PUSHRESULT=`ion-connect airgap push-artifact-url --checksum $2 --project $1 --url $3`
#echo "PR: $PUSHRESULT"
PUSHSTATUS=`echo $PUSHRESULT | jq -r '.status'`
#echo "PS: $PUSHSTATUS"
if [ "$PUSHSTATUS" = "accepted" ]; then
  PUSHID=`echo $PUSHRESULT | jq -r '.id'`
else
  echo "ERROR: Failed to post to Ion"
  echo $PUSHRESULT
  exit 1
fi

while [[ $PUSHSTATUS != "finished" ]]; do
  COUNTER=1
  if [[ $COUNTER -lt $TIMEOUT ]]; then
    sleep 1
    GETPUSHRESULT=`ion-connect airgap get-push --id $PUSHID`
    PUSHSTATUS=`echo $GETPUSHRESULT | jq -r '.status'`
  else
    echo "ERROR: ion-connect has timed out"
    exit 1
  fi
  COUNTER=COUNTER+1
done
printf "{\"scanner\":$GETSCANRESULT\n,\"airgap\":$GETPUSHRESULT\n}"
