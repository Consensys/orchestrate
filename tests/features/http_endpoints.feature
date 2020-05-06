@http-endpoints
Feature: Verify HTTP Endpoints

  Scenario: Get Chain Registry Swagger
    When I send "GET" request to "{{chain-registry}}/swagger/"
    Then the response code should be 200

  Scenario: Get Chain Registry metrics
    When I send "GET" request to "{{chain-registry-metrics}}/metrics"
    Then the response code should be 200

  Scenario: Get Chain Registry readiness
    When I send "GET" request to "{{chain-registry-metrics}}/ready"
    Then the response code should be 200

  Scenario: Get Chain Registry liveness
    When I send "GET" request to "{{chain-registry-metrics}}/live"
    Then the response code should be 200



  Scenario: Get Contract Registry Swagger
    When I send "GET" request to "{{contract-registry-http}}/swagger/"
    Then the response code should be 200

  Scenario: Get Contract Registry metrics
    When I send "GET" request to "{{contract-registry-metrics}}/metrics"
    Then the response code should be 200

  Scenario: Get Contract Registry readiness
    When I send "GET" request to "{{contract-registry-metrics}}/ready"
    Then the response code should be 200

  Scenario: Get Contract Registry liveness
    When I send "GET" request to "{{contract-registry-metrics}}/live"
    Then the response code should be 200



  Scenario: Get Envelope Store Swagger
    When I send "GET" request to "{{envelope-store-http}}/swagger/"
    Then the response code should be 200

  Scenario: Get Envelope Store metrics
    When I send "GET" request to "{{envelope-store-metrics}}/metrics"
    Then the response code should be 200

  Scenario: Get Envelope Store readiness
    When I send "GET" request to "{{envelope-store-metrics}}/ready"
    Then the response code should be 200

  Scenario: Get Envelope Store liveness
    When I send "GET" request to "{{envelope-store-metrics}}/live"
    Then the response code should be 200
