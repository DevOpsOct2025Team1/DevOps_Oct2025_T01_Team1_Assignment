Feature: Security enforcement

  Scenario: Unauthenticated cannot access admin delete
    When I send a DELETE request to "/api/admin/delete_user" with json:
      """
      {"id":"u123"}
      """
    Then the response status code should be 401

Scenario: Invalid authorization header
  When I send a DELETE request to "/api/admin/delete_user" with json:
    """
    {"id":"u123"}
    """
  Then the response status code should be 401

Scenario: Missing JWT when accessing admin delete
  When I send a DELETE request to "/api/admin/delete_user" with json:
    """
    {"id":"u123"}
    """
  Then the response status code should be 401

Scenario: Login with invalid JSON body
  When I send a POST request to "/api/login" with json:
    """
    {invalid json}
    """
  Then the response status code should be 400
