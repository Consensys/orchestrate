@tx-scheduler
@multi-tenancy
Feature: Transaction Scheduler Jobs
  As an external developer
  I want to send use transaction scheduler API to interact with the registered chains

  Background:
    Given I have the following tenants
      | alias         | tenantID |
      | tenantFoo     | foo      |
      | tenantBar     | bar      |
      | tenantDefault | _        |
    Then I register the following contracts
      | name        | artifacts        | Headers.Authorization          |
      | SimpleToken | SimpleToken.json | Bearer {{tenantDefault.token}} |
    Then I register the following chains
      | alias | Name                | URLs                         | Headers.Authorization          |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu_1.URLs}} | Bearer {{tenantDefault.token}} |
    And I have created the following accounts
      | alias    | ID              | Headers.Authorization          |
      | account1 | {{random.uuid}} | Bearer {{tenantDefault.token}} |
    Then I track the following envelopes
      | ID                  |
      | faucet-{{account1}} |
    Given I set the headers
      | Key           | Value                          |
      | Authorization | Bearer {{tenantDefault.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/transfer" with json:
  """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
      "from": "{{global.nodes.besu_1.fundedPublicKeys[0]}}",
      "to": "{{account1}}",
      "value": "100000000000000000"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}",
    	"id": "faucet-{{account1}}"
    }
}
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.sender"
    Then Envelopes should be in topic "tx.decoded"

  @besu
  Scenario: Cannot start or update other tenant jobs
    Given I register the following alias
      | alias | value              |
      | to1   | {{random.account}} |
    Then  I set the headers
      | Key           | Value                      |
      | Authorization | Bearer {{tenantFoo.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/schedules" with json:
      """
{}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias           | path |
      | scheduleOneUUID | uuid |
    When I send "POST" request to "{{global.tx-scheduler}}/jobs" with json:
      """
{
	"scheduleUUID": "{{scheduleOneUUID}}",
	"chainUUID": "{{besu.UUID}}",
	"type": "eth://ethereum/transaction",
    "transaction": {
        "from": "{{account1}}",
        "to": "{{to1}}",
        "value": "100000"
    }
}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias     | path |
      | txJobUUID | uuid |
    Then I track the following envelopes
      | ID            |
      | {{txJobUUID}} |
    Then  I set the headers
      | Key           | Value                      |
      | Authorization | Bearer {{tenantBar.token}} |
    When I send "PATCH" request to "{{global.tx-scheduler}}/jobs/{{txJobUUID}}" with json:
      """
{
    "transaction": {
        "from": "{{account1}}",
        "to": "{{to1}}",
        "value": "100000"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}"
    }
}
      """
    Then the response code should be 404
    When I send "PUT" request to "{{global.tx-scheduler}}/jobs/{{txJobUUID}}/start"
    Then the response code should be 404
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{txJobUUID}}"
    Then the response code should be 404
    Then  I set the headers
      | Key           | Value                      |
      | Authorization | Bearer {{tenantFoo.token}} |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{txJobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | uuid          | status  |
      | {{txJobUUID}} | CREATED |
    When I send "PUT" request to "{{global.tx-scheduler}}/jobs/{{txJobUUID}}/start"
    Then Envelopes should be in topic "tx.decoded"
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{txJobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | uuid          | status  |
      | {{txJobUUID}} | MINED |
