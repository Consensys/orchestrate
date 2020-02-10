@chain-registry
Feature: chain registry

  Scenario: get chain data
    When I send "GET" request to "/chains"
    Then the response code should be 200

    When I send "GET" request to "/_/chains"
    Then the response code should be 200


  Scenario: Add and remove a chain
    Given I send "DELETE" request to "/_/chains/gethTemp"
    When I send "POST" request to "/_/chains" with json:
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

    When I send "GET" request to "/chains/{{gethTempUUID}}"
    Then the response code should be 200

    When I send "GET" request to "/_/chains/gethTemp"
    Then the response code should be 200

    When I send "POST" request to "/_/chains" with json:
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

    When I send "POST" request to "/chains" with json:
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
    When I send "DELETE" request to "/_/chains/gethTemp"
    Then the response code should be 200

    When I send "GET" request to "/_/chains/gethTemp"
    Then the response code should be 500


  Scenario: Patch chain
    Given I send "DELETE" request to "/_/chains/gethTemp2"
    When I send "POST" request to "/_/chains" with json:
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

    When I send "PATCH" request to "/_/chains/gethTemp2" with json:
    """
    {
      "listener": {
        "backOffDuration": "2s",
        "externalTxEnabled": false
      }
    }
    """
    Then the response code should be 200

    When I send "PATCH" request to "/chains/{{gethTemp2UUID}}" with json:
    """
    {
      "listener": {
        "backOffDuration": "3s"
      }
    }
    """
    Then the response code should be 200

    When I send "DELETE" request to "/chains/{{gethTemp2UUID}}"
    Then the response code should be 200
