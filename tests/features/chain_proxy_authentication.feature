@chain-registry
@multi-tenancy
Feature: Chain-Proxy Authentication
  As as external developer
  I want to perform proxy calls to my chains with expected permission rules

  @geth
  Scenario: Chain-Proxy Auth with tenants
    Given I have the following tenants
      | alias    | tenantID |
      | foo      | foo      |
      | bar      | bar      |
      | wildcard | _        |
    Given I set the headers
      | Key         | Value              |
      | X-API-KEY   | {{global.api-key}} |
      | X-TENANT-ID | {{foo.tenantID}}   |
    When I send "POST" request to "{{global.api}}/chains" with json:
      """
      {
      "name": "geth-foo-{{scenarioID}}",
      "urls": {{global.nodes.geth[0].URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "chainFoo"

    Given I set the headers
      | Key       | Value              |
      | X-API-KEY | {{global.api-key}} |
    When I send "POST" request to "{{global.api}}/chains" with json:
      """
      {
      "name": "geth-default-{{scenarioID}}",
      "urls": {{global.nodes.geth[0].URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "chainDefault"

    Given I sleep "3s"

    Given I set the headers
      | Key          | Value              |
      | X-API-Key    | {{global.api-key}} |
      | Content-Type | application/json   |
    When I send "POST" request to "{{global.api}}/proxy/chains/{{chainFoo}}" with json:
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
    When I send "POST" request to "{{global.api}}/proxy/chains/{{chainDefault}}" with json:
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
    When I send "POST" request to "{{global.api}}/proxy/chains/{{chainFoo}}" with json:
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
    When I send "POST" request to "{{global.api}}/proxy/chains/{{chainDefault}}" with json:
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
      | Key          | Value              |
      | X-API-KEY    | {{global.api-key}} |
      | X-TENANT-ID  | {{foo.tenantID}}   |
      | Content-Type | application/json   |
    When I send "POST" request to "{{global.api}}/proxy/chains/{{chainFoo}}" with json:
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
    When I send "POST" request to "{{global.api}}/proxy/chains/{{chainDefault}}" with json:
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
      | Key          | Value              |
      | X-API-KEY    | {{global.api-key}} |
      | X-TENANT-ID  | {{bar.tenantID}}   |
      | Content-Type | application/json   |
    When I send "POST" request to "{{global.api}}/proxy/chains/{{chainFoo}}" with json:
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
    When I send "POST" request to "{{global.api}}/proxy/chains/{{chainDefault}}" with json:
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
      | Key          | Value              |
      | X-API-KEY    | {{global.api-key}} |
      | Content-Type | application/json   |
    When I send "POST" request to "{{global.api}}/proxy/chains/{{chainFoo}}" with json:
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
    When I send "POST" request to "{{global.api}}/proxy/chains/{{chainDefault}}" with json:
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
      | Key          | Value              |
      | X-API-KEY    | {{global.api-key}} |
      | Content-Type | application/json   |
    When I send "DELETE" request to "{{global.api}}/chains/{{chainFoo}}"
    Then the response code should be 204
    When I send "DELETE" request to "{{global.api}}/chains/{{chainDefault}}"
    Then the response code should be 204


  @besu
  Scenario: Chain-Proxy Auth with username
    Given I have the following tenants
      | alias            | tenantID  | username |
      | tenantOne        | tenantOne |          |
      | tenantOneUserOne | tenantOne | userOne  |
      | tenantOneUserTwo | tenantOne | userTwo  |
    Given I set the headers
      | Key         | Value                  |
      | X-API-KEY   | {{global.api-key}}     |
      | X-TENANT-ID | {{tenantOne.tenantID}} |
    When I send "POST" request to "{{global.api}}/chains" with json:
      """
      {
      "name": "besu-shared-{{scenarioID}}",
      "urls": {{global.nodes.besu[0].URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "sharedChain"
    Given I set the headers
      | Key         | Value                         |
      | X-API-KEY   | {{global.api-key}}            |
      | X-TENANT-ID | {{tenantOneUserOne.tenantID}} |
      | X-USERNAME  | {{tenantOneUserOne.username}} |
    When I send "POST" request to "{{global.api}}/chains" with json:
      """
      {
      "name": "besu-userOne-{{scenarioID}}",
      "urls": {{global.nodes.besu[0].URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "userOneChain"
    Given I set the headers
      | Key         | Value                         |
      | X-API-KEY   | {{global.api-key}}            |
      | X-TENANT-ID | {{tenantOneUserTwo.tenantID}} |
      | X-USERNAME  | {{tenantOneUserTwo.username}} |
    When I send "POST" request to "{{global.api}}/chains" with json:
      """
      {
      "name": "besu-userUserTwo-{{scenarioID}}",
      "urls": {{global.nodes.besu[0].URLs}}
      }
      """
    Then the response code should be 200
    Then I store the UUID as "userTwoChain"
    Given I sleep "3s"
    Given I set the headers
      | Key         | Value                         |
      | X-API-Key   | {{global.api-key}}            |
      | X-TENANT-ID | {{tenantOneUserOne.tenantID}} |
      | X-USERNAME  | {{tenantOneUserOne.username}} |
    When I send "POST" request to "{{global.api}}/proxy/chains/{{sharedChain}}" with json:
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
    When I send "POST" request to "{{global.api}}/proxy/chains/{{userOneChain}}" with json:
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
    When I send "POST" request to "{{global.api}}/proxy/chains/{{userTwoChain}}" with json:
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
    Given I set the headers
      | Key         | Value                         |
      | X-API-Key   | {{global.api-key}}            |
      | X-TENANT-ID | {{tenantOneUserTwo.tenantID}} |
      | X-USERNAME  | {{tenantOneUserTwo.username}} |
    When I send "POST" request to "{{global.api}}/proxy/chains/{{sharedChain}}" with json:
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
    When I send "POST" request to "{{global.api}}/proxy/chains/{{userOneChain}}" with json:
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
    When I send "POST" request to "{{global.api}}/proxy/chains/{{userTwoChain}}" with json:
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
      | Key         | Value                         |
      | X-API-Key   | {{global.api-key}}            |
      | X-TENANT-ID | {{tenantOneUserTwo.tenantID}} |
    When I send "POST" request to "{{global.api}}/proxy/chains/{{sharedChain}}" with json:
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
    When I send "POST" request to "{{global.api}}/proxy/chains/{{userOneChain}}" with json:
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
    When I send "POST" request to "{{global.api}}/proxy/chains/{{userTwoChain}}" with json:
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
    Given I set the headers
      | Key          | Value              |
      | X-API-KEY    | {{global.api-key}} |
      | X-USERNAME   | *                  |
      | Content-Type | application/json   |
    When I send "DELETE" request to "{{global.api}}/chains/{{sharedChain}}"
    Then the response code should be 204
    When I send "DELETE" request to "{{global.api}}/chains/{{userOneChain}}"
    Then the response code should be 204
    When I send "DELETE" request to "{{global.api}}/chains/{{userTwoChain}}"
    Then the response code should be 204
