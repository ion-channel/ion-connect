Feature: Stuff related to repositories
  Scenario: Get repository data from a dependency
     When I successfully run `ion-connect repository get-repository-by-dependency com.intel.jndn.utils:jndn-utils maven`
     Then the output should contain:
     """
     {
       "api_created_at": "2016-05-09T22:23:30.913Z",
       "api_updated_at": "2016-05-09T22:23:30.913Z",
       "description": "A collection of tools to simplify synchronous and asynchronous data transfer over the NDN network",
       "fork": false,
       "full_name": "01org/jndn-utils",
       "id": "31076872",
       "name": "jndn-utils",
       "owner_id": "1635439",
       "owner_login": "01org",
       "owner_site_admin": false,
       "owner_type": "Organization"
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
