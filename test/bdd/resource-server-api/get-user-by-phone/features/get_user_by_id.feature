Feature: Get User By ID

  Scenario: I get the user by id
    Given I have authenticated my self as a resource server
    And There is a user with phone number "251912121212"
    When I ask for a user with id
    Then I should get the user data

  Scenario: I fail to get the user by id
    Given I have authenticated my self as a resource server
    And There is a user with phone number "251912121212"
    When I ask for a user with incorrect id
    Then My request should fail with message "no user found"
