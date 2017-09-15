Feature: Projects
  Scenario: A users creates a project
    Given I have a ruleset id
    When I run the command to create a project
    Then I see a response showing the project is created

  Scenario: Get project
    Given I have a project
    When I run the command to get the project
    Then I see the projects details
