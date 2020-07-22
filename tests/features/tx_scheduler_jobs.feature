@tx-scheduler
Feature: Transaction Scheduler
  As an external developer
  I want to send use transaction scheduler API to interact with the registered chains

  Background:
    Given I have the following tenants
      | alias   | tenantID        |
      | tenant1 | {{random.uuid}} |
    Then I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | Bearer {{tenant1.token}} |
    Then I register the following chains
      | alias | Name                | URLs                         | Headers.Authorization    |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu_1.URLs}} | Bearer {{tenant1.token}} |
    And I have created the following accounts
      | alias    | ID              | Headers.Authorization    |
      | account1 | {{random.uuid}} | Bearer {{tenant1.token}} |
    Given I sign the following transactions
      | alias     | ID              | Value              | Gas   | To           | privateKey                                   | ChainUUID     | Headers.Authorization    |
      | txFaucet1 | {{random.uuid}} | 100000000000000000 | 21000 | {{account1}} | {{global.nodes.besu_1.fundedPrivateKeys[0]}} | {{besu.UUID}} | Bearer {{tenant1.token}} |
    Then I track the following envelopes
      | ID               |
      | {{txFaucet1.ID}} |
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
    Then Envelopes should be in topic "tx.decoded"

  @besu
  Scenario: Execute transfer transaction using jobs, step by step
    Given I register the following alias
      | alias | value              |
      | to1   | {{random.account}} |
      | to2   | {{random.account}} |
    Then  I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
      | X-Tenant-ID   | {{tenant1.tenantID}}     |
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
        "value": "100000"
    }
}
      """
    Then the response code should be 200
    And Response should have the following fields
      | uuid | chainUUID     | transaction.from | transaction.to | status
      | ~    | {{besu.UUID}} | {{account1}}     | {{to1}}        | CREATED
    Then I register the following response fields
      | alias        | path |
      | txOneJobUUID | uuid |
    Then I track the following envelopes
      | ID               |
      | {{txOneJobUUID}} |
    When I send "PATCH" request to "{{global.tx-scheduler}}/jobs/{{txOneJobUUID}}" with json:
      """
{
    "transaction": {
        "from": "{{account1}}",
        "to": "{{to2}}",
        "value": "100000"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}"
    }
}
      """
    Then the response code should be 200
    When I send "PUT" request to "{{global.tx-scheduler}}/jobs/{{txOneJobUUID}}/start"
    Then the response code should be 202
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{txOneJobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | uuid             | transaction.from | transaction.to | status  |
      | {{txOneJobUUID}} | {{account1}}     | {{to2}}        | STARTED |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have the following fields
      | Nonce | From         | To      |
      | 1     | {{account1}} | {{to2}} |
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have the following fields
      | Raw | TxHash |
      | ~   | ~      |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status |
      | 1              |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{txOneJobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | MINED  | CREATED        | STARTED        | PENDING        | MINED          |

  @besu
  Scenario: Execute raw transaction using jobs, step by step
    Given I sign the following transactions
      | alias | ID              | Value              | Gas   | To           | privateKey                                   | ChainUUID     | Headers.Authorization    |
      | rawTx | {{random.uuid}} | 100000000000000000 | 21000 | {{account1}} | {{global.nodes.besu_1.fundedPrivateKeys[0]}} | {{besu.UUID}} | Bearer {{tenant1.token}} |
    Then  I set the headers
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
	"type": "eth://ethereum/rawTransaction",
    "transaction": {
        "raw": "{{rawTx.Raw}}"
    }
}
      """
    Then the response code should be 200
    And Response should have the following fields
      | uuid | chainUUID     | status
      | ~    | {{besu.UUID}} | CREATED
    Then I register the following response fields
      | alias        | path |
      | txTwoJobUUID | uuid |
    Then I track the following envelopes
      | ID               |
      | {{txTwoJobUUID}} |
    Then the response code should be 200
    When I send "PUT" request to "{{global.tx-scheduler}}/jobs/{{txTwoJobUUID}}/start"
    Then the response code should be 202
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have the following fields
      | Raw | To           |
      | ~   | {{account1}} |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status |
      | 1              |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{txTwoJobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | MINED  | CREATED        | STARTED        | PENDING        | MINED          |
