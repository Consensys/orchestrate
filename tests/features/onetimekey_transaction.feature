@onetimekey-tx
Feature: One time key signature
  As an external developer
  I want to send anonymous transaction using one time usage keys

  @quorum
  @besu
  Scenario: Deploy ERC20
    Given I register the following contract
      | name        | artifacts        | tenantid                             |
      | SimpleToken | SimpleToken.json | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    When I send envelopes to topic "tx.crafter"
      | chainName | from | contractName | methodSignature | gas     | tenantid                             | id       |
      | besu      | *    | SimpleToken  | constructor()   | 2000000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 | unique-1 |
      | besu      | *    | SimpleToken  | constructor()   | 2000000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 | unique-2 |
      | quorum    | *    | SimpleToken  | constructor()   | 2000000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 | unique-1 |
      | quorum    | *    | SimpleToken  | constructor()   | 2000000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 | unique-2 |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have from set
    And Envelopes should have raw and hash set
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields:
      | receipt.status |
      | 1              |
      | 1              |
      | 1              |
      | 1              |

  @quorum
  Scenario: Deploy private ERC20 contract with Quorum and Tessera
    Given I register the following contract
      | name        | artifacts        | tenantid                             |
      | SimpleToken | SimpleToken.json | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    When I send envelopes to topic "tx.crafter"
      | chainName | from | contractName | methodSignature | gas     | privateFor                                   | privateFrom                                  | method                        | tenantid                             | id       |
      | quorum    | *    | SimpleToken  | constructor()   | 2000000 | QfeDAys9MPDs2XHExtc84jKGHxZg/aj52DTh0vtA3Xc= | BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3Bo= | ETH_SENDRAWPRIVATETRANSACTION | f30c452b-e5fb-4102-a45d-bc00a060bcc6 | unique-1 |
      | quorum    | *    | SimpleToken  | constructor()   | 2000000 | QfeDAys9MPDs2XHExtc84jKGHxZg/aj52DTh0vtA3Xc= | BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3Bo= | ETH_SENDRAWPRIVATETRANSACTION | f30c452b-e5fb-4102-a45d-bc00a060bcc6 | unique-2 |
      | quorum    | *    | SimpleToken  | constructor()   | 2000000 | QfeDAys9MPDs2XHExtc84jKGHxZg/aj52DTh0vtA3Xc= | BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3Bo= | ETH_SENDRAWPRIVATETRANSACTION | f30c452b-e5fb-4102-a45d-bc00a060bcc6 | unique-3 |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields:
      | receipt.status |
      | 1              |
      | 1              |
      | 1              |

  @besu
  Scenario: Deploy private ERC20 contract with Besu and Orion
    Given I register the following contract
      | name        | artifacts        | tenantid                             |
      | SimpleToken | SimpleToken.json | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    When I send envelopes to topic "tx.crafter"
      | chainName | from | contractName | methodSignature | privateFor                                   | privateFrom                                  | method                     | tenantid                             | id       |
      | besu      | *    | SimpleToken  | constructor()   | A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo= | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | EEA_SENDPRIVATETRANSACTION | f30c452b-e5fb-4102-a45d-bc00a060bcc6 | unique-1 |
      | besu      | *    | SimpleToken  | constructor()   | A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo= | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | EEA_SENDPRIVATETRANSACTION | f30c452b-e5fb-4102-a45d-bc00a060bcc6 | unique-2 |
      | besu      | *    | SimpleToken  | constructor()   | A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo= | Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs= | EEA_SENDPRIVATETRANSACTION | f30c452b-e5fb-4102-a45d-bc00a060bcc6 | unique-3 |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have the following fields:
      | receipt.status |
      | 1              |
      | 1              |
      | 1              |

