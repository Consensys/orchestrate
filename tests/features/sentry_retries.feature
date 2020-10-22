@tx-sentry
Feature: Deploy ERC20 contract using tx-sentry
  As an external developer
  I want to deploy a contract using tx-sentry retry feature

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
    Then I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | Bearer {{tenant1.token}} |

  Scenario: Deploy ERC20 using retry policy with zero gas increment
    Given I register the following alias
      | alias                | value           |
      | preBesuContractTxID  | {{random.uuid}} |
      | besuContractTxID     | {{random.uuid}} |
      | postBesuContractTxID | {{random.uuid}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    Then I track the following envelopes
      | ID                      |
      | {{preBesuContractTxID}} |
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
    	"id": "{{preBesuContractTxID}}"
    }
}
  """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    Then Set nonce manager records
      | Account      | ChainID          | Nonce |
      | {{account1}} | {{besu.ChainID}} | 1     |
    Then I track the following envelopes
      | ID                   |
      | {{besuContractTxID}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "contractName": "SimpleToken",
        "from": "{{account1}}",
        "gasPricePolicy": {
          "retryPolicy": {
            "interval": "1s"
          }
        }
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{besuContractTxID}}"
    }
}
  """
    Then the response code should be 202
    Then I register the following response fields
      | alias      | path         |
      | jobOneUUID | jobs[0].uuid |
    Then Envelopes should be in topic "tx.sender"
    Then I sleep "5s"
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobOneUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status  | logs[0].status | logs[1].status | logs[2].status | logs[3].status | logs[4].status |
      | PENDING | CREATED        | STARTED        | PENDING        | RESENDING      | FAILED         |
    Then Set nonce manager records
      | Account      | ChainID          | Nonce |
      | {{account1}} | {{besu.ChainID}} | 0     |
    Then I track the following envelopes
      | ID                       |
      | {{postBesuContractTxID}} |
      | {{besuContractTxID}}     |
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
    	"id": "{{postBesuContractTxID}}"
    }
}
  """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobOneUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status |
      | MINED  |

  Scenario: Deploy ERC20 using retry policy with gas increment
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    Then I track the following envelopes
      | ID                  |
      | faucet-{{account2}} |
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
    Then Envelopes should be in topic "tx.decoded"
    Given I register the following alias
      | alias                | value           |
      | preGethContractTxID  | {{random.uuid}} |
      | gethContractTxID     | {{random.uuid}} |
      | postGethContractTxID | {{random.uuid}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    Then I track the following envelopes
      | ID                      |
      | {{preGethContractTxID}} |
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
    	"id": "{{preGethContractTxID}}"
    }
}
  """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    Then Set nonce manager records
      | Account      | ChainID          | Nonce |
      | {{account2}} | {{geth.ChainID}} | 1     |
    Then I track the following envelopes
      | ID                   |
      | {{gethContractTxID}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "geth-{{scenarioID}}",
    "params": {
        "contractName": "SimpleToken",
        "from": "{{account2}}",
        "gasPricePolicy": {
          "retryPolicy": {
            "interval": "1s",
            "increment": 0.15,
            "limit": 0.45
          }
        }
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{gethContractTxID}}"
    }
}
  """
    Then the response code should be 202
    Then I register the following response fields
      | alias        | path         |
      | jobOneUUID   | jobs[0].uuid |
      | scheduleUUID | uuid         |
    Then Envelopes should be in topic "tx.sender"
    Then I sleep "5s"
    When I send "GET" request to "{{global.tx-scheduler}}/schedules/{{scheduleUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | jobs[0].status | jobs[1].status | jobs[2].status | jobs[3].status |
      | PENDING        | PENDING        | PENDING        | PENDING        |
    Then Set nonce manager records
      | Account      | ChainID          | Nonce |
      | {{account2}} | {{geth.ChainID}} | 0     |
    Then I track the following envelopes
      | ID                       |
      | {{postGethContractTxID}} |
      | {{gethContractTxID}}     |
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
    	"id": "{{postGethContractTxID}}"
    }
}
  """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    When I send "GET" request to "{{global.tx-scheduler}}/schedules/{{scheduleUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | jobs[0].status | jobs[1].status | jobs[2].status | jobs[3].status |
      | NEVER_MINED    | NEVER_MINED    | NEVER_MINED    | MINED          |


  Scenario: Deploy ERC20 using retry policy with zero gas increment to retry limit
    Given I register the following alias
      | alias               | value           |
      | preBesuContractTxID | {{random.uuid}} |
      | besuContractTxID    | {{random.uuid}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    Then I track the following envelopes
      | ID                      |
      | {{preBesuContractTxID}} |
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
    	"id": "{{preBesuContractTxID}}"
    }
}
  """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    Then Set nonce manager records
      | Account      | ChainID          | Nonce |
      | {{account1}} | {{besu.ChainID}} | 1     |
    Then I track the following envelopes
      | ID                   |
      | {{besuContractTxID}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "contractName": "SimpleToken",
        "from": "{{account1}}",
        "gasPricePolicy": {
          "retryPolicy": {
            "interval": "1s"
          }
        }
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{besuContractTxID}}"
    }
}
  """
    Then the response code should be 202
    Then I register the following response fields
      | alias      | path         |
      | jobOneUUID | jobs[0].uuid |
    Then Envelopes should be in topic "tx.sender"
    Then I sleep "15s"
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobOneUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status  | logs.length | annotations.hasBeenRetried |
      | PENDING | 23          | true                       |
