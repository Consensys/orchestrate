@account-management
Feature: Account management
  As as external developer
  I want to generate new accounts and use them to sign transactions

  Background:
    Given I have the following tenants
      | alias   | tenantID        |
      | tenant1 | {{random.uuid}} |
    Given I register the following contracts
      | name        | artifacts        | Headers.Authorization    |
      | SimpleToken | SimpleToken.json | Bearer {{tenant1.token}} |
    And I register the following chains
      | alias | Name                | URLs                         | Headers.Authorization    |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu_1.URLs}} | Bearer {{tenant1.token}} |

  Scenario: Import account and update it and sign with it
    Given I register the following alias
      | alias          | value           |
      | sendTxID       | {{random.uuid}} |
      | generateAccID  | {{random.uuid}} |
      | generateAccID2 | {{random.uuid}} |
    Given I have the following account
      | alias     |
      | importAcc |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.identity-manager}}/accounts/import" with json:
  """
{
    "alias": "{{generateAccID}}",
    "privateKey": "{{importAcc.private_key}}",
    "attributes": {
    	"scenario_id": "{{scenarioID}}"
    }
}
      """
    Then the response code should be 200
    And Response should have the following fields
      | alias             | attributes.scenario_id | address               |
      | {{generateAccID}} | {{scenarioID}}         | {{importAcc.address}} |
    When I send "PATCH" request to "{{global.identity-manager}}/accounts/{{importAcc.address}}" with json:
  """
{
    "alias": "{{generateAccID2}}", 
    "attributes": {
    	"new_attribute": "{{scenarioID}}"
    }
}
      """
    Then the response code should be 200
    Then I send "GET" request to "{{global.identity-manager}}/accounts/{{importAcc.address}}"
    Then the response code should be 200
    And Response should have the following fields
      | alias              | attributes.new_attribute |
      | {{generateAccID2}} | {{scenarioID}}           |
    Then I track the following envelopes
      | ID           |
      | {{sendTxID}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "contractName": "SimpleToken",
        "from": "{{importAcc.address}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{sendTxID}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"

  Scenario: Generate account and send transaction
    Given I register the following alias
      | alias         | value           |
      | generateAccID | {{random.uuid}} |
      | sendTxID      | {{random.uuid}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.identity-manager}}/accounts" with json:
  """
{
    "alias": "{{generateAccID}}", 
    "attributes": {
    	"scenario_id": "{{scenarioID}}"
    }
}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias            | path    |
      | generatedAccAddr | address |
    Then I track the following envelopes
      | ID           |
      | {{sendTxID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
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
    	"id": "{{sendTxID}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"

  Scenario: Sending transaction with not existed account
    Given I register the following alias
      | alias    | value              |
      | fromAcc  | {{random.account}} |
      | sendTxID | {{random.uuid}}    |
    Then I track the following envelopes
      | ID           |
      | {{sendTxID}} |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/deploy-contract" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "contractName": "SimpleToken",
        "from": "{{fromAcc}}"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "{{sendTxID}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.recover"
    And Envelopes should have the following fields
      | Errors.0.Message                 |
      | no key for account "{{fromAcc}}" |

