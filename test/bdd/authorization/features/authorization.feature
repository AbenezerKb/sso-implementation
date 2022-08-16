Feature: obtaining authorization

    As a client

    I want to obtain authorization from the resource owner

    So that I can request to get an access token or refresh token.

    Scenario Outline: Succesfull Obtaining Authorization
        Given I have the following parameters:
            | response_type   | client_id   | redirect_uri   | scope   | state   |
            | <response_type> | <client_id> | <redirect_uri> | <scope> | <state> |

        When I send a POST request
        Then I should be redirected to "<consent_uri>" with the following success parameters:
            | consentId   | state   |
            | <consentId> | <state> |
        Examples:
            | response_type | client_id     | redirect_uri           | scope  | state | consentId | state | consent_uri            |
            | code          | 3749027981234 | https://www.google.com/ | openid | 1234  | 1234      | 1234  | https://www.google.com/ |

    Scenario Outline: Unable to Obtain Authorization
        Given I have the following parameters:
            | response_type   | client_id   | redirect_uri   | scope   | state   |
            | <response_type> | <client_id> | <redirect_uri> | <scope> | <state> |

        When I send a POST request
        Then I should be redirected to "<redirect_uri>" with the following error parameters:
            | error   | error_description   | state   |
            | <error> | <error_description> | <state> |
        Examples:
            | response_type      | client_id | redirect_uri              | scope    | state | error                | error_description         |
            | code               |           | https://www.google.com/   | openid   | 1234  | invalid_request      | client_id is required.    |
            | code               | 234555    |                           | openid   | 1234  | invalid_request      | redirect_uri is required. |
            | authorization_code | 234555    | https://www.google.com/   | openid   | 1234  | invalid_request      | must be a valid value.    |
            | code               | 234555    | https://www.google.com/   | closedid | 1234  | invalid_request      | must be a valid value.    |
            | code               | 234555    | localhostts:9000/callback | openid   | 1234  | invalid_redirect_uri | invalid redirect uri      |


