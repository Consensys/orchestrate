@public-tx
Feature: Generate account
  As as external developer
  I want to generate a new account

  @account
  Scenario: Generate account
    When I send envelopes to topic "account.generator"
      | tenantid                             |
      | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    Then Envelopes should be in topic "account.generator"
    Then Envelopes should be in topic "account.generated"
    And Envelopes should have from set

  @account
  Scenario: Generate account with Faucet credit
    When I send envelopes to topic "account.generator"
      | chainName | value           | tenantid                             |
      | geth       | 100000000000000000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    Then Envelopes should be in topic "account.generator"
    Then Envelopes should be in topic "account.generated"
    And Envelopes should have from set
