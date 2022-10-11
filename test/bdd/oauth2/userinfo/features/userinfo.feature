Feature: UserInfo Endpoint

    As a Client
    I want to fetch user info
    So that i can have more information about the authenticated user

    Scenario: Successful userInfo request
        Given there is authenticated user using openid connect with following details
            | first_name | middle_name | last_name | phone        | email            | gender |
            | jon        | doe         | john      | 251923456789 | normal@gmail.com | male   |
        When I send userInfo request
        Then I should get correct userInfo response 
    Scenario Outline: Unsuccessful userInfo request
        Given there is invalid access token "<access_token>"
        When I send userInfo request
        Then the request should fail with message "<message>"

        Examples:
            | access_token | message      |
            | token        | Unauthorized |
            | tok          | Unauthorized |


