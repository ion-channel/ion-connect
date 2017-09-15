Feature: Analyze
  As a user
  I want to perform actions on projects
  So that I can manage my projects

  Scenario: A user analyzes a project
    Given I have a project
    When I run the command to analyze a project
    Then I see a response showing the project is analyzed
