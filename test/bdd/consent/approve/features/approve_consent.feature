Feature: Approve consent
  Background:
    Given I am logged in with credentials
      | email             | password |
      | consent@gmail.com | consent  |
    And There are registered scopes with the following details
      | name   | description        | resource_server_name |
      | openid | your profile info  | sso                  |
      | email  | your default email | sso                  |
    And There is a client with the following details
      | name | redirect_uris        | secret    | scopes       | client_type  | logo_url               |
      | ride | https://www.google.com | my_secret | openid email | confidential | http://logo.client.com |
    And I have a consent with the following details
      | id                                   | scope  | redirect_uri         | response_type | approved | state    | prompt  |
      | 48684fe2-43fa-46b8-ba6b-78cfc7196fb8 | openid | https://www.google.com | code          | false    | my_state | consent |
  @success
  Scenario: Consent is approved
    When I request consent approval with id "48684fe2-43fa-46b8-ba6b-78cfc7196fb8"
    Then The consent should be approved
  @failure
  Scenario Outline: consent is not approved
    When I request consent approval with id "<consent_id>"
    Then Consent approval should fail with message "<message>"
    Examples:
    | consent_id | message |
    | not_an_id | invalid consentId |