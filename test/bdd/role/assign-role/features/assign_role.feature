Feature: Assign Role
  As an admin,
  I want to assign a role for a specific user
  So that I can give some permissions for the user

  Background:
    Given I am logged in with the following credentials
      | email           | password | role             |
      | admin@gmail.com | 12345678 | update_user_role |
    And The following role is registered on the system
      | name  | permissions                               |
      | clerk | get_all_users,create_user,get_all_clients |
    And The following user is registered on the system
      | first_name | middle_name | last_name | phone        | email            | password |
      | abebe      | alemu       | rebuma    | 251923456789 | normal@gmail.com | 123456   |

  @success
  Scenario: I successfully assign the role to the user
    When I request to assign "clerk" as role for the user
    Then the role should be assigned to the user

  @failure
  Scenario Outline: I fail to assign the role to the user
    When I request to assign "<role>" as role for the user
    Then my request should fail with "<message>" and "<field_error>"

    Examples:
      | role         | message                          | field_error      |
      |              |                                  | role is required |
      | non-existing | role non-existing does not exist |                  |