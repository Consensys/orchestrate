@jwt
@multi-tenancy
Feature: JWT Authentication
  As as external developer
  I want to perform authenticate using a jwt token from an identity provider

  Background:
    Given I have the following jwt tokens
      | alias     | audience                          |
      | tenantFoo | https://orchestrate.consensys.net |

  Scenario: Create resource with jwt token
    Given I set the headers
      | Key           | Value               |
      | Authorization | {{tenantFoo.token}} |
    When I send "POST" request to "{{global.api}}/schedules" with json:
    """
    {}
    """
    Then the response code should be 200
    And Response should have the following fields
      | tenantID                                 |
      | DkwGi2hAirRT32rogzf8ntpOUCoJRulL@clients |
