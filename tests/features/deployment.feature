@public-tx
Feature: Deploy ERC20 contract
  As an external developer
  I want to deploy a contract

  Background:
    Given I have the following tenants
      | alias   | tenantID                             |
      | tenant1 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    And I register the following chains
      | alias | Name                | URLs                       | Headers.Authorization    |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu.URLs}} | Bearer {{tenant1.token}} |
      | geth  | geth-{{scenarioID}} | {{global.nodes.geth.URLs}} | Bearer {{tenant1.token}} |
    And I register the following faucets
      | Name                       | ChainRule     | CreditorAccount                         | MaxBalance          | Amount              | Cooldown | Headers.Authorization    |
      | besu-faucet-{{scenarioID}} | {{besu.UUID}} | {{global.nodes.besu.fundedAccounts[0]}} | 1000000000000000000 | 1000000000000000000 | 1s       | Bearer {{tenant1.token}} |
      | geth-faucet-{{scenarioID}} | {{geth.UUID}} | {{global.nodes.geth.fundedAccounts[0]}} | 1000000000000000000 | 1000000000000000000 | 1s       | Bearer {{tenant1.token}} |
    And I have created the following accounts
      | alias    | ChainName           | ContextLabels.faucetChildTxID | Headers.Authorization    |
      | account1 | besu-{{scenarioID}} | {{random.uuid}}               | Bearer {{tenant1.token}} |
      | account2 | geth-{{scenarioID}} | {{random.uuid}}               | Bearer {{tenant1.token}} |
    And I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | Bearer {{tenant1.token}} |

  Scenario: Deploy ERC20
    When I send envelopes to topic "tx.crafter"
      | ID              | ChainName           | From         | ContractName | MethodSignature | Gas     | Headers.Authorization    |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account1}} | SimpleToken  | constructor()   | 2000000 | Bearer {{tenant1.token}} |
      | {{random.uuid}} | geth-{{scenarioID}} | {{account2}} | SimpleToken  | constructor()   | 2000000 | Bearer {{tenant1.token}} |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status | Receipt.ContractAddress |
      | 1              | ~                       |
      | 1              | ~                       |
