@chain-registry
Feature: chain registry

  Scenario: get chain data
    Given I set authentication method "API-Key" with "with-key"
    When I send "GET" request to "{{chain-registry}}/chains"
    Then the response code should be 200

    When I send "GET" request to "{{chain-registry}}/_/chains"
    Then the response code should be 200


  Scenario: Add and remove a chain
    Given I set authentication method "API-Key" with "with-key"
    Given I send "DELETE" request to "{{chain-registry}}/_/chains/gethTemp"
    When I send "POST" request to "{{chain-registry}}/_/chains" with json:
    """
    {
      "name": "gethTemp",
      "urls": ["http://geth:8545"],
      "listener": {
        "depth": 1,
        "blockPosition": "1",
        "backOffDuration": "1s",
        "externalTxEnabled": true
      }
    }
    """
    Then the response code should be 200
    Then I store the UUID as "gethTempUUID"

    When I send "GET" request to "{{chain-registry}}/chains/{{gethTempUUID}}"
    Then the response code should be 200

    When I send "GET" request to "{{chain-registry}}/_/chains/gethTemp"
    Then the response code should be 200

    When I send "POST" request to "{{chain-registry}}/_/chains" with json:
    """
    {
      "name": "gethTemp",
      "urls": ["http://geth:8545"],
      "listener": {
        "depth": 1,
        "blockPosition": "1",
        "backOffDuration": "1s"
      }
    }
    """
    Then the response code should be 409

    When I send "POST" request to "{{chain-registry}}/chains" with json:
    """
    {
      "name": "gethTemp",
      "urls": ["http://geth:8545"],
      "listener": {
        "depth": 1,
        "blockPosition": "1",
        "backOffDuration": "1s"
      }
    }
    """
    Then the response code should be 409
    When I send "DELETE" request to "{{chain-registry}}/_/chains/gethTemp"
    Then the response code should be 200

    When I send "GET" request to "{{chain-registry}}/_/chains/gethTemp"
    Then the response code should be 500


  Scenario: Patch chain
    Given I set authentication method "API-Key" with "with-key"
    Given I send "DELETE" request to "{{chain-registry}}/_/chains/gethTemp2"
    When I send "POST" request to "{{chain-registry}}/_/chains" with json:
    """
    {
      "name": "gethTemp2",
      "urls": ["http://geth:8545"],
      "listener": {
        "depth": 1,
        "blockPosition": "1",
        "backOffDuration": "1s",
        "externalTxEnabled": true
      }
    }
    """
    Then the response code should be 200
    Then I store the UUID as "gethTemp2UUID"

    When I send "PATCH" request to "{{chain-registry}}/_/chains/gethTemp2" with json:
    """
    {
      "listener": {
        "backOffDuration": "2s",
        "externalTxEnabled": false
      }
    }
    """
    Then the response code should be 200

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
    Then the response code should be 200
