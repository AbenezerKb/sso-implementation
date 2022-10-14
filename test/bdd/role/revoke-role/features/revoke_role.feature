Feature: Revoke Role
  As an admin
  I want to revoke role from a user
  So that that user will not have that role

  Background:
    Given I am logged in with the following credentials
      | email           | password | role             |
      | admin@gmail.com | 12345678 | revoke_user_role |
    And there is a role with the following details:
      | name  | permissions               |
      | role1 | create_user,create_client |
    And the following user has the role assigned
      | first_name | middle_name | last_name | phone        | email            | password |
      | abebe      | alemu       | rebuma    | 251923456789 | normal@gmail.com | 123456   |

  @success
  Scenario: I successfully revoke the role
    When I request to revoke the role for "abebe"
    Then the user should no longer have that role assigned

  @failure
  Scenario Outline: I fail to delete the role
    When I request to revoke the role for "<user>"
    Then my request should fail with "<message>"
    Examples:
      | user    | message        |
      | no-user | user not found |
