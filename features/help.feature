Feature: Get help
 Scenario: Check version
   When I successfully run `ion-connect --version`
   Then the output should contain:
   """
ion-connect version 0.6.6
   """

 Scenario: Get Help
   When I successfully run `ion-connect --help`
   Then the output should contain:
   """
ion-connect [global options] command [command options] [arguments...]
   """
