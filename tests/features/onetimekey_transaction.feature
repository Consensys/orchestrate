@onetimekey-tx
Feature: One time key signature
  As an external developer
  I want to send anonymous transaction using one time usage keys

  Background:
    Given I have the following tenants
      | alias   | tenantID                             |
      | tenant1 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    And I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | Bearer {{tenant1.token}} |

  @quorum
  @besu
  Scenario: Deploy ERC20
    Given I register the following chains
      | alias  | Name                  | URLs                         | PrivateTxManagers                         | Headers.Authorization    |
      | besu   | besu-{{scenarioID}}   | {{global.nodes.besu.URLs}}   |                                           | Bearer {{tenant1.token}} |
      | quorum | quorum-{{scenarioID}} | {{global.nodes.quorum.URLs}} | {{global.nodes.quorum.PrivateTxManagers}} | Bearer {{tenant1.token}} |
    When I send envelopes to topic "tx.crafter"
      | ID              | ChainName             | ContractName | MethodSignature | Gas     | ContextLabels.txFrom | Headers.Authorization    |
      | {{random.uuid}} | besu-{{scenarioID}}   | SimpleToken  | constructor()   | 2000000 | one-time-key         | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu-{{scenarioID}}   | SimpleToken  | constructor()   | 2000000 | one-time-key         | Bearer {{tenant1.token}} |
      | {{random.uuid}} | quorum-{{scenarioID}} | SimpleToken  | constructor()   | 2000000 | one-time-key         | Bearer {{tenant1.token}} |
      | {{random.uuid}} | quorum-{{scenarioID}} | SimpleToken  | constructor()   | 2000000 | one-time-key         | Bearer {{tenant1.token}} |
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
      | alias  | Name                  | URLs                         | PrivateTxManagers                         | Headers.Authorization    |
      | quorum | quorum-{{scenarioID}} | {{global.nodes.quorum.URLs}} | {{global.nodes.quorum.PrivateTxManagers}} | Bearer {{tenant1.token}} |
    When I send envelopes to topic "tx.crafter"
      | ID              | ChainName             | ContractName | MethodSignature | Gas     | PrivateFor                                       | PrivateFrom                                  | Method | ContextLabels.txFrom | Headers.Authorization    |
      | {{random.uuid}} | quorum-{{scenarioID}} | SimpleToken  | constructor()   | 2000000 | ["QfeDAys9MPDs2XHExtc84jKGHxZg/aj52DTh0vtA3Xc="] | BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3Bo= | 2      | one-time-key         | Bearer {{tenant1.token}} |
      | {{random.uuid}} | quorum-{{scenarioID}} | SimpleToken  | constructor()   | 2000000 | ["QfeDAys9MPDs2XHExtc84jKGHxZg/aj52DTh0vtA3Xc="] | BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3Bo= | 2      | one-time-key         | Bearer {{tenant1.token}} |
      | {{random.uuid}} | quorum-{{scenarioID}} | SimpleToken  | constructor()   | 2000000 | ["QfeDAys9MPDs2XHExtc84jKGHxZg/aj52DTh0vtA3Xc="] | BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3Bo= | 2      | one-time-key         | Bearer {{tenant1.token}} |
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
      | alias | Name                | URLs                       | Headers.Authorization    |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu.URLs}} | Bearer {{tenant1.token}} |
    When I send envelopes to topic "tx.crafter"
      | ID              | ChainName           | ContractName | MethodSignature | PrivateFor                                       | PrivateFrom                                  | Method | ContextLabels.txFrom | Headers.Authorization    |
      | {{random.uuid}} | besu-{{scenarioID}} | SimpleToken  | constructor()   | ["A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="] | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | 3      | one-time-key         | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu-{{scenarioID}} | SimpleToken  | constructor()   | ["A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="] | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | 3      | one-time-key         | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu-{{scenarioID}} | SimpleToken  | constructor()   | ["A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="] | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | 3      | one-time-key         | Bearer {{tenant1.token}} |
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

