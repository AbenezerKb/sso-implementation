Feature: Get Permission List
  As an admin,
  I want to get the list of all possible permissions
  So that I can use them to create roles

  Background:
    Given I am logged in with the following credentials
      | email           | password | role                |
      | admin@gmail.com | 12345678 | get_permission_list |

  Scenario Outline: I get all permissions
    When I request to get all permissions with group "<group>"
    Then I should get all permissions in that group
    Examples:
      | group  |
      |        |
      | user   |
      | client |

  Scenario Outline: I fail to get all permissions
    When I request to get all permissions with group "<group>"
    Then my request should fail with message "<message>"
    Examples:
      | group    | message                      |
      | no-group | group no-group doesn't exist |
