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

Feature: Rulesets
  Scenario: Get rules
    When I successfully run with 'account_id' `ion-connect ruleset get-rules`
    Then the ion output should contain:
    """
    "category": "Code Coverage"
    """
    Then the ion output should contain:
    """
    "category": "About Dot Yaml"
    """
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
