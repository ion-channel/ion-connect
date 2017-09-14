Feature: Projects
  As a user
  I want to perform actions on projects
  So that I can manage my projects

  Scenario: A user creates a ruleset
    Given I have a set of rules
    When I run the command to create a ruleset
    Then I see a response showing the ruleset is created

  Scenario: A users creates a project
    Given I have a ruleset id
    When I run the command to create a project
    Then I see a response showing the project is created

  Scenario: Analyze the project
    Given previous output
    And a variable 'project_id' is set from the previous output from location 'id'
    And an Ion Channel team id 'test-team'
    When I successfully run with 'team_id,project_id' `./test/analyze.sh project_id team_id`
    Then the ion output should contain:
    """
    Finished about_yml scan for java-lew, valid .about.yml found.
    """
    Then the ion output should contain:
    """
    Compliance analysis completed successfully, your project at master is compliant!
    """

  Scenario: Analyze the project with branch/hash
    Given previous output
    And a variable 'project_id' is set from the previous output from location 'id'
    And an Ion Channel team id 'test-team'
    And a branch named a39b99095ddb9d6dd299f13cbcf9dd17fd5fb8c3
    When I successfully run with 'team_id,project_id,branch' `./test/analyze.sh project_id team_id branch`
    Then the ion output should contain:
    """
    a39b99095ddb9d6dd299f13cbcf9dd17fd5fb8c3
    """
    Then the ion output should contain:
    """
    Finished about_yml scan for java-lew, valid .about.yml found.
    """
    Then the ion output should contain:
    """
    Compliance analysis completed successfully, your project at a39b99095ddb9d6dd299f13cbcf9dd17fd5fb8c3 is compliant!
    """
