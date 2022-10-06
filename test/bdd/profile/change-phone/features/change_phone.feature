Feature: Change Phone

    As  a user,
    I want to update my phone number,
    So that I can get access to the system with the updated phone.
    Background:
        Given I am logged in user with the following details
            | first_name | middle_name | last_name | phone        | email            | password | gender |
            | nati       | nati        | nati      | 251923456789 | normal@gmail.com | 123456   | male   |
        And The following user is registered on the system
            | first_name | middle_name | last_name | phone        | email           | password | gender |
            | user1      | user1       | user1     | 251933333333 | user1@gmail.com | 111111   | male   |

    @success
    Scenario Outline: Successful Phone Change
        Given I fill the following details
            | phone   | otp   |
            | <phone> | <otp> |
        When I request to change my phone
        Then I should successfully change my phone

        Examples:
            | phone        | otp    |
            | 251944456789 | 123456 |
    @failure
    Scenario Outline: Phone already exists
        Given I fill the following details
            | phone   | otp   |
            | <phone> | <otp> |
        When I request to change my phone
        Then The phone changing should fail with message "<message>"

        Examples:
            | phone        | otp    | message              |
            | 251933333333 | 123456 | phone already exists |


    @failure
    Scenario Outline: Unsuccessful phone change
        Given I fill the following details with wrong info
            | phone   | otp   |
            | <phone> | <otp> |
        When I request to change my phone
        Then The phone changing should fail with message "<message>"

        Examples:
            | phone        | otp    | message              |
            | 251933333334 | 123    | invalid otp          |
            | 25193333333  | 123456 | invalid phone number |
            | 251933333334 |        | otp is required      |