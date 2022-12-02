Feature: Refresh Token
    as a client
    I want my  access token to be refreshed
    so that I do not have to authenticate every time my access token expires
    Background: There is a client registered on the system
        Given There is a regeistered user on the system:
            | first_name | middle_name | last_name | phone      | email            | password |
            | testuser1  | testuser1   | testuser1 | 0925252595 | test11@gmail.com | 1234567  |
        And There is a client on the system:
            | name | redirect_uris        | secret    | scopes       | client_type  | logo_url               |
            | ride | http://localhost.com | my_secret | openid email | confidential | http://logo.client.com |
        And The user grants access to the client:
            | refresh_token            | expires_at                          | scope        |
            | +toNc!tKC8q;,SXt7h%iu#aX | 2023-09-26T09:06:36.525293389+03:00 | openid email |

    Scenario: Refresh Token successfully
        When I refresh the access token:
            | grant_type    | refresh_token            |
            | refresh_token | +toNc!tKC8q;,SXt7h%iu#aX |
        Then I should get a new access token with the old refresh token

    Scenario Outline:missing required inputs
        When I refresh the access token:
            | grant_type   | refresh_token   |
            | <grant_type> | <refresh_token> |
        Then The request should fail with field error "<error_message>":
        Examples:
            | grant_type    | refresh_token            | error_message             |
            | refresh_token |                          | refresh_token is required |
            |               | +toNc!tKC8q;,SXt7h%iu#aX | grant_type is required    |

    Scenario Outline: Refresh Token is  expired
        Given I have an expired refresh token:
            | refresh_token            | expires_at                          | scope        |
            | 0ohn8ktKC8q;,SXt7h%iu#aX | 2022-08-26T09:06:36.525293389+03:00 | openid email |
        When I refresh the access token:
            | grant_type    | refresh_token            |
            | refresh_token | 0ohn8ktKC8q;,SXt7h%iu#aX |
        Then The request should fail with error message "<error_message>":
        Examples:
            | error_message         |
            | refresh token expired |