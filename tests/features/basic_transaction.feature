@public-tx
Feature: Send transfer transaction
  As an external developer
  I want to process a single transfer transaction

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
      | besu-faucet-{{scenarioID}} | {{besu.UUID}} | {{global.nodes.besu.fundedAccounts[7]}} | 1000000000000000000 | 1000000000000000000 | 1s       | Bearer {{tenant1.token}} |
      | geth-faucet-{{scenarioID}} | {{geth.UUID}} | {{global.nodes.geth.fundedAccounts[8]}} | 1000000000000000000 | 1000000000000000000 | 1s       | Bearer {{tenant1.token}} |
    And I have created the following accounts
      | alias    | ID              | ChainName           | ContextLabels.faucetChildTxID | Headers.Authorization    |
      | account1 | {{random.uuid}} | besu-{{scenarioID}} | {{random.uuid}}               | Bearer {{tenant1.token}} |
      | account2 | {{random.uuid}} | geth-{{scenarioID}} | {{random.uuid}}               | Bearer {{tenant1.token}} |

  Scenario: Send transfer transaction
    When I send envelopes to topic "tx.crafter"
      | ID              | ChainName           | From         | Gas   | To                                         | Value              | Headers.Authorization    |
      | {{random.uuid}} | besu-{{scenarioID}} | {{account1}} | 21000 | 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 | 100000000000000000 | Bearer {{tenant1.token}} |
      | {{random.uuid}} | geth-{{scenarioID}} | {{account2}} | 21000 | 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 | 100000000000000000 | Bearer {{tenant1.token}} |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have the following fields
      | Nonce |
      | 0     |
      | 0     |
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have the following fields
      | Raw | TxHash |
      | ~   | ~      |
      | ~   | ~      |
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields
      | Receipt.Status |
      | 1              |
      | 1              |
