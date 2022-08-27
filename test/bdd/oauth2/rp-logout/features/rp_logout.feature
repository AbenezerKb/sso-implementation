Feature: logout

    As a clinet
    I want to logout user from SSO as they request to logout from my service
    so that the end user can have a clear session on a particular device

    Background: I have given id_token
        Given I am  registered on the system
        And the user is registered on the system
        And I have id_token
    @success
    Scenario Outline: Successful Logout
        Given I have the following details:
            | post_logout_redirect_uri   | state   |
            | <post_logout_redirect_uri> | <state> |
        When I request to logout
        Then I should be redirected to "<logout_uri>" with the following query params:
            | post_logout_redirect_uri   | state   |
            | <post_logout_redirect_uri> | <state> |

        Examples:
            | post_logout_redirect_uri | state | logout_uri              |
            | https://www.google.com/  | 1234  | https://www.google.com/ |
    @failure
    Scenario Outline: Failed Logout
        Given I have the following invalid_request details:
            | id_token_hint   | post_logout_redirect_uri   | state   |
            | <id_token_hint> | <post_logout_redirect_uri> | <state> |
        When I request to logout
        Then I should be redirected to "<err_uri>" with the following failure query params:
            | error   | error_description   |
            | <error> | <error_description> |
        Examples:
            | id_token_hint                         | post_logout_redirect_uri | state | err_uri                 | error           | error_description       |
            | iielksaklcnvlajfkje.kjkladfkje.kjklad | https://www.google.com/  | 1234  | https://www.google.com/ | invalid request | no logedin user found   |
            |                                       | https://www.google.com/  | 1234  | https://www.google.com/ | invalid request | login hint is required. |

    @failure
    Scenario Outline: Successful Logout
        Given I have the following details:
            | post_logout_redirect_uri   | state   |
            | <post_logout_redirect_uri> | <state> |
        When I request to logout
        Then I should be redirected to "<err_uri>" with the following failure query params:
            | error   | error_description   |
            | <error> | <error_description> |
        Examples:
            | post_logout_redirect_uri | state | err_uri                 | error           | error_description                     |
            |                          | 1234  | https://www.google.com/ | invalid request | post logout redirect uri is required. |