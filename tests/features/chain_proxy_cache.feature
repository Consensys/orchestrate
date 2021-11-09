@chain-registry
Feature: Chain-Proxy Cache
  As as external developer
  I want to perform proxy calls to my chains

  Background:
    Given I have the following tenants
      | alias   | tenantID        |
      | tenant1 | {{random.uuid}} |

  Scenario: Chain registry should cache "eth_getBlockByNumber" request for same chainUUID
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/proxy/chains/{{chain.besu1.UUID}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": [
          "0x1",
          false
        ],
        "id": 1
      }
      """
    Then the response code should be 200
    And Response should have the following headers
      | X-Cache-Control | Content-Type     |
      | -               | application/json |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/proxy/chains/{{chain.besu1.UUID}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": [
          "0x1",
          false
        ],
        "id": 1
      }
      """
    Then the response code should be 200
    And Response should have the following headers
      | X-Cache-Control | Content-Type     |
      | ~               | application/json |
    Given I sleep "3s"
    When I send "POST" request to "{{global.api}}/proxy/chains/{{chain.besu1.UUID}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": [
          "0x1",
          false
        ],
        "id": 1
      }
      """
    Then the response code should be 200
    And Response should have the following headers
      | X-Cache-Control | Content-Type     |
      | -               | application/json |

  Scenario: Chain registry should ignore cache when user indicate it
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/proxy/chains/{{chain.besu1.UUID}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": [
          "0x3",
          false
        ],
        "id": 1
      }
      """
    Then the response code should be 200
    And Response should have the following headers
      | X-Cache-Control | Content-Type     |
      | -               | application/json |
    Given I set the headers
      | Key             | Value                |
      | X-Cache-Control | no-cache             |
      | X-API-KEY       | {{global.api-key}}   |
      | X-TENANT-ID     | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/proxy/chains/{{chain.besu1.UUID}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": [
          "0x3",
          false
        ],
        "id": 1
      }
      """
    Then the response code should be 200
    And Response should have the following headers
      | X-Cache-Control | Content-Type     |
      | -               | application/json |
