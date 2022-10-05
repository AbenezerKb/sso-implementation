Feature: Delete Scope

    As a admin
    I want to delete scope
    So that I can remove scope's that resource server's no more want

    Background: I am logged in as admin
        Given I am logged in as admin user
            | email           | password | role         |
            | admin@gmail.com | iAmAdmin | delete_scope |
        And There are scope's with the following details
            | name    | description          | resource_server_name |
            | openid  | your profile info    | sso                  |
            | email   | your default email   | sso                  |
            | profile | your default profile | sso                  |
    @success
    Scenario Outline: Successful Delete
        When I delete the scope with name "<name>"
        Then The scope should be deleted

        Examples:
            | name    |
            | email   |
            | profile |


    @failure
    Scenario Outline: no scope found
        When I delete the scope with name "<name>"
        Then The delete should fail with message "<message>"

        Examples:
            | name       | message        |
            | not_openid | no scope found |
