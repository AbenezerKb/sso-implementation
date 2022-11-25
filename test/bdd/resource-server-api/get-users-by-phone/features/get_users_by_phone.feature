Feature: Get Users By phone

  Scenario: I get users by phone
    Given I have authenticated my self as a resource server
    And There are users with phone numbers
      | phone          |
      | "251912121212" |
      | "251913131313" |
    When I ask for users with phones
      | phones                  |
      | 251912121212,0913131313 |
    Then I should get the users