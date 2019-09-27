Feature: make ERC20 transactions
  As an external developer
  I want to process multiple transactions

  Scenario: Make an transfer transaction
    Given I have the following envelope:
      | AliasContractInstance | AliasChainId | from                                       | contractName | methodSignature | gas     |
      | SimpleToken           | primary      | 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 | SimpleToken  | constructor()   | 2000000 |
    When I store the following contract
      | contractName | fileName         |
      | SimpleToken  | SimpleToken.json |
    And I send these envelopes to CoreStack
    Then I should catch their contract addresses
    Given I have the following envelope:
      | AliasChainId | from                                       | AliasTo     | methodSignature           | args                                           |
      | primary      | 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 | SimpleToken | transfer(address,uint256) | 0xdbb881a51cd4023e4400cef3ef73046743f08da3,0x1 |
      | primary      | 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 | SimpleToken | transfer(address,uint256) | 0x6009608a02a7a15fd6689d6dad560c44e9ab61ff,0x2 |
    When I send these envelopes to CoreStack
    Then CoreStack should receive envelopes
    Then the tx-crafter should set the data
    Then the tx-nonce should set the nonce
    Then the tx-signer should sign
    Then the tx-sender should send the tx
    Then the tx-listener should catch the tx
    Then the tx-decoder should decode
