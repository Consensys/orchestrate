@http-endpoints
Feature: Verify HTTP Endpoints

  Scenario: Get Chain Registry Swagger
    When I send "GET" request to "{{global.chain-registry}}/swagger/"
    Then the response code should be 200

  Scenario: Get Chain Registry Swagger spec
    When I send "GET" request to "{{global.chain-registry}}/swagger/swagger.json"
    Then the response code should be 200

  Scenario: Get Chain Registry metrics
    When I send "GET" request to "{{global.chain-registry-metrics}}/metrics"
    Then the response code should be 200

  Scenario: Get Chain Registry readiness
    When I send "GET" request to "{{global.chain-registry-metrics}}/ready"
    Then the response code should be 200

  Scenario: Get Chain Registry liveness
    When I send "GET" request to "{{global.chain-registry-metrics}}/live"
    Then the response code should be 200


  Scenario: Get Contract Registry Swagger
    When I send "GET" request to "{{global.contract-registry-http}}/swagger/"
    Then the response code should be 200

  Scenario: Get Contract Registry Swagger spec
    When I send "GET" request to "{{global.contract-registry-http}}/swagger/swagger.json"
    Then the response code should be 200

  Scenario: Get Contract Registry metrics
    When I send "GET" request to "{{global.contract-registry-metrics}}/metrics"
    Then the response code should be 200

  Scenario: Get Contract Registry readiness
    When I send "GET" request to "{{global.contract-registry-metrics}}/ready"
    Then the response code should be 200

  Scenario: Get Contract Registry liveness
    When I send "GET" request to "{{global.contract-registry-metrics}}/live"
    Then the response code should be 200


  Scenario: Get Envelope Store Swagger
    When I send "GET" request to "{{global.envelope-store-http}}/swagger/"
    Then the response code should be 200

  Scenario: Get Envelope Store Swagger spec
    When I send "GET" request to "{{global.envelope-store-http}}/swagger/swagger.json"
    Then the response code should be 200

  Scenario: Get Envelope Store metrics
    When I send "GET" request to "{{global.envelope-store-metrics}}/metrics"
    Then the response code should be 200

  Scenario: Get Envelope Store readiness
    When I send "GET" request to "{{global.envelope-store-metrics}}/ready"
    Then the response code should be 200

  Scenario: Get Envelope Store liveness
    When I send "GET" request to "{{global.envelope-store-metrics}}/live"
    Then the response code should be 200

    
  Scenario: Get Transaction Scheduler Swagger
    When I send "GET" request to "{{global.tx-scheduler}}/swagger/"
    Then the response code should be 200

  Scenario: Get Transaction Scheduler Swagger
    When I send "GET" request to "{{global.tx-scheduler}}/swagger/swagger.json"
    Then the response code should be 200

  Scenario: Get Transaction Scheduler metrics
    When I send "GET" request to "{{global.tx-scheduler-metrics}}/metrics"
    Then the response code should be 200

  Scenario: Get Transaction Scheduler readiness
    When I send "GET" request to "{{global.tx-scheduler-metrics}}/ready"
    Then the response code should be 200

  Scenario: Get Transaction Scheduler liveness
    When I send "GET" request to "{{global.tx-scheduler-metrics}}/live"
    Then the response code should be 200
