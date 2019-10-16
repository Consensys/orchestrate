Feature: Deploy private ERC20 contract
  As an external developer
  I want to deploy a private contract

  Scenario: Deploy private ERC20 contract
   Given I register the following contract
      | name         | artifacts        |
      | SimpleToken  | SimpleToken.json |
    When I send envelopes to topic "crafter"
      | chain.id      | from                                       | contract.name | method.sig    | tx.gas  | privateFor                                   | privateFrom                                   | privateTxType            | protocol |
      | chain.private | 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 | SimpleToken   | constructor() | 2000000 | QfeDAys9MPDs2XHExtc84jKGHxZg/aj52DTh0vtA3Xc= | BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3Bo=  | whatAmISupposedToSetHere | 2        |
    Then Envelopes should be in topic "crafter"
    Then Envelopes should be in topic "nonce"
    Then Envelopes should be in topic "signer"
    Then Envelopes should be in topic "sender"
    Then Envelopes should be in topic "decoded"
