@private-tx
Feature: Deploy private ERC20 contract
  As an external developer
  I want to deploy a private contract

  Background:
    Given I have the following tenants
      | alias   | tenantID                             |
      | tenant1 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    And I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | Bearer {{tenant1.token}} |

  @quorum
  Scenario: Deploy private ERC20 contract with Quorum and Tessera
    Given I register the following chains
      | alias  | Name                  | URLs                         | PrivateTxManagers                         | Headers.Authorization    |
      | quorum | quorum-{{scenarioID}} | {{global.nodes.quorum.URLs}} | {{global.nodes.quorum.PrivateTxManagers}} | Bearer {{tenant1.token}} |
    And I have created the following accounts
      | alias    | Headers.Authorization    |
      | account1 | Bearer {{tenant1.token}} |
    And I send envelopes to topic "tx.crafter"
      | ID              | ChainName             | From         | ContractName | MethodSignature | Gas     | PrivateFor                                       | PrivateFrom                                  | Method | Headers.Authorization    |
      | {{random.uuid}} | quorum-{{scenarioID}} | {{account1}} | SimpleToken  | constructor()   | 2000000 | ["QfeDAys9MPDs2XHExtc84jKGHxZg/aj52DTh0vtA3Xc="] | BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3Bo= | 2      | Bearer {{tenant1.token}} |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have the following fields
      | InternalLabels.enclaveKey |
      | ~                         |
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status |
      | 1              |

  @quorum
  Scenario: Deploy private ERC20 contract with unknown ChainName
    Given I have created the following accounts
      | alias    | Headers.Authorization    |
      | account1 | Bearer {{tenant1.token}} |
    When I send envelopes to topic "tx.crafter"
      | ID              | ChainName | From         | ContractName | MethodSignature | Gas     | PrivateFor                                       | PrivateFrom                                  | Method | Headers.Authorization    |
      | {{random.uuid}} | unknown   | {{account1}} | SimpleToken  | constructor()   | 2000000 | ["QfeDAys9MPDs2XHExtc84jKGHxZg/aj52DTh0vtA3Xc="] | BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3Bo= | 2      | Bearer {{tenant1.token}} |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.recover"
    And Envelopes should have the following fields
      | Errors[0].Message                |
      | no chain found with name unknown |

  @quorum
  Scenario: Deploy private ERC20 contract with unknown PrivateFrom
    Given I register the following chains
      | alias  | Name                  | URLs                         | PrivateTxManagers                         | Headers.Authorization    |
      | quorum | quorum-{{scenarioID}} | {{global.nodes.quorum.URLs}} | {{global.nodes.quorum.PrivateTxManagers}} | Bearer {{tenant1.token}} |
    And I have created the following accounts
      | alias    | Headers.Authorization    |
      | account1 | Bearer {{tenant1.token}} |
    When I send envelopes to topic "tx.crafter"
      | ID              | ChainName             | From         | ContractName | MethodSignature | Gas     | PrivateFor                                       | PrivateFrom  | Method | Headers.Authorization    |
      | {{random.uuid}} | quorum-{{scenarioID}} | {{account1}} | SimpleToken  | constructor()   | 2000000 | ["QfeDAys9MPDs2XHExtc84jKGHxZg/aj52DTh0vtA3Xc="] | dW5rbm93bg== | 2      | Bearer {{tenant1.token}} |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.recover"
    And Envelopes should have the following fields
      | Errors[0].Message                                                                                                                         |
      | failed to send a request to Tessera enclave: 08200@: request to '{{global.chain-registry}}/tessera/{{quorum.UUID}}/storeraw' failed - 400 |

  @besu
  Scenario: Deploy private ERC20 contract with Besu and Orion
    Given I register the following chains
      | alias | Name                | URLs                       | Headers.Authorization    |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu.URLs}} | Bearer {{tenant1.token}} |
    And I have created the following accounts
      | alias    | Headers.Authorization    |
      | account1 | Bearer {{tenant1.token}} |
    When I send envelopes to topic "tx.crafter"
      | ID              | ChainName           | From         | ContractName | MethodSignature | PrivateFor                                       | PrivateFrom                                  | Method | Headers.Authorization    |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account1}} | SimpleToken  | constructor()   | ["A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="] | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | 3      | Bearer {{tenant1.token}} |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Output | Receipt.PrivateFrom                          | Receipt.PrivateFor                               |
      | 1              | ~              | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | ["A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="] |

  @besu
  Scenario: Batch deploy private ERC20 contract with Besu and Orion with different PrivateFor
    Given I register the following chains
      | alias | Name                | URLs                       | Headers.Authorization    |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu.URLs}} | Bearer {{tenant1.token}} |
    And I have created the following accounts
      | alias    | Headers.Authorization    |
      | account1 | Bearer {{tenant1.token}} |
      | account2 | Bearer {{tenant1.token}} |
      | account3 | Bearer {{tenant1.token}} |
      | account4 | Bearer {{tenant1.token}} |
    When I send envelopes to topic "tx.crafter"
      | ID              | ChainName           | From         | ContractName | MethodSignature | PrivateFor                                       | PrivateFrom                                  | Method | Headers.Authorization    |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account1}} | SimpleToken  | constructor()   | ["A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="] | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | 3      | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account2}} | SimpleToken  | constructor()   | ["k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8="] | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | 3      | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account3}} | SimpleToken  | constructor()   | ["A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="] | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | 3      | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account4}} | SimpleToken  | constructor()   | ["k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8="] | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | 3      | Bearer {{tenant1.token}} |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Output | Receipt.PrivateFrom                          | Receipt.PrivateFor                               |
      | 1              | ~              | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | ["A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="] |
      | 1              | ~              | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | ["k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8="] |
      | 1              | ~              | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | ["A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="] |
      | 1              | ~              | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | ["k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8="] |

  @besu
  Scenario: Batch deploy private ERC20 contract with Besu and Orion with different PrivateFrom
    Given I register the following chains
      | alias  | Name                  | URLs                         | Headers.Authorization    |
      | besu   | besu-{{scenarioID}}   | {{global.nodes.besu.URLs}}   | Bearer {{tenant1.token}} |
      | besu_1 | besu_1-{{scenarioID}} | {{global.nodes.besu_1.URLs}} | Bearer {{tenant1.token}} |
    And I have created the following accounts
      | alias    | Headers.Authorization    |
      | account1 | Bearer {{tenant1.token}} |
      | account2 | Bearer {{tenant1.token}} |
      | account3 | Bearer {{tenant1.token}} |
      | account4 | Bearer {{tenant1.token}} |
    When I send envelopes to topic "tx.crafter"
      | ID              | ChainName             | From         | ContractName | MethodSignature | PrivateFor                                       | PrivateFrom                                  | Method | Headers.Authorization    |
      | {{random.uuid}} | besu-{{scenarioID}}   | {{account1}} | SimpleToken  | constructor()   | ["k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8="] | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | 3      | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu_1-{{scenarioID}} | {{account2}} | SimpleToken  | constructor()   | ["k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8="] | A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo= | 3      | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu-{{scenarioID}}   | {{account3}} | SimpleToken  | constructor()   | ["k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8="] | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | 3      | Bearer {{tenant1.token}} |
      | {{random.uuid}} | besu_1-{{scenarioID}} | {{account4}} | SimpleToken  | constructor()   | ["k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8="] | A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo= | 3      | Bearer {{tenant1.token}} |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Output | Receipt.PrivateFrom                          | Receipt.PrivateFor                               |
      | 1              | ~              | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | ["k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8="] |
      | 1              | ~              | A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo= | ["k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8="] |
      | 1              | ~              | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | ["k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8="] |
      | 1              | ~              | A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo= | ["k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8="] |

  @besu
  Scenario: Batch deploy private ERC20 for a privacy group
    Given I register the following chains
      | alias | Name                | URLs                       | Headers.Authorization    |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu.URLs}} | Bearer {{tenant1.token}} |
    And I wait "1.5s"
    And I have created the following accounts
      | alias    | Headers.Authorization    |
      | account1 | Bearer {{tenant1.token}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.chain-registry}}/{{besu.UUID}}" with json:
      """
{
    "jsonrpc": "2.0",
    "method": "priv_createPrivacyGroup",
    "params": [
        {
            "addresses": [
                "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
                "Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs=",
                "k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8="
            ],
            "name": "TestGroup",
            "description": "TestGroup"
        }
    ],
    "id": 1
}
      """
    Then the response code should be 200
    Then I store response field "result" as "privacyGroupId"
    When I send envelopes to topic "tx.crafter"
      | ID              | ChainName           | From         | ContractName | MethodSignature | PrivacyGroupID     | PrivateFrom                                  | Method | Headers.Authorization    |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account1}} | SimpleToken  | constructor()   | {{privacyGroupId}} | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | 3      | Bearer {{tenant1.token}} |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Output | Receipt.PrivateFrom                          | Receipt.PrivacyGroupId |
      | 1              | ~              | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | {{privacyGroupId}}     |
