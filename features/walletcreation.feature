Feature: generate a new wallet
  As as external developer
  I want to generate a new ethereum address

  Scenario: Make a wallet generation
    Given I have the following envelope:
    |
    |
    When I send these envelope in WalletGenerator
    Then WalletGenerator should receive them
    Then the tx-signer should set from
