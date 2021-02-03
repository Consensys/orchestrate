@faucet
Feature: Faucet funding
  As as external developer
  I want to fund accounts using registered faucet

  Background:
    Given I have the following tenants
      | alias   | tenantID        |
      | tenant1 | {{random.uuid}} |
    And I register the following chains
      | alias | Name                | URLs                          | Headers.Authorization    |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu[0].URLs}} | Bearer {{tenant1.token}} |

  Scenario: Generate account with faucet
    And I register the following faucets
      | Name                | ChainRule     | CreditorAccount                              | MaxBalance       | Amount           | Cooldown | Headers.Authorization    |
      | besu-{{scenarioID}} | {{besu.UUID}} | {{global.nodes.besu[0].fundedPublicKeys[0]}} | 1000000000000000 | 1000000000000000 | 1m       | Bearer {{tenant1.token}} |
    And I have created the following accounts
      | alias    | ID              | ChainName           | Headers.Authorization    |
      | account1 | {{random.uuid}} | besu-{{scenarioID}} | Bearer {{tenant1.token}} |
    Given I sleep "11s"
    Given I set the headers
      | Key             | Value                    |
      | Authorization   | Bearer {{tenant1.token}} |
      | X-Cache-Control | no-cache                 |
    When I send "POST" request to "{{global.api}}/proxy/chains/{{besu.UUID}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBalance",
        "params": [
          "{{account1}}",
          "latest"
        ],
        "id": 1
      }
      """
    Then the response code should be 200
    And Response should have the following fields
      | result          |
      | 0x38d7ea4c68000 |

  Scenario: Send transaction with faucet
    Given I register the following alias
      | alias         | value              |
      | toAddr        | {{random.account}} |
      | transferOneID | {{random.uuid}}    |
    And I have created the following accounts
      | alias    | ID              | ChainName           | Headers.Authorization    |
      | account1 | {{random.uuid}} | besu-{{scenarioID}} | Bearer {{tenant1.token}} |
    And I register the following faucets
      | Name                | ChainRule     | CreditorAccount                              | MaxBalance       | Amount           | Cooldown | Headers.Authorization    |
      | besu-{{scenarioID}} | {{besu.UUID}} | {{global.nodes.besu[0].fundedPublicKeys[0]}} | 1000000000000000 | 1000000000000000 | 1m       | Bearer {{tenant1.token}} |
    Then I track the following envelopes
      | ID                |
      | {{transferOneID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.api}}/transactions/transfer" with json:
      """
      {
        "chain": "besu-{{scenarioID}}",
        "params": {
          "from": "{{account1}}",
          "to": "{{toAddr}}",
          "value": "100000000000000"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{transferOneID}}"
        }
      }
      """
    Then the response code should be 202
    Then I register the following response fields
      | alias         | path         |
      | txJobUUID     | jobs[0].uuid |
      | faucetJobUUID | jobs[1].uuid |
    Then Envelopes should be in topic "tx.recover"
    When I send "GET" request to "{{global.api}}/jobs/{{txJobUUID}}"
    Then the response code should be 200
    Given I sleep "1s"
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status
      | FAILED | CREATED        | STARTED        | PENDING        | FAILED
    Given I sleep "11s"
    When I send "GET" request to "{{global.api}}/jobs/{{faucetJobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | MINED  | CREATED        | STARTED        | PENDING        | MINED          |
    Given I set the headers
      | Key             | Value                    |
      | Authorization   | Bearer {{tenant1.token}} |
      | X-Cache-Control | no-cache                 |
    When I send "POST" request to "{{global.api}}/proxy/chains/{{besu.UUID}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBalance",
        "params": [
          "{{account1}}",
          "latest"
        ],
        "id": 1
      }
      """
    Then the response code should be 200
    And Response should have the following fields
      | result          |
      | 0x38d7ea4c68000 |
