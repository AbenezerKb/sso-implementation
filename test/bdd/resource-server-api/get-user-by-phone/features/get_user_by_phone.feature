Feature: Get User By Phone

  Scenario: I get the user by phone number
    Given I have authenticated my self as a resource server
    And There is a user with phone number "251912121212"
    When I ask for a user with phone number "0912121212"
    Then I should get the user data

  Scenario: I fail to get the user by phone number
    Given I have authenticated my self as a resource server
    And There is a user with phone number "251912121212"
    When I ask for a user with phone number "0913131313"
    Then My request should fail with message "no user found"
