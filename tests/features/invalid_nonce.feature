@public-tx
Feature: Invalid Nonce
  As an external developer
  I want transaction with invalid nonce to be recovered, sent to blockchain and mined

  Scenario: Nonce Too High
    When I send envelopes to topic "tx.signer"
      | chainName | from                                       | to                                         | nonce   | gasPrice   | gas   | tenantid                             |
      | besu      | 0xa8d8db1d8919665a18212374d623fc7c0dfda410 | 0x0000000000000000000000000000000000000000 | 1000000 | 1000000000 | 21000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
      | besu      | 0xa8d8db1d8919665a18212374d623fc7c0dfda410 | 0x0000000000000000000000000000000000000000 | 1000001 | 1000000000 | 21000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
      | besu      | 0xa8d8db1d8919665a18212374d623fc7c0dfda410 | 0x0000000000000000000000000000000000000000 | 1000002 | 1000000000 | 21000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have nonce set
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"

  Scenario: Nonce Too Low
    Given I register the following contract
      | name        | artifacts        | tenantid                             |
      | SimpleToken | SimpleToken.json | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    And I have deployed contract "token"
      | chainName | from                                       | contractName | methodSignature | gas     | tenantid                             |
      | besu      | 0xbfc7137876d7ac275019d70434b0f0779824a969 | SimpleToken  | constructor()   | 2000000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    When I send envelopes to topic "tx.signer"
      | chainName | from                                       | to                                         | nonce | gasPrice   | gas   | tenantid                             |
      | besu      | 0xbfc7137876d7ac275019d70434b0f0779824a969 | 0x0000000000000000000000000000000000000000 | 0     | 1000000000 | 21000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have nonce set
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"

  Scenario: Chaotic nonce
    Given I register the following contract
      | name        | artifacts        | tenantid                             |
      | SimpleToken | SimpleToken.json | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    # Next deployment purpose is to increase account nonce to at least 1
    And I have deployed contract "token"
      | chainName | from                                       | contractName | methodSignature | gas     | tenantid                             |
      | besu      | 0x93f7274c9059e601be4512f656b57b830e019e41 | SimpleToken  | constructor()   | 2000000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    When I send envelopes to topic "tx.signer"
      | chainName | from                                       | to                                         | nonce   | gasPrice   | gas   | tenantid                             |
      | besu      | 0x93f7274c9059e601be4512f656b57b830e019e41 | 0x0000000000000000000000000000000000000000 | 1000002 | 1000000000 | 21000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
      | besu      | 0x93f7274c9059e601be4512f656b57b830e019e41 | 0x0000000000000000000000000000000000000000 | 1000000 | 1000000000 | 21000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
      | besu      | 0x93f7274c9059e601be4512f656b57b830e019e41 | 0x0000000000000000000000000000000000000000 | 0       | 1000000000 | 21000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
      | besu      | 0x93f7274c9059e601be4512f656b57b830e019e41 | 0x0000000000000000000000000000000000000000 | 1000001 | 1000000000 | 21000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have nonce set
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"

  @private-tx
  Scenario: Nonce too high with private transaction
    When I send envelopes to topic "tx.signer"
      | chainName | from                                       | to                                         | nonce   | gasPrice | gas   | tenantid                             | privateFor                                   | privateFrom                                  | method                     | tenantid                             |
      | besu      | 0xa8d8db1d8919665a18212374d623fc7c0dfda410 | 0x0000000000000000000000000000000000000000 | 1000000 | 0        | 30000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 | A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo= | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | EEA_SENDPRIVATETRANSACTION | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
      | besu      | 0xa8d8db1d8919665a18212374d623fc7c0dfda410 | 0x0000000000000000000000000000000000000000 | 1000001 | 0        | 30000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 | A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo= | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | EEA_SENDPRIVATETRANSACTION | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
      | besu      | 0xa8d8db1d8919665a18212374d623fc7c0dfda410 | 0x0000000000000000000000000000000000000000 | 1000002 | 0        | 30000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 | A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo= | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | EEA_SENDPRIVATETRANSACTION | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have nonce set
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
