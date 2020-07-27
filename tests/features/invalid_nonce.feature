@public-tx
Feature: Invalid Nonce
  As an external developer
  I want transaction with invalid nonce to be recovered, sent to blockchain and mined

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
    Given I sign the following transactions
      | alias     | ID              | Value              | Gas   | To           | privateKey                                   | ChainUUID     | Headers.Authorization    |
      | txFaucet1 | {{random.uuid}} | 100000000000000000 | 21000 | {{account1}} | {{global.nodes.besu_1.fundedPrivateKeys[0]}} | {{besu.UUID}} | Bearer {{tenant1.token}} |
      | txFaucet2 | {{random.uuid}} | 100000000000000000 | 21000 | {{account2}} | {{global.nodes.besu_1.fundedPrivateKeys[0]}} | {{besu.UUID}} | Bearer {{tenant1.token}} |
      | txFaucet3 | {{random.uuid}} | 100000000000000000 | 21000 | {{account3}} | {{global.nodes.besu_1.fundedPrivateKeys[0]}} | {{besu.UUID}} | Bearer {{tenant1.token}} |
    Then I track the following envelopes
      | ID               |
      | {{txFaucet1.ID}} |
      | {{txFaucet2.ID}} |
      | {{txFaucet3.ID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/send-raw" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "raw": "{{txFaucet1.Raw}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{txFaucet1.ID}}"
    }
}
      """
    Then the response code should be 202
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/send-raw" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "raw": "{{txFaucet2.Raw}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{txFaucet2.ID}}"
    }
}
      """
    Then the response code should be 202
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/send-raw" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "raw": "{{txFaucet3.Raw}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{txFaucet3.ID}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"

  Scenario: Nonce Too High
    Given I register the following alias
      | alias | value              |
      | to1   | {{random.account}} |
    Then I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/schedules" with json:
      """
{}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias           | path |
      | scheduleOneUUID | uuid |
    When I send "POST" request to "{{global.tx-scheduler}}/jobs" with json:
  """
{
    "scheduleUUID": "{{scheduleOneUUID}}",
	"chainUUID": "{{besu.UUID}}",
    "type": "eth://ethereum/transaction",
    "transaction": {
        "from": "{{account1}}",
        "to": "{{to1}}",
        "value": "100000",
        "nonce": "1000000",
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
    When I send "POST" request to "{{global.tx-scheduler}}/jobs" with json:
  """
{
    "scheduleUUID": "{{scheduleOneUUID}}",
	"chainUUID": "{{besu.UUID}}",
    "type": "eth://ethereum/transaction",
    "transaction": {
        "from": "{{account1}}",
        "to": "{{to1}}",
        "value": "100000",
        "nonce": "1000001",
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
    When I send "POST" request to "{{global.tx-scheduler}}/jobs" with json:
  """
{
    "scheduleUUID": "{{scheduleOneUUID}}",
	"chainUUID": "{{besu.UUID}}",
    "type": "eth://ethereum/transaction",
    "transaction": {
        "from": "{{account1}}",
        "to": "{{to1}}",
        "value": "100000",
        "nonce": "1000002",
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
      | ID                 |
      | {{txOneJobUUID}}   |
      | {{txTwoJobUUID}}   |
      | {{txThreeJobUUID}} |
    When I send "PUT" request to "{{global.tx-scheduler}}/jobs/{{txOneJobUUID}}/start"
    Then the response code should be 202
    When I send "PUT" request to "{{global.tx-scheduler}}/jobs/{{txTwoJobUUID}}/start"
    Then the response code should be 202
    When I send "PUT" request to "{{global.tx-scheduler}}/jobs/{{txThreeJobUUID}}/start"
    Then the response code should be 202
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have the following fields
      | Nonce   |
      | 1000000 |
      | 1000001 |
      | 1000002 |
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have the following fields
      | Nonce   |
      | 1000000 |
      | 1000001 |
      | 1000002 |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Nonce |
      | 1              | 0     |
      | 1              | 1     |
      | 1              | 2     |

  Scenario: Nonce Too Low
    Given I register the following alias
      | alias | value              |
      | to1   | {{random.account}} |
      | to2   | {{random.account}} |
    Then I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/schedules" with json:
      """
{}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias           | path |
      | scheduleTwoUUID | uuid |
    When I send "POST" request to "{{global.tx-scheduler}}/jobs" with json:
  """
{
    "scheduleUUID": "{{scheduleTwoUUID}}",
	"chainUUID": "{{besu.UUID}}",
    "type": "eth://ethereum/transaction",
    "transaction": {
        "from": "{{account2}}",
        "to": "{{to1}}",
        "value": "10000",
        "nonce": "0",
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
    Then I track the following envelopes
      | ID               |
      | {{txOneJobUUID}} |
    When I send "PUT" request to "{{global.tx-scheduler}}/jobs/{{txOneJobUUID}}/start"
    Then the response code should be 202
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Nonce | To      |
      | 1              | 0     | {{to1}} |
    When I send "POST" request to "{{global.tx-scheduler}}/jobs" with json:
  """
{
    "scheduleUUID": "{{scheduleTwoUUID}}",
	"chainUUID": "{{besu.UUID}}",
    "type": "eth://ethereum/transaction",
    "transaction": {
        "from": "{{account2}}",
        "to": "{{to2}}",
        "value": "10000",
        "nonce": "0",
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
    Then I track the following envelopes
      | ID               |
      | {{txOneJobUUID}} |
    When I send "PUT" request to "{{global.tx-scheduler}}/jobs/{{txOneJobUUID}}/start"
    Then the response code should be 202
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have the following fields
      | Nonce |
      | 0     |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Nonce | To      |
      | 1              | 1     | {{to2}} |


  Scenario: Chaotic nonce
    Given I register the following alias
      | alias | value              |
      | to1   | {{random.account}} |
      | to2   | {{random.account}} |
      | to3   | {{random.account}} |
    Then I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/schedules" with json:
      """
{}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias             | path |
      | scheduleThreeUUID | uuid |
    When I send "POST" request to "{{global.tx-scheduler}}/jobs" with json:
  """
{
    "scheduleUUID": "{{scheduleThreeUUID}}",
	"chainUUID": "{{besu.UUID}}",
    "type": "eth://ethereum/transaction",
    "transaction": {
        "from": "{{account3}}",
        "to": "{{to1}}",
        "value": "100000",
        "nonce": "1002",
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
    When I send "POST" request to "{{global.tx-scheduler}}/jobs" with json:
  """
{
    "scheduleUUID": "{{scheduleThreeUUID}}",
	"chainUUID": "{{besu.UUID}}",
    "type": "eth://ethereum/transaction",
    "transaction": {
        "from": "{{account3}}",
        "to": "{{to2}}",
        "value": "100000",
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
      | alias        | path |
      | txTwoJobUUID | uuid |
    When I send "POST" request to "{{global.tx-scheduler}}/jobs" with json:
  """
{
    "scheduleUUID": "{{scheduleThreeUUID}}",
	"chainUUID": "{{besu.UUID}}",
    "type": "eth://ethereum/transaction",
    "transaction": {
        "from": "{{account3}}",
        "to": "{{to3}}",
        "value": "100000",
        "nonce": "1001",
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
      | ID                 |
      | {{txOneJobUUID}}   |
      | {{txTwoJobUUID}}   |
      | {{txThreeJobUUID}} |
    When I send "PUT" request to "{{global.tx-scheduler}}/jobs/{{txOneJobUUID}}/start"
    Then the response code should be 202
    When I send "PUT" request to "{{global.tx-scheduler}}/jobs/{{txTwoJobUUID}}/start"
    Then the response code should be 202
    When I send "PUT" request to "{{global.tx-scheduler}}/jobs/{{txThreeJobUUID}}/start"
    Then the response code should be 202
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have the following fields
      | Nonce |
      | 1002  |
      | 0     |
      | 1001  |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Nonce | To      |
      | 1              | 1     | {{to1}} |
      | 1              | 2     | {{to2}} |
      | 1              | 3     | {{to3}} |
  
