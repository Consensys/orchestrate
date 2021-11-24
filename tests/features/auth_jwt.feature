@jwt
@multi-tenancy
Feature: JWT Authentication
  As as external developer
  I want to perform authenticate using a jwt token from an identity provider

  Background:
    Given I have the following jwt tokens
      | alias            | audience                                        |
      | tenantOne        | https://orchestrate.consensys.net/tenant1       |
      | tenantOneUserOne | https://orchestrate.consensys.net/tenant1:user1 |

  Scenario: Create resource with claimed tenant
    Given I set the headers
      | Key           | Value               |
      | Authorization | {{tenantOne.token}} |
    When I send "POST" request to "{{global.api}}/schedules" with json:
    """
    {}
    """
    Then the response code should be 200
    And Response should have the following fields
      | tenantID |
      | tenant1  |

  Scenario: Create resource with claimed tenant and username
    Given I set the headers
      | Key           | Value                      |
      | Authorization | {{tenantOneUserOne.token}} |
    When I send "POST" request to "{{global.api}}/schedules" with json:
    """
    {}
    """
    Then the response code should be 200
    And Response should have the following fields
      | tenantID | ownerID |
      | tenant1  | user1   |

  Scenario: Create resource impersonating default tenant
    Given I set the headers
      | Key           | Value                      |
      | Authorization | {{tenantOneUserOne.token}} |
      | X-TENANT-ID   | _                          |
    When I send "POST" request to "{{global.api}}/schedules" with json:
    """
    {}
    """
    Then the response code should be 200
    And Response should have the following fields
      | tenantID |
      | _        |

  Scenario: Fail to impersonate no default tenant
    Given I set the headers
      | Key           | Value                      |
      | Authorization | {{tenantOneUserOne.token}} |
      | X-TENANT-ID   | bar                        |
    When I send "POST" request to "{{global.api}}/schedules" with json:
    """
    {}
    """
    Then the response code should be 401

  Scenario: Fail to impersonate user
    Given I set the headers
      | Key           | Value               |
      | Authorization | {{tenantOne.token}} |
      | X-USERNAME    | user                 |
    When I send "POST" request to "{{global.api}}/schedules" with json:
    """
    {}
    """
    Then the response code should be 401
