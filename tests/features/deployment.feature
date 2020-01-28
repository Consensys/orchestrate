@public-tx
Feature: deploy ERC20 contract
  As an external developer
  I want to deploy a contract

  Scenario: Deploy ERC20
    Given I register the following contract
      | name         | artifacts        | tenantid                             |
      | SimpleToken  | SimpleToken.json | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    When I send envelopes to topic "tx.crafter"
      | chain.name | from                                       | contract.name | method.sig      | tx.gas  | tenantid                             |
      | geth       | 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 | SimpleToken   | constructor()   | 2000000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.nonce"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have nonce set
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have raw and hash set
    Then Envelopes should be in topic "tx.decoder"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have log decoded
