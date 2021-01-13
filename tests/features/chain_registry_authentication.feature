@chain-registry
@multi-tenancy
Feature: Chain-Registry Authentication
  As as external developer
  I want to register new chains and protect them under expected permission rules

  @geth
  Scenario: Create chain using X-API-Key
    Given I set the headers
      | Key       | Value              |
      | X-API-Key | {{global.api-key}} |
    When I send "POST" request to "{{global.api}}/chains" with json:
      """
      {
      "name": "geth-{{scenarioID}}",
      "urls": {{global.nodes.geth.URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "chainUUID"
    When I send "GET" request to "{{global.api}}/chains/{{chainUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | uuid          | tenantID |
      | {{chainUUID}} | _        |
    When I send "PATCH" request to "{{global.api}}/chains/{{chainUUID}}" with json:
      """
      {
        "name": "geth-new-{{scenarioID}}"
      }
      """
    Then the response code should be 200
    When I send "DELETE" request to "{{global.api}}/chains/{{chainUUID}}"
    Then the response code should be 204

  @geth
  Scenario: Create chain using X-API-Key and X-Tenant-ID
    Given I set the headers
      | Key         | Value              |
      | X-API-Key   | {{global.api-key}} |
      | X-Tenant-ID | foo                |
    When I send "POST" request to "{{global.api}}/chains" with json:
      """
      {
      "name": "geth-{{scenarioID}}",
      "urls": {{global.nodes.geth.URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "chainUUID"
    When I send "GET" request to "{{global.api}}/chains/{{chainUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | uuid          | tenantID |
      | {{chainUUID}} | foo      |
    When I send "DELETE" request to "{{global.api}}/chains/{{chainUUID}}"
    Then the response code should be 204

  @geth
  Scenario: Create chain with Wildcard JWT token
    Given I have the following tenants
      | alias    | tenantID |
      | wildcard | *        |
    Given I set the headers
      | Key           | Value                     |
      | Authorization | Bearer {{wildcard.token}} |
    When I send "POST" request to "{{global.api}}/chains" with json:
      """
      {
      "name": "geth-{{scenarioID}}",
      "urls": {{global.nodes.geth.URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "chainUUID"
    When I send "GET" request to "{{global.api}}/chains/{{chainUUID}}"
    And Response should have the following fields
      | uuid          | tenantID |
      | {{chainUUID}} | _        |
    Then the response code should be 200
    When I send "PATCH" request to "{{global.api}}/chains/{{chainUUID}}" with json:
      """
      {
        "name": "geth-new-{{scenarioID}}"
      }
      """
    Then the response code should be 200
    When I send "DELETE" request to "{{global.api}}/chains/{{chainUUID}}"
    Then the response code should be 204

  @geth
  Scenario:  Create chain with Wildcard JWT and X-Tenant-ID
    Given I have the following tenants
      | alias    | tenantID |
      | wildcard | *        |
    Given I set the headers
      | Key           | Value                     |
      | Authorization | Bearer {{wildcard.token}} |
      | X-Tenant-ID   | foo                       |
    When I send "POST" request to "{{global.api}}/chains" with json:
      """
      {
      "name": "geth-{{scenarioID}}",
      "urls": {{global.nodes.geth.URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "chainUUID"
    When I send "GET" request to "{{global.api}}/chains/{{chainUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | uuid          | tenantID |
      | {{chainUUID}} | foo      |
    When I send "DELETE" request to "{{global.api}}/chains/{{chainUUID}}"
    Then the response code should be 204

  @geth
  Scenario:  Create chain with with tenant foo and valid permissions
    Given I have the following tenants
      | alias     | tenantID |
      | tenantFoo | foo      |
      | tenantBar | bar      |
      | wildcard  | *        |
    Given I set the headers
      | Key           | Value                      |
      | Authorization | Bearer {{tenantFoo.token}} |
    When I send "POST" request to "{{global.api}}/chains" with json:
      """
      {
      "name": "geth-{{scenarioID}}",
      "urls": {{global.nodes.geth.URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "chainUUID"
    When I send "GET" request to "{{global.api}}/chains/{{chainUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | uuid          | tenantID |
      | {{chainUUID}} | foo      |
    Given I set the headers
      | Key           | Value                      |
      | Authorization | Bearer {{tenantBar.token}} |
    When I send "GET" request to "{{global.api}}/chains/{{chainUUID}}"
    Then the response code should be 404
    Given I set the headers
      | Key           | Value                     |
      | Authorization | Bearer {{wildcard.token}} |
      | X-Tenant-ID   | foo                       |
    When I send "GET" request to "{{global.api}}/chains/{{chainUUID}}"
    Then the response code should be 200
    Given I set the headers
      | Key           | Value                      |
      | Authorization | Bearer {{tenantFoo.token}} |
    When I send "DELETE" request to "{{global.api}}/chains/{{chainUUID}}"
    Then the response code should be 204

  @geth
  Scenario:  Create chain with with default tenant valid permissions
    Given I have the following tenants
      | alias     | tenantID |
      | tenantFoo | foo      |
      | tenantBar | bar      |
      | wildcard  | *        |
    Given I set the headers
      | Key           | Value                     |
      | Authorization | Bearer {{wildcard.token}} |
      | X-Tenant-ID   | _                         |
    When I send "POST" request to "{{global.api}}/chains" with json:
      """
      {
      "name": "geth-{{scenarioID}}",
      "urls": {{global.nodes.geth.URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "chainUUID"
    When I send "GET" request to "{{global.api}}/chains/{{chainUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | uuid          | tenantID |
      | {{chainUUID}} | _        |
    Given I set the headers
      | Key           | Value                      |
      | Authorization | Bearer {{tenantBar.token}} |
    When I send "GET" request to "{{global.api}}/chains/{{chainUUID}}"
    Then the response code should be 200
    Given I set the headers
      | Key           | Value                      |
      | Authorization | Bearer {{tenantFoo.token}} |
    When I send "GET" request to "{{global.api}}/chains/{{chainUUID}}"
    Then the response code should be 200
    Given I set the headers
      | Key           | Value                     |
      | Authorization | Bearer {{wildcard.token}} |
    When I send "DELETE" request to "{{global.api}}/chains/{{chainUUID}}"
    Then the response code should be 204

  @geth
  Scenario: Invalid X-API-Key
    Given I set the headers
      | Key       | Value              |
      | X-API-Key | {{global.api-key}} |
    When I send "POST" request to "{{global.api}}/chains" with json:
      """
      {
      "name": "geth-{{scenarioID}}",
      "urls": {{global.nodes.geth.URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "chainUUID"
    Given I set the headers
      | Key       | Value       |
      | X-API-Key | unknown-key |
    When I send "GET" request to "{{global.api}}/chains/{{chainUUID}}"
    Then the response code should be 401
    When I send "PATCH" request to "{{global.api}}/chains/{{chainUUID}}" with json:
      """
      {
        "name": "geth-new-{{scenarioID}}"
      }
      """
    Then the response code should be 401
    When I send "DELETE" request to "{{global.api}}/chains/{{chainUUID}}"
    Then the response code should be 401
    When I send "DELETE" request to "{{global.api}}/chains/{{chainUUID}}"
    Then the response code should be 401
    Given I set the headers
      | Key       | Value              |
      | X-API-Key | {{global.api-key}} |
    When I send "DELETE" request to "{{global.api}}/chains/{{chainUUID}}"
    Then the response code should be 204

  @geth
  Scenario: Invalid JWT token
    Given I have the following tenants
      | alias   | tenantID           |
      | tenant1 | foo-{{scenarioID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.api}}/chains" with json:
      """
      {
      "name": "geth-{{scenarioID}}",
      "urls": {{global.nodes.geth.URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "chainUUID"
    Given I set the headers
      | Key           | Value        |
      | Authorization | InvalidToken |
    When I send "GET" request to "{{global.api}}/chains/{{chainUUID}}"
    Then the response code should be 401
    When I send "DELETE" request to "{{global.api}}/chains/{{chainUUID}}"
    Then the response code should be 401
    When I send "DELETE" request to "{{global.api}}/chains/{{chainUUID}}"
    Then the response code should be 401
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "DELETE" request to "{{global.api}}/chains/{{chainUUID}}"
    Then the response code should be 204

