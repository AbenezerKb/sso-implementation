Feature: update Client Status

    As a admin
    I want to update client's status
    So that I can activate or deactivate client's as appropriate

    Background: I am logged in as admin
        Given I am logged in as admin user
            | email           | password | role          |
            | admin@gmail.com | iAmAdmin | update_client |
    @success
    Scenario Outline: Successful Client Status Update
        Given there is client with the following details:
            | name   | client_type   | redirect_uris   | scopes   | logo_url   | status   |
            | <name> | <client_type> | <redirect_uris> | <scopes> | <logo_url> | <status> |
        When I update the client's status to "<updated_status>"
        Then the client status should update to "<updated_status>"

        Examples:
            | name      | client_type  | redirect_uris       | scopes          | logo_url                               | status   | updated_status |
            | clientOne | confidential | https://google.com  | profile email   | https://ww.google.com/error-image1.png | ACTIVE   | INACTIVE       |
            | clientTwo | public       | https://youtube.com | profile balance | https://ww.google.com/error-image2.png | INACTIVE | ACTIVE         |

    @failure
    Scenario Outline: Client not found
        Given there is client with id "<id>"
        When I update the client's status to "<updated_status>"
        Then Then I should get client not found error with message "<message>"

        Examples:
            | id                          | updated_status | message           |
            | 3kjf0-kjf0afl2-afl30-afl-dk | ACTIVE         | client not found  |
            | 3kjf0-kjf0afl2-afl30-afl-dk | ACTIVE         | invalid client id |
    @failure
    Scenario Outline: Invalid Status
        Given there is client with the following details:
            | name   | client_type   | redirect_uris   | scopes   | logo_url   | status   |
            | <name> | <client_type> | <redirect_uris> | <scopes> | <logo_url> | <status> |
        When I update the client's status to "<updated_status>"
        Then Then I should get error with message "<message>"

        Examples:
            | name        | client_type  | redirect_uris        | scopes         | logo_url                               | status | updated_status | message               |
            | clientThree | public       | https://facebook.com | profile openid | https://ww.google.com/error-image3.png | ACTIVE | INACTIVED      | must be a valid value |
            | clientFour  | confidential | https://google.com   | profile        | https://ww.google.com/error-image4.png | ACTIVE |                | status is required    |

