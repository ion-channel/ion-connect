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

Feature: Get help
 Scenario: Check version
   When I successfully run `ion-connect --version`
   Then the output should contain:
   """
ion-connect version
   """

 Scenario: Get Help
   When I successfully run `ion-connect --help`
   Then the output should contain:
   """
ion-connect [global options] command [command options] [arguments...]
   """
