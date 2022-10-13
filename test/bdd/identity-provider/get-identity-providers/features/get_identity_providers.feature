Feature: Get Identity Providers

    As an Admin,
    I want to see all identity providers registered on Ride+ sso
    So that I can see and manage them easily


    Background:
        Given I am logged in as admin user
            | email           | password      | role                     |
            | admin@gmail.com | adminPassword | update_identity_provider |
        And There are identity provider with the following details
            | name | logo_uri                                       | client_id | client_secret | redirect_uri         | authorization_uri    | token_endpoint_uri | user_info_endpoint_uri |
            | ip_1 | https://www.google.com/images/errors/robot.png | client_1  | secret_1      | https://redirect.com | http://authorize.com | http://token.com   | http.userinfo.com      |
            | ip_2 | https://www.google.com/images/errors/robot.png | client_2  | secret_2      | https://redirect.com | http://authorize.com | http://token.com   | http.userinfo.com      |
            | ip_3 | https://www.google.com/images/errors/robot.png | client_3  | secret_3      | https://redirect.com | http://authorize.com | http://token.com   | http.userinfo.com      |
            | ip_4 | https://www.google.com/images/errors/robot.png | client_4  | secret_4      | https://redirect.com | http://authorize.com | http://token.com   | http.userinfo.com      |

    @success
    Scenario Outline: Successful identity providers
        When I request to get all the identity providers with the following preferences
            | page   | per_page   |
            | <page> | <per_page> |
        Then I should get the list of identity providers that pass my preferences
        Examples:
            | page | per_page |
            | 0    | 10       |
            | 0    | 3        |
            | 1    | 2        |
            | 1    | 5        |