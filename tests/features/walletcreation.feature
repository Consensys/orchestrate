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

  Scenario: Make a wallet generation with faucet credit
    Given I have the following envelope:
    | AliasChainId | value              |  
    | primary      | 100000000000000000 |
    When I send these envelope in WalletGenerator
    Then WalletGenerator should receive them
    Then the tx-signer should set from