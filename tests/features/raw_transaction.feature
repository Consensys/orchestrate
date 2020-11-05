@public-tx
@raw-tx
Feature: Send raw transfer transaction
  As an external developer
  I want to process a raw transaction using transaction scheduler API

  Background:
    Given I have the following tenants
      | alias   | tenantID        |
      | tenant1 | {{random.uuid}} |
    And I register the following chains
      | alias | Name                | URLs                         | Headers.Authorization    |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu_1.URLs}} | Bearer {{tenant1.token}} |

  Scenario: Send raw transaction
    Given I register the following alias
      | alias          | value              |
      | random_account | {{random.account}} |
    Given I sign the following transactions
      | alias  | ID              | Data | Gas   | To                 | Nonce | privateKey             | ChainUUID     | Headers.Authorization    |
      | rawTx1 | {{random.uuid}} | 0x   | 21000 | {{random_account}} | 0     | {{random.private_key}} | {{besu.UUID}} | Bearer {{tenant1.token}} |
    Then I track the following envelopes
      | ID            |
      | {{rawTx1.ID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/send-raw" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
      "raw": "{{rawTx1.Raw}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{rawTx1.ID}}"
    }
}
      """
    Then the response code should be 202
    Then I register the following response fields
      | alias   | path         |
      | jobUUID | jobs[0].uuid |
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status |
      | 1              |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | MINED  | CREATED        | STARTED        | PENDING        | MINED          |
