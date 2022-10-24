Feature: Get All Current Sessions

    As a user
    I want to see my current sessions on different devices
    So that I can manage my all my sessions in single device

    Background:
        Given I am logged in user with the following details
            | first_name | middle_name | last_name | phone        | email            | password | gender |
            | john       | doe         | jon       | 251923456789 | normal@gmail.com | 123456   | male   |

    @success
    Scenario Outline: Successful get all current sessions
        Given And I have the following sessions on the system
            | id                                   | refresh_token             | ip_address | user_agent                                                                                                                                                                                              |
            | c3a3d3f7-c1a5-4ab5-9a67-cafde0bb4721 | TXoIg917E2LdtwgAM3JbkLwT6 | 127.0.0.1  | Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36 RuxitSynthetic/1.0 v2097774569010745593 t3659150847447165606 athe94ac249 altpriv cvcv=2 smf=0 |
            | 3fb5d963-02f2-4135-a8e8-187c657b558b | TXoIg917E2LdtwgAM3JbkLwY7 | 127.0.0.1  | Mozilla/5.0 (Linux; Android 6.0.1; HTC6545LVW Build/MMB29M; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/81.0.4044.117 Mobile Safari/537.36 [FB_IAB/FB4A;FBAV/268.1.0.54.121;]         |
            | 158dffc1-cb4f-4f8b-903c-dadafe5e58b1 | TXoIg917E2LdtwgAM3JbkLwU8 | 127.0.0.1  | Mozilla/5.0 (Linux; Android 8.1.0; vivo 1808) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.116 Mobile Safari/537.36                                                                          |
        When I request to get my current sessions
        Then I should get the all my sessions