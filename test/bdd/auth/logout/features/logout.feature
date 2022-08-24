Feature: logout

    As a user
    I want to log out of the system
    so that I can have a clear session on a particular device

    @success
    Scenario Outline: Successful Logout
        Given I am a logedin  user with the following details:
            | email             | password |
            | example@email.com | 1234abcd |
        When I logout
        Then I should Successfully logout of the system
