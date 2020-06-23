@public-tx
Feature: Invalid Nonce
  As an external developer
  I want transaction with invalid nonce to be recovered, sent to blockchain and mined

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
      | besu-faucet-{{scenarioID}} | {{besu.UUID}} | {{global.nodes.besu.fundedAccounts[7]}} | 1000000000000000000 | 1000000000000000000 | 1s       | Bearer {{tenant1.token}} |
    And I have created the following accounts
      | alias    | ID              | ChainName           | ContextLabels.faucetChildTxID | Headers.Authorization    |
      | account1 | {{random.uuid}} | besu-{{scenarioID}} | {{random.uuid}}               | Bearer {{tenant1.token}} |

  Scenario: Nonce Too High
    When I send envelopes to topic "tx.signer"
      | ID              | ChainName           | From         | To                                         | Nonce   | GasPrice   | Gas   | Headers.Authorization    |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account1}} | 0x0000000000000000000000000000000000000000 | 1000000 | 1000000000 | 21000 | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account1}} | 0x0000000000000000000000000000000000000000 | 1000001 | 1000000000 | 21000 | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account1}} | 0x0000000000000000000000000000000000000000 | 1000002 | 1000000000 | 21000 | Bearer {{tenant1.token}} |
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
      | alias      | ChainName           | From         | ContractName | MethodSignature | Gas     | Headers.Authorization    |
      | besu-token | besu-{{scenarioID}} | {{account1}} | SimpleToken  | constructor()   | 2000000 | Bearer {{tenant1.token}} |
    When I send envelopes to topic "tx.signer"
      | ID              | ChainName           | From         | To                                         | Nonce | GasPrice   | Gas   | Headers.Authorization    |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account1}} | 0x0000000000000000000000000000000000000000 | 0     | 1000000000 | 21000 | Bearer {{tenant1.token}} |
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
      | alias      | ChainName           | From         | ContractName | MethodSignature | Gas     | Headers.Authorization    |
      | besu-token | besu-{{scenarioID}} | {{account1}} | SimpleToken  | constructor()   | 2000000 | Bearer {{tenant1.token}} |
    When I send envelopes to topic "tx.signer"
      | ID              | ChainName           | From         | To                                         | Nonce   | GasPrice   | Gas   | Headers.Authorization    |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account1}} | 0x0000000000000000000000000000000000000000 | 1000002 | 1000000000 | 21000 | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account1}} | 0x0000000000000000000000000000000000000000 | 1000000 | 1000000000 | 21000 | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account1}} | 0x0000000000000000000000000000000000000000 | 0       | 1000000000 | 21000 | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account1}} | 0x0000000000000000000000000000000000000000 | 1000001 | 1000000000 | 21000 | Bearer {{tenant1.token}} |
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
      | ID              | ChainName           | From         | To                                         | Nonce   | GasPrice | Gas   | PrivateFor                                       | PrivateFrom                                  | Method | Headers.Authorization    |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account1}} | 0x0000000000000000000000000000000000000000 | 1000000 | 0        | 30000 | ["A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="] | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | 3      | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account1}} | 0x0000000000000000000000000000000000000000 | 1000001 | 0        | 30000 | ["A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="] | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | 3      | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account1}} | 0x0000000000000000000000000000000000000000 | 1000002 | 0        | 30000 | ["A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="] | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | 3      | Bearer {{tenant1.token}} |
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
