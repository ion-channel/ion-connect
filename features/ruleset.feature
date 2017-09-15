Feature: Rulesets
  Scenario: Get rules
    When I successfully run with 'team_id' `ion-connect ruleset get-rules`
    Then the ion output should contain:
    """
    "category": "Code Coverage"
    """
    Then the ion output should contain:
    """
    "category": "About Dot Yaml"
    """

  Scenario: A user creates a ruleset
    Given I have a set of rules
    When I run the command to create a ruleset
    Then I see a response showing the ruleset is created

  Scenario: Get a ruleset
    Given previous output
    And a variable 'id' is set from the previous output from location 'id'
    And an Ion Channel team id 'test-team'
    When I successfully run with 'team_id,id' `ion-connect ruleset get-ruleset --team-id team_id id`
    Then the ion output should contain:
    """
    "description": "The project source is required to include a valid .about.yml file.",
    """

  Scenario: Get all rule sets for team
    Given an Ion Channel team id 'test-team'
    When I successfully run with 'team_id' `ion-connect ruleset get-rulesets --team-id team_id`
    Then the ion output should contain:
    """
    "description": "The project source is required to include a valid .about.yml file.",
    """
