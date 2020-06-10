@multi-tenancy
Feature: Chain-Proxy Authentication
  Scenario: Chain-Proxy Auth
    Given I have the following tenants
      | alias    | tenantID |
      | foo      |   foo    |
      | bar      |   bar    |
      | wildcard |    *     |
    Given I set the headers
      | Key       | Value    |
      | Authorization | Bearer {{foo.token}} |
    When I send "POST" request to "{{global.chain-registry}}/chains" with json:
      """
      {
        "name": "geth-foo-{{scenarioUID}}",
        "urls": {{global.nodes.geth.URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "chainFoo"

    Given I set the headers
      | Key       | Value    |
      | Authorization | Bearer {{foo.token}} |
      | X-Tenant-ID   |            _             |
    When I send "POST" request to "{{global.chain-registry}}/chains" with json:
      """
      {
        "name": "geth-default-{{scenarioUID}}",
        "urls": {{global.nodes.geth.URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "chainDefault"

    Given I sleep "3s"

    Given I set the headers
      | Key       | Value    |
      | X-API-Key | with-key |
      | Content-Type | application/json |
    When I send "POST" request to "{{global.chain-registry}}/{{chainFoo}}" with json:
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
    When I send "POST" request to "{{global.chain-registry}}/{{chainDefault}}" with json:
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
      | Key       | Value       |
      | X-API-Key | unknown-key |
      | Content-Type | application/json |
      When I send "POST" request to "{{global.chain-registry}}/{{chainFoo}}" with json:
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
    When I send "POST" request to "{{global.chain-registry}}/{{chainDefault}}" with json:
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
      | Key       | Value    |
      | Authorization | Bearer {{foo.token}} |
      | Content-Type | application/json |
    When I send "POST" request to "{{global.chain-registry}}/{{chainFoo}}" with json:
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
    When I send "POST" request to "{{global.chain-registry}}/{{chainDefault}}" with json:
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
      | Key       | Value    |
      | Authorization | Bearer {{bar.token}} |
      | Content-Type | application/json |
    When I send "POST" request to "{{global.chain-registry}}/{{chainFoo}}" with json:
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
    When I send "POST" request to "{{global.chain-registry}}/{{chainDefault}}" with json:
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
      | Key       | Value    |
      | Authorization | Bearer {{wildcard.token}} |
      | Content-Type | application/json |
    When I send "POST" request to "{{global.chain-registry}}/{{chainFoo}}" with json:
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
    When I send "POST" request to "{{global.chain-registry}}/{{chainDefault}}" with json:
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
      | Key       | Value    |
      | Authorization | Bearer {{wildcard.token}} |
      | Content-Type | application/json |
    When I send "DELETE" request to "{{global.chain-registry}}/chains/{{chainFoo}}"
    Then the response code should be 204
    When I send "DELETE" request to "{{global.chain-registry}}/chains/{{chainDefault}}"
    Then the response code should be 204
