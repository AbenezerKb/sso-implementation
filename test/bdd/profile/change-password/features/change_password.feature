Feature: Change Password

    As  a User,
    I want to change my password
    So that I will have more secure password

    Background:
        Given I am logged in user with the following details
            | first_name | middle_name | last_name | phone        | email            | password | gender |
            | nati       | nati        | nati      | 251923456789 | normal@gmail.com | 123456   | male   |

    @success
    Scenario Outline: Successful Password Change
        Given I fill the following details
            | old_password   | new_password   |
            | <old_password> | <new_password> |
        When I request to change my password
        Then I should successfully change my password

        Examples:
            | old_password | new_password |
            | 123456       | 654321       |
            
    @failure
    Scenario Outline: Invalid Credential
        Given I fill the following details
            | old_password   | new_password   |
            | <old_password> | <new_password> |
        When I request to change my password
        Then The password changing should fail with message "<message>"

        Examples:
            | old_password | new_password | message            |
            | 123457       | 654321       | invalid credential |

    @failure
    Scenario Outline: Unsuccessful Password Change
        Given I fill the following details
            | old_password   | new_password   |
            | <old_password> | <new_password> |
        When I request to change my password
        Then The password changing should fail with field error message "<message>"

        Examples:
            | old_password | new_password | message            |
            | 12345        | 654321       | password too short |



