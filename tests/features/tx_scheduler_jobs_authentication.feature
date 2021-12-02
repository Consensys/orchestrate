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
      | name        | artifacts        | API-KEY            |
      | SimpleToken | SimpleToken.json | {{global.api-key}} |
    Then I register the following chains
      | alias | Name                | URLs                          | API-KEY            |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu[0].URLs}} | {{global.api-key}} |
    And I have created the following accounts
      | alias    | ID              | API-KEY            |
      | account1 | {{random.uuid}} | {{global.api-key}} |
    Then I track the following envelopes
      | ID                  |
      | faucet-{{account1}} |
    Given I set the headers
      | Key       | Value              |
      | X-API-KEY | {{global.api-key}} |
    When I send "POST" request to "{{global.api}}/transactions/transfer" with json:
      """
      {
        "chain": "besu-{{scenarioID}}",
        "params": {
          "from": "{{global.nodes.besu[0].fundedPublicKeys[0]}}",
          "to": "{{account1}}",
          "value": "0x16345785D8A0000"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}",
          "id": "faucet-{{account1}}"
        }
      }
      """
    Then the response code should be 202
    Then Envelopes should be in topic "tx.decoded"

  @besu
  Scenario: Cannot start or update other tenant jobs
    Given I register the following alias
      | alias | value              |
      | to1   | {{random.account}} |
    Then  I set the headers
      | Key         | Value                  |
      | X-API-KEY   | {{global.api-key}}     |
      | X-TENANT-ID | {{tenantFoo.tenantID}} |
    When I send "POST" request to "{{global.api}}/schedules" with json:
      """
      {}
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias           | path |
      | scheduleOneUUID | uuid |
    When I send "POST" request to "{{global.api}}/jobs" with json:
      """
      {
        "scheduleUUID": "{{scheduleOneUUID}}",
        "chainUUID": "{{besu.UUID}}",
        "type": "eth://ethereum/transaction",
        "transaction": {
          "from": "{{account1}}",
          "to": "{{to1}}",
          "value": "0x186A0"
        }
      }
      """
    Then the response code should be 200
    Then I register the following response fields
      | alias     | path |
      | txJobUUID | uuid |
    Then I track the following envelopes
      | ID                  |
      | {{scheduleOneUUID}} |
    Then  I set the headers
      | Key         | Value                  |
      | X-API-KEY   | {{global.api-key}}     |
      | X-TENANT-ID | {{tenantBar.tenantID}} |
    When I send "PATCH" request to "{{global.api}}/jobs/{{txJobUUID}}" with json:
      """
      {
        "transaction": {
          "from": "{{account1}}",
          "to": "{{to1}}",
          "value": "0x186A0"
        },
        "labels": {
          "scenario.id": "{{scenarioID}}"
        }
      }
      """
    Then the response code should be 404
    When I send "PUT" request to "{{global.api}}/jobs/{{txJobUUID}}/start"
    Then the response code should be 404
    When I send "GET" request to "{{global.api}}/jobs/{{txJobUUID}}"
    Then the response code should be 404
    Then  I set the headers
      | Key         | Value                  |
      | X-API-KEY   | {{global.api-key}}     |
      | X-TENANT-ID | {{tenantFoo.tenantID}} |
    When I send "GET" request to "{{global.api}}/jobs/{{txJobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | uuid          | status  |
      | {{txJobUUID}} | CREATED |
    When I send "PUT" request to "{{global.api}}/jobs/{{txJobUUID}}/start"
    Then Envelopes should be in topic "tx.decoded"
    When I send "GET" request to "{{global.api}}/jobs/{{txJobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | uuid          | status |
      | {{txJobUUID}} | MINED  |
