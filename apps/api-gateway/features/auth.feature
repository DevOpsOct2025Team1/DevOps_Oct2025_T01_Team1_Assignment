Feature: Authentication

  Scenario: Login success
    When I send a POST request to "/api/login" with json:
      """
      {"username":"testuser","password":"password123"}
      """
    Then the response status code should be 200

  Scenario: Login missing username
    When I send a POST request to "/api/login" with json:
      """
      {"password":"password123"}
      """
    Then the response status code should be 400

Scenario: Login with wrong password
  When I send a POST request to "/api/login" with json:
    """
    {"username":"testuser","password":"wrongpass"}
    """
  Then the response status code should be 500

Scenario: Login with empty body
  When I send a POST request to "/api/login" with json:
    """
    {}
    """
  Then the response status code should be 400
