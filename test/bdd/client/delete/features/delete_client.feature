Feature: Delete Client

    As a admin
    I want to delete client
    So that I can remove unwanted clients

    Background: I am logged in as admin
        Given I am logged in as admin user
            | email           | password | role          |
            | admin@gmail.com | iAmAdmin | delete_client |
        And There is a client with the following details
            | name | redirect_uris        | secret    | scopes       | client_type  | logo_url               |
            | ride | http://localhost.com | my_secret | openid email | confidential | http://logo.client.com |
    @success
    Scenario: Successful Delete
        When I delete the client
        Then The client should be deleted

    @failure
    Scenario Outline: no client
        When I delete the client with id "<id>"
        Then The delete should fail with message "<message>"

        Examples:
            | id                                   | message          |
            | 60d56419-c2e9-4ee4-951f-04644d245ee3 | client not found |

    @failure
    Scenario Outline: invalid client id
        When I delete the client with id "<id>"
        Then The delete should fail with error message "<message>"

        Examples:
            | id                                 | message          |
            | 60d56419-c2e9-4ee4-951f-04644d245e | client not found |
            | 4                                  | client not found |

