@public-tx
Feature: Generate account
  As as external developer
  I want to generate a new account

  Background:
    Given I have the following tenants
      | alias   | tenantID        |
      | tenant1 | {{random.uuid}} |

  @account
  Scenario: Generate account
    When I send envelopes to topic "account.generator"
      | Headers.Authorization    |
      | Bearer {{tenant1.token}} |
    Then Envelopes should be in topic "account.generator"
    Then Envelopes should be in topic "account.generated"
    And Envelopes should have the following fields
      | From |
      | ~    |


#  @account-faucet
#  Scenario: Generate account with faucet
#    And I register the following chains
#      | alias  | Name                  | URLs                         | Headers.Authorization    |
#      | besu_1 | besu_1-{{scenarioID}} | {{global.nodes.besu_1.URLs}} | Bearer {{tenant1.token}} |
#    And I have created the following accounts
#      | alias    | ID              | Headers.Authorization    |
#      | account1 | {{random.uuid}} | Bearer {{tenant1.token}} |
#    Given I sign the following transactions
#      | alias     | ID              | Value              | Gas   | To           | privateKey                                   | ChainUUID       | Headers.Authorization    |
#      | txFaucet1 | {{random.uuid}} | 150000000000000000 | 21000 | {{account1}} | {{global.nodes.besu_1.fundedPrivateKeys[0]}} | {{besu_1.UUID}} | Bearer {{tenant1.token}} |
#    Then I track the following envelopes
#      | ID               |
#      | {{txFaucet1.ID}} |
#    Given I set the headers
#      | Key           | Value                    |
#      | Authorization | Bearer {{tenant1.token}} |
#      | Content-Type  | application/json         |
#    When I send "POST" request to "{{global.tx-scheduler}}/transactions/send-raw" with json:
#  """
#{
#    "chain": "besu_1-{{scenarioID}}",
#    "params": {
#        "raw": "{{txFaucet1.Raw}}"
#    },
#    "labels": {
#    	"scenario.id": "{{scenarioID}}",
#    	"id": "{{txFaucet1.ID}}"
#    }
#}
#      """
#    Then the response code should be 202
#    Then Envelopes should be in topic "tx.decoded"
#    And I register the following faucets
#      | Name                         | ChainRule       | CreditorAccount | MaxBalance       | Amount           | Cooldown | Headers.Authorization    |
#      | besu_1-faucet-{{scenarioID}} | {{besu_1.UUID}} | {{account1}}    | 1000000000000000 | 1000000000000000 | 1s       | Bearer {{tenant1.token}} |
#    And I have created the following accounts
#      | alias    | ID              | ChainName             | ContextLabels.faucetChildTxID | Headers.Authorization    |
#      | account2 | {{random.uuid}} | besu_1-{{scenarioID}} | {{random.uuid}}               | Bearer {{tenant1.token}} |
#    When I send "POST" request to "{{global.chain-registry}}/{{besu_1.UUID}}" with json:
#      """
#      {
#        "jsonrpc": "2.0",
#        "method": "eth_getBalance",
#        "params": [
#          "{{account2}}",
#          "latest"
#        ],
#        "id": 1
#      }
#      """
#    Then the response code should be 200
#    And Response should have the following fields
#      | result          |
#      | 0x38d7ea4c68000 |
