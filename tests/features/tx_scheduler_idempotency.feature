@tx-scheduler
Feature: Transaction Scheduler Idempotency
  As an external developer
  I want to use idempotency-key to interact with the transaction scheduler API

  Background:
    Given I have the following tenants
      | alias   | tenantID        |
      | tenant1 | {{random.uuid}} |
    Then I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | {{tenant1.token}} |
    And I have created the following accounts
      | alias    | ID              | Headers.Authorization    |
      | account1 | {{random.uuid}} | {{tenant1.token}} |
      | account2 | {{random.uuid}} | {{tenant1.token}} |

  Scenario: Send twice same transaction using same idempotency key
    Given I register the following alias
      | alias          | value           |
      | idempotencykey | {{random.uuid}} |
      | deployTxOneID  | {{random.uuid}} |
      | deployTxTwoID  | {{random.uuid}} |
    Then I track the following envelopes
      | ID                |
      | {{deployTxOneID}} |
    Then  I set the headers
      | Key               | Value                    |
      | Authorization     | {{tenant1.token}} |
      | X-Idempotency-Key | {{idempotencykey}}       |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{account1}}",
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{deployTxOneID}}"
        }
      }
      """
    Then the response code should be 202
    Then  I set the headers
      | Key               | Value                    |
      | Authorization     | {{tenant1.token}} |
      | X-Idempotency-Key | {{idempotencykey}}       |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{account1}}",
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{deployTxTwoID}}"
        }
      }
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    When I send "GET" request to "{{global.api}}/transactions?idempotency_keys={{idempotencykey}}"
    Then the response code should be 200
    And Response should have the following fields
      | 0.idempotencyKey   | 0.jobs[0].status | 0.jobs[0].labels.id |
      | {{idempotencykey}} | MINED            | {{deployTxOneID}}   |

  Scenario: Send twice different transaction using same idempotency key
    Given I register the following alias
      | alias          | value           |
      | idempotencykey | {{random.uuid}} |
      | deployTxOneID  | {{random.uuid}} |
      | deployTxTwoID  | {{random.uuid}} |
    Then I track the following envelopes
      | ID                |
      | {{deployTxOneID}} |
    Then  I set the headers
      | Key               | Value                    |
      | Authorization     | {{tenant1.token}} |
      | X-Idempotency-Key | {{idempotencykey}}       |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{account1}}",
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{deployTxOneID}}"
        }
      }
      """
    Then the response code should be 202
    Then  I set the headers
      | Key               | Value                    |
      | Authorization     | {{tenant1.token}} |
      | X-Idempotency-Key | {{idempotencykey}}       |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "oneTimeKey": true,
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{deployTxTwoID}}"
        }
      }
      """
    Then the response code should be 409
    And Response should have the following fields
      | message                                                                                                          |
      | DB101@use-cases.send-tx: transaction request with the same idempotency key and different params already exists |
    Then Envelopes should be in topic "tx.decoded"
    When I send "GET" request to "{{global.api}}/transactions?idempotency_keys={{idempotencykey}}"
    Then the response code should be 200
    And Response should have the following fields
      | 0.idempotencyKey   | 0.jobs[0].status | 0.jobs[0].labels.id |
      | {{idempotencykey}} | MINED            | {{deployTxOneID}}   |
