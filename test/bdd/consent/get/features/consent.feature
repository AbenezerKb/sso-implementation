Feature: Get Conset By ID
    As a user
    I want to get a conset by ID
    So that I can see the conset details

    Background: I have consent with the following details
        Given I have a consent with the following details
            | client_id                            | client_name | scopes | description | client_type  | redirect_uri     | response_type | user_id                              | consent_id                           |
            | 48684fe2-43fa-46b8-ba6b-78cfc7196fb8 | ride        | openid | description | confidential | http://localhost | code          | 48684fe2-43fa-46b8-ba6b-78cfc7196fb8 | 48684fe2-43fa-46b8-ba6b-78cfc7196fb8 |
    @success
    Scenario: Valid Conset ID
        Given I have a consent with ID "<consent_id>"
        And user with ID "<user_id>"
        When I request consent Data
        Then I should get valid consent data
        Examples:
            | consent_id                           | user_id                              |
            | 48684fe2-43fa-46b8-ba6b-78cfc7196fb8 | 48684fe2-43fa-46b8-ba6b-78cfc7196fb8 |
    @error
    Scenario Outline: Invalid Conset ID
        Given I have a consent with ID "<consent_id>"
        And user with ID "<user_id>"
        When I request consent Data
        Then I should get error "<error_message>"

        Examples:
            | consent_id                           | user_id                              | error_message     |
            | nati4fe2-34af-46b8-ba6b-78cfc7196fo8 | 48684fe2-43fa-46b8-ba6b-78cfc7196fb8 | consent not found |
    @error
    Scenario Outline: Invalid User ID
        Given I have a consent with ID "<consent_id>"
        And Invalid user ID "<user_id>"
        When I request consent Data
        Then I should get error "<error_message>"
        Examples:
            | consent_id                           | user_id | error_message  |
            | 48684fe2-43fa-46b8-ba6b-78cfc7196fb8 | 2       | user not found |
