Feature: Generate wallet
  As as external developer
  I want to generate a new wallet

  Scenario: Generate wallet
    When I send envelopes to topic "wallet.generator"
    |
    |
    Then Envelopes should be in topic "wallet.generator"
    Then Envelopes should be in topic "wallet.generated"
    And Envelopes should have from set

  Scenario: Generate wallet with Faucet credit
     When I send envelopes to topic "wallet.generator"
    | chain.id           | tx.value           |  
    | chain.primary      | 100000000000000000 |
    Then Envelopes should be in topic "wallet.generator"
    Then Envelopes should be in topic "wallet.generated"
    And Envelopes should have from set
