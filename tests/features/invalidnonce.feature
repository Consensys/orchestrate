@public-tx
Feature: Invalid Nonce
  As an external developer
  I want transaction with invalid nonce to be recovered, sent to blockchain and mined

  Scenario: Nonce Too High
    When I send envelopes to topic "tx.signer"
      | chain.nodeName | from                                       | tx.nonce  | tx.gasPrice | tx.gas | tenantid                             |
      | geth           | 0xa8d8db1d8919665a18212374d623fc7c0dfda410 | 1000000   | 1000000000  | 21000  | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
      | geth           | 0xa8d8db1d8919665a18212374d623fc7c0dfda410 | 1000001   | 1000000000  | 21000  | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
      | geth           | 0xa8d8db1d8919665a18212374d623fc7c0dfda410 | 1000002   | 1000000000  | 21000  | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.nonce"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have nonce set
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoder"
    Then Envelopes should be in topic "tx.decoded"

  Scenario: Nonce Too Low
    Given I register the following contract
      | name         | artifacts        | tenantid                             |
      | SimpleToken  | SimpleToken.json | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    And I have deployed contract "token"
      | chain.nodeName | from                                       | contract.name | method.sig    | tx.gas  | tenantid                             |
      | geth           | 0xbfc7137876d7ac275019d70434b0f0779824a969 | SimpleToken   | constructor() | 2000000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    When I send envelopes to topic "tx.signer"
      | chain.nodeName | from                                       | tx.nonce  | tx.gasPrice | tx.gas | tenantid                             |
      | geth           | 0xbfc7137876d7ac275019d70434b0f0779824a969 | 0         | 1000000000  | 21000  | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.nonce"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have nonce set
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoder"
    Then Envelopes should be in topic "tx.decoded"

  Scenario: Chaotic nonce
    When I send envelopes to topic "tx.signer"
      | chain.nodeName | from                                       | tx.nonce  | tx.gasPrice | tx.gas | tenantid                             |
      | geth           | 0x93f7274c9059e601be4512f656b57b830e019e41 | 1000002   | 1000000000  | 21000  | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
      | geth           | 0x93f7274c9059e601be4512f656b57b830e019e41 | 2         | 1000000000  | 21000  | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
      | geth           | 0x93f7274c9059e601be4512f656b57b830e019e41 | 1000000   | 1000000000  | 21000  | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
      | geth           | 0x93f7274c9059e601be4512f656b57b830e019e41 | 1         | 1000000000  | 21000  | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
      | geth           | 0x93f7274c9059e601be4512f656b57b830e019e41 | 1000001   | 1000000000  | 21000  | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
      | geth           | 0x93f7274c9059e601be4512f656b57b830e019e41 | 3         | 1000000000  | 21000  | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.nonce"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have nonce set
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoder"
    Then Envelopes should be in topic "tx.decoded"
