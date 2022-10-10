Feature: Get Scope By Name

    As an admin
    I want to fetch a particular scope by name
    So that I can see the details and manage that particular scope

    Background:
        Given The following scopes are registered on the system
            | name   | description        | resource_server_name |
            | openid | your profile info  | sso                  |
            | email  | your default email | sso                  |

        And I am logged in as admin user
            | email           | password      | role      |
            | admin@gmail.com | adminPassword | get_scope |
    @success
    Scenario Outline: Successful get scope by name
        When I request to get a scope by "<filled_name>"
        Then I should get the following scope
            | name   | description   | resource_server_name   |
            | <name> | <description> | <resource_server_name> |

        Examples:
            | filled_name | name   | description       | resource_server_name |
            | openid      | openid | your profile info | sso                  |

    @failure
    Scenario Outline: Failed get scope by name
        When I request to get a scope by "<filled_name>"
        Then my request should fail with "<message>"

        Examples:
            | filled_name | message         |
            | scope       | scope not found |
