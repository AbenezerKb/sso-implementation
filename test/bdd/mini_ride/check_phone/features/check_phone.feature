Feature: Check phone for mini-ride

    As a MiniRide
    I want to check the existence of a phone number on the sso
    So that I can update/create updated/new drivers and stream the event's back to sso and ridePlus

    Background: populate user's
        Given I am authenticated with following credential
            | username | password |
            | username | password |
        And they are the following user's on sso
            | id                                   | first_name | middle_name | last_name | phone        | profile_picture                                                                                                                       | status |
            | 06eb340a-862a-4dd0-8a3f-5e4c1f767d3d | abebe      | kebede      | teshome   | 251944123345 | image                                                                                                                                 | ACTIVE |
            | 495f6800-dd63-49e2-9809-107076ed2c72 | Surafel    | Zerihun     | Surafel   | 251967968549 | https://onde-images.s3.amazonaws.com/profile/2021-06-08/0333b19d-9a8e-4597-95f2-cd2379504c36-bfe7b669-40d9-4d03-9a8f-4d78feb93708.png | ACTIVE |

    Scenario Outline: Successfull check
        When I request to check users with the following phone "<check_phone>"
        Then I should get the following response
            | id   | first_name   | middle_name   | last_name   | phone   | profile_picture   | status   | exists   |
            | <id> | <first_name> | <middle_name> | <last_name> | <phone> | <profile_picture> | <status> | <exists> |

        Examples:
            | id                                   | first_name | middle_name | last_name | phone        | profile_picture | status | exists | check_phone  |
            | 06eb340a-862a-4dd0-8a3f-5e4c1f767d3d | abebe      | kebede      | teshome   | 251944123345 | image           | ACTIVE | true   | 251944123345 |
            | 06eb340a-862a-4dd0-8a3f-5e4c1f767d3d |            |             |           |              |                 |        | false  | 251967968579 |

    Scenario Outline: Unsuccessfull check
        When I request to check users with the following phone "<check_phone>"
        Then I should get error with message "<message>"
        Examples:
            | check_phone   | message              |
            | invalid_phone | invalid phone number |