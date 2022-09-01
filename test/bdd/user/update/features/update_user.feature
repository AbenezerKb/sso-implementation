Feature: update User Profile

    As a user
    I want to update my profile
    So that I will have upto date information on my profile

    Background: I am logged in user
        Given I am logged in user with the following details:
            | first_name | middle_name | last_name | phone     | email           |
            | nati       | nati        | nati      | 092345678 | admin@gmail.com |

    @success
    Scenario Outline: Successful  Profile Update
        Given I fill the form with the following details:
            | first_name   | middle_name   | last_name   | phone   | email   |
            | <first_name> | <middle_name> | <last_name> | <phone> | <email> |

        When I update my profile
        Then my profile should be updated

        Examples:
            | first_name | middle_name | last_name | phone      | email           |
            | testuser1  | testuser1   | testuser1 | 0925252525 | test1@gmail.com |

    @failure
    Scenario Outline: Failed Profile Update
        Given I fill the form with the following details:
            | first_name   | middle_name   | last_name   | phone   | email   |
            | <first_name> | <middle_name> | <last_name> | <phone> | <email> |

        When I update my profile
        Then The update should fail with message "<message>"

        Examples:
            | first_name | middle_name | last_name | phone      | email           | message              |
            | testuser1  | testuser1   | testuser1 | 0925252525 | test1gmail.com  | email is not valid   |
            | testuser1  | testuser1   | testuser1 | 33333333   | test1@gmail.com | invalid phone number |

