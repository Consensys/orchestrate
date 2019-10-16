Feature: Send transfer transaction
  As an external developer
  I want to process a single transfer transaction

  Scenario: Send transfer transaction
    When I send envelopes to topic "crafter"
      | chain.id           | from                                       | tx.gas | tx.to                                      | tx.value            |
      | chain.primary      | 0xdbb881a51cd4023e4400cef3ef73046743f08da3 | 21000  | 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 | 1000000000000000000 |
    Then Envelopes should be in topic "crafter"
    Then Envelopes should be in topic "nonce"
    Then Envelopes should be in topic "signer"
    And Envelopes should have nonce set
    Then Envelopes should be in topic "sender"
    And Envelopes should have raw and hash set
    Then Envelopes should be in topic "decoded"
    And Envelopes should have log decoded
