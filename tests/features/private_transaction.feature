@private-tx
Feature: Private transactions
  As an external developer
  I want to deploy a private contract using transaction scheduler

  Background:
    Given I have the following tenants
      | alias   | tenantID        |
      | tenant1 | {{random.uuid}} |
    And I have created the following accounts
      | alias    | API-KEY            | Tenant               |
      | account1 | {{global.api-key}} | {{tenant1.tenantID}} |
    And I register the following contracts
      | name        | artifacts        | API-KEY            | Tenant               |
      | SimpleToken | SimpleToken.json | {{global.api-key}} | {{tenant1.tenantID}} |

  @go-quorum
  Scenario: Deploy private ERC20 contract and send transaction to it with Quorum and Tessera
    Given I register the following alias
      | alias                    | value           |
      | quorumDeployContractTxID | {{random.uuid}} |
      | quorumSentTxContractTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                           |
      | {{quorumDeployContractTxID}} |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.quorum0.Name}}",
        "params": {
          "from": "{{account1}}",
          "protocol": "Tessera",
          "privateFor": [
            "{{global.nodes.quorum[1].privateAddress[0]}}"
          ],
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{quorumDeployContractTxID}}"
        }
      }
      """
    Then the response code should be 202
    Then I register the following response fields
      | alias           | path         |
      | jobPrivTxOne    | jobs[0].uuid |
      | jobMarkingTxOne | jobs[1].uuid |
      | evlpID          | uuid         |
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.ContractAddress |
      | 1              | ~                       |
    And I register the following envelope fields
      | id         | alias               | path                    |
      | {{evlpID}} | counterContractAddr | Receipt.ContractAddress |
    When I send "GET" request to "{{global.api}}/jobs/{{jobPrivTxOne}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status |
      | STORED | CREATED        | STARTED        | STORED         |
    When I send "GET" request to "{{global.api}}/jobs/{{jobMarkingTxOne}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | MINED  | CREATED        | STARTED        | PENDING        | MINED          |
    Then I track the following envelopes
      | ID                           |
      | {{quorumSentTxContractTxID}} |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/send" with json:
      """
      {
        "chain": "{{chain.quorum0.Name}}",
        "params": {
          "from": "{{account1}}",
          "to": "{{counterContractAddr}}",
          "methodSignature": "transfer(address,uint256)",
          "contractName": "SimpleToken",
          "args": [
            "0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff",
            "0x2"
          ],
          "protocol": "Tessera",
          "privateFrom": "{{global.nodes.quorum[0].privateAddress[0]}}",
          "privateFor": [
            "{{global.nodes.quorum[1].privateAddress[0]}}"
          ]
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{quorumSentTxContractTxID}}"
        }
      }
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"

  @go-quorum
  Scenario: Deploy private ERC20 contract with Quorum and Tessera using privacy enhancement
    Given I register the following alias
      | alias                       | value           |
      | quorumDeployContractTxIDTwo | {{random.uuid}} |
    Then I track the following envelopes
      | ID                              |
      | {{quorumDeployContractTxIDTwo}} |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.quorum0.Name}}",
        "params": {
          "from": "{{account1}}",
          "protocol": "Tessera",
          "privateFrom": "{{global.nodes.quorum[0].privateAddress[0]}}",
          "privateFor": [
            "{{global.nodes.quorum[1].privateAddress[0]}}","{{global.nodes.quorum[2].privateAddress[0]}}"
          ],
          "mandatoryFor": [
            "{{global.nodes.quorum[1].privateAddress[0]}}"
          ],
          "privacyFlag": 1,
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{quorumDeployContractTxIDTwo}}"
        }
      }
      """
    Then the response code should be 202
    Then I register the following response fields
      | alias           | path         |
      | jobPrivTxOne    | jobs[0].uuid |
      | jobMarkingTxOne | jobs[1].uuid |
      | evlpID          | uuid         |
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.ContractAddress |
      | 1              | ~                       |

  @go-quorum
  Scenario: Fail to deploy private ERC20 contract with unknown ChainName
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "UnknownChain",
        "params": {
          "from": "{{account1}}",
          "protocol": "Tessera",
          "privateFrom": "{{global.nodes.quorum[0].privateAddress[0]}}",
          "privateFor": [
            "{{global.nodes.quorum[1].privateAddress[0]}}"
          ],
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}"
        }
      }
      """
    Then the response code should be 422
    And Response should have the following fields
      | code   | message |
      | 271360 | ~       |

  @go-quorum
  Scenario: Fail to deploy private ERC20 contract with invalid private from
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.quorum0.Name}}",
        "params": {
          "from": "{{account1}}",
          "protocol": "Tessera",
          "privateFrom": "{{global.nodes.quorum[1].privateAddress[0]}}",
          "privateFor": [
            "{{global.nodes.quorum[0].privateAddress[0]}}"
          ],
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}"
        }
      }
      """
    Then the response code should be 202
    Then I register the following response fields
      | alias           | path         |
      | jobPrivTxOne    | jobs[0].uuid |
      | jobMarkingTxOne | jobs[1].uuid |
      | evlpID          | uuid         |
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.recover"

  @go-quorum
  Scenario: Force not correlative nonce for private and public txs in Quorum/Tessera
    Given I register the following alias
      | alias                           | value           |
      | publicQuorumDeployContractTxID  | {{random.uuid}} |
      | privateQuorumDeployContractTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                                 |
      | {{publicQuorumDeployContractTxID}} |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.quorum0.Name}}",
        "params": {
          "from": "{{account1}}",
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{publicQuorumDeployContractTxID}}"
        }
      }
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    Then I track the following envelopes
      | ID                                  |
      | {{privateQuorumDeployContractTxID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.quorum0.Name}}",
        "params": {
          "from": "{{account1}}",
          "protocol": "Tessera",
          "privateFrom": "{{global.nodes.quorum[0].privateAddress[0]}}",
          "privateFor": [
            "{{global.nodes.quorum[1].privateAddress[0]}}"
          ],
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{privateQuorumDeployContractTxID}}"
        }
      }
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.ContractAddress |
      | 1              | ~                       |

  @go-quorum
  Scenario: Fail to deploy private ERC20 contract with unknown PrivateFor
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.quorum0.Name}}",
        "params": {
          "from": "{{account1}}",
          "protocol": "Tessera",
          "privateFrom": "{{global.nodes.quorum[0].privateAddress[0]}}",
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}"
        }
      }
      """
    Then the response code should be 400
    And Response should have the following fields
      | code   | message |
      | 271104 | ~       |

  @go-quorum
  Scenario: Fail to deploy private ERC20 contract with not authorized chain
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "quorum_2-{{scenarioID}}",
        "params": {
          "from": "{{account1}}",
          "protocol": "Tessera",
          "privateFrom": "{{global.nodes.quorum[0].privateAddress[0]}}",
          "privateFor": [
            "{{global.nodes.quorum[1].privateAddress[0]}}"
          ],
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}"
        }
      }
      """
    Then the response code should be 422
    And Response should have the following fields
      | code   | message |
      | 271360 | ~       |

  @besu
  Scenario: Deploy private ERC20 contract with Besu/EEA
    Given I register the following alias
      | alias                  | value           |
      | besuDeployContractTxID | {{random.uuid}} |
      | besuSentTxContractTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                         |
      | {{besuDeployContractTxID}} |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{account1}}",
          "protocol": "EEA",
          "privateFrom": "{{global.nodes.besu[0].privateAddress[0]}}",
          "privateFor": [
            "{{global.nodes.besu[1].privateAddress[0]}}"
          ],
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{besuDeployContractTxID}}"
        }
      }
      """
    Then the response code should be 202
    Then I register the following response fields
      | alias           | path         |
      | jobPrivTxTwo    | jobs[0].uuid |
      | jobMarkingTxTwo | jobs[1].uuid |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Output | Receipt.ContractAddress | Receipt.PrivateFrom                        | Receipt.PrivateFor                             |
      | 1              | ~              | ~                       | {{global.nodes.besu[0].privateAddress[0]}} | ["{{global.nodes.besu[1].privateAddress[0]}}"] |
    And I register the following envelope fields
      | id                         | alias               | path                    |
      | {{besuDeployContractTxID}} | counterContractAddr | Receipt.ContractAddress |
    When I send "GET" request to "{{global.api}}/jobs/{{jobPrivTxTwo}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | transaction.from |
      | STORED | CREATED        | STARTED        | STORED         | {{account1}}     |
    When I send "GET" request to "{{global.api}}/jobs/{{jobMarkingTxTwo}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status | annotations.oneTimeKey |
      | MINED  | CREATED        | STARTED        | PENDING        | MINED          | true                   |
    Then I track the following envelopes
      | ID                         |
      | {{besuSentTxContractTxID}} |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/send" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{account1}}",
          "to": "{{counterContractAddr}}",
          "methodSignature": "transfer(address,uint256)",
          "contractName": "SimpleToken",
          "args": [
            "0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff",
            "0x2"
          ],
          "protocol": "EEA",
          "privateFrom": "{{global.nodes.besu[0].privateAddress[0]}}",
          "privateFor": [
            "{{global.nodes.besu[1].privateAddress[0]}}"
          ]
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{besuSentTxContractTxID}}"
        }
      }
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status |
      | 1              |

  @besu
  Scenario: Batch deploy private ERC20 contract with Besu/EEA with different PrivateFor
    Given I register the following alias
      | alias                       | value           |
      | besuDeployContractTxOneID   | {{random.uuid}} |
      | besuDeployContractTxTwoID   | {{random.uuid}} |
      | besuDeployContractTxThreeID | {{random.uuid}} |
      | besuDeployContractTxFourID  | {{random.uuid}} |
    And I have created the following accounts
      | alias    | API-KEY            | Tenant               |
      | account2 | {{global.api-key}} | {{tenant1.tenantID}} |
      | account3 | {{global.api-key}} | {{tenant1.tenantID}} |
      | account4 | {{global.api-key}} | {{tenant1.tenantID}} |
    Then I track the following envelopes
      | ID                              |
      | {{besuDeployContractTxOneID}}   |
      | {{besuDeployContractTxTwoID}}   |
      | {{besuDeployContractTxThreeID}} |
      | {{besuDeployContractTxFourID}}  |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{account1}}",
          "protocol": "EEA",
          "privateFrom": "{{global.nodes.besu[0].privateAddress[0]}}",
          "privateFor": [
            "{{global.nodes.besu[1].privateAddress[0]}}"
          ],
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{besuDeployContractTxOneID}}"
        }
      }
      """
    Then the response code should be 202
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{account2}}",
          "protocol": "EEA",
          "privateFrom": "{{global.nodes.besu[0].privateAddress[0]}}",
          "privateFor": [
            "{{global.nodes.besu[2].privateAddress[0]}}"
          ],
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{besuDeployContractTxTwoID}}"
        }
      }
      """
    Then the response code should be 202
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{account3}}",
          "protocol": "EEA",
          "privateFrom": "{{global.nodes.besu[0].privateAddress[0]}}",
          "privateFor": [
            "{{global.nodes.besu[1].privateAddress[0]}}"
          ],
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{besuDeployContractTxThreeID}}"
        }
      }
      """
    Then the response code should be 202
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{account4}}",
          "protocol": "EEA",
          "privateFrom": "{{global.nodes.besu[0].privateAddress[0]}}",
          "privateFor": [
            "{{global.nodes.besu[2].privateAddress[0]}}"
          ],
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{besuDeployContractTxFourID}}"
        }
      }
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Output | Receipt.PrivateFrom                        | Receipt.PrivateFor                             |
      | 1              | ~              | {{global.nodes.besu[0].privateAddress[0]}} | ["{{global.nodes.besu[1].privateAddress[0]}}"] |
      | 1              | ~              | {{global.nodes.besu[0].privateAddress[0]}} | ["{{global.nodes.besu[2].privateAddress[0]}}"] |
      | 1              | ~              | {{global.nodes.besu[0].privateAddress[0]}} | ["{{global.nodes.besu[1].privateAddress[0]}}"] |
      | 1              | ~              | {{global.nodes.besu[0].privateAddress[0]}} | ["{{global.nodes.besu[2].privateAddress[0]}}"] |

  @besu
  Scenario: Batch deploy private ERC20 contract with Bes/EEA with different PrivateFrom
    Given I register the following alias
      | alias                       | value           |
      | besuDeployContractTxOneID   | {{random.uuid}} |
      | besuDeployContractTxTwoID   | {{random.uuid}} |
      | besuDeployContractTxThreeID | {{random.uuid}} |
      | besuDeployContractTxFourID  | {{random.uuid}} |
    And I have created the following accounts
      | alias    | API-KEY            | Tenant               |
      | account2 | {{global.api-key}} | {{tenant1.tenantID}} |
      | account3 | {{global.api-key}} | {{tenant1.tenantID}} |
      | account4 | {{global.api-key}} | {{tenant1.tenantID}} |
    Then I track the following envelopes
      | ID                              |
      | {{besuDeployContractTxOneID}}   |
      | {{besuDeployContractTxTwoID}}   |
      | {{besuDeployContractTxThreeID}} |
      | {{besuDeployContractTxFourID}}  |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{account1}}",
          "protocol": "EEA",
          "privateFrom": "{{global.nodes.besu[0].privateAddress[0]}}",
          "privateFor": [
            "{{global.nodes.besu[2].privateAddress[0]}}"
          ],
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{besuDeployContractTxOneID}}"
        }
      }
      """
    Then the response code should be 202
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu1.Name}}",
        "params": {
          "from": "{{account2}}",
          "protocol": "EEA",
          "privateFrom": "{{global.nodes.besu[1].privateAddress[0]}}",
          "privateFor": [
            "{{global.nodes.besu[2].privateAddress[0]}}"
          ],
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{besuDeployContractTxTwoID}}"
        }
      }
      """
    Then the response code should be 202
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{account3}}",
          "protocol": "EEA",
          "privateFrom": "{{global.nodes.besu[0].privateAddress[0]}}",
          "privateFor": [
            "{{global.nodes.besu[2].privateAddress[0]}}"
          ],
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{besuDeployContractTxThreeID}}"
        }
      }
      """
    Then the response code should be 202
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu1.Name}}",
        "params": {
          "from": "{{account4}}",
          "protocol": "EEA",
          "privateFrom": "{{global.nodes.besu[1].privateAddress[0]}}",
          "privateFor": [
            "{{global.nodes.besu[2].privateAddress[0]}}"
          ],
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{besuDeployContractTxFourID}}"
        }
      }
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Output | Receipt.PrivateFrom                        | Receipt.PrivateFor                             |
      | 1              | ~              | {{global.nodes.besu[0].privateAddress[0]}} | ["{{global.nodes.besu[2].privateAddress[0]}}"] |
      | 1              | ~              | {{global.nodes.besu[1].privateAddress[0]}} | ["{{global.nodes.besu[2].privateAddress[0]}}"] |
      | 1              | ~              | {{global.nodes.besu[0].privateAddress[0]}} | ["{{global.nodes.besu[2].privateAddress[0]}}"] |
      | 1              | ~              | {{global.nodes.besu[1].privateAddress[0]}} | ["{{global.nodes.besu[2].privateAddress[0]}}"] |

  @besu
  Scenario: Deploy private ERC20 for a privacy group
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    Given I sleep "2s"
    When I send "POST" request to "{{global.api}}/proxy/chains/{{chain.besu0.UUID}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "priv_createPrivacyGroup",
        "params": [
          {
            "addresses": [
              "{{global.nodes.besu[0].privateAddress[0]}}",
              "{{global.nodes.besu[1].privateAddress[0]}}",
              "{{global.nodes.besu[2].privateAddress[0]}}"
            ],
            "name": "TestGroup",
            "description": "TestGroup"
          }
        ],
        "id": 1
      }
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias          | path   |
      | privacyGroupId | result |
    Given I register the following alias
      | alias                       | value           |
      | besuDeployContractTxGroupID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                              |
      | {{besuDeployContractTxGroupID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{account1}}",
          "protocol": "EEA",
          "privateFrom": "{{global.nodes.besu[0].privateAddress[0]}}",
          "privacyGroupId": "{{privacyGroupId}}",
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{besuDeployContractTxGroupID}}"
        }
      }
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Output | Receipt.PrivateFrom                        | Receipt.PrivacyGroupId |
      | 1              | ~              | {{global.nodes.besu[0].privateAddress[0]}} | {{privacyGroupId}}     |


  @besu
  Scenario: Deploy private ERC20 for a privacy group
    Given I register the following alias
      | alias               | value           |
      | privacyGroupNameOne | {{random.uuid}} |
    Then I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    Then I sleep "2s"
    When I send "POST" request to "{{global.api}}/proxy/chains/{{chain.besu0.UUID}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "priv_createPrivacyGroup",
        "params": [
          {
            "addresses": [
              "{{global.nodes.besu[0].privateAddress[0]}}",
              "{{global.nodes.besu[1].privateAddress[0]}}",
              "{{global.nodes.besu[2].privateAddress[0]}}"
            ],
            "name": "{{privacyGroupNameOne}}",
            "description": "TestGroup"
          }
        ],
        "id": 1
      }
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias          | path   |
      | privacyGroupId | result |
    Given I register the following alias
      | alias                       | value           |
      | besuDeployContractTxGroupID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                              |
      | {{besuDeployContractTxGroupID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{account1}}",
          "protocol": "EEA",
          "privateFrom": "{{global.nodes.besu[0].privateAddress[0]}}",
          "privacyGroupId": "{{privacyGroupId}}",
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{besuDeployContractTxGroupID}}"
        }
      }
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Output | Receipt.ContractAddress | Receipt.PrivateFrom                        | Receipt.PrivacyGroupId |
      | 1              | ~              | ~                       | {{global.nodes.besu[0].privateAddress[0]}} | {{privacyGroupId}}     |

  @besu
  Scenario: Fail to deploy private ERC20 with a privacy group and privateFor
    Given I register the following alias
      | alias               | value           |
      | privacyGroupNameTwo | {{random.uuid}} |
    Then I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    Then I sleep "2s"
    When I send "POST" request to "{{global.api}}/proxy/chains/{{chain.besu0.UUID}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "priv_createPrivacyGroup",
        "params": [
          {
            "addresses": [
              "{{global.nodes.besu[0].privateAddress[0]}}",
              "{{global.nodes.besu[1].privateAddress[0]}}",
              "{{global.nodes.besu[2].privateAddress[0]}}"
            ],
            "name": "{{privacyGroupNameTwo}}",
            "description": "TestGroup"
          }
        ],
        "id": 1
      }
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias          | path   |
      | privacyGroupId | result |
    Given I register the following alias
      | alias                       | value           |
      | besuDeployContractTxGroupID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                              |
      | {{besuDeployContractTxGroupID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{account1}}",
          "protocol": "EEA",
          "privateFrom": "{{global.nodes.besu[0].privateAddress[0]}}",
          "privacyGroupId": "{{privacyGroupId}}",
          "privateFor": [
            "{{global.nodes.besu[2].privateAddress[0]}}"
          ],
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{besuDeployContractTxGroupID}}"
        }
      }
      """
    Then the response code should be 400
    And Response should have the following fields
      | code   | message |
      | 271104 | ~       |

  @besu
  @oneTimeKey
  Scenario: Deploy private ERC20 with one-time-key
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    Given I register the following alias
      | alias                     | value           |
      | besuDeployContractTxOTKID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                            |
      | {{besuDeployContractTxOTKID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "oneTimeKey": true,
          "protocol": "EEA",
          "privateFrom": "{{global.nodes.besu[0].privateAddress[0]}}",
          "privateFor": [
            "{{global.nodes.besu[2].privateAddress[0]}}"
          ],
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{besuDeployContractTxOTKID}}"
        }
      }
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Output | Receipt.ContractAddress | Receipt.PrivateFrom                        |
      | 1              | ~              | ~                       | {{global.nodes.besu[0].privateAddress[0]}} |

  @besu
  Scenario: Force not correlative nonce for private and public txs
    Given I register the following alias
      | alias                         | value           |
      | publicBesuDeployContractTxID  | {{random.uuid}} |
      | privateBesuDeployContractTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                               |
      | {{publicBesuDeployContractTxID}} |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
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
          "id": "{{publicBesuDeployContractTxID}}"
        }
      }
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.ContractAddress |
      | 1              | ~                       |
    Then I track the following envelopes
      | ID                                |
      | {{privateBesuDeployContractTxID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "from": "{{account1}}",
          "protocol": "EEA",
          "privateFrom": "{{global.nodes.besu[0].privateAddress[0]}}",
          "privateFor": [
            "{{global.nodes.besu[1].privateAddress[0]}}",
            "{{global.nodes.besu[2].privateAddress[0]}}"
          ],
          "contractName": "SimpleToken"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{privateBesuDeployContractTxID}}"
        }
      }
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Output | Receipt.ContractAddress | Receipt.PrivateFrom                        | Receipt.PrivateFor                             |
      | 1              | ~              | ~                       | {{global.nodes.besu[0].privateAddress[0]}} | ["{{global.nodes.besu[1].privateAddress[0]}}"] |

  @besu
  Scenario: Private Transaction using Job and too high nonce
    Given I register the following alias
      | alias | value              |
      | to1   | {{random.account}} |
    Then I set the headers
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
    When I send "POST" request to "{{global.api}}/jobs" with json:
      """
      {
        "scheduleUUID": "{{scheduleUUID}}",
        "chainUUID": "{{chain.besu0.UUID}}",
        "type": "eth://eea/markingTransaction",
        "transaction": {
          "from": "{{account1}}"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}"
        }
      }
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias              | path |
      | txMarkingTxJobUUID | uuid |
    When I send "POST" request to "{{global.api}}/jobs" with json:
      """
      {
        "scheduleUUID": "{{scheduleUUID}}",
        "chainUUID": "{{chain.besu0.UUID}}",
        "nextJobUUID": "{{txMarkingTxJobUUID}}",
        "type": "eth://eea/privateTransaction",
        "transaction": {
          "from": "{{account1}}",
          "to": "{{to1}}",
          "nonce": 1000001,
          "privateFrom": "{{global.nodes.besu[0].privateAddress[0]}}",
          "privateFor": [
            "{{global.nodes.besu[1].privateAddress[0]}}"
          ]
        },
        "labels": {
          "scenario.id": "{{scenarioID}}"
        }
      }
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias         | path |
      | txPrivJobUUID | uuid |
    Then I track the following envelopes
      | ID               |
      | {{scheduleUUID}} |
    When I send "PUT" request to "{{global.api}}/jobs/{{txPrivJobUUID}}/start"
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status |
      | 2              |
