Feature: Update Profile Picture

    As a User,
    I want to update my profile
    So that I will have up to dated picture in my profile

    Background: I am logged in user
        Given I am logged in user with the following details:
            | first_name | middle_name | last_name | phone        | email           | password | gender | profile_picture   |
            | jon        | dou         | john      | 251923456789 | admin@gmail.com | 123456   | male   | <profile_picture> |

    @success
    Scenario: Successful Update
        Given I selected this picture "<picture>"
        When I update my profile picture
        Then my profile picture should be updated

        Examples:
            | picture              |
            | ./assets/girl.jpeg   |
            | ./assets/hacker.jpeg |

    @failure
    Scenario: UnUnsuccessful Update
        Given I selected this picture "<picture>"
        When I update my profile picture
        Then The update should fail with message "<message>"

        Examples:
            | picture             | message                                                 |
            | ./assets/links.txt  | profile_picture must be one of types (png,jpg,jpeg)     |
            | ./assets/sample.pdf | profile_picture must be one of types (png,jpg,jpeg)     |
            | ./assets/big.jpg    | profile_picture size must be less than 1MB              |
