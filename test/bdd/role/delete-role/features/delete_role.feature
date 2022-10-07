Feature: Delete Role
  As an admin
  I want to delete specific role
  So that no user will have that role

  Background:
    Given I am logged in with the following credentials
      | email           | password | role        |
      | admin@gmail.com | 12345678 | delete_role |
    And there is a role with the following details:
      | name  | permissions               | status |
      | role1 | create_user,create_client | ACTIVE |
    And the following user has the role assigned
      | first_name | middle_name | last_name | phone        | email            | password |
      | abebe      | alemu       | rebuma    | 251923456789 | normal@gmail.com | 123456   |

  @success
  Scenario: I successfully delete the role
    When I request to delete the role "role1"
    Then the role should be deleted
    And the user should no longer have that role assigned

  @failure
  Scenario Outline: I fail to delete the role
    When I request to delete the role "<role>"
    Then my request should fail with "<message>" and "<field_error>"
    Examples:
      | role    | message        | field_error      |
      |         |                | role is required |
      | no-role | role not found |                  |