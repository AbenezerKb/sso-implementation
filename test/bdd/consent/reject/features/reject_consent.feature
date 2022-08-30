Feature: Reject consent

  Background:
    Given I am logged in with credentials
      | email             | password |
      | consent@gmail.com | consent  |
    And There are registered scopes with the following details
      | name   | description        | resource_server_name |
      | openid | your profile info  | sso                  |
      | email  | your default email | sso                  |
    And There is a client with the following details
      | name | redirect_uris          | secret    | scopes       | client_type  | logo_url               |
      | ride | https://www.google.com | my_secret | openid email | confidential | http://logo.client.com |
    And I have a consent with the following details
      | id                                   | scope  | redirect_uri           | response_type | approved | state    | prompt  |
      | 48684fe2-43fa-46b8-ba6b-78cfc7196fb8 | openid | https://www.google.com | code          | false    | my_state | consent |

  @success
  Scenario Outline: Consent is rejected
    When I request consent rejection with id "48684fe2-43fa-46b8-ba6b-78cfc7196fb8" and message "<message>"
    Then The consent should be rejected
    Examples:
      | message               |
      |                       |
      | login is required     |
      | user rejected consent |

  @failure
  Scenario Outline: consent is not rejected correctly
    When I request consent rejection with id "48684fe2-43fa-46b8-ba6b-78cfc7196fb8" and message "<message>"
    Then Consent rejection should fail with message "<error>"
    Examples:
      | message           | error             |
      | login is required | invalid consentId |