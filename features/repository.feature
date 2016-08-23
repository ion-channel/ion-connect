Feature: Stuff related to repositories
  Scenario: Get repository data from a dependency
     When I successfully run `ion-connect repository get-repository-by-dependency com.intel.jndn.utils:jndn-utils maven`
     Then the output should contain:
     """
     {
       "api_created_at": "2016-05-09T22:23:30.913Z",
       "api_updated_at": "2016-05-09T22:23:30.913Z",
       "created_at": "2015-02-20T17:49:42Z",
       "default_branch": "master",
       "description": "A collection of tools to simplify synchronous and asynchronous data transfer over the NDN network",
       "fork": false,
       "forks_count": 4,
       "full_name": "01org/jndn-utils",
       "has_downloads": true,
       "has_issues": true,
       "has_pages": true,
       "has_wiki": true,
       "homepage": "",
       "id": "31076872",
       "language": "Java",
       "name": "jndn-utils",
       "network_count": 4,
       "open_issues_count": 0,
       "owner_id": "1635439",
       "owner_login": "01org",
       "owner_site_admin": false,
       "owner_type": "Organization",
       "pushed_at": "2015-09-04T22:33:30Z",
       "size": 537,
       "stargazers_count": 0,
       "subscribers_count": 4,
       "updated_at": "2015-05-07T18:06:08Z",
       "watchers_count": 0
     }
     """
  Scenario: Get repository data from a dependency
    When I successfully run `ion-connect repository get-repository ruby/ruby`
    Then the output should contain:
    """
    {
      "api_created_at": "2016-05-09T22:07:10.563Z",
      "api_updated_at": "2016-05-09T22:07:10.563Z",
      "description": "The Ruby Programming Language",
      "fork": false,
      "full_name": "ruby/ruby",
      "id": "538746",
      "name": "ruby",
      "owner_id": "210414",
      "owner_login": "ruby",
      "owner_site_admin": false,
      "owner_type": "Organization"
    }
    """
