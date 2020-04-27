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
      | chainName | from                                       | contractName | methodSignature | gasPrice | privateFor                                   | privateFrom                                  | method                     | tenantid                             |
      | besu      | 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 | SimpleToken  | constructor()   | 0        | A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo= | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | EEA_SENDPRIVATETRANSACTION | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields:
      | receipt.status |
      | 1              |

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
