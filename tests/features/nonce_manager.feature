@nonce
Feature: Nonce manager
  As an external developer
  I want transaction with empty nonce to be calibrated, sent to blockchain and mined

  Background:
    Given I have the following tenants
      | alias   | tenantID        |
      | tenant1 | {{random.uuid}} |
    And I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | Bearer {{tenant1.token}} |
    And I register the following chains
      | alias | Name                | URLs                         | Headers.Authorization    |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu_1.URLs}} | Bearer {{tenant1.token}} |
    And I have created the following accounts
      | alias    | ID              | Headers.Authorization    |
      | account1 | {{random.uuid}} | Bearer {{tenant1.token}} |
      | account2 | {{random.uuid}} | Bearer {{tenant1.token}} |
      | account3 | {{random.uuid}} | Bearer {{tenant1.token}} |
    Then I track the following envelopes
      | ID                  |
      | faucet-{{account1}} |
      | faucet-{{account2}} |
      | faucet-{{account3}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.api}}/transactions/transfer" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
      "from": "{{global.nodes.besu_1.fundedPublicKeys[0]}}",
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
    When I send "POST" request to "{{global.api}}/transactions/transfer" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
      "from": "{{global.nodes.besu_1.fundedPublicKeys[0]}}",
      "to": "{{account2}}",
      "value": "100000000000000000"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "faucet-{{account2}}"
    }
}
      """
    Then the response code should be 202
    When I send "POST" request to "{{global.api}}/transactions/transfer" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
      "from": "{{global.nodes.besu_1.fundedPublicKeys[0]}}",
      "to": "{{account3}}",
      "value": "100000000000000000"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "faucet-{{account3}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    
  Scenario: Nonce recalibrating
    Given I register the following alias
      | alias | value              |
      | to1   | {{random.account}} |
      | to2   | {{random.account}} |
      | to3   | {{random.account}} |
    Then I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.api}}/schedules" with json:
      """
{}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias           | path |
      | scheduleOneUUID | uuid |
    When I send "POST" request to "{{global.api}}/schedules" with json:
      """
{}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias           | path |
      | scheduleTwoUUID | uuid |
    When I send "POST" request to "{{global.api}}/schedules" with json:
      """
{}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias             | path |
      | scheduleThreeUUID | uuid |
    When I send "POST" request to "{{global.api}}/jobs" with json:
  """
{
    "scheduleUUID": "{{scheduleOneUUID}}",
	"chainUUID": "{{besu.UUID}}",
    "type": "eth://ethereum/transaction",
    "transaction": {
        "from": "{{account1}}",
        "to": "{{to1}}",
        "value": "140000",
        "gas": "21000"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}"
    }
}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias        | path |
      | txOneJobUUID | uuid |
    When I send "POST" request to "{{global.api}}/jobs" with json:
  """
{
    "scheduleUUID": "{{scheduleTwoUUID}}",
	"chainUUID": "{{besu.UUID}}",
    "type": "eth://ethereum/transaction",
    "transaction": {
        "from": "{{account1}}",
        "to": "{{to2}}",
        "value": "130000",
        "nonce": "2",
        "gas": "21000"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}"
    }
}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias        | path |
      | txTwoJobUUID | uuid |
    When I send "POST" request to "{{global.api}}/jobs" with json:
  """
{
    "scheduleUUID": "{{scheduleThreeUUID}}",
	"chainUUID": "{{besu.UUID}}",
    "type": "eth://ethereum/transaction",
    "transaction": {
        "from": "{{account1}}",
        "to": "{{to3}}",
        "value": "109000",
        "gas": "21000"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}"
    }
}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias          | path |
      | txThreeJobUUID | uuid |
    Then I track the following envelopes
      | ID                    |
      | {{scheduleOneUUID}}   |
      | {{scheduleTwoUUID}}   |
      | {{scheduleThreeUUID}} |
    When I send "PUT" request to "{{global.api}}/jobs/{{txOneJobUUID}}/start"
    Then the response code should be 202
    When I send "PUT" request to "{{global.api}}/jobs/{{txTwoJobUUID}}/start"
    Then the response code should be 202
    When I send "PUT" request to "{{global.api}}/jobs/{{txThreeJobUUID}}/start"
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Nonce |
      | 1              | 0     |
      | 1              | 2     |
      | 1              | 1     |

  Scenario: Chaotic nonce
    Given I register the following alias
      | alias | value              |
      | to1   | {{random.account}} |
      | to2   | {{random.account}} |
      | to3   | {{random.account}} |
    Then I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.api}}/schedules" with json:
      """
{}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias           | path |
      | scheduleOneUUID | uuid |
    When I send "POST" request to "{{global.api}}/schedules" with json:
      """
{}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias           | path |
      | scheduleTwoUUID | uuid |
    When I send "POST" request to "{{global.api}}/schedules" with json:
      """
{}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias             | path |
      | scheduleThreeUUID | uuid |
    When I send "POST" request to "{{global.api}}/jobs" with json:
  """
{
    "scheduleUUID": "{{scheduleOneUUID}}",
	"chainUUID": "{{besu.UUID}}",
    "type": "eth://ethereum/transaction",
    "transaction": {
        "from": "{{account3}}",
        "to": "{{to1}}",
        "value": "100000",
        "gasPrice": "1000000000",
        "gas": "21000"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}"
    }
}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias        | path |
      | txOneJobUUID | uuid |
    When I send "POST" request to "{{global.api}}/jobs" with json:
  """
{
    "scheduleUUID": "{{scheduleTwoUUID}}",
	"chainUUID": "{{besu.UUID}}",
    "type": "eth://ethereum/transaction",
    "transaction": {
        "from": "{{account3}}",
        "to": "{{to2}}",
        "value": "100000",
        "nonce": "1",
        "gasPrice": "1000000000",
        "gas": "21000"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}"
    }
}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias        | path |
      | txTwoJobUUID | uuid |
    When I send "POST" request to "{{global.api}}/jobs" with json:
  """
{
    "scheduleUUID": "{{scheduleThreeUUID}}",
	"chainUUID": "{{besu.UUID}}",
    "type": "eth://ethereum/transaction",
    "transaction": {
        "from": "{{account3}}",
        "to": "{{to3}}",
        "value": "100100",
        "nonce": "0",
        "gasPrice": "1000000000",
        "gas": "21000"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}"
    }
}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias          | path |
      | txThreeJobUUID | uuid |
    Then I track the following envelopes
      | ID                    |
      | {{scheduleOneUUID}}   |
      | {{scheduleTwoUUID}}   |
      | {{scheduleThreeUUID}} |
    When I send "PUT" request to "{{global.api}}/jobs/{{txOneJobUUID}}/start"
    Then the response code should be 202
    When I send "PUT" request to "{{global.api}}/jobs/{{txTwoJobUUID}}/start"
    Then the response code should be 202
    When I send "PUT" request to "{{global.api}}/jobs/{{txThreeJobUUID}}/start"
    Then the response code should be 202
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Nonce | To      |
      | 1              | 0     | {{to1}} |
      | 1              | 1     | {{to2}} |
      | 1              | 2     | {{to3}} |
    When I send "GET" request to "{{global.api}}/jobs/{{txThreeJobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status | logs[4].status | logs[5].status |
      | MINED  | CREATED        | STARTED        | PENDING        | RECOVERING     | PENDING        | MINED          |

  @private-tx
  Scenario: Private Transaction, invalid nonce order
    Given I register the following alias
      | alias | value              |
      | to1   | {{random.account}} |
    Then I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.api}}/schedules" with json:
      """
{}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias           | path |
      | scheduleOneUUID | uuid |
    When I send "POST" request to "{{global.api}}/schedules" with json:
      """
{}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias           | path |
      | scheduleTwoUUID | uuid |
    When I send "POST" request to "{{global.api}}/schedules" with json:
      """
{}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias             | path |
      | scheduleThreeUUID | uuid |
    When I send "POST" request to "{{global.api}}/jobs" with json:
  """
{
    "scheduleUUID": "{{scheduleOneUUID}}",
	"chainUUID": "{{besu.UUID}}",
    "type": "eth://orion/eeaTransaction",
    "transaction": {
        "from": "{{account1}}",
        "to": "{{to1}}",
        "nonce": "0",
        "privateFrom": "{{global.nodes.besu_1.privateAddress}}",
        "privateFor": ["{{global.nodes.besu_2.privateAddress}}"]
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}"
    }
}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias        | path |
      | txOneJobUUID | uuid |
    When I send "POST" request to "{{global.api}}/jobs" with json:
  """
{
    "scheduleUUID": "{{scheduleTwoUUID}}",
	"chainUUID": "{{besu.UUID}}",
    "type": "eth://orion/eeaTransaction",
    "transaction": {
        "from": "{{account1}}",
        "to": "{{to1}}",
        "privateFrom": "{{global.nodes.besu_1.privateAddress}}",
        "privateFor": ["{{global.nodes.besu_2.privateAddress}}"]
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}"
    }
}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias        | path |
      | txTwoJobUUID | uuid |
    When I send "POST" request to "{{global.api}}/jobs" with json:
  """
{
    "scheduleUUID": "{{scheduleThreeUUID}}",
	"chainUUID": "{{besu.UUID}}",
    "type": "eth://orion/eeaTransaction",
    "transaction": {
        "from": "{{account1}}",
        "to": "{{to1}}",
        "privateFrom": "{{global.nodes.besu_1.privateAddress}}",
        "privateFor": ["{{global.nodes.besu_2.privateAddress}}"]
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}"
    }
}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias          | path |
      | txThreeJobUUID | uuid |
    Then I track the following envelopes
      | ID                    |
      | {{scheduleOneUUID}}   |
      | {{scheduleTwoUUID}}   |
      | {{scheduleThreeUUID}} |
    When I send "PUT" request to "{{global.api}}/jobs/{{txOneJobUUID}}/start"
    Then the response code should be 202
    When I send "PUT" request to "{{global.api}}/jobs/{{txTwoJobUUID}}/start"
    Then the response code should be 202
    When I send "PUT" request to "{{global.api}}/jobs/{{txThreeJobUUID}}/start"
    Then the response code should be 202
    Given I sleep "5s"
    When I send "GET" request to "{{global.api}}/jobs/{{txOneJobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | transaction.nonce |
      | STORED | CREATED        | STARTED        | STORED         | 0                 |
    When I send "GET" request to "{{global.api}}/jobs/{{txTwoJobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | transaction.nonce |
      | STORED | CREATED        | STARTED        | STORED         | 1                 |
    When I send "GET" request to "{{global.api}}/jobs/{{txThreeJobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | transaction.nonce |
      | STORED | CREATED        | STARTED        | STORED         | 2                 |
