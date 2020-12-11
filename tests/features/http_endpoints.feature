@http-endpoints
Feature: Verify HTTP Endpoints
  ###################
  # CHAIN REGISTRY
  ###################
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

  ###################
  # CONTRACT REGISTRY
  ###################
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


  ###################
  # TRANSACTION SCHEDULER
  ###################
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


  ###################
  # TX-CRAFTER
  ###################
  Scenario: Get tx-crafter readiness
    When I send "GET" request to "{{global.tx-crafter-metrics}}/ready?full=1"
    Then the response code should be 200
    And Response should have the following fields
      | chain-registry | redis | transaction-scheduler | kafka |
      | OK             | OK    | OK                    | OK    |

  Scenario: Get tx-crafter liveness
    When I send "GET" request to "{{global.tx-crafter-metrics}}/live"


  ###################
  # TX-SIGNER
  ###################
  Scenario: Get tx-signer readiness
    When I send "GET" request to "{{global.tx-signer-metrics}}/ready?full=1"
    Then the response code should be 200
    And Response should have the following fields
      | kafka | transaction-scheduler
      | OK    | OK

  Scenario: Get tx-signer liveness
    When I send "GET" request to "{{global.tx-signer-metrics}}/live"
    Then the response code should be 200


  ###################
  # TX-SENDER
  ###################
  Scenario: Get tx-sender readiness
    When I send "GET" request to "{{global.tx-sender-metrics}}/ready?full=1"
    Then the response code should be 200
    And Response should have the following fields
      | transaction-scheduler | kafka | redis |
      | OK                    | OK    | OK    |

  Scenario: Get tx-sender liveness
    When I send "GET" request to "{{global.tx-sender-metrics}}/live"
    Then the response code should be 200


  ###################
  # TX-LISTENER
  ###################
  Scenario: Get tx-listener readiness
    When I send "GET" request to "{{global.tx-listener-metrics}}/ready?full=1"
    Then the response code should be 200
    And Response should have the following fields
      | chain-registry | transaction-scheduler | kafka |
      | OK             | OK                    | OK    |

  Scenario: Get tx-listener liveness
    When I send "GET" request to "{{global.tx-listener-metrics}}/live"
    Then the response code should be 200
    
  ###################
  # Key Manager
  ###################
  Scenario: Get key-manager liveness
    When I send "GET" request to "{{global.key-manager-metrics}}/live"
    Then the response code should be 200

  Scenario: Get key-manager readiness
    When I send "GET" request to "{{global.key-manager-metrics}}/ready?full=1"
    Then the response code should be 200
    And Response should have the following fields
      | vault |
      | OK    |

  ###################
  # Identity Manager
  ###################
  Scenario: Get Identity Manager Swagger JSON file
    When I send "GET" request to "{{global.identity-manager}}/swagger/swagger.json"
    Then the response code should be 200

  Scenario: Get Identity Manager metrics
    When I send "GET" request to "{{global.identity-manager-metrics}}/metrics"
    Then the response code should be 200

  Scenario: Get Identity Manager readiness
    When I send "GET" request to "{{global.identity-manager-metrics}}/ready?full=1"
    Then the response code should be 200
    And Response should have the following fields
      | chain-registry | database |
      | OK             | OK       |

  Scenario: Get Identity Manager liveness
    When I send "GET" request to "{{global.identity-manager-metrics}}/live"
    Then the response code should be 200
