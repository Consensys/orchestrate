@account-management
@multi-tenancy
Feature: Account management
  As as external developer
  I want to generate new accounts and use them to sign transactions

  Background:
    Given I have the following tenants
      | alias         | tenantID | username |
      | tenantFoo     | foo      |          |
      | tenantFooUser | foo      | userFoo  |
      | tenantBar     | bar      |          |
      | wildcard      | _        |          |
    Given I register the following contracts
      | name        | artifacts        | API-KEY            | Tenant                |
      | SimpleToken | SimpleToken.json | {{global.api-key}} | {{wildcard.tenantID}} |

  Scenario: Accounts own by default tenant can be used by other authorized tenants
    Given I register the following alias
      | alias            | value           |
      | generateAccID    | {{random.uuid}} |
      | generateAccID2   | {{random.uuid}} |
      | fooSendTxID      | {{random.uuid}} |
      | wildcardSendTxID | {{random.uuid}} |
    Given I set the headers
      | Key       | Value              |
      | X-API-KEY | {{global.api-key}} |
    When I send "POST" request to "{{global.api}}/accounts" with json:
  """
{
    "alias": "{{generateAccID}}", 
    "attributes": {
    	"scenario_id": "{{scenarioID}}"
    }
}
      """
    Then the response code should be 200
    And Response should have the following fields
      | alias             | attributes.scenario_id | tenantID |
      | {{generateAccID}} | {{scenarioID}}         | _        |
    Then I register the following response fields
      | alias            | path    |
      | generatedAccAddr | address |
    Then I track the following envelopes
      | ID              |
      | {{fooSendTxID}} |
    Given I set the headers
      | Key         | Value                  |
      | X-API-KEY   | {{global.api-key}}     |
      | X-TENANT-ID | {{tenantFoo.tenantID}} |
    When I send "PATCH" request to "{{global.api}}/accounts/{{generatedAccAddr}}" with json:
  """
{
    "alias": "{{generateAccID2}}"
}
      """
    Then the response code should be 200
    Then I send "GET" request to "{{global.api}}/accounts/{{generatedAccAddr}}"
    Then the response code should be 200
    And Response should have the following fields
      | alias              |
      | {{generateAccID2}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
  """
{
    "chain": "{{chain.besu0.Name}}",
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
      | Key         | Value                  |
      | X-API-KEY   | {{global.api-key}}     |
      | X-TENANT-ID | {{tenantFoo.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
  """
{
    "chain": "{{chain.besu0.Name}}",
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
      | Key         | Value                 |
      | X-API-KEY   | {{global.api-key}}    |
      | X-TENANT-ID | {{wildcard.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
  """
{
    "chain": "{{chain.besu0.Name}}",
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

  Scenario: Accounts own by tenant foo can be access only by tenant foo
    Given I register the following alias
      | alias            | value           |
      | generateAccID    | {{random.uuid}} |
      | generateAccID2   | {{random.uuid}} |
      | fooSendTxID      | {{random.uuid}} |
      | barSendTxID      | {{random.uuid}} |
      | wildcardSendTxID | {{random.uuid}} |
    Given I set the headers
      | Key         | Value                  |
      | X-API-KEY   | {{global.api-key}}     |
      | X-TENANT-ID | {{tenantFoo.tenantID}} |
    When I send "POST" request to "{{global.api}}/accounts" with json:
  """
{
    "alias": "{{generateAccID}}", 
    "attributes": {
    	"scenario_id": "{{scenarioID}}"
    }
}
      """
    Then the response code should be 200
    And Response should have the following fields
      | alias             | attributes.scenario_id | tenantID               |
      | {{generateAccID}} | {{scenarioID}}         | {{tenantFoo.tenantID}} |
    Then I register the following response fields
      | alias            | path    |
      | generatedAccAddr | address |
    Then I track the following envelopes
      | ID              |
      | {{fooSendTxID}} |
    Given I set the headers
      | Key         | Value                  |
      | X-API-KEY   | {{global.api-key}}     |
      | X-TENANT-ID | {{tenantFoo.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
  """
{
    "chain": "{{chain.besu0.Name}}",
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
      | Key         | Value                  |
      | X-API-KEY   | {{global.api-key}}     |
      | X-TENANT-ID | {{tenantBar.tenantID}} |
    When I send "PATCH" request to "{{global.api}}/accounts/{{generatedAccAddr}}" with json:
  """
{
    "alias": "{{generateAccID2}}"
}
      """
    Then the response code should be 404
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
  """
{
    "chain": "{{chain.besu0.Name}}",
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
    Then the response code should be 422
    Then I track the following envelopes
      | ID                   |
      | {{wildcardSendTxID}} |
    Given I set the headers
      | Key         | Value                  |
      | X-API-KEY   | {{global.api-key}}     |
      | X-TENANT-ID | {{tenantFoo.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
  """
{
    "chain": "{{chain.besu0.Name}}",
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

  Scenario: Accounts own by an specific user  can be access only by itself
    Given I register the following alias
      | alias          | value           |
      | userSendTxID   | {{random.uuid}} |
      | tenantSendTxID | {{random.uuid}} |
    Given I register the following alias
      | alias         | value           |
      | generateAccID | {{random.uuid}} |
    Given I set the headers
      | Key         | Value                      |
      | X-API-KEY   | {{global.api-key}}         |
      | X-TENANT-ID | {{tenantFooUser.tenantID}} |
      | X-USERNAME  | {{tenantFooUser.username}} |
    When I send "POST" request to "{{global.api}}/accounts" with json:
  """
{
    "alias": "{{generateAccID}}", 
    "attributes": {
    	"scenario_id": "{{scenarioID}}"
    }
}
      """
    Then the response code should be 200
    And Response should have the following fields
      | alias             | attributes.scenario_id | tenantID                   | ownerID                    |
      | {{generateAccID}} | {{scenarioID}}         | {{tenantFooUser.tenantID}} | {{tenantFooUser.username}} |
    Then I register the following response fields
      | alias            | path    |
      | generatedAccAddr | address |
    Then I track the following envelopes
      | ID               |
      | {{userSendTxID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
  """
{
    "chain": "{{chain.besu0.Name}}",
    "params": {
        "contractName": "SimpleToken",
        "from": "{{generatedAccAddr}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{userSendTxID}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"
    Then I track the following envelopes
      | ID                 |
      | {{tenantSendTxID}} |
    Given I set the headers
      | Key         | Value                  |
      | X-API-KEY   | {{global.api-key}}     |
      | X-TENANT-ID | {{tenantFoo.tenantID}} |
    When I send "POST" request to "{{global.api}}/transactions/deploy-contract" with json:
  """
{
    "chain": "{{chain.besu0.Name}}",
    "params": {
        "contractName": "SimpleToken",
        "from": "{{generatedAccAddr}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{tenantSendTxID}}"
    }
}
      """
    Then the response code should be 422
