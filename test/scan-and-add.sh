PROJECT_ID=c11d5730e1725dcea5b591c27035e70b

#RUN Analysis
ANALYSIS_ID=$(IONCHANNEL_ENDPOINT_URL=https://api.test.ionchannel.io/ ion-connect scanner analyze-project --account-id account_id --project-id $PROJECT_ID $BUILD_NUMBER | jq -r .id)

#Get the Analysis
TIMEOUT=120
COUNTER=10
STATUS=$(IONCHANNEL_ENDPOINT_URL=https://api.test.ionchannel.io/ ion-connect scanner get-analysis-status --account-id account_id --project-id $PROJECT_ID $ANALYSIS_ID | jq -r .status)
while [[ $STATUS = "accepted" ]]; do
  COUNTER=$((COUNTER+10))
  if [[ $COUNTER -lt $TIMEOUT ]]; then
    sleep 5
    STATUS=$(IONCHANNEL_ENDPOINT_URL=https://api.test.ionchannel.io/ ion-connect scanner get-analysis-status --account-id account_id --project-id $PROJECT_ID $ANALYSIS_ID | jq -r .status)
  else
    echo "ERROR: ion-connect has timed out waiting for analysis to finish"
    exit 1
  fi
done
echo "Analysis has finished"

#Parse the coverage value
URL="https://sonar.geointservices.io/api/resources?resource=ion:mage&metrics=coverage"
VALUE=$(curl -s -u "$sonarqube_ro_token:" "$URL" | jq -r '.[].msr[].val')
echo "Retreived latest code coverage result value ($VALUE)"


#Add the value to the analysis
IONCHANNEL_ENDPOINT_URL=https://api.test.ionchannel.io/ ion-connect --debug scanner add-scan-result --account-id account_id --project-id $PROJECT_ID --analysis-id $ANALYSIS_ID finished "{\"value\":$VALUE}" coverage
echo "Added code coverage result to analyis"

#Get the new scan id
SCAN=$(IONCHANNEL_ENDPOINT_URL=https://api.test.ionchannel.io/ ion-connect scanner get-analysis-status --account-id account_id --project-id $PROJECT_ID $ANALYSIS_ID | jq -r '.scan_status[] | select(.name | contains("coverage")) | .id')
COUNTER=1
while [[ -z "$SCAN_ID" ]]; do
  COUNTER=$((COUNTER+1))
  if [[ $COUNTER -lt $TIMEOUT ]]; then
    sleep 1
    SCAN_ID=$(IONCHANNEL_ENDPOINT_URL=https://api.test.ionchannel.io/ ion-connect scanner get-analysis-status --account-id account_id --project-id $PROJECT_ID $ANALYSIS_ID | jq -r '.scan_status[] | select(.name | contains("coverage")) | .id')
  else
    echo "ERROR: ion-connect has timed out waiting for scan to finish"
    exit 1
  fi
done
echo "Retreived external scan id from analysis"

#Get the results of the scan
SCAN_RESULT=$(IONCHANNEL_ENDPOINT_URL=https://api.test.ionchannel.io/ ion-connect animal get-scan --account-id account_id --project-id $PROJECT_ID $ANALYSIS_ID $SCAN_ID)
echo "Retreived external scan data from analysis"

#Apply the rules to the scan
RULESET=ruleset2
IONCHANNEL_ENDPOINT_URL=https://api.test.ionchannel.io/ ion-connect ruleset get-ruleset --account-id account_id $RULESET
SUMMARY=$(IONCHANNEL_ENDPOINT_URL=https://api.test.ionchannel.io/ ion-connect ruleset apply-ruleset --account-id account_id $RULESET "[$SCAN_RESULT]" | jq -r .summary)
echo "Applied ruleset to result of coverage scan"
echo $SUMMARY


RULESET_DATA=$(IONCHANNEL_ENDPOINT_URL=https://api.test.ionchannel.io/ ion-connect ruleset get-ruleset --account-id account_id $RULESET)

COMPARATOR=$(echo $RULESET_DATA | jq -r ".rules[].comparator")
ALLOWED_VALUE=$(echo $RULESET_DATA | jq -r ".rules[].value")

echo "Applied rule set to external scan with result $VALUE $COMPARATOR $ALLOWED_VALUE ($SUMMARY)"

if [ "$SUMMARY" == "fail" ]; then
  exit 1
fi
