@chain-registry
Feature: Chain registry
  As as external developer
  I want to register new chains

  Scenario: get chain data with API key
    Given I set the headers
      | Key       | Value       |
      | X-API-Key | {{global.api-key}} |
    When I send "GET" request to "{{global.chain-registry}}/chains"
    Then the response code should be 200

  Scenario: get chain data with JWT
    Given I have the following tenants
      | alias   |
      | tenant1 |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "GET" request to "{{global.chain-registry}}/chains"
    Then the response code should be 200

  Scenario: Add and remove a chain with API key
    Given I set the headers
      | Key       | Value       |
      | X-API-Key | {{global.api-key}} |
    When I send "POST" request to "{{global.chain-registry}}/chains" with json:
      """
      {
        "name": "gethTemp-{{scenarioID}}",
        "urls": {{global.nodes.geth.URLs}},
        "listener": {
          "depth": 1,
          "fromBlock": "1",
          "backOffDuration": "1s",
          "externalTxEnabled": true
        }
      }
      """
    Then the response code should be 200
    Then I store the UUID as "gethTempUUID"

    When I send "GET" request to "{{global.chain-registry}}/chains/{{gethTempUUID}}"
    Then the response code should be 200

    When I send "POST" request to "{{global.chain-registry}}/chains" with json:
      """
      {
        "name": "gethTemp-{{scenarioID}}",
        "urls": {{global.nodes.geth.URLs}},
        "listener": {
          "depth": 1,
          "fromBlock": "1",
          "backOffDuration": "1s"
        }
      }
      """
    Then the response code should be 409

    When I send "DELETE" request to "{{global.chain-registry}}/chains/{{gethTempUUID}}"
    Then the response code should be 204

    When I send "GET" request to "{{global.chain-registry}}/chains/{{gethTempUUID}}"
    Then the response code should be 404

  Scenario: Add and remove a chain with JWT token
    Given I have the following tenants
      | alias   |
      | tenant1 |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.chain-registry}}/chains" with json:
      """
      {
        "name": "gethTemp-{{scenarioID}}",
        "urls": {{global.nodes.geth.URLs}},
        "listener": {
          "depth": 1,
          "fromBlock": "1",
          "backOffDuration": "1s",
          "externalTxEnabled": true
        }
      }
      """
    Then the response code should be 200
    Then I store the UUID as "gethTempUUID"

    When I send "GET" request to "{{global.chain-registry}}/chains/{{gethTempUUID}}"
    Then the response code should be 200

    When I send "POST" request to "{{global.chain-registry}}/chains" with json:
      """
      {
        "name": "gethTemp-{{scenarioID}}",
        "urls": {{global.nodes.geth.URLs}},
        "listener": {
          "depth": 1,
          "fromBlock": "1",
          "backOffDuration": "1s"
        }
      }
      """
    Then the response code should be 409

    When I send "DELETE" request to "{{global.chain-registry}}/chains/{{gethTempUUID}}"
    Then the response code should be 204

    When I send "GET" request to "{{global.chain-registry}}/chains/{{gethTempUUID}}"
    Then the response code should be 404

  Scenario: Register chain with API key
    Given I set the headers
      | Key       | Value       |
      | X-API-Key | {{global.api-key}} |
    When I send "POST" request to "{{global.chain-registry}}/chains" with json:
      """
      {
        "name": "gethTemp2-{{scenarioID}}",
        "urls": {{global.nodes.geth.URLs}},
        "listener": {
          "depth": 1,
          "fromBlock": "1",
          "backOffDuration": "1s",
          "externalTxEnabled": true
        }
      }
      """
    Then the response code should be 200
    Then I store the UUID as "gethTemp2UUID"

    When I send "PATCH" request to "{{global.chain-registry}}/chains/{{gethTemp2UUID}}" with json:
      """
      {
        "listener": {
          "backOffDuration": "1000"
        }
      }
      """
    Then the response code should be 400

    When I send "PATCH" request to "{{global.chain-registry}}/chains/{{gethTemp2UUID}}" with json:
      """
      {
        "urls": [
          "&£$&£$%"
        ]
      }
      """
    Then the response code should be 400

    When I send "PATCH" request to "{{global.chain-registry}}/chains/{{gethTemp2UUID}}" with json:
      """
      {
        "listener": {
          "backOffDuration": "3s"
        }
      }
      """
    Then the response code should be 200

    When I send "DELETE" request to "{{global.chain-registry}}/chains/{{gethTemp2UUID}}"
    Then the response code should be 204

  Scenario: Update chain with JWT
    Given I have the following tenants
      | alias   |
      | tenant1 |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.chain-registry}}/chains" with json:
      """
      {
        "name": "gethTemp2-{{scenarioID}}",
        "urls": {{global.nodes.geth.URLs}},
        "listener": {
          "depth": 1,
          "fromBlock": "1",
          "backOffDuration": "1s",
          "externalTxEnabled": true
        }
      }
      """
    Then the response code should be 200
    Then I store the UUID as "gethTemp2UUID"

    When I send "PATCH" request to "{{global.chain-registry}}/chains/{{gethTemp2UUID}}" with json:
      """
      {
        "listener": {
          "backOffDuration": "3s"
        }
      }
      """
    Then the response code should be 200

    When I send "DELETE" request to "{{global.chain-registry}}/chains/{{gethTemp2UUID}}"
    Then the response code should be 204

  Scenario: Fail to register chains with invalid values
    Given I have the following tenants
      | alias   |
      | tenant1 |
    Given I set the headers
      | Key           | Value                    |
      | Authorization | Bearer {{tenant1.token}} |
    When I send "POST" request to "{{global.chain-registry}}/chains" with json:
      """
      {
        "name": "gethInvalid-{{scenarioID}}",
        "urls": {{global.nodes.geth.URLs}},
        "listener": {
          "depth": 1,
          "fromBlock": "1",
          "backOffDuration": "1000"
        }
      }
      """
    Then the response code should be 400

    When I send "POST" request to "{{global.chain-registry}}/chains" with json:
      """
      {
        "name": "gethInvalid-{{scenarioID}}",
        "urls": [
          "&£$&£$%"
        ],
        "listener": {
          "depth": 1,
          "fromBlock": "1",
          "backOffDuration": "1s"
        }
      }
      """
    Then the response code should be 400

    When I send "POST" request to "{{global.chain-registry}}/chains" with json:
      """
      {
        "name": "gethInvalid-{{scenarioID}}",
        "urls": [],
        "listener": {
          "depth": 1,
          "fromBlock": "1",
          "backOffDuration": "1s"
        }
      }
      """
    Then the response code should be 400
