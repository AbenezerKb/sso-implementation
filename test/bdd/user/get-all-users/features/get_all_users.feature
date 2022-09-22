Feature: Get All Users
  As an admin,
  I want to get all the list of users
  So that I can manage details of each user

  Background:
    Given The following users are registered on the system
      | first_name | middle_name | last_name | phone        | email           | password |
      | user1      | user1       | user1     | 251911111111 | user1@gmail.com | 111111   |
      | user2      | user2       | user2     | 251922222222 | user2@gmail.com | 222222   |
      | user3      | user3       | user3     | 251933333333 | user3@gmail.com | 333333   |
      | user4      | user4       | user4     | 251944444444 | user4@gmail.com | 444444   |
      | user5      | user5       | user5     | 251955555555 | user5@gmail.com | 555555   |
    And I am logged in as admin user
      | email           | password      | role       |
      | admin@gmail.com | adminPassword | super-user |

  @success
  Scenario Outline: I get all the users
    When I request to get all the users with the following preferences
      | page   | per_page   |
      | <page> | <per_page> |
    Then I should get the list of users that pass my preferences
    Examples:
      | page | per_page |
      | 0    | 10       |
      | 0    | 3        |
      | 1    | 2        |
      | 1    | 5        |

  Scenario Outline: I fail to get all the users due to invalid request
    When I request to get all the users with the following preferences
      | page   | per_page   |
      | <page> | <per_page> |
    Then I should get error message "<message>"
    Examples:
      | page | per_page | message |
#      | hello | 10       | invalid filter params |
#      | 1     | hello    | invalid filter params |
