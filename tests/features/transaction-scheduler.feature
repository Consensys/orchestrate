@tx-scheduler
Feature: Transaction Scheduler
  As an external developer
  I want to send use transaction scheduler API to interact with the registered chains

  Background:
    Given I have the following tenants
      | alias   |
      | tenant1 |
    Given I register the following chains
      | alias | Name                | URLs                       | Headers.Authorization    |
      | besu  | besu-{{scenarioID}} | {{global.nodes.besu.URLs}} | Bearer {{tenant1.token}} |
    And I wait "1.5s"

  Scenario: Send contract transaction and start a job
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/send" with json:
      """
{
    "chain": "besu-{{scenarioID}}",
    "params": {
        "from": "0x93f7274c9059e601be4512f656b57b830e019e41",
        "to": "0x93f7274c9059e601be4512f656b57b830e019e23",
        "methodSignature": "constructor()"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}"
    }
}
      """
    Then the response code should be 202
    Then I store response field "schedule.jobs.0.uuid" as "jobUUID"
    Then Register new envelope tracker "jobUUID"
    And Response should have the following fields
      | params.methodSignature | schedule.uuid | schedule.jobs.0.uuid | schedule.jobs[0].status
      | constructor()          | ~             | ~                    | STARTED

  Scenario: New JOB started step by step
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/schedules" with json:
      """
{}
      """
    Then the response code should be 200
    Then I store response field "uuid" as "scheduleUUID"
    When I send "POST" request to "{{global.tx-scheduler}}/jobs" with json:
      """
{
	"scheduleUUID": "{{scheduleUUID}}",
	"chainUUID": "{{besu.UUID}}",
	"type": "ETH_SENDRAWTRANSACTION",
    "transaction": {
        "from": "0x93f7274c9059e601be4512f656b57b830e019e41",
        "to": "0x93f7274c9059e601be4512f656b57b830e019e42"
    }
}
      """
    Then the response code should be 200
    Then I store response field "uuid" as "jobUUID"
    And Response should have the following fields
      | uuid | chainUUID     | transaction.from                           | transaction.to                             | status
      | ~    | {{besu.UUID}} | 0x93f7274c9059e601be4512f656b57b830e019e41 | 0x93f7274c9059e601be4512f656b57b830e019e42 | CREATED
    When I send "PATCH" request to "{{global.tx-scheduler}}/jobs/{{jobUUID}}" with json:
      """
{
    "transaction": {
        "from": "0x6009608a02a7a15fd6689d6dad560c44e9ab61ff",
        "to": "0x93f7274c9059e601be4512f656b57b830e019e44"
    },
    "labels": {
    	"scenario.id": "{{scenarioID}}"
    },
    "status": "PENDING"
}
      """
    Then the response code should be 200
    When I send "PUT" request to "{{global.tx-scheduler}}/jobs/{{jobUUID}}/start"
    Then the response code should be 202
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | uuid | transaction.from                           | transaction.to                             | status
      | ~    | 0x6009608a02a7a15fd6689d6dad560c44e9ab61ff | 0x93f7274c9059e601be4512f656b57b830e019e44 | PENDING
