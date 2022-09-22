Feature: update User Profile

    As a user
    I want to update my profile
    So that I will have upto date information on my profile

    Background: I am logged in user
        Given I am logged in user with the following details:
            | first_name | middle_name | last_name | phone        | email           | password | gender |
            | nati       | nati        | nati      | 251923456789 | admin@gmail.com | 123456   | male   |

    @success
    Scenario Outline: Successful  Profile Update
        Given I fill the form with the following details:
            | first_name   | middle_name   | last_name   | gender   |
            | <first_name> | <middle_name> | <last_name> | <gender> |

        When I update my profile
        Then my profile should be updated

        Examples:
            | first_name | middle_name | last_name | gender |
            | testuser1  | testuser1   | testuser1 | male   |

    @failure
    Scenario Outline: Failed Profile Update
        Given I fill the form with the following details:
            | first_name   | middle_name   | last_name   | gender   |
            | <first_name> | <middle_name> | <last_name> | <gender> |

        When I update my profile
        Then The update should fail with message "<message>"

        Examples:
            | first_name | middle_name | last_name | gender | message                 |
            |            | testuser1   | testuser1 | male   | first name is required  |
            | testuser1  |             | testuser1 | female | middle name is required |
            | testuser1  | testuser1   |           | female | last name is required   |
            | testuser1  | testuser1   | testuser1 |        | gender is required      |
