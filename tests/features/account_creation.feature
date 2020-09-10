@public-tx
Feature: Generate account
  As as external developer
  I want to generate a new account and generate account with faucet

  Background:
    Given I have the following tenants
      | alias   | tenantID        |
      | tenant1 | {{random.uuid}} |

  @account
  Scenario: Generate account
    When I send envelopes to topic "account.generator"
      | ID              | Headers.Authorization    |
      | {{random.uuid}} | Bearer {{tenant1.token}} |
    Then Envelopes should be in topic "account.generator"
    Then Envelopes should be in topic "account.generated"
    And Envelopes should have the following fields
      | From |
      | ~    |

  Scenario: Generate account with faucet
    And I register the following chains
      | alias | Name                | URLs                         | Headers.Authorization    |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu_1.URLs}} | Bearer {{tenant1.token}} |
    And I register the following faucets
      | Name                | ChainRule     | CreditorAccount                             | MaxBalance       | Amount           | Cooldown | Headers.Authorization    |
      | besu-{{scenarioID}} | {{besu.UUID}} | {{global.nodes.besu_1.fundedPublicKeys[0]}} | 1000000000000000 | 1000000000000000 | 1s       | Bearer {{tenant1.token}} |
    And I have created the following accounts
      | alias    | ID              | ChainName           | Headers.Authorization    |
      | account1 | {{random.uuid}} | besu-{{scenarioID}} | Bearer {{tenant1.token}} |
    Given I sleep "3s"
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.chain-registry}}/{{besu.UUID}}" with json:
      """
      {
        "jsonrpc": "2.0",
        "method": "eth_getBalance",
        "params": [
          "{{account1}}",
          "latest"
        ],
        "id": 1
      }
      """
    Then the response code should be 200
    And Response should have the following fields
      | result          |
      | 0x38d7ea4c68000 |
