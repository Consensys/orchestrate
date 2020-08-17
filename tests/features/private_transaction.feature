@private-tx
Feature: Private transactions
  As an external developer
  I want to deploy a private contract using transaction scheduler

  Background:
    Given I have the following tenants
      | alias   | tenantID        |
      | tenant1 | {{random.uuid}} |
    And I have created the following accounts
      | alias    | Headers.Authorization    |
      | account1 | Bearer {{tenant1.token}} |
      | account2 | Bearer {{tenant1.token}} |
      | account3 | Bearer {{tenant1.token}} |
      | account4 | Bearer {{tenant1.token}} |
    And I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | Bearer {{tenant1.token}} |
    And I register the following chains
      | alias  | Name                  | URLs                           | PrivateTxManagers                           | Headers.Authorization    |
      | quorum | quorum-{{scenarioID}} | {{global.nodes.quorum_1.URLs}} | {{global.nodes.quorum_1.PrivateTxManagers}} | Bearer {{tenant1.token}} |
    And I register the following chains
      | alias  | Name                  | URLs                         | Headers.Authorization    |
      | besu   | besu-{{scenarioID}}   | {{global.nodes.besu_1.URLs}} | Bearer {{tenant1.token}} |
      | besu_2 | besu_2-{{scenarioID}} | {{global.nodes.besu_2.URLs}} | Bearer {{tenant1.token}} |

  @quorum
  Scenario: Deploy private ERC20 contract with Quorum and Tessera
    Given I register the following alias
      | alias                    | value           |
      | quorumDeployContractTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                           |
      | {{quorumDeployContractTxID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "quorum-{{scenarioID}}",
    "params": {
        "from": "{{account1}}",
        "protocol": "Tessera",
        "privateFrom": "{{global.nodes.quorum_1.privateAddress}}",
        "privateFor": ["{{global.nodes.quorum_2.privateAddress}}"],
        "contractName": "SimpleToken"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{quorumDeployContractTxID}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have the following fields
      | InternalLabels.enclaveKey |
      | ~                         |
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.ContractAddress |
      | 1              | ~                       |

  @quorum
  Scenario: Fail to deploy private ERC20 contract with unknown ChainName
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "UnknownChain",
    "params": {
        "from": "{{account1}}",
        "protocol": "Tessera",
        "privateFrom": "{{global.nodes.quorum_1.privateAddress}}",
        "privateFor": ["{{global.nodes.quorum_2.privateAddress}}"],
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

  @quorum
  Scenario: Force not correlative nonce for private and public txs in Quorum/Tessera
    Given I register the following alias
      | alias                           | value           |
      | publicQuorumDeployContractTxID  | {{random.uuid}} |
      | privateQuorumDeployContractTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                                 |
      | {{publicQuorumDeployContractTxID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "quorum-{{scenarioID}}",
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
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    Then I track the following envelopes
      | ID                                  |
      | {{privateQuorumDeployContractTxID}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "quorum-{{scenarioID}}",
    "params": {
        "from": "{{account1}}",
        "protocol": "Tessera",
        "privateFrom": "{{global.nodes.quorum_1.privateAddress}}",
        "privateFor": ["{{global.nodes.quorum_2.privateAddress}}"],
        "contractName": "SimpleToken"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{privateQuorumDeployContractTxID}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have the following fields
      | InternalLabels.enclaveKey |
      | ~                         |
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.ContractAddress |
      | 1              | ~                       |

  @quorum
  Scenario: Fail to deploy private ERC20 contract with unknown PrivateFor
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "quorum-{{scenarioID}}",
    "params": {
        "from": "{{account1}}",
        "protocol": "Tessera",
        "privateFrom": "{{global.nodes.quorum_1.privateAddress}}",
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

  @quorum
  Scenario: Fail to deploy private ERC20 contract with not authorized chain
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "quorum_2-{{scenarioID}}",
    "params": {
        "from": "{{account1}}",
        "protocol": "Tessera",
        "privateFrom": "{{global.nodes.quorum_1.privateAddress}}",
        "privateFor": ["{{global.nodes.quorum_2.privateAddress}}"],
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
  Scenario: Deploy private ERC20 contract with Besu and Orion
    Given I register the following alias
      | alias                  | value           |
      | besuDeployContractTxID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                         |
      | {{besuDeployContractTxID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "from": "{{account2}}",
        "protocol": "Orion",
        "privateFrom": "{{global.nodes.besu_1.privateAddress}}",
        "privateFor": ["{{global.nodes.besu_2.privateAddress}}"],
        "contractName": "SimpleToken"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{besuDeployContractTxID}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Output | Receipt.ContractAddress | Receipt.PrivateFrom                    | Receipt.PrivateFor                         |
      | 1              | ~              | ~                       | {{global.nodes.besu_1.privateAddress}} | ["{{global.nodes.besu_2.privateAddress}}"] |

  @besu
  Scenario: Batch deploy private ERC20 contract with Besu and Orion with different PrivateFor
    Given I register the following alias
      | alias                       | value           |
      | besuDeployContractTxOneID   | {{random.uuid}} |
      | besuDeployContractTxTwoID   | {{random.uuid}} |
      | besuDeployContractTxThreeID | {{random.uuid}} |
      | besuDeployContractTxFourID  | {{random.uuid}} |
    Then I track the following envelopes
      | ID                              |
      | {{besuDeployContractTxOneID}}   |
      | {{besuDeployContractTxTwoID}}   |
      | {{besuDeployContractTxThreeID}} |
      | {{besuDeployContractTxFourID}}  |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "from": "{{account1}}",
        "protocol": "Orion",
        "privateFrom": "{{global.nodes.besu_1.privateAddress}}",
        "privateFor": ["{{global.nodes.besu_2.privateAddress}}"],
        "contractName": "SimpleToken"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{besuDeployContractTxOneID}}"
    }
}
      """
    Then the response code should be 202
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "from": "{{account2}}",
        "protocol": "Orion",
        "privateFrom": "{{global.nodes.besu_1.privateAddress}}",
        "privateFor": ["{{global.nodes.besu_3.privateAddress}}"],
        "contractName": "SimpleToken"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{besuDeployContractTxTwoID}}"
    }
}
      """
    Then the response code should be 202
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "from": "{{account3}}",
        "protocol": "Orion",
        "privateFrom": "{{global.nodes.besu_1.privateAddress}}",
        "privateFor": ["{{global.nodes.besu_2.privateAddress}}"],
        "contractName": "SimpleToken"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{besuDeployContractTxThreeID}}"
    }
}
      """
    Then the response code should be 202
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "from": "{{account4}}",
        "protocol": "Orion",
        "privateFrom": "{{global.nodes.besu_1.privateAddress}}",
        "privateFor": ["{{global.nodes.besu_3.privateAddress}}"],
        "contractName": "SimpleToken"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{besuDeployContractTxFourID}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Output | Receipt.PrivateFrom                    | Receipt.PrivateFor                         |
      | 1              | ~              | {{global.nodes.besu_1.privateAddress}} | ["{{global.nodes.besu_2.privateAddress}}"] |
      | 1              | ~              | {{global.nodes.besu_1.privateAddress}} | ["{{global.nodes.besu_3.privateAddress}}"] |
      | 1              | ~              | {{global.nodes.besu_1.privateAddress}} | ["{{global.nodes.besu_2.privateAddress}}"] |
      | 1              | ~              | {{global.nodes.besu_1.privateAddress}} | ["{{global.nodes.besu_3.privateAddress}}"] |

  @besu
  Scenario: Batch deploy private ERC20 contract with Besu and Orion with different PrivateFrom
    Given I register the following alias
      | alias                       | value           |
      | besuDeployContractTxOneID   | {{random.uuid}} |
      | besuDeployContractTxTwoID   | {{random.uuid}} |
      | besuDeployContractTxThreeID | {{random.uuid}} |
      | besuDeployContractTxFourID  | {{random.uuid}} |
    Then I track the following envelopes
      | ID                              |
      | {{besuDeployContractTxOneID}}   |
      | {{besuDeployContractTxTwoID}}   |
      | {{besuDeployContractTxThreeID}} |
      | {{besuDeployContractTxFourID}}  |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "from": "{{account1}}",
        "protocol": "Orion",
        "privateFrom": "{{global.nodes.besu_1.privateAddress}}",
        "privateFor": ["{{global.nodes.besu_3.privateAddress}}"],
        "contractName": "SimpleToken"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{besuDeployContractTxOneID}}"
    }
}
      """
    Then the response code should be 202
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu_2-{{scenarioID}}",
    "params": {
        "from": "{{account2}}",
        "protocol": "Orion",
        "privateFrom": "{{global.nodes.besu_2.privateAddress}}",
        "privateFor": ["{{global.nodes.besu_3.privateAddress}}"],
        "contractName": "SimpleToken"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{besuDeployContractTxTwoID}}"
    }
}
      """
    Then the response code should be 202
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "from": "{{account3}}",
        "protocol": "Orion",
        "privateFrom": "{{global.nodes.besu_1.privateAddress}}",
        "privateFor": ["{{global.nodes.besu_3.privateAddress}}"],
        "contractName": "SimpleToken"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{besuDeployContractTxThreeID}}"
    }
}
      """
    Then the response code should be 202
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu_2-{{scenarioID}}",
    "params": {
        "from": "{{account4}}",
        "protocol": "Orion",
        "privateFrom": "{{global.nodes.besu_2.privateAddress}}",
        "privateFor": ["{{global.nodes.besu_3.privateAddress}}"],
        "contractName": "SimpleToken"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{besuDeployContractTxFourID}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Output | Receipt.PrivateFrom                    | Receipt.PrivateFor                         |
      | 1              | ~              | {{global.nodes.besu_1.privateAddress}} | ["{{global.nodes.besu_3.privateAddress}}"] |
      | 1              | ~              | {{global.nodes.besu_2.privateAddress}} | ["{{global.nodes.besu_3.privateAddress}}"] |
      | 1              | ~              | {{global.nodes.besu_1.privateAddress}} | ["{{global.nodes.besu_3.privateAddress}}"] |
      | 1              | ~              | {{global.nodes.besu_2.privateAddress}} | ["{{global.nodes.besu_3.privateAddress}}"] |

  @besu
  Scenario: Deploy private ERC20 for a privacy group
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    Given I sleep "2s"
    When I send "POST" request to "{{global.chain-registry}}/{{besu.UUID}}" with json:
      """
{
    "jsonrpc": "2.0",
    "method": "priv_createPrivacyGroup",
    "params": [
        {
            "addresses": [
                "{{global.nodes.besu_1.privateAddress}}",
                "{{global.nodes.besu_2.privateAddress}}",
                "{{global.nodes.besu_3.privateAddress}}"
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
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "from": "{{account1}}",
        "protocol": "Orion",
        "privateFrom": "{{global.nodes.besu_1.privateAddress}}",
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
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Output | Receipt.PrivateFrom                    | Receipt.PrivacyGroupId |
      | 1              | ~              | {{global.nodes.besu_1.privateAddress}} | {{privacyGroupId}}     |


  @besu
  Scenario: Deploy private ERC20 for a privacy group
    Given I register the following alias
      | alias               | value           |
      | privacyGroupNameOne | {{random.uuid}} |
    Then I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    Then I sleep "2s"
    When I send "POST" request to "{{global.chain-registry}}/{{besu.UUID}}" with json:
      """
{
    "jsonrpc": "2.0",
    "method": "priv_createPrivacyGroup",
    "params": [
        {
            "addresses": [
                "{{global.nodes.besu_1.privateAddress}}",
                "{{global.nodes.besu_2.privateAddress}}",
                "{{global.nodes.besu_3.privateAddress}}"
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
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "from": "{{account1}}",
        "protocol": "Orion",
        "privateFrom": "{{global.nodes.besu_1.privateAddress}}",
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
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Output | Receipt.ContractAddress | Receipt.PrivateFrom                    | Receipt.PrivacyGroupId |
      | 1              | ~              | ~                       | {{global.nodes.besu_1.privateAddress}} | {{privacyGroupId}}     |

  @besu
  Scenario: Fail to deploy private ERC20 with a privacy group and privateFor
    Given I register the following alias
      | alias               | value           |
      | privacyGroupNameTwo | {{random.uuid}} |
    Then I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    Then I sleep "2s"
    When I send "POST" request to "{{global.chain-registry}}/{{besu.UUID}}" with json:
      """
{
    "jsonrpc": "2.0",
    "method": "priv_createPrivacyGroup",
    "params": [
        {
            "addresses": [
                "{{global.nodes.besu_1.privateAddress}}",
                "{{global.nodes.besu_2.privateAddress}}",
                "{{global.nodes.besu_3.privateAddress}}"
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
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "from": "{{account1}}",
        "protocol": "Orion",
        "privateFrom": "{{global.nodes.besu_1.privateAddress}}",
        "privacyGroupId": "{{privacyGroupId}}",
        "privateFor": ["{{global.nodes.besu_3.privateAddress}}"],
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
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    Given I register the following alias
      | alias                     | value           |
      | besuDeployContractTxOTKID | {{random.uuid}} |
    Then I track the following envelopes
      | ID                            |
      | {{besuDeployContractTxOTKID}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "annotations": { "oneTimeKey": true },
        "protocol": "Orion",
        "privateFrom": "{{global.nodes.besu_1.privateAddress}}",
        "privateFor": ["{{global.nodes.besu_3.privateAddress}}"],
        "contractName": "SimpleToken"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{besuDeployContractTxOTKID}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Output | Receipt.ContractAddress | Receipt.PrivateFrom                    |
      | 1              | ~              | ~                       | {{global.nodes.besu_1.privateAddress}} |

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
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
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
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.ContractAddress |
      | 1              | ~                       |
    Then I track the following envelopes
      | ID                                |
      | {{privateBesuDeployContractTxID}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "from": "{{account1}}",
        "protocol": "Orion",
        "privateFrom": "{{global.nodes.besu_1.privateAddress}}",
        "privateFor": ["{{global.nodes.besu_2.privateAddress}}","{{global.nodes.besu_3.privateAddress}}"],
        "contractName": "SimpleToken"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{privateBesuDeployContractTxID}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.Output | Receipt.ContractAddress | Receipt.PrivateFrom                    | Receipt.PrivateFor                         |
      | 1              | ~              | ~                       | {{global.nodes.besu_1.privateAddress}} | ["{{global.nodes.besu_2.privateAddress}}"] |
