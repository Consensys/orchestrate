@tx-scheduler
Feature: Transaction Scheduler
  As an external developer
  I want to send use transaction scheduler API to interact with the registered chains

  Background:
    Given I have the following tenants
      | alias   | tenantID                             |
      | tenant1 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    And I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | Bearer {{tenant1.token}} |
    And I register the following chains
      | alias | Name                | URLs                       | Headers.Authorization    |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu.URLs}} | Bearer {{tenant1.token}} |
    And I register the following faucets
      | Name                       | ChainRule     | CreditorAccount                         | MaxBalance          | Amount              | Cooldown | Headers.Authorization    |
      | besu-faucet-{{scenarioID}} | {{besu.UUID}} | {{global.nodes.besu.fundedAccounts[1]}} | 1000000000000000000 | 1000000000000000000 | 1s       | Bearer {{tenant1.token}} |
    And I have created the following accounts
      | alias    | ID              | ChainName           | ContextLabels.faucetChildTxID | Headers.Authorization    |
      | account1 | {{random.uuid}} | besu-{{scenarioID}} | {{random.uuid}}               | Bearer {{tenant1.token}} |
    And I have deployed the following contracts
      | alias      | ChainName           | From         | ContractName | MethodSignature | Gas     | Headers.Authorization    |
      | token-besu | besu-{{scenarioID}} | {{account1}} | SimpleToken  | constructor()   | 2000000 | Bearer {{tenant1.token}} |
    Given I register the following alias
      | alias   | value           |
      | labelID | {{random.uuid}} |
    Then I track the following envelopes
      | ID          |
      | {{labelID}} |

  Scenario: Send contract transaction and start a job
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
    	"id": "{{labelID}}"
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
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status | logs[4].status |
      | MINED  | CREATED        | STARTED        | PENDING        | SENT           | MINED          |

  Scenario: Send contract transaction with unknown from address
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
      | ID          | Nonce | Data | Gas | GasPrice | From                                       | To             |
      | {{jobUUID}} | ~     | ~    | ~   | ~        | 0x931D387731bBbC988B312206c74F77D004D6B84b | {{token-besu}} |
    Then Envelopes should be in topic "tx.recover"
    And Envelopes should have the following fields
      | Errors[0].Message                                               |
      | no key for account "0x931D387731bBbC988B312206c74F77D004D6B84b" |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | status | logs[0].status | logs[1].status | logs[2].status |
      | FAILED | CREATED        | STARTED        | FAILED         |

  Scenario: New JOB started step by step
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
      | status | logs[0].status | logs[1].status | logs[2].status | logs[3].status | logs[4].status |
      | MINED  | CREATED        | STARTED        | PENDING        | SENT           | MINED          |
