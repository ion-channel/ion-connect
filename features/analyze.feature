# Copyright [2016] [Selection Pressure]
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

Feature: Projects
  Scenario: Create a ruleset
    Given an Ion Channel team id 'test-team'
    When I successfully run with 'team_id' `ion-connect ruleset create-ruleset --team-id team_id test-ruleset "this is a test ruleset" '["c30b917956c3040daa2c571ef31dbe3a"]'`
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
