Feature: List Users

  Scenario: Admin can list all users
    Given I am authenticated as "admin"
    When I send a GET request to "/api/admin/list_users"
    Then the response status code should be 200

  Scenario: Admin can filter users by role
    Given I am authenticated as "admin"
    When I send a GET request to "/api/admin/list_users?role=admin"
    Then the response status code should be 200

  Scenario: Admin can search users by username
    Given I am authenticated as "admin"
    When I send a GET request to "/api/admin/list_users?username=john"
    Then the response status code should be 200

  Scenario: Admin can combine role filter and username search
    Given I am authenticated as "admin"
    When I send a GET request to "/api/admin/list_users?role=user&username=test"
    Then the response status code should be 200

  Scenario: Regular user cannot list users
    Given I am authenticated as "user"
    When I send a GET request to "/api/admin/list_users"
    Then the response status code should be 401

  Scenario: Unauthenticated user cannot list users
    When I send a GET request to "/api/admin/list_users"
    Then the response status code should be 401

  Scenario: Invalid role filter returns 400
    Given I am authenticated as "admin"
    When I send a GET request to "/api/admin/list_users?role=invalid"
    Then the response status code should be 400