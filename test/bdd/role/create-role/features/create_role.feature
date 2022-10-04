Feature: Create Role
  As an Admin,
  I want to create a role
  So that I can assign it to a user

  Background:
    Given I am logged in with the following credentials
      | email           | password | role        |
      | admin@gmail.com | 12345678 | create_role |

  Scenario: I successfully create a role
    When I request to create a role with the following permissions
      | role_name | permissions                                     |
      | my_role   | create_user,create_client,update_user,get_scope |
    Then the role should successfully be created

  Scenario Outline:
    When I request to create a role with the following permissions
      | role_name   | permissions   |
      | <role_name> | <permissions> |
    Then my request should fail with "<message>" and "<field_error>"
    Examples:
      | role_name | permissions               | message                          | field_error             |
      |           | create_user,create_client |                                  | name is required        |
      | my_role   |                           |                                  | permissions is required |
      | my_role   | unknown                   | permission unknown doesn't exist |                         |