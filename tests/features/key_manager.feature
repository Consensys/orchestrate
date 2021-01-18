@key-manager
Feature: Key Manager
  As as external developer
  I want to sign transactions using keys store in the Vault

  @zk-snarks
  Scenario: Sign zk-snarks account, sign and verify it
    Given I have the following tenants
      | alias   | tenantID        |
      | tenant1 | {{random.uuid}} |
    Then I send "POST" request to "{{global.key-manager}}/zk-snarks/accounts" with json:
      """
      {
        "namespace": "{{tenant1.tenantID}}"
      }
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias      | path      |
      | accountOne | publicKey |
    Then I send "GET" request to "{{global.key-manager}}/zk-snarks/accounts/{{accountOne}}?namespace={{tenant1.tenantID}}"
    Then the response code should be 200
    And Response should have the following fields
      | publicKey      | curve | signingAlgorithm |
      | {{accountOne}} | bn256 | eddsa            |
    Then I send "POST" request to "{{global.key-manager}}/zk-snarks/accounts/{{accountOne}}/sign" with json:
      """
      {
        "data": "data to sign",
        "namespace": "{{tenant1.tenantID}}"
      }
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias        | path |
      | signatureOne | .    |
    Then I send "POST" request to "{{global.key-manager}}/zk-snarks/accounts/verify-signature" with json:
      """
      {
        "data": "data to sign",
        "signature": "{{signatureOne}}",
        "publicKey": "{{accountOne}}"
      }
      """
    Then the response code should be 204

    
    
