Feature: Invalid Nonce
  As an external developer
  I want transaction with invalid nonce to be recovered, sent to blockchain and mined

  Scenario: Nonce Too High
    When I send envelopes to topic "signer"
      | chain.id      | from                                       | tx.nonce  | tx.gasPrice | tx.gas |
      | chain.primary | 0xa8d8db1d8919665a18212374d623fc7c0dfda410 | 1000000   | 1000000000  | 21000  |
      | chain.primary | 0xa8d8db1d8919665a18212374d623fc7c0dfda410 | 1000001   | 1000000000  | 21000  |
      | chain.primary | 0xa8d8db1d8919665a18212374d623fc7c0dfda410 | 1000002   | 1000000000  | 21000  |
    Then Envelopes should be in topic "signer"
    Then Envelopes should be in topic "sender"
    Then Envelopes should be in topic "nonce"
    Then Envelopes should be in topic "signer"
    And Envelopes should have nonce set
    Then Envelopes should be in topic "sender"
    Then Envelopes should be in topic "decoded"

  Scenario: Nonce Too Low
    When I send envelopes to topic "signer"
      | chain.id      | from                                       | tx.nonce  | tx.gasPrice | tx.gas |
      | chain.primary | 0xa8d8db1d8919665a18212374d623fc7c0dfda410 | 0   | 1000000000  | 21000  |
    Then Envelopes should be in topic "signer"
    Then Envelopes should be in topic "sender"
    Then Envelopes should be in topic "nonce"
    Then Envelopes should be in topic "signer"
    And Envelopes should have nonce set
    Then Envelopes should be in topic "sender"
    Then Envelopes should be in topic "decoded"
