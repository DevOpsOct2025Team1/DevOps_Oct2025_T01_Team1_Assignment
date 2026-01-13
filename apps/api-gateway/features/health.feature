Feature: Health Check
  As a system administrator
  I want to check the health status of the API Gateway
  So that I can monitor the service availability

  Scenario: Health endpoint returns OK status
    When I send a GET request to "/health"
    Then the response status code should be 200
    And the response should contain "status" with value "ok"

  Scenario: Health endpoint responds quickly
    When I send a GET request to "/health"
    Then the response status code should be 200
    And the response time should be less than 1000 milliseconds