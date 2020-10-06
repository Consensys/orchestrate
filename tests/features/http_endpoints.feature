@http-endpoints
Feature: Verify HTTP Endpoints

  Scenario: Get Chain Registry Swagger
    When I send "GET" request to "{{global.chain-registry}}/swagger/"
    Then the response code should be 200

  Scenario: Get Chain Registry Swagger JSON file
    When I send "GET" request to "{{global.chain-registry}}/swagger/swagger.json"
    Then the response code should be 200

  Scenario: Get Chain Registry metrics
    When I send "GET" request to "{{global.chain-registry-metrics}}/metrics"
    Then the response code should be 200

  Scenario: Get Chain Registry readiness
    When I send "GET" request to "{{global.chain-registry-metrics}}/ready?full=1"
    Then the response code should be 200
    And Response should have the following fields
      | database |
      | OK       |

  Scenario: Get Chain Registry liveness
    When I send "GET" request to "{{global.chain-registry-metrics}}/live"
    Then the response code should be 200

  Scenario: Get Contract Registry Swagger
    When I send "GET" request to "{{global.contract-registry-http}}/swagger/"
    Then the response code should be 200

  Scenario: Get Contract Registry Swagger JSON file
    When I send "GET" request to "{{global.contract-registry-http}}/swagger/swagger.json"
    Then the response code should be 200

  Scenario: Get Contract Registry metrics
    When I send "GET" request to "{{global.contract-registry-metrics}}/metrics"
    Then the response code should be 200

  Scenario: Get Contract Registry readiness
    When I send "GET" request to "{{global.contract-registry-metrics}}/ready?full=1"
    Then the response code should be 200
    And Response should have the following fields
      | database |
      | OK       |

  Scenario: Get Contract Registry liveness
    When I send "GET" request to "{{global.contract-registry-metrics}}/live"
    Then the response code should be 200


  Scenario: Get Transaction Scheduler Swagger
    When I send "GET" request to "{{global.tx-scheduler}}/swagger/"
    Then the response code should be 200

  Scenario: Get Transaction Scheduler Swagger JSON file
    When I send "GET" request to "{{global.tx-scheduler}}/swagger/swagger.json"
    Then the response code should be 200

  Scenario: Get Transaction Scheduler metrics
    When I send "GET" request to "{{global.tx-scheduler-metrics}}/metrics"
    Then the response code should be 200

  Scenario: Get Transaction Scheduler readiness
    When I send "GET" request to "{{global.tx-scheduler-metrics}}/ready?full=1"
    Then the response code should be 200
    And Response should have the following fields
      | chain-registry | contract-registry | database | kafka |
      | OK             | OK                | OK       | OK    |

  Scenario: Get Transaction Scheduler liveness
    When I send "GET" request to "{{global.tx-scheduler-metrics}}/live"
    Then the response code should be 200


  Scenario: Get tx-crafter readiness
    When I send "GET" request to "{{global.tx-crafter-metrics}}/ready?full=1"
    Then the response code should be 200
    And Response should have the following fields
      | chain-registry | redis | transaction-scheduler | kafka |
      | OK             | OK    | OK                    | OK    |

  Scenario: Get tx-crafter liveness
    When I send "GET" request to "{{global.tx-crafter-metrics}}/live"


  Scenario: Get tx-signer readiness
    When I send "GET" request to "{{global.tx-signer-metrics}}/ready?full=1"
    Then the response code should be 200
    And Response should have the following fields
      | secret-store | kafka |
      | OK           | OK    |

  Scenario: Get tx-signer liveness
    When I send "GET" request to "{{global.tx-signer-metrics}}/live"
    Then the response code should be 200


  Scenario: Get tx-sender readiness
    When I send "GET" request to "{{global.tx-sender-metrics}}/ready?full=1"
    Then the response code should be 200
    And Response should have the following fields
      | transaction-scheduler | kafka | redis |
      | OK                    | OK    | OK    |

  Scenario: Get tx-sender liveness
    When I send "GET" request to "{{global.tx-sender-metrics}}/live"
    Then the response code should be 200


  Scenario: Get tx-listener readiness
    When I send "GET" request to "{{global.tx-listener-metrics}}/ready?full=1"
    Then the response code should be 200
    And Response should have the following fields
      | chain-registry | transaction-scheduler | kafka |
      | OK             | OK                    | OK    |

  Scenario: Get tx-listener liveness
    When I send "GET" request to "{{global.tx-listener-metrics}}/live"
    Then the response code should be 200
