Feature: Get OpenID Authorized Clients
  As a user,
  I want to get all clients that are logged in with my SSO account,
  So that I can view details of each client.

  Background:
    Given I am logged in as the following user
      | email          | password     |
      | user@gmail.com | userPassword |
    And I have given authorization for the following clients
      | name        | client_type  | redirect_uris        | scopes               | logo_url               | granted_scopes |
      | clientOne   | confidential | https://google.com   | profile email        | https://www.google.com | profile email  |
      | clientTwo   | public       | https://facebook.com | profile email openid | https://www.google.com | email openid   |
      | clientThree | confidential | https://yahoo.com    | openid               | https://www.google.com | openid         |

  @success
  Scenario: I get openid authorized clients
    When I request to get openid authorized clients
    Then I should get the list of openid authorized clients
