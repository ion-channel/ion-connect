Feature: Rulesets
  Scenario: Create a ruleset
    Given an Ion Channel account id 'test-account'
    When I successfully run with 'account_id' `ion-connect ruleset create-ruleset --account-id account_id test-ruleset "this is a test ruleset" '["0239f0f8c5223fc47f32ebdf6636f4f0","c30b917956c3040daa2c571ef31dbe3a"]'`
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

  Scenario: Get a ruleset
    Given previous output
    And a variable 'id' is set from the previous output from location 'id'
    And an Ion Channel account id 'test-account'
    When I successfully run with 'account_id,id' `ion-connect ruleset get-ruleset --account-id account_id id`
    Then the ion output should contain:
    """
    "description": "The project source is required to include a valid .about.yml file.",
    """

  Scenario: Get all rule sets for account
    Given an Ion Channel account id 'test-account'
    When I successfully run with 'account_id' `ion-connect ruleset get-rulesets --account-id account_id`
    Then the ion output should contain:
    """
    "description": "The project source is required to include a valid .about.yml file.",
    """

  
