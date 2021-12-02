@faucet
@multi-tenancy
Feature: Faucet funding
  As as external developer
  I want to fund accounts using registered faucet and multiple tenants

  Background:
    Given I have the following tenants
      | alias     | tenantID |
      | tenantFoo | foo      |
      | tenantBar | bar      |

  Scenario: Generate account with faucet and different tenant
    And I register the following faucets
      | Name                  | ChainRule            | CreditorAccount                              | MaxBalance      | Amount          | Cooldown | API-KEY            | Tenant                 |
      | faucet-{{scenarioID}} | {{chain.besu0.UUID}} | {{global.nodes.besu[0].fundedPublicKeys[0]}} | 0x38D7EA4C68000 | 0x38D7EA4C68000 | 1m       | {{global.api-key}} | {{tenantFoo.tenantID}} |
    And I have created the following accounts
      | alias    | ID              | ChainName            | API-KEY            | Tenant                 |
      | account1 | {{random.uuid}} | {{chain.besu0.Name}} | {{global.api-key}} | {{tenantBar.tenantID}} |
    Given I sleep "5s"
    Given I set the headers
      | Key             | Value                  |
      | X-API-KEY       | {{global.api-key}}     |
      | X-TENANT-ID     | {{tenantBar.tenantID}} |
      | X-Cache-Control | no-cache               |
    When I send "POST" request to "{{global.api}}/proxy/chains/{{chain.besu0.UUID}}" with json:
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
      | result |
      | 0x0    |

  Scenario: Send transaction with faucet and different tenant
    Given I register the following alias
      | alias         | value              |
      | toAddr        | {{random.account}} |
      | transferOneID | {{random.uuid}}    |
    And I have created the following accounts
      | alias    | ID              | ChainName            | API-KEY            | Tenant                 |
      | account1 | {{random.uuid}} | {{chain.besu0.Name}} | {{global.api-key}} | {{tenantBar.tenantID}} |
    And I register the following faucets
      | Name                  | ChainRule            | CreditorAccount                              | MaxBalance      | Amount          | Cooldown | API-KEY            | Tenant                 |
      | faucet-{{scenarioID}} | {{chain.besu0.UUID}} | {{global.nodes.besu[0].fundedPublicKeys[0]}} | 0x38D7EA4C68000 | 0x38D7EA4C68000 | 1m       | {{global.api-key}} | {{tenantFoo.tenantID}} |
    Then I track the following envelopes
      | ID                |
      | {{transferOneID}} |
    Given I set the headers
      | Key         | Value                  |
      | X-API-KEY   | {{global.api-key}}     |
      | X-TENANT-ID | {{tenantBar.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/transfer" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{account1}}",
          "to": "{{toAddr}}",
          "value": "0x16345785D8A0000"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{transferOneID}}"
        }
      }
      """
    Then the response code should be 202
    And Response should have the following fields
      | jobs.length |
      | 1           |
    Then I register the following response fields
      | alias     | path         |
      | txJobUUID | jobs[0].uuid |
    Then Envelopes should be in topic "tx.recover"
    When I send "GET" request to "{{global.api}}/jobs/{{txJobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | FAILED | CREATED        | STARTED        | PENDING        | FAILED         |
    Given I set the headers
      | Key             | Value                  |
      | X-API-KEY       | {{global.api-key}}     |
      | X-TENANT-ID     | {{tenantBar.tenantID}} |
      | X-Cache-Control | no-cache               |
    When I send "POST" request to "{{global.api}}/proxy/chains/{{chain.besu0.UUID}}" with json:
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
      | result |
      | 0x0    |
