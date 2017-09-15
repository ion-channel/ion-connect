#!/bin/bash
# Copyright [2016] [Selection Pressure]
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PROJECT_ID=$1
TEAM_ID=$2

if [ -z "$3" ]; then
  BRANCH=''
else
  BRANCH=" --branch $3"
fi

echo $BRANCH

function cool_echo {
# local x=$1
# for (( i=0 ; i < ${#x} ; i++ )); do python3 -c "print('${x:$i:1}', end='')"; done
echo $1
echo
}

#RUN Analysis
cool_echo "Begining compliance analysis of project ($PROJECT_ID)."
ANALYSIS_ID=$(ion-connect scanner analyze-project --team-id $TEAM_ID --project-id $PROJECT_ID $BRANCH | jq -r .id)
echo "Analysis requested the id is $ANALYSIS_ID"
#Get the Analysis
TIMEOUT=1200
COUNTER=10
STATUS=$(ion-connect scanner get-analysis-status --team-id $TEAM_ID --project-id $PROJECT_ID $ANALYSIS_ID | jq -r .status)
while [[ $STATUS = "accepted" ]]; do
echo -n '.'
COUNTER=$((COUNTER+10))
if [[ $COUNTER -lt $TIMEOUT ]]; then
sleep 5
STATUS=$(ion-connect scanner get-analysis-status --team-id $TEAM_ID --project-id $PROJECT_ID $ANALYSIS_ID | jq -r .status)
else
cool_echo "ERROR: ion-connect has timed out waiting for analysis to finish"
exit 1
fi
done
echo

cool_echo "All project scans have finished."
cool_echo "Evaluating analysis for compliance."

#Get the results of the analysis
ANALYSIS=$(ion-connect report get-analysis --team-id $TEAM_ID --project-id $PROJECT_ID $ANALYSIS_ID)
PASSED="$(echo $ANALYSIS | jq -r .passed)"
ANALYZED_BRANCH="$(echo $ANALYSIS | jq -r .branch)"
OUTPUT=$(echo $ANALYSIS | jq -r .scan_summaries[].summary)
while read -r LINE; do
cool_echo "$LINE"
sleep 1s
done <<< "$OUTPUT"

if [ "$PASSED" == "false" ]; then
cool_echo "Compliance analysis failed, your project is not compliant :("
exit 1
fi
cool_echo "Compliance analysis completed successfully, your project at $ANALYZED_BRANCH is compliant!"
