Feature: update User Profile

    As a user
    I want to update my profile
    So that I will have upto date information on my profile

    Background: I am logged in user
        Given I am logged in user with the following details:
            | first_name | middle_name | last_name | phone      | email           | password |
            | nati       | nati        | nati      | 251923456789 | admin@gmail.com | 123456   |

    @success
    Scenario Outline: Successful  Profile Update
        Given I fill the form with the following details:
            | first_name   | middle_name   | last_name   | phone   | email   |
            | <first_name> | <middle_name> | <last_name> | <phone> | <email> |

        When I update my profile
        Then my profile should be updated

        Examples:
            | first_name | middle_name | last_name | phone      | email           |
            | testuser1  | testuser1   | testuser1 | 251925252525 | test1@gmail.com |

    @failure
    Scenario Outline: Failed Profile Update
        Given I fill the form with the following details:
            | first_name   | middle_name   | last_name   | phone   | email   |
            | <first_name> | <middle_name> | <last_name> | <phone> | <email> |

        When I update my profile
        Then The update should fail with message "<message>"

        Examples:
            | first_name | middle_name | last_name | phone      | email           | message              |
            | testuser1  | testuser1   | testuser1 | 251925252525 | test1gmail.com  | email is not valid   |
            | testuser1  | testuser1   | testuser1 | 33333333   | test1@gmail.com | invalid phone number |
    @failure
    Scenario Outline: Dublicated value
        Given there is user with following details:
            | first_name     | middle_name      | last_name      | phone      | email               |
            | dublicate_name | dublicate_middle | dublicate_last | 251912345678 | dublicate@gmail.com |
        And I fill the form with the following details:
            | first_name   | middle_name   | last_name   | phone   | email   |
            | <first_name> | <middle_name> | <last_name> | <phone> | <email> |

        When I update my profile
        Then The update should fail with error message "<message>"

        Examples:
            | first_name | middle_name | last_name | phone      | email           | message                             |
            | testuser1  | testuser1   | testuser1 | 251925252525 | dublicate@gmail.com | user with this email already exists |
            | testuser1  | testuser1   | testuser1 | 251912345678 | test1@gmail.com | user with this phone already exists |


