Feature: Admin actions

  Scenario: Admin can delete user
    Given I am authenticated as "admin"
    When I send a DELETE request to "/api/admin/delete_user" with json:
      """
      {"id":"u123"}
      """
    Then the response status code should be 200

  Scenario: Admin delete user with missing id
    Given I am authenticated as "admin"
    When I send a DELETE request to "/api/admin/delete_user" with json:
      """
      {}
      """
    Then the response status code should be 400


  Scenario: Regular user cannot delete user
    Given I am authenticated as "user"
    When I send a DELETE request to "/api/admin/delete_user" with json:
      """
      {"id":"u123"}
      """
    Then the response status code should be 401

