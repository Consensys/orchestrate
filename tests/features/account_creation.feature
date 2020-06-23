@public-tx
Feature: Generate account
  As as external developer
  I want to generate a new account

  Background:
    Given I have the following tenants
      | alias   | tenantID                             |
      | tenant1 | f30c452b-e5fb-4102-a45d-bc00a060bcc6 |

  @account
  Scenario: Generate account
    When I send envelopes to topic "account.generator"
      | Headers.Authorization    |
      | Bearer {{tenant1.token}} |
    Then Envelopes should be in topic "account.generator"
    Then Envelopes should be in topic "account.generated"
    And Envelopes should have the following fields
      | From |
      | ~    |

