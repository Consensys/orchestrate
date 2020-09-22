@chain-registry
Feature: Chain-Proxy Cache
  As as external developer
  I want to perform proxy calls to my chains

  Background:
    Given I have the following tenants
      | alias   | tenantID        |
      | tenant1 | {{random.uuid}} |
      | tenant2 | {{random.uuid}} |
    Then I register the following chains
      | alias     | Name                 | URLs                         | Headers.Authorization    |
      | besuOne   | besu-{{scenarioID}}  | {{global.nodes.besu_2.URLs}} | Bearer {{tenant1.token}} |
      | besuTwo   | besu2-{{scenarioID}} | {{global.nodes.besu_2.URLs}} | Bearer {{tenant1.token}} |
    Given I sleep "3s"

  Scenario: Chain registry should cache "eth_getBlockByNumber" request for same chainUUID
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.chain-registry}}/{{besuOne.UUID}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": ["0x1",false],
        "id": 1
      }
      """
    Then the response code should be 200
    And Response should have the following headers
      | X-Cache-Control | Content-Type     |
      | -               | application/json |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.chain-registry}}/{{besuOne.UUID}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": ["0x1",false],
        "id": 1
      }
      """
    Then the response code should be 200
    And Response should have the following headers
      | X-Cache-Control | Content-Type     |
      | ~               | application/json |
    Given I sleep "3s"
    When I send "POST" request to "{{global.chain-registry}}/{{besuOne.UUID}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": ["0x1",false],
        "id": 1
      }
      """
    Then the response code should be 200
    And Response should have the following headers
      | X-Cache-Control | Content-Type     |
      | -               | application/json |

  Scenario: Chain registry should cache "eth_getBlockByNumber" request for same chainID
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.chain-registry}}/{{besuOne.UUID}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": ["0x2",false],
        "id": 1
      }
      """
    Then the response code should be 200
    And Response should have the following headers
      | X-Cache-Control | Content-Type     |
      | -               | application/json |
    When I send "POST" request to "{{global.chain-registry}}/{{besuTwo.UUID}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": ["0x2",false],
        "id": 1
      }
      """
    Then the response code should be 200
    And Response should have the following headers
      | X-Cache-Control | Content-Type     |
      | ~               | application/json |

  Scenario: Chain registry should ignore cache when user indicate it
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.chain-registry}}/{{besuOne.UUID}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": ["0x3",false],
        "id": 1
      }
      """
    Then the response code should be 200
    And Response should have the following headers
      | X-Cache-Control | Content-Type     |
      | -               | application/json |
    Given I set the headers
      | Key             | Value                    |
      | X-Cache-Control | no-cache                 |
      | Authorization   | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.chain-registry}}/{{besuOne.UUID}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": ["0x3",false],
        "id": 1
      }
      """
    Then the response code should be 200
    And Response should have the following headers
      | X-Cache-Control | Content-Type     |
      | -               | application/json |
