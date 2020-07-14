@tx-scheduler
Feature: Transaction Scheduler
  As an external developer
  I want to send use transaction scheduler API to interact with the registered chains

  Background:
    Given I have the following tenants
      | alias   | tenantID        |
      | tenant1 | {{random.uuid}} |

  @besu
  Scenario: Send contract transaction and start a job
    # Prepare Orchestrate and Blockchain state
    And I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | Bearer {{tenant1.token}} |
    And I register the following chains
      | alias | Name                | URLs                       | Headers.Authorization    |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu.URLs}} | Bearer {{tenant1.token}} |
    # Create new account and fund it
    And I have created the following accounts
      | alias    | ID              | Headers.Authorization    |
      | account1 | {{random.uuid}} | Bearer {{tenant1.token}} |
    Given I sign the following transactions
      | Value | Gas   | To           | privateKey                                 | ChainUUID     | Headers.Authorization    | alias |
      | 10000 | 21000 | {{account1}} | {{global.nodes.besu.fundedPrivateKeys[0]}} | {{besu.UUID}} | Bearer {{tenant1.token}} | tx1   |
    Given I register the following alias
      | alias      | value           |
      | faucetTxID | {{random.uuid}} |
      | sendTxID   | {{random.uuid}} |
    Then I track the following envelopes
      | ID             |
      | {{faucetTxID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/send-raw" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "raw": "{{tx1.Raw}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{faucetTxID}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    # Start scenario
    And I have deployed the following contracts
      | alias      | ChainName           | From         | ContractName | MethodSignature | Gas     | Headers.Authorization    |
      | token-besu | besu-{{scenarioID}} | {{account1}} | SimpleToken  | constructor()   | 2000000 | Bearer {{tenant1.token}} |
    Then I track the following envelopes
      | ID           |
      | {{sendTxID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/send" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "from": "{{account1}}",
        "to": "{{token-besu}}",
        "methodSignature": "transfer(address,uint256)",
        "args": ["0xdbb881a51CD4023E4400CEF3ef73046743f08da3",1]
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{sendTxID}}"
    }
}
      """
    Then the response code should be 202
    And Response should have the following fields
      | params.methodSignature    | schedule.uuid | schedule.jobs[0].uuid | schedule.jobs[0].status |
      | transfer(address,uint256) | ~             | ~                     | STARTED                 |
    Then I register the following response fields
      | alias   | path                  |
      | jobUUID | schedule.jobs[0].uuid |
    Then Envelopes should be in topic "tx.crafter"
    And Envelopes should have the following fields
      | ID          | Nonce |
      | {{jobUUID}} | -     |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status  |
      | STARTED |
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have the following fields
      | ID          | Nonce | Data | Gas | GasPrice | From         | To             |
      | {{jobUUID}} | 1     | ~    | ~   | ~        | {{account1}} | {{token-besu}} |
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have the following fields
      | Raw | TxHash |
      | ~   | ~      |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Logs[0].Event             | Receipt.Logs[0].DecodedData.from | Receipt.Logs[0].DecodedData.to             | Receipt.Logs[0].DecodedData.value |
      | 1              | Transfer(address,address,uint256) | {{account1}}                     | 0xdbb881a51CD4023E4400CEF3ef73046743f08da3 | 1                                 |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | MINED  | CREATED        | STARTED        | PENDING        | MINED          |

  @besu
  Scenario: Send contract transaction with unknown from address
    And I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | Bearer {{tenant1.token}} |
    And I register the following chains
      | alias | Name                | URLs                       | Headers.Authorization    |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu.URLs}} | Bearer {{tenant1.token}} |
    Given I register the following alias
      | alias   | value           |
      | labelID | {{random.uuid}} |
    Then I track the following envelopes
      | ID          |
      | {{labelID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "contractName": "SimpleToken",
        "from": "0x931D387731bBbC988B312206c74F77D004D6B84b"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{labelID}}"
    }
}
      """
    Then the response code should be 202
    And Response should have the following fields
      | schedule.uuid | schedule.jobs[0].uuid | schedule.jobs[0].status |
      | ~             | ~                     | STARTED                 |
    Then I register the following response fields
      | alias   | path                  |
      | jobUUID | schedule.jobs[0].uuid |
    Then Envelopes should be in topic "tx.crafter"
    And Envelopes should have the following fields
      | ID          |
      | {{jobUUID}} |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status  |
      | STARTED |
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have the following fields
      | ID          | Nonce | Data | Gas | GasPrice | From                                       |
      | {{jobUUID}} | ~     | ~    | ~   | ~        | 0x931D387731bBbC988B312206c74F77D004D6B84b |
    Then Envelopes should be in topic "tx.recover"
    And Envelopes should have the following fields
      | Errors[0].Message                                               |
      | no key for account "0x931D387731bBbC988B312206c74F77D004D6B84b" |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status |
      | FAILED | CREATED        | STARTED        | FAILED         |

  @besu
  Scenario: New JOB started step by step
    And I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | Bearer {{tenant1.token}} |
    And I register the following chains
      | alias | Name                | URLs                       | Headers.Authorization    |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu.URLs}} | Bearer {{tenant1.token}} |
    And I have created the following accounts
      | alias    | ID              | Headers.Authorization    |
      | account1 | {{random.uuid}} | Bearer {{tenant1.token}} |
    And I have deployed the following contracts
      | alias      | ChainName           | From         | ContractName | MethodSignature | Gas     | Headers.Authorization    |
      | token-besu | besu-{{scenarioID}} | {{account1}} | SimpleToken  | constructor()   | 2000000 | Bearer {{tenant1.token}} |
    Given I register the following alias
      | alias   | value           |
      | labelID | {{random.uuid}} |
    Then I track the following envelopes
      | ID          |
      | {{labelID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
      | X-Tenant-ID   | {{tenant1.tenantID}}     |
    When I send "POST" request to "{{global.tx-scheduler}}/schedules" with json:
      """
{}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias        | path |
      | scheduleUUID | uuid |
    When I send "POST" request to "{{global.tx-scheduler}}/jobs" with json:
      """
{
	"scheduleUUID": "{{scheduleUUID}}",
	"chainUUID": "{{besu.UUID}}",
	"type": "ETH_SENDRAWTRANSACTION",
    "transaction": {
        "from": "0x93f7274c9059e601be4512f656b57b830e019e41",
        "to": "0x93f7274c9059e601be4512f656b57b830e019e42"
    }
}
      """
    Then the response code should be 200
    And Response should have the following fields
      | uuid | chainUUID     | transaction.from                           | transaction.to                             | status
      | ~    | {{besu.UUID}} | 0x93f7274c9059e601be4512f656b57b830e019e41 | 0x93f7274c9059e601be4512f656b57b830e019e42 | CREATED
    Then I register the following response fields
      | alias   | path |
      | jobUUID | uuid |
    # data corresponds to a "transfer(address,uint256)" as methodSignature and ["0xdbb881a51CD4023E4400CEF3ef73046743f08da3",2] as args
    When I send "PATCH" request to "{{global.tx-scheduler}}/jobs/{{jobUUID}}" with json:
      """
{
    "transaction": {
        "from": "{{account1}}",
        "to": "{{token-besu}}",
        "data": "0xa9059cbb000000000000000000000000dbb881a51cd4023e4400cef3ef73046743f08da30000000000000000000000000000000000000000000000000000000000000002"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{labelID}}"
    }
}
      """
    Then the response code should be 200
    When I send "PUT" request to "{{global.tx-scheduler}}/jobs/{{jobUUID}}/start"
    Then the response code should be 202
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | uuid        | transaction.from | transaction.to | status  |
      | {{jobUUID}} | {{account1}}     | {{token-besu}} | STARTED |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have the following fields
      | Nonce | Data                                                                                                                                       | From         | To             |
      | 1     | 0xa9059cbb000000000000000000000000dbb881a51cd4023e4400cef3ef73046743f08da30000000000000000000000000000000000000000000000000000000000000002 | {{account1}} | {{token-besu}} |
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have the following fields
      | Raw | TxHash |
      | ~   | ~      |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Logs[0].Event             | Receipt.Logs[0].DecodedData.from | Receipt.Logs[0].DecodedData.to             | Receipt.Logs[0].DecodedData.value |
      | 1              | Transfer(address,address,uint256) | {{account1}}                     | 0xdbb881a51CD4023E4400CEF3ef73046743f08da3 | 2                                 |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | MINED  | CREATED        | STARTED        | PENDING        | MINED          |

  @besu
  Scenario: Send deploy contract transaction
    And I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | Bearer {{tenant1.token}} |
    And I register the following chains
      | alias | Name                | URLs                       | Headers.Authorization    |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu.URLs}} | Bearer {{tenant1.token}} |
    # Create new account and fund it
    And I have created the following accounts
      | alias    | ID              | Headers.Authorization    |
      | account1 | {{random.uuid}} | Bearer {{tenant1.token}} |
    Given I sign the following transactions
      | alias | Value | Gas   | To           | privateKey                                 | ChainUUID     | Headers.Authorization    |
      | tx1   | 10000 | 21000 | {{account1}} | {{global.nodes.besu.fundedPrivateKeys[0]}} | {{besu.UUID}} | Bearer {{tenant1.token}} |
    Given I register the following alias
      | alias              | value           |
      | faucetTxID         | {{random.uuid}} |
      | deployContractTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID             |
      | {{faucetTxID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/send-raw" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "raw": "{{tx1.Raw}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{faucetTxID}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    Then I track the following envelopes
      | ID                     |
      | {{deployContractTxID}} |
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
    	"id": "{{deployContractTxID}}"
    }
}
      """
    Then the response code should be 202
    And Response should have the following fields
      | schedule.uuid | schedule.jobs[0].uuid | schedule.jobs[0].status |
      | ~             | ~                     | STARTED                 |
    Then I register the following response fields
      | alias   | path                  |
      | jobUUID | schedule.jobs[0].uuid |
    Then Envelopes should be in topic "tx.crafter"
    And Envelopes should have the following fields
      | ID          |
      | {{jobUUID}} |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status  |
      | STARTED |
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have the following fields
      | ID          | Nonce | Data | Gas | GasPrice | From         |
      | {{jobUUID}} | ~     | ~    | ~   | ~        | {{account1}} |
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have the following fields
      | Raw | TxHash |
      | ~   | ~      |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Logs[0].Event             |
      | 1              | Transfer(address,address,uint256) |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | MINED  | CREATED        | STARTED        | PENDING        | MINED          |

  @besu
  Scenario: Send transfer transaction
    And I register the following chains
      | alias | Name                | URLs                       | Headers.Authorization    |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu.URLs}} | Bearer {{tenant1.token}} |
        # Create new account and fund it
    And I have created the following accounts
      | alias    | ID              | Headers.Authorization    |
      | account1 | {{random.uuid}} | Bearer {{tenant1.token}} |
    Given I sign the following transactions
      | alias | Value   | Gas   | To           | privateKey                                 | ChainUUID     | Headers.Authorization    |
      | tx1   | 1000000 | 21000 | {{account1}} | {{global.nodes.besu.fundedPrivateKeys[0]}} | {{besu.UUID}} | Bearer {{tenant1.token}} |
    Given I register the following alias
      | alias         | value              |
      | faucetTxID    | {{random.uuid}}    |
      | transfertTxID | {{random.uuid}}    |
      | recipient     | {{random.account}} |
    Then I track the following envelopes
      | ID             |
      | {{faucetTxID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/send-raw" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "raw": "{{tx1.Raw}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{faucetTxID}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    Then I track the following envelopes
      | ID                |
      | {{transfertTxID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/transfer" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "from": "{{account1}}",
        "to": "{{recipient}}",
        "value": "12345"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{transfertTxID}}"
    }
}
      """
    Then the response code should be 202
    And Response should have the following fields
      | schedule.uuid | schedule.jobs[0].uuid | schedule.jobs[0].status |
      | ~             | ~                     | STARTED                 |
    Then I register the following response fields
      | alias   | path                  |
      | jobUUID | schedule.jobs[0].uuid |
    Then Envelopes should be in topic "tx.crafter"
    And Envelopes should have the following fields
      | ID          |
      | {{jobUUID}} |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status  |
      | STARTED |
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have the following fields
      | ID          | Nonce | Data | Gas   | GasPrice | From         |
      | {{jobUUID}} | ~     | -    | 21000 | ~        | {{account1}} |
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have the following fields
      | Raw | TxHash |
      | ~   | ~      |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status |
      | 1              |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | MINED  | CREATED        | STARTED        | PENDING        | MINED          |
    When I send "POST" request to "{{global.chain-registry}}/{{besu.UUID}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBalance",
        "params": [
          "{{recipient}}",
          "latest"
        ],
        "id": 1
      }
      """
    Then the response code should be 200
    And Response should have the following fields
      | result |
      | 0x3039 |

  @besu
  Scenario: Send raw transaction
    Given I have the following tenants
      | alias   |
      | tenant1 |
    And I register the following chains
      | alias | Name                | URLs                       | Headers.Authorization    |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu.URLs}} | Bearer {{tenant1.token}} |
    Given I register the following alias
      | alias     | value              |
      | labelID   | {{random.uuid}}    |
      | recipient | {{random.account}} |
    Then I track the following envelopes
      | ID          |
      | {{labelID}} |
    Given I sign the following transactions
      | alias | Value | Gas   | To            | privateKey                                 | ChainUUID     | Headers.Authorization    |
      | tx1   | 10000 | 21000 | {{recipient}} | {{global.nodes.besu.fundedPrivateKeys[0]}} | {{besu.UUID}} | Bearer {{tenant1.token}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/send-raw" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "raw": "{{tx1.Raw}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{labelID}}"
    }
}
      """
    Then the response code should be 202
    And Response should have the following fields
      | schedule.uuid | schedule.jobs[0].uuid | schedule.jobs[0].status |
      | ~             | ~                     | STARTED                 |
    Then I register the following response fields
      | alias   | path                  |
      | jobUUID | schedule.jobs[0].uuid |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status |
      | 1              |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | MINED  | CREATED        | STARTED        | PENDING        | MINED          |
    When I send "POST" request to "{{global.chain-registry}}/{{besu.UUID}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBalance",
        "params": [
          "{{recipient}}",
          "latest"
        ],
        "id": 1
      }
      """
    Then the response code should be 200
    And Response should have the following fields
      | result |
      | 0x2710 |

  @besu
  @private-tx
  Scenario: Send a private Orion deploy contract transaction
    And I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | Bearer {{tenant1.token}} |
    Given I register the following chains
      | alias | Name                | URLs                       | Headers.Authorization    |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu.URLs}} | Bearer {{tenant1.token}} |
    And I have created the following accounts
      | alias    | ID              | Headers.Authorization    |
      | account1 | {{random.uuid}} | Bearer {{tenant1.token}} |
    Given I register the following alias
      | alias              | value           |
      | deployContractTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                     |
      | {{deployContractTxID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
      | Content-Type  | application/json         |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "from": "{{account1}}",
        "protocol": "Orion",
        "privateFrom": "Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs=",
        "privateFor": ["k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8="],
        "contractName": "SimpleToken"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{deployContractTxID}}"
    }
}
      """
    Then the response code should be 202
    And Response should have the following fields
      | schedule.uuid | schedule.jobs[0].uuid | schedule.jobs[0].status |
      | ~             | ~                     | STARTED                 |
    Then I register the following response fields
      | alias   | path                  |
      | jobUUID | schedule.jobs[0].uuid |
    Then Envelopes should be in topic "tx.crafter"
    And Envelopes should have the following fields
      | ID          |
      | {{jobUUID}} |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status  |
      | STARTED |
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have the following fields
      | ID          | Nonce | Data | Gas | GasPrice | From         |
      | {{jobUUID}} | ~     | ~    | ~   | ~        | {{account1}} |
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have the following fields
      | Raw | TxHash |
      | ~   | ~      |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Logs[0].Event             |
      | 1              | Transfer(address,address,uint256) |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[3].status |
      | MINED  | CREATED        | STARTED        | MINED          |

  @quorum
  @private-tx
  Scenario: Send a private Tessera deploy contract transaction
    And I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | Bearer {{tenant1.token}} |
    Given I register the following chains
      | alias  | Name                  | URLs                         | PrivateTxManagers                         | Headers.Authorization    |
      | quorum | quorum-{{scenarioID}} | {{global.nodes.quorum.URLs}} | {{global.nodes.quorum.PrivateTxManagers}} | Bearer {{tenant1.token}} |
    And I have created the following accounts
      | alias    | ID              | Headers.Authorization    |
      | account1 | {{random.uuid}} | Bearer {{tenant1.token}} |
    Given I register the following alias
      | alias              | value           |
      | deployContractTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                     |
      | {{deployContractTxID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
      | Content-Type  | application/json         |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "quorum-{{scenarioID}}",
    "params": {
        "from": "{{account1}}",
        "protocol": "Tessera",
        "privateFrom": "BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3Bo=",
        "privateFor": ["QfeDAys9MPDs2XHExtc84jKGHxZg/aj52DTh0vtA3Xc="],
        "contractName": "SimpleToken"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{deployContractTxID}}"
    }
}
      """
    Then the response code should be 202
    And Response should have the following fields
      | schedule.uuid | schedule.jobs[0].uuid | schedule.jobs[0].status |
      | ~             | ~                     | STARTED                 |
    Then I register the following response fields
      | alias   | path                  |
      | jobUUID | schedule.jobs[0].uuid |
    Then Envelopes should be in topic "tx.crafter"
    And Envelopes should have the following fields
      | ID          |
      | {{jobUUID}} |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status  |
      | STARTED |
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have the following fields
      | ID          | Nonce | Data | Gas | GasPrice | From         |
      | {{jobUUID}} | ~     | ~    | ~   | ~        | {{account1}} |
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have the following fields
      | Raw | TxHash |
      | ~   | ~      |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Logs[0].Event             |
      | 1              | Transfer(address,address,uint256) |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status |
      | MINED  | CREATED        | STARTED        | PENDING        | MINED          |
