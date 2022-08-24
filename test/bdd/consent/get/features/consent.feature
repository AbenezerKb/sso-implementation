Feature: Get Consent By ID
  As a user
  I want to get a consent by ID
  So that I can see the consent details

  Background: I have consent with the following details
    Given I am logged in with credentials
      | email             | password |
      | consent@gmail.com | consent  |
    And There are registered scopes with the following details
      | name   | description        | resource_server_name |
      | openid | your profile info  | sso                  |
      | email  | your default email | sso                  |
    And There is a client with the following details
      | name | redirect_uris        | secret    | scopes       | client_type  | logo_url               |
      | ride | http://localhost.com | my_secret | openid email | confidential | http://logo.client.com |
    And I have a consent with the following details
      | id                                   | scope  | redirect_uri         | response_type | approved | state    | prompt  |
      | 48684fe2-43fa-46b8-ba6b-78cfc7196fb8 | openid | http://localhost.com | code          | false    | my_state | consent |

  @success
  Scenario: Valid request
    Given I have a consent with ID "48684fe2-43fa-46b8-ba6b-78cfc7196fb8"
    When I request consent Data
    Then I should get valid consent data

  @failure
  Scenario Outline: invalid request
    Given I have a consent with ID "<consent_id>"
    When I request consent Data
    Then I should get error with message "<message>" and field error "<field_error>"
    Examples:
      | consent_id                           | message           | field_error        |
      | invalid_consent_id                   | consent not found |                    |
      | 12684fe2-43fa-46b8-ba6b-78cfc7196fb8 | consent not found |                    |