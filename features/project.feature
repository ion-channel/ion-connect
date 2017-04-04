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
    When I successfully run with 'team_id' `ion-connect ruleset create-ruleset --team-id team_id test-ruleset "this is a test ruleset" '["0239f0f8c5223fc47f32ebdf6636f4f0","c30b917956c3040daa2c571ef31dbe3a"]'`
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
    And an Ion Channel team id 'test-team'
    When I successfully run with 'team_id,ruleset_id' `ion-connect project create-project --team-id team_id --ruleset-id ruleset_id --active sonar-auth-geoaxis "https://gitlab.devops.geointservices.io/DevOps/sonar-auth-geoaxis.git" "Sonar Plugin for auth with geoaxis"`
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

  Scenario: Get project
    Given previous output
    And a variable 'id' is set from the previous output from location 'id'
    And an Ion Channel team id 'test-team'
    When I successfully run with 'team_id,id' `ion-connect project get-project --team-id team_id id`
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
