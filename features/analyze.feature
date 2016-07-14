Feature: Projects
  Scenario: Create a ruleset
    Given an Ion Channel account id 'test-account'
    When I successfully run with 'account_id' `ion-connect ruleset create-ruleset --account-id account_id test-ruleset "this is a test ruleset" '["c30b917956c3040daa2c571ef31dbe3a"]'`
    Then the ion output should contain:
    """
    rules
    """
    Then the ion output should contain:
    """
    "name": "test-ruleset"
    """
    Then the ion output should contain:
    """
    "description": "The project source is required to include a valid .about.yml file.",
    """

  Scenario: Create a project
    Given previous output
    And a variable 'ruleset_id' is set from the previous output from location 'id'
    And an Ion Channel account id 'test-account'
    When I successfully run with 'account_id,ruleset_id' `ion-connect project create-project --account-id account_id --ruleset-id ruleset_id --active sonar-auth-geoaxis "https://gitlab.devops.geointservices.io/DevOps/sonar-auth-geoaxis.git" "Sonar Plugin for auth with geoaxis"`
    Then the ion output should contain:
    """
    "active": true
    """
    Then the ion output should contain:
    """
    "branch": "master"
    """
    Then the ion output should contain:
    """
    "source": "https://gitlab.devops.geointservices.io/DevOps/sonar-auth-geoaxis.git"
    """

  Scenario: Analyze the project
    Given previous output
    And a variable 'project_id' is set from the previous output from location 'id'
    And an Ion Channel account id 'test-account'
    When I successfully run with 'account_id,project_id' `./test/analyze.sh project_id account_id`
    Then the ion output should contain:
    """
    Finished about_yml scan for sonar-auth-geoaxis, valid .about.yml found.
    """
    Then the ion output should contain:
    """
    Compliance analysis completed successfully, your project is compliant!
    """
