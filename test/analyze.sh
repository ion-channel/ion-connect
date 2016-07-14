#!/bin/bash

PROJECT_ID=$1
ACCOUNT_ID=$2


function cool_echo {
# local x=$1
# for (( i=0 ; i < ${#x} ; i++ )); do python3 -c "print('${x:$i:1}', end='')"; done
echo $1
echo
}

#RUN Analysis
cool_echo "Begining compliance analysis of project ($PROJECT_ID)."
ANALYSIS_ID=$(ion-connect scanner analyze-project --account-id $ACCOUNT_ID --project-id $PROJECT_ID 1 | jq -r .id)
echo "Analysis requested the id is $ANALYSIS_ID"
#Get the Analysis
TIMEOUT=1200
COUNTER=10
STATUS=$(ion-connect scanner get-analysis-status --account-id $ACCOUNT_ID --project-id $PROJECT_ID $ANALYSIS_ID | jq -r .status)
while [[ $STATUS = "accepted" ]]; do
echo -n '.'
COUNTER=$((COUNTER+10))
if [[ $COUNTER -lt $TIMEOUT ]]; then
sleep 5
STATUS=$(ion-connect scanner get-analysis-status --account-id $ACCOUNT_ID --project-id $PROJECT_ID $ANALYSIS_ID | jq -r .status)
else
cool_echo "ERROR: ion-connect has timed out waiting for analysis to finish"
exit 1
fi
done
echo

cool_echo "All project scans have finished."
cool_echo "Evaluating analysis for compliance."

#Get the results of the analysis
ANALYSIS=$(ion-connect analysis get-analysis --account-id $ACCOUNT_ID --project-id $PROJECT_ID $ANALYSIS_ID)
PASSED="$(echo $ANALYSIS | jq -r .passed)"
OUTPUT=$(echo $ANALYSIS | jq -r .scan_summaries[].summary)
while read -r LINE; do
cool_echo "$LINE"
sleep 1s
done <<< "$OUTPUT"

if [ "$PASSED" == "false" ]; then
cool_echo "Compliance analysis failed, your project is not compliant :("
exit 1
fi
cool_echo "Compliance analysis completed successfully, your project is compliant!"
