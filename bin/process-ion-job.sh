#!/bin/bash
## process-ion-job.sh
## usage: process-ion-job.sh name hash url
##
## Copyright (C) 2015 Selection Pressure LLC
##
## This software may be modified and distributed under the terms
## of the MIT license.  See the LICENSE file for details.

function error {
  echo "Error: process-ion-job.sh name hash url"
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

airgap_id=`ion-connect airgap push-artifact-url --checksum $2 --name $1 --url $3 | jq -r '.airgap_id'`

STATUS="started"
while [[ $STATUS != "finished" ]]; do
  sleep 2
  RESULT=`ion-connect airgap get-push --airgapid $airgap_id`
  STATUS=`echo $RESULT | jq -r '.scan_status'`
done
echo "$RESULT"
