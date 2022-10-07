Feature: Update Role
  As an admin
  I want to update a role
  So that I can change the permissions given to that role

  Background:
    Given I am logged in with the following credentials
      | email           | password | role        |
      | admin@gmail.com | 12345678 | update_role |
    And there is a role with the following details
      | name  | permissions               |
      | role1 | create_user,create_client |

  @success
  Scenario: I successfully update role
    When I request to update "role1" with the following permissions
      | permissions                            |
      | create_user, update_scope, create_role |
    Then the role should be updated

  @failure
  Scenario Outline: I fail to update the role
    When I request to update "<role>" with the following permissions
      | permissions   |
      | <permissions> |
    Then my request should fail with "<message>" and "<field_error>"
    Examples:
      | role  | permissions         | message                           | field_error             |
      | role1 |                     |                                   | permissions is required |
      | role1 | no_perm,create_user | permission no_perm does not exist |                         |

  Scenario Outline: I fail to update the role
    When I request to update "<role>" with the following permissions
      | permissions   |
      | <permissions> |
    Then my request should fail with no role found "<message>"
    Examples:
      | role    | permissions | message        |
      | no-role | create_user | role not found |
