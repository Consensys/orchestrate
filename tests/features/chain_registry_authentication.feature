@multi-tenancy
Feature: Chain-Registry Authentication
  Scenario: Valid X-API-Key and X-Tenant-ID unset
    Given I set the headers
        | Key       | Value    |
        | X-API-Key | with-key |
    When I send "POST" request to "{{global.chain-registry}}/chains" with json:
      """
      {
        "name": "geth-{{scenarioID}}",
        "urls": {{global.nodes.geth.URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "chainUID"
    When I send "GET" request to "{{global.chain-registry}}/chains/{{chainUID}}"
    Then the response code should be 200
    When I send "PATCH" request to "{{global.chain-registry}}/chains/{{chainUID}}" with json:
      """
      {
        "name": "geth-new-{{scenarioID}}"
      }
      """
    Then the response code should be 200
    When I send "DELETE" request to "{{global.chain-registry}}/chains/{{chainUID}}"
    Then the response code should be 204
    
  Scenario: Invalid X-API-Key
    Given I set the headers
        | Key       | Value    |
        | X-API-Key | with-key |
    When I send "POST" request to "{{global.chain-registry}}/chains" with json:
      """
      {
        "name": "geth-{{scenarioID}}",
        "urls": {{global.nodes.geth.URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "chainUID"
    Given I set the headers
      | Key       | Value       |
      | X-API-Key | unknown-key |
    When I send "GET" request to "{{global.chain-registry}}/chains/{{chainUID}}"
    Then the response code should be 401
    When I send "PATCH" request to "{{global.chain-registry}}/chains/{{chainUID}}" with json:
      """
      {
        "name": "geth-new-{{scenarioID}}"
      }
      """
    Then the response code should be 401
    When I send "DELETE" request to "{{global.chain-registry}}/chains/{{chainUID}}"
    Then the response code should be 401
    When I send "DELETE" request to "{{global.chain-registry}}/chains/{{chainUID}}"
    Then the response code should be 401
    Given I set the headers
      | Key         |   Value     |
      | X-API-Key   | with-key    |
    When I send "DELETE" request to "{{global.chain-registry}}/chains/{{chainUID}}"
    Then the response code should be 204

  Scenario: Valid X-API-Key and X-Tenant-ID
    Given I set the headers
      | Key         | Value    |
      | X-API-Key   | with-key |
      | X-Tenant-ID |    foo   |
    When I send "POST" request to "{{global.chain-registry}}/chains" with json:
      """
      {
        "name": "geth-{{scenarioID}}",
        "urls": {{global.nodes.geth.URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "chainUID"
    When I send "GET" request to "{{global.chain-registry}}/chains/{{chainUID}}"
    Then the response code should be 200
    When I send "PATCH" request to "{{global.chain-registry}}/chains/{{chainUID}}" with json:
      """
      {
        "name": "geth-new-{{scenarioID}}"
      }
      """
    Then the response code should be 200
    When I send "DELETE" request to "{{global.chain-registry}}/chains/{{chainUID}}"
    Then the response code should be 204

  Scenario: Valid X-API-Key and invalid X-Tenant-ID
    Given I set the headers
      | Key         | Value    |
      | X-API-Key   | with-key |
      | X-Tenant-ID |    foo   |
    When I send "POST" request to "{{global.chain-registry}}/chains" with json:
      """
      {
        "name": "geth-{{scenarioID}}",
        "urls": {{global.nodes.geth.URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "chainUID"
    Given I set the headers
      | Key         | Value    |
      | X-API-Key   | with-key |
      | X-Tenant-ID |    bar   |
    When I send "GET" request to "{{global.chain-registry}}/chains/{{chainUID}}"
    Then the response code should be 404
    When I send "PATCH" request to "{{global.chain-registry}}/chains/{{chainUID}}" with json:
      """
      {
        "name": "geth-new-{{scenarioID}}"
      }
      """
    Then the response code should be 404
    When I send "DELETE" request to "{{global.chain-registry}}/chains/{{chainUID}}"
    Then the response code should be 404
    When I send "DELETE" request to "{{global.chain-registry}}/chains/{{chainUID}}"
    Then the response code should be 404
    Given I set the headers
      | Key         |   Value     |
      | X-API-Key   | with-key    |
      | X-Tenant-ID |    foo      |
     When I send "DELETE" request to "{{global.chain-registry}}/chains/{{chainUID}}"
    Then the response code should be 204

  Scenario: JWT token and X-Tenant-ID unset
    Given I have the following tenants
      | alias   | tenantID                 |
      | tenant1 |   foo-{{scenarioID}}    |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.chain-registry}}/chains" with json:
      """
      {
        "name": "geth-{{scenarioID}}",
        "urls": {{global.nodes.geth.URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "chainUID"
    When I send "GET" request to "{{global.chain-registry}}/chains/{{chainUID}}"
    Then the response code should be 200
    When I send "PATCH" request to "{{global.chain-registry}}/chains/{{chainUID}}" with json:
      """
      {
        "name": "geth-new-{{scenarioID}}"
      }
      """
    Then the response code should be 200
    When I send "DELETE" request to "{{global.chain-registry}}/chains/{{chainUID}}"
    Then the response code should be 204

  Scenario: JWT token
    Given I have the following tenants
      | alias   | tenantID                 |
      | tenant1 |   foo-{{scenarioID}}    |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
      | X-Tenant-ID   |            _             |
    When I send "POST" request to "{{global.chain-registry}}/chains" with json:
      """
      {
        "name": "geth-{{scenarioID}}",
        "urls": {{global.nodes.geth.URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "chainUID"
    When I send "GET" request to "{{global.chain-registry}}/chains/{{chainUID}}"
    Then the response code should be 200
    When I send "PATCH" request to "{{global.chain-registry}}/chains/{{chainUID}}" with json:
      """
      {
        "name": "geth-new-{{scenarioID}}"
      }
      """
    Then the response code should be 200
    When I send "DELETE" request to "{{global.chain-registry}}/chains/{{chainUID}}"
    Then the response code should be 204

  Scenario: Invalid JWT token
    Given I have the following tenants
      | alias   | tenantID                 |
      | tenant1 |   foo-{{scenarioID}}    |
      | tenant2 |   bar-{{scenarioID}}    |
      | tenant3 |   *                      |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.chain-registry}}/chains" with json:
      """
      {
        "name": "geth-{{scenarioID}}",
        "urls": {{global.nodes.geth.URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "chainUID"
    Given I set the headers
      | Key           | Value                      |
      | Authorization | Bearer {{tenant2.token}}   |
    When I send "GET" request to "{{global.chain-registry}}/chains/{{chainUID}}"
    Then the response code should be 404
    When I send "PATCH" request to "{{global.chain-registry}}/chains/{{chainUID}}" with json:
      """
      {
        "name": "geth-new-{{scenarioID}}"
      }
      """
    Then the response code should be 404
    When I send "DELETE" request to "{{global.chain-registry}}/chains/{{chainUID}}"
    Then the response code should be 404
    When I send "DELETE" request to "{{global.chain-registry}}/chains/{{chainUID}}"
    Then the response code should be 404
    Given I set the headers
      | Key         |   Value     |
      | Authorization | Bearer {{tenant3.token}}   |
    When I send "DELETE" request to "{{global.chain-registry}}/chains/{{chainUID}}"
    Then the response code should be 204
  
  Scenario: Wildcard * JWT token
    Given I have the following tenants
      | alias   | tenantID                 |
      | tenant1 |   foo-{{scenarioID}}    |
      | tenant2 |   *   |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.chain-registry}}/chains" with json:
      """
      {
        "name": "geth-{{scenarioID}}",
        "urls": {{global.nodes.geth.URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "chainUID"
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant2.token}} |
    When I send "GET" request to "{{global.chain-registry}}/chains/{{chainUID}}"
    Then the response code should be 200
    When I send "PATCH" request to "{{global.chain-registry}}/chains/{{chainUID}}" with json:
      """
      {
        "name": "geth-new-{{scenarioID}}"
      }
      """
    Then the response code should be 200
    When I send "DELETE" request to "{{global.chain-registry}}/chains/{{chainUID}}"
    Then the response code should be 204
