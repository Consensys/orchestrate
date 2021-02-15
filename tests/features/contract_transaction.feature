@public-tx
Feature: Send contract transactions
  As an external developer
  I want to send multiple contract transactions using transaction-scheduler API

  Background:
    Given I have the following tenants
      | alias   | tenantID        |
      | tenant1 | {{random.uuid}} |
    And I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | Bearer {{tenant1.token}} |
      | Counter     | Counter.json     | Bearer {{tenant1.token}} |
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
    When I send "POST" request to "{{global.api}}/transactions/transfer" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{global.nodes.besu[0].fundedPublicKeys[0]}}",
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
        "chain": "{{chain.geth0.Name}}",
        "params": {
          "from": "{{global.nodes.geth[0].fundedPublicKeys[0]}}",
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
    And Envelopes should have the following fields
      | Receipt.Status |
      | 1              |
      | 1              |
    Given I register the following alias
      | alias            | value           |
      | besuContractTxID | {{random.uuid}} |
      | gethContractTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                   |
      | {{besuContractTxID}} |
      | {{gethContractTxID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
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
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.geth0.Name}}",
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
    Then Envelopes should be in topic "tx.decoded"
    And I register the following envelope fields
      | id                   | alias            | path                    |
      | {{besuContractTxID}} | besuContractAddr | Receipt.ContractAddress |
      | {{gethContractTxID}} | gethContractAddr | Receipt.ContractAddress |

  @geth
  Scenario: Send contract transactions
    Given I register the following alias
      | alias             | value           |
      | besuSendTxOneID   | {{random.uuid}} |
      | besuSendTxTwoID   | {{random.uuid}} |
      | besuSendTxThreeID | {{random.uuid}} |
      | gethSendTxOneID   | {{random.uuid}} |
      | gethSendTxTwoID   | {{random.uuid}} |
      | gethSendTxThreeID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                    |
      | {{besuSendTxOneID}}   |
      | {{besuSendTxTwoID}}   |
      | {{besuSendTxThreeID}} |
      | {{gethSendTxOneID}}   |
      | {{gethSendTxTwoID}}   |
      | {{gethSendTxThreeID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.api}}/transactions/send" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{account1}}",
          "to": "{{besuContractAddr}}",
          "methodSignature": "transfer(address,uint256)",
          "args": [
            "0xdbb881a51CD4023E4400CEF3ef73046743f08da3",
            "1"
          ]
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{besuSendTxOneID}}"
        }
      }
      """
    Then the response code should be 202
    When I send "POST" request to "{{global.api}}/transactions/send" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{account1}}",
          "to": "{{besuContractAddr}}",
          "methodSignature": "transfer(address,uint256)",
          "args": [
            "0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff",
            "0x2"
          ]
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{besuSendTxTwoID}}"
        }
      }
      """
    Then the response code should be 202
    When I send "POST" request to "{{global.api}}/transactions/send" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{account1}}",
          "to": "{{besuContractAddr}}",
          "methodSignature": "transfer(address,uint256)",
          "gas": "100000",
          "args": [
            "0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff",
            "0x8"
          ]
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{besuSendTxThreeID}}"
        }
      }
      """
    Then the response code should be 202
    When I send "POST" request to "{{global.api}}/transactions/send" with json:
      """
      {
        "chain": "{{chain.geth0.Name}}",
        "params": {
          "from": "{{account2}}",
          "to": "{{gethContractAddr}}",
          "methodSignature": "transfer(address,uint256)",
          "args": [
            "0xdbb881a51CD4023E4400CEF3ef73046743f08da3",
            "0x1"
          ]
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{gethSendTxOneID}}"
        }
      }
      """
    Then the response code should be 202
    When I send "POST" request to "{{global.api}}/transactions/send" with json:
      """
      {
        "chain": "{{chain.geth0.Name}}",
        "params": {
          "from": "{{account2}}",
          "to": "{{gethContractAddr}}",
          "methodSignature": "transfer(address,uint256)",
          "args": [
            "0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff",
            "0x2"
          ]
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{gethSendTxTwoID}}"
        }
      }
      """
    Then the response code should be 202
    When I send "POST" request to "{{global.api}}/transactions/send" with json:
      """
      {
        "chain": "{{chain.geth0.Name}}",
        "params": {
          "from": "{{account2}}",
          "to": "{{gethContractAddr}}",
          "methodSignature": "transfer(address,uint256)",
          "gas": "100000",
          "args": [
            "0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff",
            "2"
          ]
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{gethSendTxThreeID}}"
        }
      }
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Logs[0].Event             | Receipt.Logs[0].DecodedData.from | Receipt.Logs[0].DecodedData.to             | Receipt.Logs[0].DecodedData.value |
      | 1              | Transfer(address,address,uint256) | {{account1}}                     | 0xdbb881a51CD4023E4400CEF3ef73046743f08da3 | 1                                 |
      | 1              | Transfer(address,address,uint256) | {{account1}}                     | 0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff | 2                                 |
      | 1              | Transfer(address,address,uint256) | {{account1}}                     | 0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff | 8                                 |
      | 1              | Transfer(address,address,uint256) | {{account2}}                     | 0xdbb881a51CD4023E4400CEF3ef73046743f08da3 | 1                                 |
      | 1              | Transfer(address,address,uint256) | {{account2}}                     | 0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff | 2                                 |
      | 1              | Transfer(address,address,uint256) | {{account2}}                     | 0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff | 2                                 |

  @oneTimeKey
  Scenario: Send contract transactions with one-time-key
    Given I register the following alias
      | alias             | value           |
      | counterDeployTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                    |
      | {{counterDeployTxID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "oneTimeKey": true,
          "contractName": "Counter"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{counterDeployTxID}}"
        }
      }
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    And I register the following envelope fields
      | id                    | alias               | path                    |
      | {{counterDeployTxID}} | counterContractAddr | Receipt.ContractAddress |
    Given I register the following alias
      | alias       | value           |
      | sendOTKTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID              |
      | {{sendOTKTxID}} |
    When I send "POST" request to "{{global.api}}/transactions/send" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "oneTimeKey": true,
          "to": "{{counterContractAddr}}",
          "methodSignature": "increment(uint256)",
          "args": [
            1
          ]
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{sendOTKTxID}}"
        }
      }
      """
    Then the response code should be 202
    Then I register the following response fields
      | alias      | path         |
      | jobOTKUUID | jobs[0].uuid |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Logs[0].Event        | Receipt.Logs[0].DecodedData.from |
      | 1              | Incremented(address,uint256) | ~                                |
    When I send "GET" request to "{{global.api}}/jobs/{{jobOTKUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | MINED  | CREATED        | STARTED        | PENDING        | MINED          |


  Scenario: Fail to send contract transactions with invalid args
    Given I register the following alias
      | alias             | value           |
      | counterDeployTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                    |
      | {{counterDeployTxID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{account1}}",
          "contractName": "Counter"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{counterDeployTxID}}"
        }
      }
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    And I register the following envelope fields
      | id                    | alias               | path                    |
      | {{counterDeployTxID}} | counterContractAddr | Receipt.ContractAddress |
    Given I register the following alias
      | alias    | value           |
      | sendTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID           |
      | {{sendTxID}} |
    When I send "POST" request to "{{global.api}}/transactions/send" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{account1}}",
          "to": "{{counterContractAddr}}",
          "methodSignature": "increment(uint256)",
          "args": [
            "string"
          ]
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{sendTxID}}"
        }
      }
      """
    Then the response code should be 422
    And Response should have the following fields
      | code   | message |
      | 271360 | ~       |

  Scenario: Fail to send contract transactions with invalid methodSignature
    Given I register the following alias
      | alias             | value           |
      | counterDeployTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                    |
      | {{counterDeployTxID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{account1}}",
          "contractName": "Counter"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{counterDeployTxID}}"
        }
      }
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    And I register the following envelope fields
      | id                    | alias               | path                    |
      | {{counterDeployTxID}} | counterContractAddr | Receipt.ContractAddress |
    Given I register the following alias
      | alias    | value           |
      | sendTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID           |
      | {{sendTxID}} |
    When I send "POST" request to "{{global.api}}/transactions/send" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{account1}}",
          "to": "{{counterContractAddr}}",
          "methodSignature": "increment(uint256,uint256,uint256)",
          "args": [
            1,
            2,
            "3"
          ]
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{sendTxID}}"
        }
      }
      """
    Then the response code should be 202
    Then I register the following response fields
      | alias         | path         |
      | jobFailedUUID | jobs[0].uuid |
    Then Envelopes should be in topic "tx.recover"
    When I send "GET" request to "{{global.api}}/jobs/{{jobFailedUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status |
      | FAILED | CREATED        | STARTED        | FAILED         |
