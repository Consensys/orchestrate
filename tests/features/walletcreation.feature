@public-tx
Feature: Generate wallet
  As as external developer
  I want to generate a new wallet

  @wallet
  Scenario: Generate wallet
    When I send envelopes to topic "wallet.generator"
    | tenantid                             |
    | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    Then Envelopes should be in topic "wallet.generator"
    Then Envelopes should be in topic "wallet.generated"
    And Envelopes should have from set

  @wallet
  Scenario: Generate wallet with Faucet credit
     When I send envelopes to topic "wallet.generator"
    | chain.name | tx.value           | tenantid                             |
    | geth       | 100000000000000000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    Then Envelopes should be in topic "wallet.generator"
    Then Envelopes should be in topic "wallet.generated"
    And Envelopes should have from set
