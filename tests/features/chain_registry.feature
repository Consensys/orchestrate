@chain-registry
Feature: chain registry

  Scenario: get chain data with API key
    Given I set authentication method "API-Key" with "with-key"
    When I send "GET" request to "{{chain-registry}}/chains"
    Then the response code should be 200

  Scenario: get chain data with JWT
    Given I set authentication method "JWT" with "f30c452b-e5fb-4102-a45d-bc00a060bcc6"
    When I send "GET" request to "{{chain-registry}}/chains"
    Then the response code should be 200

  Scenario: Add and remove a chain with API key
    Given I set authentication method "API-Key" with "with-key"
    When I send "POST" request to "{{chain-registry}}/chains" with json:
      """
      {
        "name": "gethTemp",
        "urls": [
          "http://geth:8545"
        ],
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

    When I send "GET" request to "{{chain-registry}}/chains/{{gethTempUUID}}"
    Then the response code should be 200

    When I send "POST" request to "{{chain-registry}}/chains" with json:
      """
      {
        "name": "gethTemp",
        "urls": [
          "http://geth:8545"
        ],
        "listener": {
          "depth": 1,
          "fromBlock": "1",
          "backOffDuration": "1s"
        }
      }
      """
    Then the response code should be 409

    When I send "DELETE" request to "{{chain-registry}}/chains/{{gethTempUUID}}"
    Then the response code should be 204

    When I send "GET" request to "{{chain-registry}}/chains/{{gethTempUUID}}"
    Then the response code should be 404

  Scenario: Add and remove a chain with JWT token
    Given I set authentication method "JWT" with "f30c452b-e5fb-4102-a45d-bc00a060bcc6"
    When I send "POST" request to "{{chain-registry}}/chains" with json:
      """
      {
        "name": "gethTemp",
        "urls": [
          "http://geth:8545"
        ],
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

    When I send "GET" request to "{{chain-registry}}/chains/{{gethTempUUID}}"
    Then the response code should be 200

    When I send "POST" request to "{{chain-registry}}/chains" with json:
      """
      {
        "name": "gethTemp",
        "urls": [
          "http://geth:8545"
        ],
        "listener": {
          "depth": 1,
          "fromBlock": "1",
          "backOffDuration": "1s"
        }
      }
      """
    Then the response code should be 409

    When I send "DELETE" request to "{{chain-registry}}/chains/{{gethTempUUID}}"
    Then the response code should be 204

    When I send "GET" request to "{{chain-registry}}/chains/{{gethTempUUID}}"
    Then the response code should be 404

  Scenario: Patch chain with API key
    Given I set authentication method "API-Key" with "with-key"
    When I send "POST" request to "{{chain-registry}}/chains" with json:
      """
      {
        "name": "gethTemp2",
        "urls": [
          "http://geth:8545"
        ],
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
    
    When I send "PATCH" request to "{{chain-registry}}/chains/{{gethTemp2UUID}}" with json:
      """
      {
        "listener": {
          "backOffDuration": "1000"
        }
      }
      """
    Then the response code should be 400
    
    When I send "PATCH" request to "{{chain-registry}}/chains/{{gethTemp2UUID}}" with json:
      """
      {
        "urls": [
          "&£$&£$%"
        ]
      }
      """
    Then the response code should be 400

    When I send "PATCH" request to "{{chain-registry}}/chains/{{gethTemp2UUID}}" with json:
      """
      {
        "listener": {
          "backOffDuration": "3s"
        }
      }
      """
    Then the response code should be 200

    When I send "DELETE" request to "{{chain-registry}}/chains/{{gethTemp2UUID}}"
    Then the response code should be 204

  Scenario: Patch chain with JWT
    Given I set authentication method "JWT" with "f30c452b-e5fb-4102-a45d-bc00a060bcc6"
    When I send "POST" request to "{{chain-registry}}/chains" with json:
      """
      {
        "name": "gethTemp2",
        "urls": [
          "http://geth:8545"
        ],
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

    When I send "PATCH" request to "{{chain-registry}}/chains/{{gethTemp2UUID}}" with json:
      """
      {
        "listener": {
          "backOffDuration": "3s"
        }
      }
      """
    Then the response code should be 200

    When I send "DELETE" request to "{{chain-registry}}/chains/{{gethTemp2UUID}}"
    Then the response code should be 204
  
  Scenario: Fail to register chains with invalid values
    When I send "POST" request to "{{chain-registry}}/chains" with json:
      """
      {
        "name": "gethInvalid",
        "urls": [
          "http://geth:8545"
        ],
        "listener": {
          "depth": 1,
          "fromBlock": "1",
          "backOffDuration": "1000"
        }
      }
      """
    Then the response code should be 400
    
    When I send "POST" request to "{{chain-registry}}/chains" with json:
      """
      {
        "name": "gethInvalid",
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
    
    When I send "POST" request to "{{chain-registry}}/chains" with json:
      """
      {
        "name": "gethInvalid",
        "urls": [],
        "listener": {
          "depth": 1,
          "fromBlock": "1",
          "backOffDuration": "1s"
        }
      }
      """
    Then the response code should be 400
