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

  Scenario: Analyze the project
    Given previous output
    And a variable 'project_id' is set from the previous output from location 'id'
    And an Ion Channel team id 'test-team'
    When I successfully run with 'team_id,project_id' `./test/analyze.sh project_id team_id`
    Then the ion output should contain:
    """
    Finished about_yml scan for sonar-auth-geoaxis, valid .about.yml found.
    """
    Then the ion output should contain:
    """
    Compliance analysis completed successfully, your project at master is compliant!
    """

  Scenario: Analyze the project with branch/hash
    Given previous output
    And a variable 'project_id' is set from the previous output from location 'id'
    And an Ion Channel team id 'test-team'
    And a branch named b979b868ab320e0236b1e7c5f1530ae2401083ab
    When I successfully run with 'team_id,project_id,branch' `./test/analyze.sh project_id team_id branch`
    Then the ion output should contain:
    """
    b979b868ab320e0236b1e7c5f1530ae2401083ab
    """
    Then the ion output should contain:
    """
    Finished about_yml scan for sonar-auth-geoaxis, valid .about.yml found.
    """
    Then the ion output should contain:
    """
    Compliance analysis completed successfully, your project at b979b868ab320e0236b1e7c5f1530ae2401083ab is compliant!
    """
