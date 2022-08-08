Feature: Login

    @success
    Scenario: Successful Login
        Given I fill the following details
            | phone   | password   | otp   |
            | <phone> | <password> | <otp> |
        When I submit the registration form
        Then I will be logged in securely to my account
        Examples:
            | phone        | password | otp  |
            | 251911121314 | password | 1234 |

    @invalid
    Scenario: Failed Login
        Given I fill the following details
            | email   | password   |
            | <email> | <password> |
        When I submit the registration form
        Then the login should fail with "<message>"
        Examples:
            | email             | password  | message            |
            | notexample.2f.com | 123456    | invalid credential |
            | example.2f.com    | not123456 | invalid credential |