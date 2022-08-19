Feature: Code Grant Flow

    Background: I have an authorization code
        Given Their is a client
        And Their is a user
        And The user granted access to the client:
            | code                                 |
            | 002a05af-15ae-4c21-8e4a-d7bde55f6aff |
    @success
    Scenario: Access Token successfully issued
        Given I have the following parameters:
            | grant_type         | code                                 | redirect_uri           |
            | authorization_code | 002a05af-15ae-4c21-8e4a-d7bde55f6aff | https://www.google.com |
        When The client request for token
        Then Token should successfully be issued

    @failure
    Scenario Outline: Issuing Access Token Failed
        Given I have the following parameters:
            | grant_type   | code   | redirect_uri   |
            | <grant_type> | <code> | <redirect_uri> |
        When The client request for token
        Then The request should fail with field error "<field_error>" and message "<error_message>"
        Examples:
            | grant_type         | code                                 | redirect_uri            | field_error              | error_message           |
            | authorization_code |                                      | https://www.google.com/ | code is required         |                         |
            | authorization_code | 002a05af-15ae-4c21-8e4a-d7bde55f6aff |                         | redirect_uri is required |                         |
            | authorization_code | 002a05af-15ae-4c21-8e4a-d7bde55f6aff | something               | invalid redirect_uri     |                         |
            |                    | 002a05af-15ae-4c21-8e4a-d7bde55f6aff | https://www.google.com/ | grant_type is required   |                         |
            | authorization_code | 002a05af-15ae-4c21-8e4a-d7bde55f6afz | https://www.google.com/ |                          | no record of code found |

