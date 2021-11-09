@external-tx
Feature: Listen to external transactions
  As an external developer
  I want to listen to transactions external to Orchestrate

  Background:
    Given I have the following tenants
      | alias   | tenantID        |
      | tenant1 | {{random.uuid}} |
    And I register the following chains
      | alias | Name                | URLs                          | API-KEY            | Tenant               | Listener.ExternalTxEnabled |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu[2].URLs}} | {{global.api-key}} | {{tenant1.tenantID}} | true                       |

  Scenario: Listen to external tx
    Given I register the following alias
      | alias          | value              |
      | random_account | {{random.account}} |
    Given I sign the following transactions
      | alias | ID              | Data | Gas   | To                 | Nonce | privateKey             | ChainUUID     | API-KEY            | Tenant               |
      | rawTx | {{random.uuid}} | 0x   | 21000 | {{random_account}} | 0     | {{random.private_key}} | {{besu.UUID}} | {{global.api-key}} | {{tenant1.tenantID}} |
    Then I track the following envelopes
      | ID                           |
      | {{global.external-tx-label}} |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/proxy/chains/{{besu.UUID}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_sendRawTransaction",
        "params": ["{{rawTx.Raw}}"],
        "id": 1
      }
      """
    Then the response code should be 200
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status |
      | 1              |
