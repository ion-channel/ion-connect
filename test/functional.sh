#! /bin/bash

export IONCHANNEL_ENDPOINT_URL=https://api.test.ionchannel.io/

echo "Using $IONCHANNEL_ENDPOINT_URL for testing"

if ! type "jq" &> /dev/null; then
  echo "Missing required command 'jq' please install before continuing"
  exit 2
fi

echo "Running function test suite"

echo "It should install"
OUTPUT=$(go install ./...)
if [ "$?" != "0" ]; then
  echo "Failed"
  exit 1
fi

echo "It should provide help"
OUTPUT=$(ion-connect help)
if [ "$?" != "0" ]; then
  echo "Failed"
  exit 1
fi

echo "It should geocode locations"
OUTPUT=$(ion-connect --insecure metadata get-locations "portland")
if [ "$?" != "0" ]; then
  echo "Failed - $OUTPUT"
  exit 1
fi

echo "It should get vulnerabilities for a project"
OUTPUT=$(ion-connect --insecure vulnerability get-vulnerabilities --limit 1 --offset 0 "solr")
if [ "$?" != "0" ]; then
  echo "Failed - $OUTPUT"
  exit 1
fi

echo "It should get vulnerabilities for text"
OUTPUT=$(ion-connect --insecure vulnerability get-vulnerabilities --text --limit 12 --offset 10 "testing")
if [ "$?" != "0" ]; then
  echo "Failed - $OUTPUT"
  exit 1
fi

echo "It should get sentiment for some text"
OUTPUT=$(ion-connect --insecure metadata get-sentiment "I love Ion Channel")
if [ "$?" != "0" ]; then
  echo "Failed - $OUTPUT"
  exit 1
fi

echo "It should handle not finding data"
OUTPUT=$(ion-connect --insecure scanner get-scan notreallyascan)
if [ "$?" != "1" ]; then
  echo "Failed - $OUTPUT"
  exit 1
fi
if [ "$OUTPUT" != "Item with id (notreallyascan) not found" ]; then
  echo "Failed - $OUTPUT"
  exit 1
fi

echo "It should scan and push an artifact"
OUTPUT=$(./bin/process-ion-job.sh ion-channel/ion-connect c88d59be7c087bd379af09953693ecf7 "https://github.com/ion-channel/agmockapp/archive/agmockapp-0.0.1.tar.gz")
if [ "$?" != "0" ]; then
  echo "Failed - $OUTPUT"
  exit 1
fi

echo "It should scan and push a local artifact"
OUTPUT=$(./bin/process-ion-job.sh ion-channel/ion-connect 02754f7539c7db341000387fd6437f9377931a37 "file://./LICENSE.txt")
if [ "$?" != "0" ]; then
  echo "Failed - $OUTPUT"
  exit 1
fi

STATUS=$(echo $OUTPUT | jq -r .scanner.status)
if [ "$STATUS" != "finished" ]; then
  echo "Failed - $STATUS - $OUTPUT"
  exit 1
fi

STATUS=$(echo $OUTPUT | jq -r .airgap.status)
if [ "$STATUS" != "finished" ]; then
  echo "Failed - $STATUS - $OUTPUT"
  exit 1
fi

echo "It should scan for dependencies"
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
OUTPUT=$(ion-connect dependency resolve-dependencies-in-file --flatten $DIR/package.json)
if [ "$?" != "0" ]; then
  echo "Failed - $OUTPUT"
  exit 1
fi

echo "function test suite completed"
