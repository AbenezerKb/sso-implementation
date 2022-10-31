Feature: Get User Permissions

  As a user
  I want to see my permissions
  So that I can recognize my allowed permissions

  Background:
    Given I am logged in with the following credentials
      | email          | password | role                                                                    |
      | user@gmail.com | 12345678 | create_user,get_all_roles,update_client,change_role_status,delete_scope |

  @success
  Scenario: I successfully get my permissions
    When I request to get my permissions
    Then I should get all my permissions
