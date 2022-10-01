Feature: Get All Scopes
    As an Admin,
    I want to get all registered scopes
    So that I can view details of each scope

    Background:
        Given The following scopes are registered on the system
            | name   | description        | resource_server_name |
            | openid | your profile info  | sso                  |
            | email  | your default email | sso                  |
        And I am logged in as admin user
            | email           | password      | role       |
            | admin@gmail.com | adminPassword | super-user |

    @success
    Scenario Outline: I get all the scopes
        When I request to get all the scopes with the following preferences
            | page   | per_page   |
            | <page> | <per_page> |
        Then I should get the list of scopes that pass my preferences
        Examples:
            | page | per_page |
            | 0    | 10       |
            | 0    | 3        |
            | 1    | 2        |
            | 1    | 5        |