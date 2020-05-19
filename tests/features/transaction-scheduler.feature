@tx-scheduler
Feature: Transaction Scheduler
  As an external developer
  I want to send use transaction scheduler API to interact with the registered chains

  Scenario: Send contract transaction and start a job
    Given I set authentication method "JWT" with "f30c452b-e5fb-4102-a45d-bc00a060bcc6"
    When I send "GET" request to "{{chain-registry}}/chains?name=geth"
    Then the response code should be 200
    Then I store response field "0.uuid" as "gethUUID"
    When I send "POST" request to "{{tx-scheduler-http}}/transactions/{{gethUUID}}/send" with json:
      """
{
    "idempotencyKey": "test6",
    "params": {
        "from": "0x93f7274c9059e601be4512f656b57b830e019e41",
        "to": "0x93f7274c9059e601be4512f656b57b830e019e23",
        "methodSignature": "constructor()"
    }
}
      """
    Then the response code should be 202
    Then I store response field "schedule.UUID" as "scheduleUUID"
    And Response should have the following fields:
      | idempotencyKey | params.methodSignature | schedule.uuid | schedule.chainUUID | schedule.jobs.0.uuid | schedule.jobs[0].status
      | test6          | constructor()          | ~             | {{gethUUID}}       | ~                    | STARTED


  Scenario: New JOB started step by step
    Given I set authentication method "JWT" with "f30c452b-e5fb-4102-a45d-bc00a060bcc6"
    When I send "GET" request to "{{chain-registry}}/chains?name=geth"
    Then the response code should be 200
    Then I store response field "0.uuid" as "gethUUID"
    When I send "POST" request to "{{tx-scheduler-http}}/schedules" with json:
      """
{
    "chainUUID": "{{gethUUID}}"
}
      """
    Then the response code should be 200
    Then I store response field "uuid" as "scheduleUUID"
    When I send "POST" request to "{{tx-scheduler-http}}/jobs" with json:
      """
{
	"scheduleUUID": "{{scheduleUUID}}",
	"type": "ETH_SENDRAWTRANSACTION",
    "transaction": {
        "from": "0x93f7274c9059e601be4512f656b57b830e019e41",
        "to": "0x93f7274c9059e601be4512f656b57b830e019e42"
    }
}
      """
    Then the response code should be 200
    Then I store response field "uuid" as "jobUUID"
    And Response should have the following fields:
      | uuid | transaction.from                           | transaction.to                             | status
      | ~    | 0x93f7274c9059e601be4512f656b57b830e019e41 | 0x93f7274c9059e601be4512f656b57b830e019e42 | CREATED
    When I send "PATCH" request to "{{tx-scheduler-http}}/jobs/{{jobUUID}}" with json:
      """
{
    "transaction": {
        "from": "0x93f7274c9059e601be4512f656b57b830e019e43",
        "to": "0x93f7274c9059e601be4512f656b57b830e019e44"
    }
}
      """
    Then the response code should be 200
    When I send "PUT" request to "{{tx-scheduler-http}}/jobs/{{jobUUID}}/start"
    Then the response code should be 202
    When I send "GET" request to "{{tx-scheduler-http}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields:
      | uuid | transaction.from                           | transaction.to                             | status
      | ~    | 0x93f7274c9059e601be4512f656b57b830e019e43 | 0x93f7274c9059e601be4512f656b57b830e019e44 | STARTED
