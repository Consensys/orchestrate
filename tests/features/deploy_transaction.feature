@deploy-contract
Feature: Deploy contracts
  As an external developer
  I want to deploy a contract using transaction scheduler API

  Background:
    Given I have the following tenants
      | alias   | tenantID        |
      | tenant1 | {{random.uuid}} |
    And I have created the following accounts
      | alias    | ID              | API-KEY            | Tenant               |
      | account1 | {{random.uuid}} | {{global.api-key}} | {{tenant1.tenantID}} |
    Given I register the following contracts
      | name        | artifacts        | API-KEY            | Tenant               |
      | SimpleToken | SimpleToken.json | {{global.api-key}} | {{tenant1.tenantID}} |
      | ERC20       | ERC20.json       | {{global.api-key}} | {{tenant1.tenantID}} |
    Then I track the following envelopes
      | ID                  |
      | faucet-{{account1}} |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/transfer" with json:
      """
      {
        "chain": "{{chain.geth0.Name}}",
        "params": {
          "from": "{{global.nodes.geth[0].fundedPublicKeys[0]}}",
          "to": "{{account1}}",
          "value": "0x16345785D8A0000"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "faucet-{{account1}}"
        }
      }
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"

  @besu
  Scenario: Deploy SimpleToken in Besu
    Given I register the following alias
      | alias            | value           |
      | besuContractTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                   |
      | {{besuContractTxID}} |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
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
    Then I register the following response fields
      | alias      | path         |
      | jobOneUUID | jobs[0].uuid |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.ContractAddress |
      | 1              | ~                       |
    When I send "GET" request to "{{global.api}}/jobs/{{jobOneUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | MINED  | CREATED        | STARTED        | PENDING        | MINED          |

  @besu
  Scenario: Deploy ERC20 in Besu
    Given I register the following alias
      | alias     | value           |
      | erc20TxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID            |
      | {{erc20TxID}} |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "contractName": "ERC20",
          "from": "{{account1}}",
          "args":["WindToken", "WIND"]
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{erc20TxID}}"
        }
      }
      """
    Then the response code should be 202
    Then I register the following response fields
      | alias      | path         |
      | jobOneUUID | jobs[0].uuid |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.ContractAddress |
      | 1              | ~                       |
    When I send "GET" request to "{{global.api}}/jobs/{{jobOneUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | MINED  | CREATED        | STARTED        | PENDING        | MINED          |

  @geth
  Scenario: Deploy SimpleToken in Geth (dynamic_fee)
    Given I register the following alias
      | alias            | value           |
      | gethContractTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                   |
      | {{gethContractTxID}} |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.geth0.Name}}",
        "params": {
          "contractName": "SimpleToken",
          "from": "{{account1}}"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{gethContractTxID}}"
        }
      }
      """
    Then the response code should be 202
    Then I register the following response fields
      | alias   | path         |
      | jobUUID | jobs[0].uuid |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.ContractAddress | Receipt.ContractName | Receipt.ContractTag |
      | 1              | ~                       | SimpleToken          | latest              |
    When I send "GET" request to "{{global.api}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status | transaction.hash |
      | MINED  | CREATED        | STARTED        | PENDING        | MINED          | ~                |

  @geth
  Scenario: Deploy SimpleToken in Geth (legacy)
    Given I register the following alias
      | alias            | value           |
      | gethContractTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                   |
      | {{gethContractTxID}} |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.geth0.Name}}",
        "params": {
          "contractName": "SimpleToken",
          "from": "{{account1}}",
          "transactionType": "legacy"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "{{gethContractTxID}}"
        }
      }
      """
    Then the response code should be 202
    Then I register the following response fields
      | alias   | path         |
      | jobUUID | jobs[0].uuid |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.ContractAddress |
      | 1              | ~                       |
    When I send "GET" request to "{{global.api}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | MINED  | CREATED        | STARTED        | PENDING        | MINED          |

  @oneTimeKey @besu
  Scenario: Deploy SimpleToken with one-time-key
    Given I register the following contracts
      | name        | artifacts        | API-KEY            | Tenant               |
      | SimpleToken | SimpleToken.json | {{global.api-key}} | {{tenant1.tenantID}} |
    Given I register the following alias
      | alias            | value           |
      | besuContractTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                   |
      | {{besuContractTxID}} |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
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
      | alias      | path         |
      | jobOTKUUID | jobs[0].uuid |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.ContractAddress |
      | 1              | ~                       |
    When I send "GET" request to "{{global.api}}/jobs/{{jobOTKUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | MINED  | CREATED        | STARTED        | PENDING        | MINED          |

  @besu
  Scenario: Fail to deploy SimpleToken with too low gas
    Given I register the following contracts
      | name        | artifacts        | API-KEY            | Tenant               |
      | SimpleToken | SimpleToken.json | {{global.api-key}} | {{tenant1.tenantID}} |
    Given I register the following alias
      | alias            | value           |
      | besuContractTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                   |
      | {{besuContractTxID}} |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
        "params": {
          "contractName": "SimpleToken",
          "oneTimeKey": true,
          "gas": 1
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
      | jobOTKUUID | jobs[0].uuid |
    Then Envelopes should be in topic "tx.recover"
    And Envelopes should have the following fields
      | Errors.0.Message                                        |
      | code: -32003 - message: Intrinsic gas exceeds gas limit |
    When I send "GET" request to "{{global.api}}/jobs/{{jobOTKUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | FAILED | CREATED        | STARTED        | PENDING        | FAILED         |

  @besu
  Scenario: Fail to deploy SimpleToken with invalid contract tag
    Given I register the following contracts
      | name        | artifacts        | API-KEY            | Tenant               |
      | SimpleToken | SimpleToken.json | {{global.api-key}} | {{tenant1.tenantID}} |
    Given I register the following alias
      | alias            | value           |
      | besuContractTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                   |
      | {{besuContractTxID}} |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
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

  @besu
  Scenario: Fail to deploy SimpleToken with missing contractName
    Given I register the following contracts
      | name        | artifacts        | API-KEY            | Tenant               |
      | SimpleToken | SimpleToken.json | {{global.api-key}} | {{tenant1.tenantID}} |
    Given I register the following alias
      | alias            | value           |
      | besuContractTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                   |
      | {{besuContractTxID}} |
    Given I set the headers
      | Key         | Value                |
      | X-API-KEY   | {{global.api-key}}   |
      | X-TENANT-ID | {{tenant1.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
      """
      {
        "chain": "{{chain.besu0.Name}}",
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
