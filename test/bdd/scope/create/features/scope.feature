Feature: Create Scope

    AS a admin user
    I want to create a scope
    So that I can create a scope for a resource servers

    Background: I am logged in as admin
        Given I am logged in as admin user
            | email           | password | role         |
            | admin@gmail.com | iAmAdmin | create_scope |
    @success
    Scenario : Successful  Scope Creation
        Given I fill the form with following fields:
            | name    | description  | resource_server |
            | profile | test profile | test_server     |
        When I create the scope
        Then I should have new scope

    @failure
    Scenario Outline: Unsuccessful Scope Creation
        Given I fill the form with following fields:
            | name   | description   | resource_server   |
            | <name> | <description> | <resource_server> |
        When I create the scope
        Then The creation should fail with "<message>"
        Examples:
            | name    | description | resource_server | message                     |
            |         | openid      | test_server     | name is required            |
            | openid  |             | test_server     | description is required     |
            | open id | openid      | test_server     | name can not contain spaces |


