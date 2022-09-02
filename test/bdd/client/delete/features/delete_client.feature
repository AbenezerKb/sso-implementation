Feature: Delete Client

    As a admin
    I want to delete client
    So that I can remove unwanted clients

    Background: I am logged in as admin
        Given I am logged in as admin user
            | email           | password | role          |
            | admin@gmail.com | iAmAdmin | create_client |
        And There is a client with the following details
            | name | redirect_uris        | secret    | scopes       | client_type  | logo_url               |
            | ride | http://localhost.com | my_secret | openid email | confidential | http://logo.client.com |
    Scenario: Successful Delete
        When I delete the client
        Then The client should be deleted

    Scenario Outline: Failed Delete
        When I delete the client with id <"id">
        Then The delete should fail with message "<message>"

        Examples:
            | id      | message          |
            | Value 1 | client not found |
