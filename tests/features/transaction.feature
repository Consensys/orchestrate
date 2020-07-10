@public-tx
Feature: Multiple transactions
  As an external developer
  I want to process multiple transactions

  Background:
    Given I have the following tenants
      | alias   | tenantID        |
      | tenant1 | {{random.uuid}} |
    And I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | Bearer {{tenant1.token}} |
    And I register the following chains
      | alias | Name                | URLs                       | Headers.Authorization    |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu.URLs}} | Bearer {{tenant1.token}} |
      | geth  | geth-{{scenarioID}} | {{global.nodes.geth.URLs}} | Bearer {{tenant1.token}} |
    And I have created the following accounts
      | alias    | ID              | Headers.Authorization    |
      | account1 | {{random.uuid}} | Bearer {{tenant1.token}} |
      | account2 | {{random.uuid}} | Bearer {{tenant1.token}} |
    Given I sign the following transactions
      | alias     | ID              | Value              | Gas   | To           | privateKey                                 | ChainUUID     | Headers.Authorization    |
      | txFaucet1 | {{random.uuid}} | 100000000000000000 | 21000 | {{account1}} | {{global.nodes.besu.fundedPrivateKeys[0]}} | {{besu.UUID}} | Bearer {{tenant1.token}} |
      | txFaucet2 | {{random.uuid}} | 100000000000000000 | 21000 | {{account2}} | {{global.nodes.geth.fundedPrivateKeys[0]}} | {{geth.UUID}} | Bearer {{tenant1.token}} |
    Then I track the following envelopes
      | ID               |
      | {{txFaucet1.ID}} |
      | {{txFaucet2.ID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
      | Content-Type  | application/json         |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/send-raw" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
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
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/send-raw" with json:
  """
{
    "chain": "geth-{{scenarioID}}",
    "params": {
        "raw": "{{txFaucet2.Raw}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{txFaucet2.ID}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
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
    Given I register the following alias
      | alias     | value              |
      | recipient | {{random.account}} |
    When I send envelopes to topic "tx.crafter"
      | ID              | ChainName           | From         | To            | Value     | Headers.Authorization    |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account1}} | {{recipient}} | 100000000 | Bearer {{tenant1.token}} |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have the following fields
      | Value     | From         |
      | 100000000 | {{account1}} |
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have the following fields
      | Value     | From         |
      | 100000000 | {{account1}} |
    And Envelopes should have the following fields
      | Raw | TxHash |
      | ~   | ~      |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status |
      | 1              |
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
      | result    |
      | 0x5f5e100 |
