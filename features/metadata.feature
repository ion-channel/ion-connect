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

Feature: Do some metadata parsing
 Scenario: Get locations
   When I successfully run `ion-connect metadata get-locations "portland"`
   Then the output should contain:
   """
   [
     {
       "boundingbox": [],
       "city": "",
       "class": "city, village, etc.",
       "country": "United States",
       "country_code": "US",
       "county": "",
       "display_name": "Portland",
       "id": "5746545",
       "importance": 1,
       "latitude": 45.52345,
       "longitude": -122.67621,
       "state": "",
       "type": ""
     }
   ]
   """

   Scenario: Get languages
     When I successfully run `ion-connect metadata get-languages "Hola como estas"`
     Then the output should contain:
     """
     [
       {
         "code": "es",
         "name": "Spanish",
         "reliable": true
       }
     ]
     """

   Scenario: Get sentiment
     When I successfully run `ion-connect metadata get-sentiment "I love cucumber"`
     Then the output should contain:
     """
     {
       "score": 0.925,
       "sentiment": "positive"
     }
     """

   # TODO: deploy fossology
   Scenario: Get licenses
     Given a file named "LICENSE.txt" with:
     """
     The MIT License (MIT)

     Copyright (c) 2015 Selection Pressure LLC

     Permission is hereby granted, free of charge, to any person obtaining a copy
     of this software and associated documentation files (the "Software"), to deal
     in the Software without restriction, including without limitation the rights
     to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
     copies of the Software, and to permit persons to whom the Software is
     furnished to do so, subject to the following conditions:

     The above copyright notice and this permission notice shall be included in all
     copies or substantial portions of the Software.

     THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
     IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
     FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
     AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
     LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
     OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
     SOFTWARE.
     """
     When I successfully run ion `ion-connect metadata get-licenses "$(cat LICENSE.txt)"`
     Then the ion output should contain:
     """
     [
       {
         "name": "mit"
       }
     ]
     """
