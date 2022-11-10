Feature: Update Client

    As an Admin
    I want to update client
    So that that client will have up to dated information

    Background:
        Given I am logged in as admin user
            | email           | password      | role       |
            | admin@gmail.com | adminPassword | super-user |

    @success
    Scenario: Successful Client Update
        Given There is client with the following details
            | name      | client_type  | redirect_uris      | scopes        | logo_url                                       |
            | clientOne | confidential | https://google.com | profile email | https://www.google.com/images/errors/robot.png |
        And I fill the form with the following details
            | name      | client_type  | redirect_uris      | scopes        | logo_url                                       |
            | clientOne | confidential | https://google.com | profile email | https://www.google.com/images/errors/robot.png |

        When I update the client
        Then The client should be updated

    @failure
    Scenario Outline: Failed Client Update
        Given There is client with the following details
            | name      | client_type  | redirect_uris      | scopes        | logo_url                               |
            | clientOne | confidential | https://google.com | profile email | https://ww.google.com/error-image1.png |
        And I fill the form with the following details
            | name   | client_type   | redirect_uris   | scopes   | logo_url   |
            | <name> | <client_type> | <redirect_uris> | <scopes> | <logo_url> |
        When I update the client
        Then The client updated should fail with message "<message>"

        Examples:
            | name      | client_type  | redirect_uris           | scopes        | logo_url                                       | message                                           |
            |           | confidential | https://google.com      | profile email | https://www.google.com/images/errors/robot.png | name is required                                  |
            | clientOne | confidential |                         | profile email | https://www.google.com/images/errors/robot.png | redirect_uris is required                         |
            | clientOne | confidential | https://google.com      | profile email |                                                | logo_url is required                              |
            | clientOne |              | https://google.com      | profile email | https://www.google.com/images/errors/robot.png | client_type is required                           |
            | newClient | my_type      | https://google.com      | profile email | https://www.google.com/images/errors/robot.png | client type must be either confidential or public |
            | newClient | confidential | https://google.com      | profile email | my-logo-url                                    | invalid logo_url                                  |
            | newClient | confidential | https://google.com      | profile email | http://hello-there.com/logo.png                | logo not found                                    |
