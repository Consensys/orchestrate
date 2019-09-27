Feature: Invalid Nonce
  As an external developer
  I want transaction with invalid nonce to be recovered, sent to blockchain and mined

  Scenario: Nonce Too High
    Given I have the following envelope:
      | AliasChainId | from                                       | nonce     |  gas  | gasPrice   |
      | primary      | 0xa8d8db1d8919665a18212374d623fc7c0dfda410 | 1000000   | 21000 | 1000000000 |
      | primary      | 0xa8d8db1d8919665a18212374d623fc7c0dfda410 | 1000001   | 21000 | 1000000000 |
      | primary      | 0xa8d8db1d8919665a18212374d623fc7c0dfda410 | 1000002   | 21000 | 1000000000 |
    When I send this envelope to tx-signer
    Then the tx-signer should sign
    Then the tx-listener should catch the tx

  Scenario: Nonce Too Low
    Given I have the following envelope:
      | AliasChainId | from                                       | nonce     |  gas  | gasPrice | 
      | primary      | 0xa8d8db1d8919665a18212374d623fc7c0dfda410 | 0   | 21000 | 1000000000 |
    When I send this envelope to tx-signer
    Then the tx-signer should sign
    Then the tx-listener should catch the tx