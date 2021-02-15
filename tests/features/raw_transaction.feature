@public-tx
@raw-tx
Feature: Send raw transfer transaction
  As an external developer
  I want to process a raw transaction using transaction scheduler API

  Background:
    Given I have the following tenants
      | alias   | tenantID        |
      | tenant1 | {{random.uuid}} |

  Scenario: Send raw transaction
    Given I register the following alias
      | alias          | value              |
      | random_account | {{random.account}} |
      | idempotencykey | {{random.uuid}}    |
    Given I sign the following transactions
      | alias | ID              | Data | Gas   | To                 | Nonce | privateKey             | ChainUUID            | Headers.Authorization    |
      | rawTx | {{random.uuid}} | 0x   | 21000 | {{random_account}} | 0     | {{random.private_key}} | {{chain.besu0.UUID}} | Bearer {{tenant1.token}} |
    Then I track the following envelopes
      | ID           |
      | {{rawTx.ID}} |
    Given I set the headers
      | Key               | Value                    |
      | Authorization     | Bearer {{tenant1.token}} |
      | X-Idempotency-Key | {{idempotencykey}}       |
    When I send "POST" request to "{{global.api}}/transactions/send-raw" with json:
  """
{
    "chain": "{{chain.besu0.Name}}",
    "params": {
      "raw": "{{rawTx.Raw}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{rawTx.ID}}"
    }
}
      """
    Then the response code should be 202
    Then I register the following response fields
      | alias   | path         |
      | jobUUID | jobs[0].uuid |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status |
      | 1              |
    When I send "GET" request to "{{global.api}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | MINED  | CREATED        | STARTED        | PENDING        | MINED          |

  Scenario: Send same raw transaction twice
    Given I register the following alias
      | alias          | value              |
      | random_account | {{random.account}} |
    Given I sign the following transactions
      | alias | ID              | Data | Gas   | To                 | Nonce | privateKey             | ChainUUID            | Headers.Authorization    |
      | rawTx | {{random.uuid}} | 0x   | 21000 | {{random_account}} | 0     | {{random.private_key}} | {{chain.besu0.UUID}} | Bearer {{tenant1.token}} |
    Then I track the following envelopes
      | ID           |
      | {{rawTx.ID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.api}}/transactions/send-raw" with json:
  """
{
    "chain": "{{chain.besu0.Name}}",
    "params": {
      "raw": "{{rawTx.Raw}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{rawTx.ID}}"
    }
}
      """
    Then the response code should be 202
    Then I register the following response fields
      | alias   | path         |
      | jobUUID | jobs[0].uuid |
    Then Envelopes should be in topic "tx.decoded"
    When I send "GET" request to "{{global.api}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | MINED  | CREATED        | STARTED        | PENDING        | MINED          |
    When I send "POST" request to "{{global.api}}/transactions/send-raw" with json:
  """
{
    "chain": "{{chain.besu0.Name}}",
    "params": {
      "raw": "{{rawTx.Raw}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{rawTx.ID}}"
    }
}
      """
    Then the response code should be 202
    Then I register the following response fields
      | alias   | path         |
      | jobUUID | jobs[0].uuid |
    Then Envelopes should be in topic "tx.recover"
    When I send "GET" request to "{{global.api}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | FAILED | CREATED        | STARTED        | PENDING        | FAILED         |
