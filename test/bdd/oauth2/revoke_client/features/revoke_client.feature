Feature: Revoke client access
  As a user,
  I want to to revoke access given to a client
  So that the client won’t access anything on my behalf

  Background:
    Given I am logged in with the following credentials
      | email             | password   |
      | myEmail@gmail.com | myPassword |
    And I have given access to the following client
      | name      | client_type  | redirect_uris      | scopes        | logo_url                                       | secret |
      | newClient | confidential | https://google.com | profile email | https://www.google.com/images/errors/robot.png | secret |

  @success
  Scenario: I successfully revoke access to the client
    When I request to revoke access to the client
    Then The client should no longer have access to my data
    And My action should be recorded

  @failure
  Scenario Outline: I fail to revoke access to the client with invalid request
    When I request to revoke access to the client with id "<client_id>"
    Then My request fails with field error "<message>"
    Examples:
      | client_id      | message           |
      | not-correct-id | invalid client_id |

  @failure
  Scenario Outline: I fail to revoke access to the client with valid request
    When I request to revoke access to the client with id "<client_id>"
    Then My request fails with error message "<message>"
    Examples:
      | client_id                            | message                |
      | 9fb7169c-735c-4638-bb24-7a01a345b0ac | no client access found |
