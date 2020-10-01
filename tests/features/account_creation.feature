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

  Scenario: Generate account
    Given I register the following alias
      | alias         | value           |
      | generateAccID | {{random.uuid}} |
      | sendTxID      | {{random.uuid}} |
    When I send envelopes to topic "account.generator"
      | ID                | Headers.Authorization    |
      | {{generateAccID}} | Bearer {{tenant1.token}} |
    Then Envelopes should be in topic "account.generator"
    Then Envelopes should be in topic "account.generated"
    And Envelopes should have the following fields
      | From | ID                |
      | ~    | {{generateAccID}} |
    And I register the following envelope fields
      | id                | alias            | path |
      | {{generateAccID}} | generatedAccAddr | From |
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
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"

  Scenario: Account not found
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

