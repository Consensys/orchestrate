@private-tx
Feature: Deploy private ERC20 contract
  As an external developer
  I want to deploy a private contract

  Scenario: Deploy private ERC20 contract
    Given I register the following contract
      | name        | artifacts        | tenantid                             |
      | SimpleToken | SimpleToken.json | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    When I send envelopes to topic "tx.crafter"
      | chainName | from                                       | contractName | methodSignature    | gas  | privateFor                                   | privateFrom                                  | privateTxType            | method | tenantid                             |
      | quorum     | 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 | SimpleToken   | constructor() | 2000000 | QfeDAys9MPDs2XHExtc84jKGHxZg/aj52DTh0vtA3Xc= | BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3Bo= | whatAmISupposedToSetHere | ETH_SENDRAWPRIVATETRANSACTION        | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"
