Feature: Get Client By ID

    As an Admin,
    I want to get a specific client by id
    So that I can see the details of that particular  client

    Background: I am logged in as admin
        Given I am logged in as admin user
            | email           | password | role       |
            | admin@gmail.com | iAmAdmin | get_client |

        And there is client with the following details:
            | name      | client_type  | redirect_uris      | scopes        | logo_url                               | status |
            | clientOne | confidential | https://google.com | profile email | https://ww.google.com/error-image1.png | active |

    @success
    Scenario: Successful Get client
        Given I have client id
        When I Get the client
        Then I should successfully get the client

    @failure
    Scenario Outline: invalid id
        Given I have client with id "<id>"
        When I Get the client
        Then Then I should get error with message "<message>"

        Examples:
            | id                                   | message           |
            | 06f7404f-5402-4832-8b0b-53da2cdd7efc | no client found   |
            | 3kjf0-kljkla0-afl30-afl-dk           | invalid client id |
