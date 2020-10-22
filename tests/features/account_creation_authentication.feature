@account-management
@multi-tenancy
Feature: Account management
  As as external developer
  I want to generate new accounts and use them to sign transactions

  Background:
    Given I have the following tenants
      | alias     | tenantID |
      | tenantFoo | foo      |
      | tenantBar | bar      |
      | wildcard  | *        |
    Given I register the following contracts
      | name        | artifacts        | Headers.Authorization     |
      | SimpleToken | SimpleToken.json | Bearer {{wildcard.token}} |
    And I register the following chains
      | alias | Name                | URLs                         | Headers.Authorization     |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu_1.URLs}} | Bearer {{wildcard.token}} |


  Scenario: Generate account as default tenant
    Given I register the following alias
      | alias            | value           |
      | generateAccID    | {{random.uuid}} |
      | fooSendTxID      | {{random.uuid}} |
      | wildcardSendTxID | {{random.uuid}} |
    And I have created the following accounts
      | alias            | ID              | Headers.Authorization     |
      | generatedAccAddr | {{random.uuid}} | Bearer {{wildcard.token}} |
    Then I track the following envelopes
      | ID              |
      | {{fooSendTxID}} |
    Given I set the headers
      | Key           | Value                      |
      | Authorization | Bearer {{tenantFoo.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "contractName": "SimpleToken",
        "from": "{{generatedAccAddr}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{fooSendTxID}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    Then I track the following envelopes
      | ID              |
      | {{fooSendTxID}} |
    Given I set the headers
      | Key           | Value                      |
      | Authorization | Bearer {{tenantFoo.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "contractName": "SimpleToken",
        "from": "{{generatedAccAddr}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{fooSendTxID}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    Then I track the following envelopes
      | ID                   |
      | {{wildcardSendTxID}} |
    Given I set the headers
      | Key           | Value                     |
      | Authorization | Bearer {{wildcard.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "contractName": "SimpleToken",
        "from": "{{generatedAccAddr}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{wildcardSendTxID}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    
  Scenario: Generate account as tenant foo
    Given I register the following alias
      | alias            | value           |
      | generateAccID    | {{random.uuid}} |
      | fooSendTxID      | {{random.uuid}} |
      | barSendTxID      | {{random.uuid}} |
      | wildcardSendTxID | {{random.uuid}} |
    And I have created the following accounts
      | alias            | ID              | Headers.Authorization      |
      | generatedAccAddr | {{random.uuid}} | Bearer {{tenantFoo.token}} |
    Then I track the following envelopes
      | ID              |
      | {{fooSendTxID}} |
    Given I set the headers
      | Key           | Value                      |
      | Authorization | Bearer {{tenantFoo.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "contractName": "SimpleToken",
        "from": "{{generatedAccAddr}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{fooSendTxID}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    Then I track the following envelopes
      | ID              |
      | {{barSendTxID}} |
    Given I set the headers
      | Key           | Value                      |
      | Authorization | Bearer {{tenantBar.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "contractName": "SimpleToken",
        "from": "{{generatedAccAddr}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{barSendTxID}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.recover"
    And Envelopes should have the following fields
      | Errors.0.Message                          |
      | no key for account "{{generatedAccAddr}}" |
    Then I track the following envelopes
      | ID                   |
      | {{wildcardSendTxID}} |
    Given I set the headers
      | Key           | Value                     |
      | Authorization | Bearer {{wildcard.token}} |
      | X-Tenant-ID   | foo                       |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "contractName": "SimpleToken",
        "from": "{{generatedAccAddr}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{wildcardSendTxID}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"

