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

Feature: Dependency
  Scenario: Get a list of resolved dependencies for a project
    When I successfully run `ion-connect dependency get-resolved-dependencies mocha npmjs`
    Then the output should contain:
    """
    {
      "mocha@2.4.5": {
        "commander@2.3.0": {},
        "debug@2.2.0": {
          "ms@0.7.1": {}
        },
        "diff@1.4.0": {},
        "escape-string-regexp@1.0.2": {},
        "glob@3.2.3": {
          "graceful-fs@2.0.3": {},
          "inherits@2.0.1": {},
          "minimatch@0.2.14": {
            "lru-cache@2.7.3": {},
            "sigmund@1.0.1": {}
          }
        },
        "growl@1.8.1": {},
        "jade@0.26.3": {
          "commander@0.6.1": {},
          "mkdirp@0.3.0": {}
        },
        "mkdirp@0.5.1": {
          "minimist@0.0.8": {}
        },
        "supports-color@1.2.0": {}
      }
    }
    """
  Scenario: Get a list of dependencies for a project
    When I successfully run `ion-connect dependency get-dependencies fart rubygems`
    Then the output should contain:
    """
    {
      "_uniqueKey": "[\"fart\",\"0.0.2\"]",
      "authors": "Manish Das",
      "built_at": "2012-04-13T00:00:00Z",
      "checksum": "12hk/pgNufuuwq45tgr2lXYvALb7959N7LMfPpAlrIU=",
      "created_at": "2012-04-13T07:37:32Z",
      "dependencies": [
        [
          "thor",
          "~> 0.14"
        ]
      ],
      "description": "Simple usage of command line interface",
      "downloads_count": 1928,
      "latest": true,
      "metadata": "{}",
      "name": "fart",
      "number": "0.0.2",
      "platform": "ruby",
      "prerelease": false,
      "sha": "d76864fe980db9fbaec2ae39b60af695762f00b6fbf79f4decb31f3e9025ac85",
      "summary": "Fart Noise"
    }
    """
  Scenario: Get the latest dependency version
    When I successfully run `ion-connect dependency get-latest-version-for-dependency rails rubygems`
    Then the output should contain:
    """
    {
      "version": "4.2.6"
    }
    """
  Scenario: Get the latest dependency version
    When I successfully run `ion-connect dependency get-latest-versions-for-dependency fart rubygems`
    Then the output should contain:
    """
    [
      "0.0.1",
      "0.0.2"
    ]
    """
  Scenario: Get dependencies from a file
    When I successfully run `ion-connect dependency resolve-dependencies-in-file --type npmjs --flatten "../../test/package.json"`
    Then the output should contain:
    """
    "bson"
    """
