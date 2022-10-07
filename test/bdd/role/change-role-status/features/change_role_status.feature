Feature: Change Role Status
  As an admin
  I want to change a status of a role
  So that I can disable/enable the role

  Background: I am logged in as admin
    Given I am logged in as admin user
      | email           | password | role               |
      | admin@gmail.com | iAmAdmin | change_role_status |

  @success
  Scenario Outline: Successful Role Status Update
    Given there is a role with the following details:
      | name   | permissions   | status   |
      | <name> | <permissions> | <status> |
    When I update the role's status to "<updated_status>"
    Then the role status should update to "<updated_status>"
    Examples:
      | name  | permissions               | status   | updated_status |
      | role1 | create_user,create_client | ACTIVE   | INACTIVE       |
      | role2 | get_all_users,create_role | INACTIVE | ACTIVE         |

  @failure
  Scenario Outline: role not found
    Given there is role with name "<name>"
    When I update the role's status to "<updated_status>"
    Then Then I should get role not found error with message "<message>"

    Examples:
      | name    | updated_status | message        |
      | no-role | ACTIVE         | role not found |

  @failure
  Scenario Outline: Invalid Status
    Given there is a role with the following details:
      | name   | permissions   | status   |
      | <name> | <permissions> | <status> |
    When I update the role's status to "<updated_status>"
    Then Then I should get error with message "<message>"

    Examples:
      | name  | permissions               | status   | updated_status | message            |
      | role1 | create_user,create_client | ACTIVE   | INACTIVE       | invalid status     |
      | role2 | get_all_users,create_role | INACTIVE | ACTIVE         | status is required |

