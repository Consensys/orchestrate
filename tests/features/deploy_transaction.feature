@public-tx
Feature: Deploy ERC20 contract
  As an external developer
  I want to deploy a contract using transaction scheduler API

  Background:
    Given I have the following tenants
      | alias   | tenantID        |
      | tenant1 | {{random.uuid}} |
    And I register the following chains
      | alias | Name                | URLs                         | Headers.Authorization    |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu_1.URLs}} | Bearer {{tenant1.token}} |
      | geth  | geth-{{scenarioID}} | {{global.nodes.geth.URLs}}   | Bearer {{tenant1.token}} |
    And I have created the following accounts
      | alias    | ID              | Headers.Authorization    |
      | account1 | {{random.uuid}} | Bearer {{tenant1.token}} |
      | account2 | {{random.uuid}} | Bearer {{tenant1.token}} |
    Then I track the following envelopes
      | ID                  |
      | faucet-{{account1}} |
      | faucet-{{account2}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/transfer" with json:
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
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/transfer" with json:
  """
{
    "chain": "geth-{{scenarioID}}",
    "params": {
      "from": "{{global.nodes.geth.fundedPublicKeys[0]}}",
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
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"

  @besu @geth
  Scenario: Deploy ERC20
    Given I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | Bearer {{tenant1.token}} |
    Given I register the following alias
      | alias            | value           |
      | besuContractTxID | {{random.uuid}} |
      | gethContractTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                   |
      | {{besuContractTxID}} |
      | {{gethContractTxID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "contractName": "SimpleToken",
        "from": "{{account1}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{besuContractTxID}}"
    }
}
      """
    Then the response code should be 202
    Then I register the following response fields
      | alias      | path                  |
      | jobOneUUID | jobs[0].uuid |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "geth-{{scenarioID}}",
    "params": {
        "contractName": "SimpleToken",
        "from": "{{account2}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{gethContractTxID}}"
    }
}
      """
    Then the response code should be 202
    Then I register the following response fields
      | alias      | path                  |
      | jobTwoUUID | jobs[0].uuid |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.ContractAddress |
      | 1              | ~                       |
      | 1              | ~                       |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobOneUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | MINED  | CREATED        | STARTED        | PENDING        | MINED          |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobTwoUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | MINED  | CREATED        | STARTED        | PENDING        | MINED          |


  @oneTimeKey
  Scenario: Deploy ERC20 with one-time-key
    Given I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | Bearer {{tenant1.token}} |
    Given I register the following alias
      | alias            | value           |
      | besuContractTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                   |
      | {{besuContractTxID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "contractName": "SimpleToken",
        "oneTimeKey": true
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{besuContractTxID}}"
    }
}
      """
    Then the response code should be 202
    Then I register the following response fields
      | alias      | path                  |
      | jobOTKUUID | jobs[0].uuid |
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.ContractAddress |
      | 1              | ~                       |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobOTKUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | MINED  | CREATED        | STARTED        | PENDING        | MINED          |
    
  Scenario: Fail to deploy ERC20 with too low gas
    Given I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | Bearer {{tenant1.token}} |
    Given I register the following alias
      | alias            | value           |
      | besuContractTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                   |
      | {{besuContractTxID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "contractName": "SimpleToken",
        "oneTimeKey": true,
        "gas": "1"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{besuContractTxID}}"
    }
}
      """
    Then the response code should be 202
    Then I register the following response fields
      | alias      | path                  |
      | jobOTKUUID | jobs[0].uuid |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.recover"
    And Envelopes should have the following fields
      | Errors.0.Message                                        |
      | code: -32003 - message: Intrinsic gas exceeds gas limit |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobOTKUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | FAILED | CREATED        | STARTED        | PENDING        | FAILED         |

  Scenario: Fail to deploy ERC20 with invalid contract tag
    Given I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | Bearer {{tenant1.token}} |
    Given I register the following alias
      | alias            | value           |
      | besuContractTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                   |
      | {{besuContractTxID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "contractName": "SimpleToken",
        "contractTag": "invalid",
        "oneTimeKey": true
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{besuContractTxID}}"
    }
}
      """
    Then the response code should be 422
    And Response should have the following fields
      | code   | message |
      | 271360 | ~       |


  Scenario: Fail to deploy ERC20 with missing contractName
    Given I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | Bearer {{tenant1.token}} |
    Given I register the following alias
      | alias            | value           |
      | besuContractTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                   |
      | {{besuContractTxID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "from": "{{account1}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{besuContractTxID}}"
    }
}
      """
    Then the response code should be 400
    And Response should have the following fields
      | code   | message |
      | 271104 | ~       |
