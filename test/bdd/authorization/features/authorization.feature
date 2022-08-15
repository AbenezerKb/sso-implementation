Feature: obtaining authorization

    As a client

    I want to obtain authorization from the resource owner

    So that I can request to get an access token or refresh token.

    Scenario Outline: Succesfull Obtaining Authorization
        Given I am I have the following parameters:
            | response_type   | client_id   | redirect_uri   | scope   | state   |
            | <response_type> | <client_id> | <redirect_uri> | <scope> | <state> |

        When I send a POST request
        Then I should be redirected to "<redirect_uri>" with the following parameters:
            | code   | state   |
            | <code> | <state> |
        Examples:
            | response_type | client_id | redirect_uri     | scope           | state | code | state |
            | code          | 234555    | http://localhost | openid userinfo | 1234  | 1234 | 1234  |

    Scenario Outline: Unable to Obtain Authorization
        Given I am I have the following parameters:
            | response_type   | client_id   | redirect_uri   | scope   | state   |
            | <response_type> | <client_id> | <redirect_uri> | <scope> | <state> |

        When I send a POST request
        Then I should be redirected to "<redirect_uri>" with the following parameters:
            | error   | error_description   | state   |
            | <error> | <error_description> | <state> |
        Examples:
            | response_type      | client_id | redirect_uri     | scope  | state | error                     | error_description         |
            | authorization_code | 234555    | http://localhost | openid | 1234  | access_denied             | access_denied             |
            | authorization_code | 234555    | http://localhost | openid | 1234  | invalid_request           | invalid_request           |
            | authorization_code | 234555    | http://localhost | openid | 1234  | server_error              | server_error              |
            | authorization_code | 234555    | http://localhost | openid | 1234  | unauthorized_client       | unauthorized_client       |
            | authorization_code | 234555    | http://localhost | openid | 1234  | unsupported_response_type | unsupported_response_type |
            | authorization_code | 234555    | http://localhost | openid | 1234  | unsupported_grant_type    | unsupported_grant_type    |
            | authorization_code | 234555    | http://localhost | openid | 1234  | invalid_scope             | invalid_scope             |
            | authorization_code | 234555    | http://localhost | openid | 1234  | invalid_client            | invalid_client            |


