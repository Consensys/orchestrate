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
      | alias  | Name                  | URLs                         | Headers.Authorization    |
      | besu_1 | besu_1-{{scenarioID}} | {{global.nodes.besu_1.URLs}} | Bearer {{tenant1.token}} |
    And I have created the following accounts
      | alias    | ID              | Headers.Authorization    |
      | account1 | {{random.uuid}} | Bearer {{tenant1.token}} |
    Given I sign the following transactions
      | alias     | ID              | Value              | Gas   | To           | privateKey                                   | ChainUUID     | Headers.Authorization    |
      | txFaucet1 | {{random.uuid}} | 100000000000000000 | 21000 | {{account1}} | {{global.nodes.besu_1.fundedPrivateKeys[0]}} | {{besu_1.UUID}} | Bearer {{tenant1.token}} |
    Then I track the following envelopes
      | ID               |
      | {{txFaucet1.ID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
      | Content-Type  | application/json         |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/send-raw" with json:
  """
{
    "chain": "besu_1-{{scenarioID}}",
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

  Scenario: Nonce Too High
    When I send envelopes to topic "tx.signer"
      | ID              | ChainName             | From         | To                                         | Nonce   | GasPrice   | Gas   | Headers.Authorization    |
      | {{random.uuid}} | besu_1-{{scenarioID}} | {{account1}} | 0x0000000000000000000000000000000000000000 | 1000000 | 1000000000 | 21000 | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu_1-{{scenarioID}} | {{account1}} | 0x0000000000000000000000000000000000000000 | 1000001 | 1000000000 | 21000 | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu_1-{{scenarioID}} | {{account1}} | 0x0000000000000000000000000000000000000000 | 1000002 | 1000000000 | 21000 | Bearer {{tenant1.token}} |
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
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have the following fields
      | Nonce |
      | 0     |
      | 1     |
      | 2     |
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"

  Scenario: Nonce Too Low
    Given I have deployed the following contracts
      | alias      | ChainName             | From         | ContractName | MethodSignature | Gas     | Headers.Authorization    |
      | besu-token | besu_1-{{scenarioID}} | {{account1}} | SimpleToken  | constructor()   | 2000000 | Bearer {{tenant1.token}} |
    When I send envelopes to topic "tx.signer"
      | ID              | ChainName             | From         | To                                         | Nonce | GasPrice   | Gas   | Headers.Authorization    |
      | {{random.uuid}} | besu_1-{{scenarioID}} | {{account1}} | 0x0000000000000000000000000000000000000000 | 0     | 1000000000 | 21000 | Bearer {{tenant1.token}} |
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have the following fields
      | Nonce |
      | 0     |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have the following fields
      | Nonce |
      | 1     |
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"

  Scenario: Chaotic nonce
    # Next deployment purpose is to increase account nonce to at least 1
    Given I have deployed the following contracts
      | alias      | ChainName             | From         | ContractName | MethodSignature | Gas     | Headers.Authorization    |
      | besu-token | besu_1-{{scenarioID}} | {{account1}} | SimpleToken  | constructor()   | 2000000 | Bearer {{tenant1.token}} |
    When I send envelopes to topic "tx.signer"
      | ID              | ChainName             | From         | To                                         | Nonce   | GasPrice   | Gas   | Headers.Authorization    |
      | {{random.uuid}} | besu_1-{{scenarioID}} | {{account1}} | 0x0000000000000000000000000000000000000000 | 1000002 | 1000000000 | 21000 | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu_1-{{scenarioID}} | {{account1}} | 0x0000000000000000000000000000000000000000 | 1000000 | 1000000000 | 21000 | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu_1-{{scenarioID}} | {{account1}} | 0x0000000000000000000000000000000000000000 | 0       | 1000000000 | 21000 | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu_1-{{scenarioID}} | {{account1}} | 0x0000000000000000000000000000000000000000 | 1000001 | 1000000000 | 21000 | Bearer {{tenant1.token}} |
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have the following fields
      | Nonce   |
      | 1000002 |
      | 1000000 |
      | 0       |
      | 1000001 |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have the following fields
      | Nonce |
      | 1     |
      | 2     |
      | 3     |
      | 4     |
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"

  @private-tx
  Scenario: Nonce too high with private transaction
    # TODO: Able to parse enums like METHOD - shoud be able to pass ETH_SENDRAWTRANSACTION instead of 1
    When I send envelopes to topic "tx.signer"
      | ID              | ChainName             | From         | To                                         | Nonce   | GasPrice | Gas   | PrivateFor                                 | PrivateFrom                            | Method | Headers.Authorization    |
      | {{random.uuid}} | besu_1-{{scenarioID}} | {{account1}} | 0x0000000000000000000000000000000000000000 | 1000000 | 0        | 30000 | ["{{global.nodes.besu_2.privateAddress}}"] | {{global.nodes.besu_1.privateAddress}} | 3      | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu_1-{{scenarioID}} | {{account1}} | 0x0000000000000000000000000000000000000000 | 1000001 | 0        | 30000 | ["{{global.nodes.besu_2.privateAddress}}"] | {{global.nodes.besu_1.privateAddress}} | 3      | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu_1-{{scenarioID}} | {{account1}} | 0x0000000000000000000000000000000000000000 | 1000002 | 0        | 30000 | ["{{global.nodes.besu_2.privateAddress}}"] | {{global.nodes.besu_1.privateAddress}} | 3      | Bearer {{tenant1.token}} |
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have the following fields
      | Nonce   | Raw |
      | 1000000 | ~   |
      | 1000001 | ~   |
      | 1000002 | ~   |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have the following fields
      | Nonce |
      | 0     |
      | 1     |
      | 2     |
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
