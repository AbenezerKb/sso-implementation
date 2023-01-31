Feature: Forget Password

  As  a User,
  I want to get a new generated password when I forget it
  So that I will not be locked out of the system

  Background:
    Given I have a user account with the following details
      | first_name | middle_name | last_name | phone        | email            | password | gender |
      | nati       | nati        | nati      | 251923456789 | normal@gmail.com | 123456   | male   |

  @success
  Scenario Outline: Successful password reset
    Given I fill my phone number as "<phone_number>"
    When I request to have forgotten my password
    Then I should successfully get a change password request code
    And I should successfully change my password using the request code

    Examples:
      | phone_number |
      | 251923456789 |

  @failure
  Scenario Outline: Unsuccessful password reset
    Given I fill my phone number as "<phone_number>"
    When I request to have forgotten my password
    Then I should successfully get a change password request code
    And I should fail change my password using an incorrect request code

    Examples:
      | phone_number |
      | 251923456789 |
