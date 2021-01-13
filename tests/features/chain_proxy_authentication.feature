@chain-registry
@multi-tenancy
Feature: Chain-Proxy Authentication
  As as external developer
  I want to perform proxy calls to my chains with expected permission rules

  @geth
  Scenario: Chain-Proxy Auth
    Given I have the following tenants
      | alias    | tenantID |
      | foo      | foo      |
      | bar      | bar      |
      | wildcard | *        |
    Given I set the headers
      | Key           | Value                |
      | Authorization | Bearer {{foo.token}} |
    When I send "POST" request to "{{global.api}}/chains" with json:
      """
      {
      "name": "geth-foo-{{scenarioID}}",
      "urls": {{global.nodes.geth.URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "chainFoo"

    Given I set the headers
      | Key           | Value                |
      | Authorization | Bearer {{foo.token}} |
      | X-Tenant-ID   | _                    |
    When I send "POST" request to "{{global.api}}/chains" with json:
      """
      {
      "name": "geth-default-{{scenarioID}}",
      "urls": {{global.nodes.geth.URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "chainDefault"

    Given I sleep "3s"

    Given I set the headers
      | Key          | Value              |
      | X-API-Key    | {{global.api-key}} |
      | Content-Type | application/json   |
    When I send "POST" request to "{{global.api}}/{{chainFoo}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": [
          "latest",
          false
        ],
        "id": 1
      }
      """
    Then the response code should be 200
    When I send "POST" request to "{{global.api}}/{{chainDefault}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": [
          "latest",
          false
        ],
        "id": 1
      }
      """
    Then the response code should be 200

    Given I set the headers
      | Key          | Value            |
      | X-API-Key    | unknown-key      |
      | Content-Type | application/json |
    When I send "POST" request to "{{global.api}}/{{chainFoo}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": [
          "latest",
          false
        ],
        "id": 1
      }
      """
    Then the response code should be 401
    When I send "POST" request to "{{global.api}}/{{chainDefault}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": [
          "latest",
          false
        ],
        "id": 1
      }
      """
    Then the response code should be 401

    Given I set the headers
      | Key           | Value                |
      | Authorization | Bearer {{foo.token}} |
      | Content-Type  | application/json     |
    When I send "POST" request to "{{global.api}}/{{chainFoo}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": [
          "latest",
          false
        ],
        "id": 1
      }
      """
    Then the response code should be 200
    When I send "POST" request to "{{global.api}}/{{chainDefault}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": [
          "latest",
          false
        ],
        "id": 1
      }
      """
    Then the response code should be 200

    Given I set the headers
      | Key           | Value                |
      | Authorization | Bearer {{bar.token}} |
      | Content-Type  | application/json     |
    When I send "POST" request to "{{global.api}}/{{chainFoo}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": [
          "latest",
          false
        ],
        "id": 1
      }
      """
    Then the response code should be 404
    When I send "POST" request to "{{global.api}}/{{chainDefault}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": [
          "latest",
          false
        ],
        "id": 1
      }
      """
    Then the response code should be 200

    Given I set the headers
      | Key           | Value                     |
      | Authorization | Bearer {{wildcard.token}} |
      | Content-Type  | application/json          |
    When I send "POST" request to "{{global.api}}/{{chainFoo}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": [
          "latest",
          false
        ],
        "id": 1
      }
      """
    Then the response code should be 200
    When I send "POST" request to "{{global.api}}/{{chainDefault}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": [
          "latest",
          false
        ],
        "id": 1
      }
      """
    Then the response code should be 200

    Given I set the headers
      | Key           | Value                     |
      | Authorization | Bearer {{wildcard.token}} |
      | Content-Type  | application/json          |
    When I send "DELETE" request to "{{global.api}}/chains/{{chainFoo}}"
    Then the response code should be 204
    When I send "DELETE" request to "{{global.api}}/chains/{{chainDefault}}"
    Then the response code should be 204
