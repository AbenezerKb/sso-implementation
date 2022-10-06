Feature: Get Roles
  As an admin
  I want to fetch list of roles
  So that I can manage roles

  Background:
    Given The following roles are registered on the system
      | name  | permissions                       |
      | role1 | create_user,create_client         |
      | role2 | delete_user,create_client         |
      | role3 | create_user,delete_client         |
      | role4 | get_all_permissions,create_client |
      | role5 | create_role,create_user           |
    And I am logged in as admin user
      | email           | password      | role          |
      | admin@gmail.com | adminPassword | get_all_roles |

  @success
  Scenario Outline: I get all the roles
    When I request to get all the roles with the following preferences
      | page   | per_page   |
      | <page> | <per_page> |
    Then I should get the list of roles that pass my preferences
    Examples:
      | page | per_page |
      | 0    | 10       |
      | 0    | 3        |
      | 1    | 2        |
      | 1    | 5        |

  Scenario Outline: I fail to get all the roles due to invalid request
    When I request to get all the roles with the following preferences
      | page   | per_page   |
      | <page> | <per_page> |
    Then I should get error message "<message>"
    Examples:
      | page | per_page | message |
#      | hello | 10       | invalid filter params |
#      | 1     | hello    | invalid filter params |
