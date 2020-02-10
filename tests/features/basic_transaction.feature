@public-tx
Feature: Send transfer transaction
  As an external developer
  I want to process a single transfer transaction

  Scenario: Send transfer transaction
    When I send envelopes to topic "tx.crafter"
      | chainName | from                                       | gas | to                                      | value            | tenantid                             |
      | besu       | 0xdbb881a51cd4023e4400cef3ef73046743f08da3 | 21000  | 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 | 1000000000000000000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
      | geth       | 0xdbb881a51cd4023e4400cef3ef73046743f08da3 | 21000  | 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 | 1000000000000000000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.nonce"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have nonce set
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have raw and hash set
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have log decoded
