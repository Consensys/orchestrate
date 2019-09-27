Feature: deploy a private ERC20 contract
  As an external developer
  I want to deploy a private contract

  Scenario: Create a private instance of ERC20
    Given I have the following envelope:
      | chainId | from                                       | contractName | methodSignature | gas     | privateFor                                   | privateFrom                                   | privateTxType            | protocol |
      | 10      | 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 | SimpleToken  | constructor()   | 2000000 | QfeDAys9MPDs2XHExtc84jKGHxZg/aj52DTh0vtA3Xc= | BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3Bo=  | whatAmISupposedToSetHere | 2        |
    When I send these envelopes to CoreStack
    Then CoreStack should receive envelopes
    Then the tx-crafter should set the data
    Then the tx-nonce should set the nonce
    Then the tx-signer should sign
    Then the tx-sender should send the tx
    Then the tx-listener should catch the tx
    Then the tx-decoder should decode
