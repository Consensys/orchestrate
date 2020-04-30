@private-tx
Feature: Deploy private ERC20 contract
  As an external developer
  I want to deploy a private contract

  @quorum
  Scenario: Deploy private ERC20 contract with Quorum and Tessera
    Given I register the following contract
      | name        | artifacts        | tenantid                             |
      | SimpleToken | SimpleToken.json | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    When I send envelopes to topic "tx.crafter"
      | chainName | from                                       | contractName | methodSignature | gas     | privateFor                                   | privateFrom                                  | method                        | tenantid                             |
      | quorum    | 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 | SimpleToken  | constructor()   | 2000000 | QfeDAys9MPDs2XHExtc84jKGHxZg/aj52DTh0vtA3Xc= | BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3Bo= | ETH_SENDRAWPRIVATETRANSACTION | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have the following fields:
      | internalLabels.enclaveKey |
      | ~                         |
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields:
      | receipt.status |
      | 1              |

  @besu
  Scenario: Deploy private ERC20 contract with Besu and Orion
    Given I register the following contract
      | name        | artifacts        | tenantid                             |
      | SimpleToken | SimpleToken.json | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    When I send envelopes to topic "tx.crafter"
      | chainName | from                                       | contractName | methodSignature | privateFor                                   | privateFrom                                  | method                     | tenantid                             |
      | besu      | 0x93f7274c9059e601be4512f656b57b830e019e41 | SimpleToken  | constructor()   | A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo= | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | EEA_SENDPRIVATETRANSACTION | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields:
      | receipt.status | receipt.output | receipt.privateFrom                          | receipt.privateFor                             |
      | 1              | ~              | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | [A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=] |

  @besu
  Scenario: Batch deploy private ERC20 contract with Besu and Orion with different privateFor
    Given I register the following contract
      | name        | artifacts        | tenantid                             |
      | SimpleToken | SimpleToken.json | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    When I send envelopes to topic "tx.crafter"
      | chainName | from                                       | contractName | methodSignature | privateFor                                   | privateFrom                                  | method                     | tenantid                             |
      | besu      | 0x93f7274c9059e601be4512f656b57b830e019e41 | SimpleToken  | constructor()   | A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo= | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | EEA_SENDPRIVATETRANSACTION | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
      | besu      | 0xa8d8db1d8919665a18212374d623fc7c0dfda410 | SimpleToken  | constructor()   | k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8= | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | EEA_SENDPRIVATETRANSACTION | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
      | besu      | 0xdbb881a51cd4023e4400cef3ef73046743f08da3 | SimpleToken  | constructor()   | A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo= | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | EEA_SENDPRIVATETRANSACTION | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
      | besu      | 0xbfc7137876d7ac275019d70434b0f0779824a969 | SimpleToken  | constructor()   | k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8= | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | EEA_SENDPRIVATETRANSACTION | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields:
      | receipt.status | receipt.output | receipt.privateFrom                          | receipt.privateFor                             |
      | 1              | ~              | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | [A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=] |
      | 1              | ~              | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | [k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8=] |
      | 1              | ~              | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | [A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=] |
      | 1              | ~              | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | [k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8=] |

  @besu
  Scenario: Batch deploy private ERC20 contract with Besu and Orion with different privateFrom
    Given I register the following contract
      | name        | artifacts        | tenantid |
      | SimpleToken | SimpleToken.json | _        |
    When I send envelopes to topic "tx.crafter"
      | chainName | from                                       | contractName | methodSignature | privateFor                                   | privateFrom                                  | method                     | tenantid |
      | besu      | 0xffbba394def3ff1df0941c6429887107f58d4e9b | SimpleToken  | constructor()   | k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8= | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | EEA_SENDPRIVATETRANSACTION | _        |
      | besu_1    | 0x6009608a02a7a15fd6689d6dad560c44e9ab61ff | SimpleToken  | constructor()   | k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8= | A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo= | EEA_SENDPRIVATETRANSACTION | _        |
      | besu      | 0xff778b716fc07d98839f48ddb88d8be583beb684 | SimpleToken  | constructor()   | k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8= | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | EEA_SENDPRIVATETRANSACTION | _        |
      | besu_1    | 0xf5956eb46b377ae41b41bda94e6270208d8202bb | SimpleToken  | constructor()   | k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8= | A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo= | EEA_SENDPRIVATETRANSACTION | _        |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields:
      | receipt.status | receipt.output | receipt.privateFrom                          | receipt.privateFor                             |
      | 1              | ~              | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | [k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8=] |
      | 1              | ~              | A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo= | [k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8=] |
      | 1              | ~              | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | [k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8=] |
      | 1              | ~              | A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo= | [k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8=] |


  @quorum
  Scenario: Deploy private ERC20 contract with unknown chainName
    Given I register the following contract
      | name        | artifacts        | tenantid                             |
      | SimpleToken | SimpleToken.json | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    When I send envelopes to topic "tx.crafter"
      | chainName | from                                       | contractName | methodSignature | gas     | privateFor                                   | privateFrom                                  | method                        | tenantid                             |
      | unknown   | 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 | SimpleToken  | constructor()   | 2000000 | QfeDAys9MPDs2XHExtc84jKGHxZg/aj52DTh0vtA3Xc= | BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3Bo= | ETH_SENDRAWPRIVATETRANSACTION | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.recover"
    And Envelopes should have the following errors:
      | errors                           |
      | no chain found with name unknown |

  @quorum
  Scenario: Deploy private ERC20 contract with unknown privateFrom
    Given I register the following contract
      | name        | artifacts        | tenantid                             |
      | SimpleToken | SimpleToken.json | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    When I send envelopes to topic "tx.crafter"
      | chainName | from                                       | contractName | methodSignature | gas     | privateFor                                   | privateFrom  | method                        | tenantid                             |
      | quorum    | 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 | SimpleToken  | constructor()   | 2000000 | QfeDAys9MPDs2XHExtc84jKGHxZg/aj52DTh0vtA3Xc= | dW5rbm93bg== | ETH_SENDRAWPRIVATETRANSACTION | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.recover"
    And Envelopes should have the following errors:
      | errors                                      |
      | failed to send a request to Tessera enclave |
