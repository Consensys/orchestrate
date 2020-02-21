@public-tx
Feature: Multiple transactions
  As an external developer
  I want to process multiple transactions

  Scenario: Send transactions
    Given I register the following contract
      | name        | artifacts        | tenantid                             |
      | SimpleToken | SimpleToken.json | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    And I have deployed contract "token-besu"
      | chainName | from                                       | contractName | methodSignature | gas     | tenantid                             |
      | besu      | 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 | SimpleToken  | constructor()   | 2000000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    And I have deployed contract "token-geth"
      | chainName | from                                       | contractName | methodSignature | gas     | tenantid                             |
      | geth      | 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 | SimpleToken  | constructor()   | 2000000 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    When I send envelopes to topic "tx.crafter"
      | chainName | from                                       | to         | methodSignature           | args                                           | tenantid                             |
      | besu      | 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 | token-besu | transfer(address,uint256) | 0xdbb881a51cd4023e4400cef3ef73046743f08da3,0x1 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
      | geth      | 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 | token-geth | transfer(address,uint256) | 0xdbb881a51cd4023e4400cef3ef73046743f08da3,0x1 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
      | besu      | 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 | token-besu | transfer(address,uint256) | 0x6009608a02a7a15fd6689d6dad560c44e9ab61ff,0x2 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
      | geth      | 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 | token-geth | transfer(address,uint256) | 0x6009608a02a7a15fd6689d6dad560c44e9ab61ff,0x2 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |
    Then Envelopes should be in topic "tx.crafter"
    Then Envelopes should be in topic "tx.signer"
    And Envelopes should have nonce set
    Then Envelopes should be in topic "tx.sender"
    And Envelopes should have raw and hash set
    Then Envelopes should be in topic "tx.decoded"
    And Envelopes should have log decoded

  Scenario: Send raw transactions
    # Send raw transaction with address:0xf5956Eb46b377Ae41b41BDa94e6270208d8202bb, nonce:0, gasLimit:21000, to:0x88a5C2d9919e46F883EB62F7b8Dd9d0CC45bc290, value:0,5ETH
    When I send envelopes to topic "tx.sender"
      | chainName | contextLabels.txMode | raw                                                                                                                                                                                                                  |
      | geth      | raw                  | 0xf86780808252089488a5c2d9919e46f883eb62f7b8dd9d0cc45bc2908806f05b59d3b20000801ca09fd94be4942219541b1fd100341706e2e4caa365c926cde48d8c7aac8c5a0f69a034a00bd00f4ef680a3586208c462280d99cb35e8a89479494af09a7228fdd46a |
    Then Envelopes should be in topic "tx.decoded"
