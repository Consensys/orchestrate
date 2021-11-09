@tx-sentry
Feature: Send transactions using tx-sentry
  As an external developer
  I want to send transactions using tx-sentry retry feature

  Background:
    Given I have the following tenants
      | alias   | tenantID        |
      | tenant1 | {{random.uuid}} |
    And I have created the following accounts
      | alias    | ID              | API-KEY            | Tenant               |
      | account1 | {{random.uuid}} | {{global.api-key}} | {{tenant1.tenantID}} |
      | account2 | {{random.uuid}} | {{global.api-key}} | {{tenant1.tenantID}} |
    Then I register the following contracts
      | name        | artifacts        | API-KEY            | Tenant               |
      | SimpleToken | SimpleToken.json | {{global.api-key}} | {{tenant1.tenantID}} |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    Then I track the following envelopes
      | ID                  |
      | faucet-{{account1}} |
    When I send "POST" request to "{{global.api}}/transactions/transfer" with json:
      """
      {
        "chain": "{{chain.geth0.Name}}",
        "params": {
          "from": "{{global.nodes.geth[0].fundedPublicKeys[0]}}",
          "to": "{{account1}}",
          "value": "100000000000000000"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "faucet-{{account1}}"
        }
      }
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"

  Scenario: Retry transaction with zero gas increment
    Given I register the following alias
      | alias   | value           |
      | txOneID | {{random.uuid}} |
      | txTwoID | {{random.uuid}} |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/schedules" with json:
      """
      {}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias        | path |
      | scheduleUUID | uuid |
    Then I track the following envelopes
      | ID          |
      | {{txTwoID}} |
    When I send "POST" request to "{{global.api}}/jobs" with json:
      """
      {
        "scheduleUUID": "{{scheduleUUID}}",
        "chainUUID": "{{chain.besu0.UUID}}",
        "type": "eth://ethereum/transaction",
        "transaction": {
          "from": "{{account1}}",
          "to": "{{random.account}}",
          "data": "0x",
          "gas": "21000",
          "nonce": "1"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{txTwoID}}"
        },
        "annotations": {
          "gasPricePolicy": {
            "retryPolicy": {
              "interval": "1s"
            }
          }
        }
      }
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias     | path |
      | txTwoUUID | uuid |
    When I send "PUT" request to "{{global.api}}/jobs/{{txTwoUUID}}/start"
    Then the response code should be 202
    Then I sleep "5s"
    When I send "POST" request to "{{global.api}}/jobs" with json:
      """
      {
        "scheduleUUID": "{{scheduleUUID}}",
        "chainUUID": "{{chain.besu0.UUID}}",
        "type": "eth://ethereum/transaction",
        "transaction": {
          "from": "{{account1}}",
          "to": "{{random.account}}",
          "data": "0x",
          "gas": "21000",
          "nonce": "0"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{txOneID}}"
        }
      }
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias     | path |
      | txOneUUID | uuid |
    When I send "PUT" request to "{{global.api}}/jobs/{{txOneUUID}}/start"
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    When I send "GET" request to "{{global.api}}/jobs/{{txTwoUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | MINED  | CREATED        | STARTED        | PENDING        | RESENDING      |

  @geth
  Scenario: Retry transaction with gas increment
    Given I register the following alias
      | alias   | value           |
      | txOneID | {{random.uuid}} |
      | txTwoID | {{random.uuid}} |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/schedules" with json:
      """
      {}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias        | path |
      | scheduleUUID | uuid |
    Then I track the following envelopes
      | ID          |
      | {{txTwoID}} |
    When I send "POST" request to "{{global.api}}/jobs" with json:
      """
      {
        "scheduleUUID": "{{scheduleUUID}}",
        "chainUUID": "{{chain.geth0.UUID}}",
        "type": "eth://ethereum/transaction",
        "transaction": {
          "from": "{{account1}}",
          "to": "{{random.account}}",
          "data": "0x",
          "gas": "21000",
          "nonce": "1"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{txTwoID}}"
        },
        "annotations": {
          "gasPricePolicy": {
            "retryPolicy": {
              "interval": "1s",
              "increment": 0.15,
              "limit": 0.45
            }
          }
        }
      }
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias     | path |
      | txTwoUUID | uuid |
    When I send "PUT" request to "{{global.api}}/jobs/{{txTwoUUID}}/start"
    Then the response code should be 202
    Then I sleep "5s"
    When I send "GET" request to "{{global.api}}/schedules/{{scheduleUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | jobs[0].status | jobs[1].status | jobs[2].status | jobs[3].status |
      | PENDING        | PENDING        | PENDING        | PENDING        |
    When I send "POST" request to "{{global.api}}/jobs" with json:
      """
      {
        "scheduleUUID": "{{scheduleUUID}}",
        "chainUUID": "{{chain.geth0.UUID}}",
        "type": "eth://ethereum/transaction",
        "transaction": {
          "from": "{{account1}}",
          "to": "{{random.account}}",
          "data": "0x",
          "gas": "21000",
          "nonce": "0"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{txOneID}}"
        }
      }
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias     | path |
      | txOneUUID | uuid |
    When I send "PUT" request to "{{global.api}}/jobs/{{txOneUUID}}/start"
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    When I send "GET" request to "{{global.api}}/schedules/{{scheduleUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | jobs[0].status | jobs[1].status | jobs[2].status | jobs[3].status |
      | NEVER_MINED    | NEVER_MINED    | NEVER_MINED    | MINED          |

  Scenario: Send transaction using retry policy with zero gas increment to retry limit
    Given I register the following alias
      | alias   | value           |
      | txOneID | {{random.uuid}} |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/schedules" with json:
      """
      {}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias        | path |
      | scheduleUUID | uuid |
    Then I track the following envelopes
      | ID          |
      | {{txOneID}} |
    When I send "POST" request to "{{global.api}}/jobs" with json:
      """
      {
        "scheduleUUID": "{{scheduleUUID}}",
        "chainUUID": "{{chain.besu0.UUID}}",
        "type": "eth://ethereum/transaction",
        "transaction": {
          "from": "{{account2}}",
          "to": "{{random.account}}",
          "data": "0x",
          "gas": "21000",
          "nonce": "1"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{txOneID}}"
        },
        "annotations": {
          "gasPricePolicy": {
            "retryPolicy": {
              "interval": "1s"
            }
          }
        }
      }
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias     | path |
      | txOneUUID | uuid |
    When I send "PUT" request to "{{global.api}}/jobs/{{txOneUUID}}/start"
    Then the response code should be 202
    Then I sleep "15s"
    When I send "GET" request to "{{global.api}}/jobs/{{txOneUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status  | annotations.hasBeenRetried |
      | PENDING | true                       |
