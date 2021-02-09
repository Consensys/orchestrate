#@http-endpoints
#Feature: Verify HTTP Endpoints
#  ###################
#  # API
#  ###################
#  Scenario: Get API Swagger JSON file
#    When I send "GET" request to "{{global.api}}/swagger/swagger.json"
#    Then the response code should be 200
#
#  Scenario: Get API metrics
#    When I send "GET" request to "{{global.api-metrics}}/metrics"
#    Then the response code should be 200
#
#  Scenario: Get API readiness
#    When I send "GET" request to "{{global.api-metrics}}/ready?full=1"
#    Then the response code should be 200
#    And Response should have the following fields
#      | database | kafka |
#      | OK       | OK    |
#
#  Scenario: Get API liveness
#    When I send "GET" request to "{{global.api-metrics}}/live"
#    Then the response code should be 200
#
#  ###################
#  # TX-SENDER
#  ###################
#  Scenario: Get tx-sender readiness
#    When I send "GET" request to "{{global.tx-sender-metrics}}/ready?full=1"
#    Then the response code should be 200
#    And Response should have the following fields
#      | api | kafka |
#      | OK  | OK    |
#
#  Scenario: Get tx-sender liveness
#    When I send "GET" request to "{{global.tx-sender-metrics}}/live"
#    Then the response code should be 200
#
#  ###################
#  # TX-LISTENER
#  ###################
#  Scenario: Get tx-listener readiness
#    When I send "GET" request to "{{global.tx-listener-metrics}}/ready?full=1"
#    Then the response code should be 200
#    And Response should have the following fields
#      | api | kafka |
#      | OK  | OK    |
#
#  Scenario: Get tx-listener liveness
#    When I send "GET" request to "{{global.tx-listener-metrics}}/live"
#    Then the response code should be 200
#
#  ###################
#  # Key Manager
#  ###################
#  Scenario: Get key-manager liveness
#    When I send "GET" request to "{{global.key-manager-metrics}}/live"
#    Then the response code should be 200
#
#  Scenario: Get key-manager readiness
#    When I send "GET" request to "{{global.key-manager-metrics}}/ready?full=1"
#    Then the response code should be 200
#    And Response should have the following fields
#      | vault |
#      | OK    |
