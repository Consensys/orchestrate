@multi-tenancy
Feature: Authentication on API

  Scenario: get authentication with API Key on chain-registry
    Given I set authentication method "API-Key" with "with-key"
    When I send "GET" request to "{{chain-registry}}/chains"
    Then the response code should be 200

  Scenario: try to authentication with bad API Key on chain-registry
    Given I set authentication method "API-Key" with "bad-key"
    When I send "GET" request to "{{chain-registry}}/chains"
    Then the response code should be 401

  Scenario: get authentication with JWT on chain-registry
    Given I set authentication method "JWT" with "f30c452b-e5fb-4102-a45d-bc00a060bcc6"
    When I send "GET" request to "{{chain-registry}}/chains"
    Then the response code should be 200

  Scenario: get authentication with bad JWT on chain-registry
    Given I set authentication method "JWT" with "Hello!"
    When I send "GET" request to "{{chain-registry}}/chains"
    Then the response code should be 200

  Scenario: try to authentication without auth method on chain-registry
    When I send "GET" request to "{{chain-registry}}/chains"
    Then the response code should be 401


  Scenario: get authentication with API Key on contract-registry
    Given I set authentication method "API-Key" with "with-key"
    When I send "GET" request to "{{contract-registry-http}}/contracts"
    Then the response code should be 200

  # TODO: Should return 401 and returns 500 now
  # Scenario: try to authentication with bad API Key on contract-registry
  #   Given I set authentication method "API-Key" with "bad-key"
  #   When I send "GET" request to "{{contract-registry-http}}/contracts"
  #   Then the response code should be 401

  Scenario: get authentication with JWT on contract-registry
    Given I set authentication method "JWT" with "f30c452b-e5fb-4102-a45d-bc00a060bcc6"
    When I send "GET" request to "{{contract-registry-http}}/contracts"
    Then the response code should be 200

  Scenario: get authentication with bad JWT on contract-registry
    Given I set authentication method "JWT" with "Hello!"
    When I send "GET" request to "{{contract-registry-http}}/contracts"
    Then the response code should be 200

# TODO: Should return 401 and returns 500 now
# Scenario: try to authentication without auth method on contract-registry
#   When I send "GET" request to "{{contract-registry-http}}/contracts"
#   Then the response code should be 401
