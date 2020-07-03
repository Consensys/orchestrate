@public-tx
Feature: Multiple transactions
  As an external developer
  I want to process multiple transactions

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
      | geth  | geth-{{scenarioID}} | {{global.nodes.geth.URLs}} | Bearer {{tenant1.token}} |
    And I register the following faucets
      | Name                       | ChainRule     | CreditorAccount                         | MaxBalance          | Amount              | Cooldown | Headers.Authorization    |
      | besu-faucet-{{scenarioID}} | {{besu.UUID}} | {{global.nodes.besu.fundedAccounts[0]}} | 1000000000000000000 | 1000000000000000000 | 1s       | Bearer {{tenant1.token}} |
      | geth-faucet-{{scenarioID}} | {{geth.UUID}} | {{global.nodes.geth.fundedAccounts[0]}} | 1000000000000000000 | 1000000000000000000 | 1s       | Bearer {{tenant1.token}} |
    And I have created the following accounts
      | alias    | ID              | ChainName           | ContextLabels.faucetChildTxID | Headers.Authorization    |
      | account1 | {{random.uuid}} | besu-{{scenarioID}} | {{random.uuid}}               | Bearer {{tenant1.token}} |
      | account2 | {{random.uuid}} | geth-{{scenarioID}} | {{random.uuid}}               | Bearer {{tenant1.token}} |
    And I have deployed the following contracts
      | alias      | ChainName           | From         | ContractName | MethodSignature | Gas     | Headers.Authorization    |
      | token-besu | besu-{{scenarioID}} | {{account1}} | SimpleToken  | constructor()   | 2000000 | Bearer {{tenant1.token}} |
      | token-geth | geth-{{scenarioID}} | {{account2}} | SimpleToken  | constructor()   | 3000000 | Bearer {{tenant1.token}} |

  Scenario: Send contract transactions
    When I send envelopes to topic "tx.crafter"
      | ID              | ChainName           | From         | To             | MethodSignature           | Gas    | Args                                                 | Headers.Authorization    |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account1}} | {{token-besu}} | transfer(address,uint256) |        | ["0xdbb881a51CD4023E4400CEF3ef73046743f08da3","1"]   | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account1}} | {{token-besu}} | transfer(address,uint256) |        | ["0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff","0x2"] | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account1}} | {{token-besu}} | transfer(address,uint256) | 100000 | ["0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff","0x8"] | Bearer {{tenant1.token}} |
      | {{random.uuid}} | geth-{{scenarioID}} | {{account2}} | {{token-geth}} | transfer(address,uint256) |        | ["0xdbb881a51CD4023E4400CEF3ef73046743f08da3","0x1"] | Bearer {{tenant1.token}} |
      | {{random.uuid}} | geth-{{scenarioID}} | {{account2}} | {{token-geth}} | transfer(address,uint256) |        | ["0x6009608a02a7a15fd6689d6dad560c44e9ab61ff","0x2"] | Bearer {{tenant1.token}} |
      | {{random.uuid}} | geth-{{scenarioID}} | {{account2}} | {{token-geth}} | transfer(address,uint256) |        | ["0x6009608a02a7a15fd6689d6dad560c44e9ab61ff","2"]   | Bearer {{tenant1.token}} |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have the following fields
      | Nonce | Data | Gas    | GasPrice | From         |
      | ~     | ~    | ~      | ~        | {{account1}} |
      | ~     | ~    | ~      | ~        | {{account1}} |
      | ~     | ~    | 100000 | ~        | {{account1}} |
      | ~     | ~    | ~      | ~        | {{account2}} |
      | ~     | ~    | ~      | ~        | {{account2}} |
      | ~     | ~    | ~      | ~        | {{account2}} |
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have the following fields
      | Raw | TxHash |
      | ~   | ~      |
      | ~   | ~      |
      | ~   | ~      |
      | ~   | ~      |
      | ~   | ~      |
      | ~   | ~      |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Logs[0].Event             | Receipt.Logs[0].DecodedData.from | Receipt.Logs[0].DecodedData.to             | Receipt.Logs[0].DecodedData.value |
      | 1              | Transfer(address,address,uint256) | {{account1}}                     | 0xdbb881a51CD4023E4400CEF3ef73046743f08da3 | 1                                 |
      | 1              | Transfer(address,address,uint256) | {{account1}}                     | 0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff | 2                                 |
      | 1              | Transfer(address,address,uint256) | {{account1}}                     | 0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff | 8                                 |
      | 1              | Transfer(address,address,uint256) | {{account2}}                     | 0xdbb881a51CD4023E4400CEF3ef73046743f08da3 | 1                                 |
      | 1              | Transfer(address,address,uint256) | {{account2}}                     | 0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff | 2                                 |
      | 1              | Transfer(address,address,uint256) | {{account2}}                     | 0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff | 2                                 |

  Scenario: Send transfer transaction
    When I send envelopes to topic "tx.crafter"
      | ID              | ChainName           | From         | To                                         | Value               | Headers.Authorization    |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account1}} | 0x0000000000000000000000000000000000000000 | 1000000000000000000 | Bearer {{tenant1.token}} |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have the following fields
      | Value               | From         |
      | 1000000000000000000 | {{account1}} |
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have the following fields
      | Value               | From         |
      | 1000000000000000000 | {{account1}} |
    And Envelopes should have the following fields
      | Raw | TxHash |
      | ~   | ~      |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status |
      | 1              |

 # TODO: this scenario cannot run twice on the same network (sending twice the same transaction)
 #  Scenario: Send raw transactions
  # Send raw transaction with address:0xf5956Eb46b377Ae41b41BDa94e6270208d8202bb, nonce:0, gasLimit:21000, to:0x88a5C2d9919e46F883EB62F7b8Dd9d0CC45bc290, value:0,5ETH
   # When I send envelopes to topic "tx.sender"
   #   | chainName | contextLabels.txMode | tenantID                             | raw                                                                                                                                                                                                                  |
   #   | besu      | raw                  | {{tenant1.token}}| 0xf86780808252089488a5c2d9919e46f883eb62f7b8dd9d0cc45bc2908806f05b59d3b20000801ca09fd94be4942219541b1fd100341706e2e4caa365c926cde48d8c7aac8c5a0f69a034a00bd00f4ef680a3586208c462280d99cb35e8a89479494af09a7228fdd46a |
   # Then Envelopes should be in topic "tx.decoded"
