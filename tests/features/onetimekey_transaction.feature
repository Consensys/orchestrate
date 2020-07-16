@onetimekey-tx
Feature: One time key signature
  As an external developer
  I want to send anonymous transaction using one time usage keys

  Background:
    Given I have the following tenants
      | alias   | tenantID        |
      | tenant1 | {{random.uuid}} |
    And I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | Bearer {{tenant1.token}} |

  @quorum
  @besu
  Scenario: Deploy ERC20
    Given I register the following chains
      | alias    | Name                    | URLs                           | PrivateTxManagers                           | Headers.Authorization    |
      | besu_1   | besu_1-{{scenarioID}}   | {{global.nodes.besu_1.URLs}}   |                                             | Bearer {{tenant1.token}} |
      | quorum_1 | quorum_1-{{scenarioID}} | {{global.nodes.quorum_1.URLs}} | {{global.nodes.quorum_1.PrivateTxManagers}} | Bearer {{tenant1.token}} |
    When I send envelopes to topic "tx.crafter"
      | ID              | ChainName               | ContractName | MethodSignature | Gas     | ContextLabels.txFrom | Headers.Authorization    |
      | {{random.uuid}} | besu_1-{{scenarioID}}   | SimpleToken  | constructor()   | 2000000 | one-time-key         | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu_1-{{scenarioID}}   | SimpleToken  | constructor()   | 2000000 | one-time-key         | Bearer {{tenant1.token}} |
      | {{random.uuid}} | quorum_1-{{scenarioID}} | SimpleToken  | constructor()   | 2000000 | one-time-key         | Bearer {{tenant1.token}} |
      | {{random.uuid}} | quorum_1-{{scenarioID}} | SimpleToken  | constructor()   | 2000000 | one-time-key         | Bearer {{tenant1.token}} |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have the following fields
      | From | Raw | TxHash |
      | ~    | ~   | ~      |
      | ~    | ~   | ~      |
      | ~    | ~   | ~      |
      | ~    | ~   | ~      |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status |
      | 1              |
      | 1              |
      | 1              |
      | 1              |

  @quorum
  Scenario: Deploy private ERC20 contract with Quorum and Tessera
    Given I register the following chains
      | alias  | Name                    | URLs                           | PrivateTxManagers                           | Headers.Authorization    |
      | quorum | quorum_1-{{scenarioID}} | {{global.nodes.quorum_1.URLs}} | {{global.nodes.quorum_1.PrivateTxManagers}} | Bearer {{tenant1.token}} |
    When I send envelopes to topic "tx.crafter"
      | ID              | ChainName               | ContractName | MethodSignature | Gas     | PrivateFor                                   | PrivateFrom                              | Method | ContextLabels.txFrom | Headers.Authorization    |
      | {{random.uuid}} | quorum_1-{{scenarioID}} | SimpleToken  | constructor()   | 2000000 | ["{{global.nodes.quorum_2.privateAddress}}"] | {{global.nodes.quorum_1.privateAddress}} | 2      | one-time-key         | Bearer {{tenant1.token}} |
      | {{random.uuid}} | quorum_1-{{scenarioID}} | SimpleToken  | constructor()   | 2000000 | ["{{global.nodes.quorum_2.privateAddress}}"] | {{global.nodes.quorum_1.privateAddress}} | 2      | one-time-key         | Bearer {{tenant1.token}} |
      | {{random.uuid}} | quorum_1-{{scenarioID}} | SimpleToken  | constructor()   | 2000000 | ["{{global.nodes.quorum_2.privateAddress}}"] | {{global.nodes.quorum_1.privateAddress}} | 2      | one-time-key         | Bearer {{tenant1.token}} |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have the following fields
      | From | Raw | TxHash |
      | ~    | ~   | ~      |
      | ~    | ~   | ~      |
      | ~    | ~   | ~      |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status |
      | 1              |
      | 1              |
      | 1              |

  @besu
  Scenario: Deploy private ERC20 contract with Besu and Orion
    Given I register the following chains
      | alias  | Name                  | URLs                         | Headers.Authorization    |
      | besu_1 | besu_1-{{scenarioID}} | {{global.nodes.besu_1.URLs}} | Bearer {{tenant1.token}} |
    When I send envelopes to topic "tx.crafter"
      | ID              | ChainName             | ContractName | MethodSignature | PrivateFor                                 | PrivateFrom                            | Method | ContextLabels.txFrom | Headers.Authorization    |
      | {{random.uuid}} | besu_1-{{scenarioID}} | SimpleToken  | constructor()   | ["{{global.nodes.besu_2.privateAddress}}"] | {{global.nodes.besu_1.privateAddress}} | 3      | one-time-key         | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu_1-{{scenarioID}} | SimpleToken  | constructor()   | ["{{global.nodes.besu_2.privateAddress}}"] | {{global.nodes.besu_1.privateAddress}} | 3      | one-time-key         | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu_1-{{scenarioID}} | SimpleToken  | constructor()   | ["{{global.nodes.besu_2.privateAddress}}"] | {{global.nodes.besu_1.privateAddress}} | 3      | one-time-key         | Bearer {{tenant1.token}} |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have the following fields
      | From | Raw | TxHash |
      | ~    | ~   | ~      |
      | ~    | ~   | ~      |
      | ~    | ~   | ~      |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status |
      | 1              |
      | 1              |
      | 1              |

