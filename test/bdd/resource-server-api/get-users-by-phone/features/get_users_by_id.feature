Feature: Get Users By ID

  Scenario: I get users by id
    Given I have authenticated my self as a resource server
    And There are users with phone numbers
      | phone          |
      | "251912121212" |
      | "251913131313" |
    When I ask for users with ids
    Then I should get the users