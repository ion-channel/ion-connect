#! /bin/bash

#RUN Analysis
ANALYSIS_ID=$(IONCHANNEL_ENDPOINT_URL=https://api.test.ionchannel.io/ ion-connect scanner analyze-project --account-id account_id --project-id 044d931cd9056898a1fd755e34ab0cb6 8 | jq -r .id)

#Get the Analysis
STATUS=$(IONCHANNEL_ENDPOINT_URL=https://api.test.ionchannel.io/ ion-connect scanner get-analysis-status --account-id account_id --project-id 044d931cd9056898a1fd755e34ab0cb6 $ANALYSIS_ID | jq -r .status)

#Parse the coverage value
VALUE=$(cat test/sonarresponse.json | jq -r '.[].msr[].val')

#Add the value to the analysis
IONCHANNEL_ENDPOINT_URL=https://api.test.ionchannel.io/ ion-connect --debug scanner add-scan-result --account-id account_id --project-id 044d931cd9056898a1fd755e34ab0cb6 --analysis-id $ANALYSIS_ID finished "{\"value\":$VALUE}" coverage

#Get the new scan id
SCAN_ID=$(IONCHANNEL_ENDPOINT_URL=https://api.test.ionchannel.io/ ion-connect scanner get-analysis-status --account-id account_id --project-id 044d931cd9056898a1fd755e34ab0cb6 $ANALYSIS_ID | jq -r '.scan_status[] | select(.name | contains("coverage")) | .id')

#Get the results of the scan
SCAN_RESULT=$(IONCHANNEL_ENDPOINT_URL=https://api.test.ionchannel.io/ ion-connect animal get-scan --account-id account_id --project-id 044d931cd9056898a1fd755e34ab0cb6 $ANALYSIS_ID $SCAN_ID)

#Apply the rules to the scan
SUMMARY=$(IONCHANNEL_ENDPOINT_URL=https://api.test.ionchannel.io/ ion-connect ruleset apply-ruleset --account-id account_id ruleset1 "[$SCAN_RESULT]" | jq -r .summary)

#If it's 'fail' fail the build
# exit 1
