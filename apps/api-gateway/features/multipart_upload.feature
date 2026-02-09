Feature: Multipart Upload

  Scenario: Initiate multipart upload
    Given I am authenticated as "user"
    When I send a POST request to "/api/files/multipart/initiate" with json:
      """
      {"filename":"large-file.mp4","content_type":"video/mp4","total_size":104857600}
      """
    Then the response status code should be 200
    And the response should contain "upload_id" with value "test-upload-id"

  Scenario: Initiate multipart upload missing filename
    Given I am authenticated as "user"
    When I send a POST request to "/api/files/multipart/initiate" with json:
      """
      {"total_size":104857600}
      """
    Then the response status code should be 400

  Scenario: Initiate multipart upload missing total_size
    Given I am authenticated as "user"
    When I send a POST request to "/api/files/multipart/initiate" with json:
      """
      {"filename":"large-file.mp4"}
      """
    Then the response status code should be 400

  Scenario: Initiate multipart upload unauthenticated
    When I send a POST request to "/api/files/multipart/initiate" with json:
      """
      {"filename":"large-file.mp4","content_type":"video/mp4","total_size":104857600}
      """
    Then the response status code should be 401

  Scenario: Upload part
    Given I am authenticated as "user"
    When I send a multipart form POST to "/api/files/multipart/test-upload-id/part/1" with file "chunk" containing "chunkdata"
    Then the response status code should be 200

  Scenario: Upload part unauthenticated
    When I send a multipart form POST to "/api/files/multipart/test-upload-id/part/1" with file "chunk" containing "chunkdata"
    Then the response status code should be 401

  Scenario: Complete multipart upload
    Given I am authenticated as "user"
    When I send a POST request to "/api/files/multipart/test-upload-id/complete" with json:
      """
      {"parts":[{"part_number":1,"etag":"\"etag1\""},{"part_number":2,"etag":"\"etag2\""}]}
      """
    Then the response status code should be 200

  Scenario: Complete multipart upload missing parts
    Given I am authenticated as "user"
    When I send a POST request to "/api/files/multipart/test-upload-id/complete" with json:
      """
      {}
      """
    Then the response status code should be 400

  Scenario: Complete multipart upload unauthenticated
    When I send a POST request to "/api/files/multipart/test-upload-id/complete" with json:
      """
      {"parts":[{"part_number":1,"etag":"\"etag1\""}]}
      """
    Then the response status code should be 401

  Scenario: Abort multipart upload
    Given I am authenticated as "user"
    When I send a DELETE request to "/api/files/multipart/test-upload-id" with json:
      """
      """
    Then the response status code should be 200

  Scenario: Abort multipart upload unauthenticated
    When I send a DELETE request to "/api/files/multipart/test-upload-id" with json:
      """
      """
    Then the response status code should be 401
