#! /bin/bash

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
OUTPUT=$(ion-connect metadata get-locations "portland")
if [ "$?" != "0" ]; then
  echo "Failed"
  exit 1
fi

echo "It should get vulnerabilities for a project"
OUTPUT=$(ion-connect vulnerabilities get-vulnerabilities --limit 1 --offset 0 "solr")
if [ "$?" != "0" ]; then
  echo "Failed"
  exit 1
fi

echo "It should get vulnerabilities for text"
OUTPUT=$(ion-connect vulnerabilities get-vulnerabilities --text --limit 12 --offset 10 "testing")
if [ "$?" != "0" ]; then
  echo "Failed"
  exit 1
fi

echo "It should get sentiment for some text"
OUTPUT=$(ion-connect metadata get-sentiment "I love Ion Channel")
if [ "$?" != "0" ]; then
  echo "Failed"
  exit 1
fi

echo "It should handle not finding data"
OUTPUT=$(ion-connect scanner get-scan notreallyascan)
if [ "$?" != "1" ]; then
  echo "Failed"
  exit 1
fi
if [ "$OUTPUT" != "Item with id (notreallyascan) not found" ]; then
  echo "Failed"
  exit 1
fi

echo "It should scan and push an artifact"
OUTPUT=$(./bin/process-ion-job.sh ion-channel/ion-connect c88d59be7c087bd379af09953693ecf7 "https://github.com/ion-channel/agmockapp/archive/agmockapp-0.0.1.tar.gz")
if [ "$?" != "0" ]; then
  echo "Failed"
  exit 1
fi

STATUS=$(echo $OUTPUT | jq -r .scanner.status)
if [ "$STATUS" != "finished" ]; then
  echo "Failed"
  exit 1
fi

STATUS=$(echo $OUTPUT | jq -r .airgap.status)
if [ "$STATUS" != "finished" ]; then
  echo "Failed"
  exit 1
fi

echo "It should scan for dependencies"
OUTPUT=$(./bin/dependency-scan-job.sh ./test/Gemfile gemfile)
if [ "$?" != "0" ]; then
  echo "Failed"
  exit 1
fi

STATUS=$(echo $OUTPUT | jq -r .status)
if [ "$STATUS" != "finished" ]; then
  echo "Failed"
  exit 1
fi

echo "function test suite completed"
