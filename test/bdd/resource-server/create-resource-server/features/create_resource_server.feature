Feature: Create Resource Server
  As an admin,
  I want to be able to create resource servers,
  So that I can manage access to resource data with the SSO

  Background:
    Given I am logged in with the following credentials
      | email           | password | role                   |
      | admin@gmail.com | 12345678 | create_resource_server |

  @success
  Scenario: I create a resource server successfully
    Given I have filled resource server name "resource_server" and the following scopes
      | name    | description     |
      | scope_1 | this is scope 1 |
      | scope_2 | this is scope 2 |
    When I submit to create a resource server
    Then the resource server should be created

  @failure
  Scenario Outline: I fail to create a resource server
    Given the resource server "existing_server" is registered
    And I have filled resource server name "<server_name>" and the following scopes
      | name           | description     |
      | <scope_name_1> | <description_1> |
      | <scope_name_2> | <description_2> |
    When I submit to create a resource server
    Then the request should fail with "<message>" and "<field_error>"
    Examples:
      | server_name     | scope_name_1 | description_1   | scope_name_2 | description_2   | field_error                   | message                   |
      |                 | scope_1      | this is scope 1 | scope_2      | this is scope 2 | server name is required       |                           |
      | resource_server |              | this is scope 1 | scope_2      | this is scope 2 | scope name is required        |                           |
      | resource_server | scope_1      |                 | scope_2      | this is scope 2 | scope description is required |                           |
      | resource_server | scope_1      | this is scope 1 | scope_1      | this is scope 2 | scope name must be unique     |                           |
      | existing_server | scope_1      | this is scope 1 | scope_2      | this is scope 2 |                               | this server name is taken |