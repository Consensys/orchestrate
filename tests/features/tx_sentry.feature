@tx-sentry
Feature: Tx Sentry

  Scenario: Retry settings are persisted
    Given I have the following tenants
      | alias   | tenantID        |
      | tenant1 | {{random.uuid}} |
    Then I register the following chains
      | alias  | Name                  | URLs                         | Headers.Authorization    |
      | besu_1 | besu_1-{{scenarioID}} | {{global.nodes.besu_1.URLs}} | Bearer {{tenant1.token}} |
    Then I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    And I have created the following accounts
      | alias    | ID              | Headers.Authorization    |
      | account1 | {{random.uuid}} | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.tx-scheduler}}/transactions/transfer" with json:
  """
{
    "chain": "besu_1-{{scenarioID}}",
    "params": {
      "from": "{{account1}}",
      "to": "0x0000000000000000000000000000000000000000",
      "value": "100000",
      "retry": {
        "interval": "1m",
        "gasPriceIncrementLevel": "low",
        "gasPriceLimit": 1.2
      }
    }
}
      """
    Then the response code should be 202
    Then I register the following response fields
      | alias   | path                  |
      | jobUUID | schedule.jobs[0].uuid |
    When I send "GET" request to "{{global.tx-scheduler}}/jobs/{{jobUUID}}"
    Then the response code should be 200
    And Response should have the following fields
      | annotations.retry.interval | annotations.retry.gasPriceIncrementLevel | annotations.retry.gasPriceLimit |
      | 1m                         | low                                      | 1.2                             |
