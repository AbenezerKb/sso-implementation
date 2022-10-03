Feature: Get All Resource Servers
  As an Admin,
  I want to get all registered resource servers
  So that I can view details of each resource server

  Background:
    Given The following resource servers are registered on the system
      | name              |
      | resource_server_1 |
      | resource_server_2 |
      | resource_server_3 |
      | resource_server_4 |
      | resource_server_5 |
    And the resource servers have the following scopes
      | name    | description     | resource_server_name |
      | scope_1 | this is scope 1 | resource_server_1    |
      | scope_2 | this is scope 2 | resource_server_2    |
      | scope_3 | this is scope 3 | resource_server_2    |
      | scope_4 | this is scope 4 | resource_server_3    |
      | scope_5 | this is scope 5 | resource_server_4    |
    And I am logged in as admin user
      | email           | password      | role       |
      | admin@gmail.com | adminPassword | super-user |

  @success
  Scenario Outline: I get all the resource servers
    When I request to get all the resource servers with the following preferences
      | page   | per_page   |
      | <page> | <per_page> |
    Then I should get the list of resource servers that pass my preferences
    Examples:
      | page | per_page |
      | 0    | 10       |
      | 0    | 3        |
      | 1    | 2        |
      | 1    | 5        |

  Scenario Outline: I fail to get all the resource servers due to invalid request
    When I request to get all the resource servers with the following preferences
      | page   | per_page   |
      | <page> | <per_page> |
    Then I should get error message "<message>"
    Examples:
      | page | per_page | message |
#      | hello | 10       | invalid filter params |
#      | 1     | hello    | invalid filter params |
