#!/bin/bash
## process-ion-job.sh
## usage: process-ion-job.sh name hash url timeout(optional, default: 120s)
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

POSTRESULT=`ion-connect airgap push-artifact-url --checksum $2 --name $1 --url $3` 
ID=`echo $POSTRESULT | jq -r '.airgap_id'`

if [ "$ID" = "null" ]; then
  echo "ERROR: Failed to post to Ion"
  echo $POSTRESULT
  exit 1
fi

STATUS="started"
while [[ $STATUS != "finished" ]]; do
  COUNTER=1
  if [[ $COUNTER -lt $TIMEOUT ]]; then
    sleep 1
    RESULT=`ion-connect airgap get-push --airgapid $ID`
    STATUS=`echo $RESULT | jq -r '.scan_status'`
  else
    exit 1
  fi
  COUNTER=COUNTER+1
done
echo "$RESULT"
