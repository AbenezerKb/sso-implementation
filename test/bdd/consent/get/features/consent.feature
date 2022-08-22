Feature: Get Conset By ID
  As a user
  I want to get a conset by ID
  So that I can see the conset details

  Background: I have consent with the following details
    Given I am logged in with credentials
      | email             | password | role    |
      | consent@gmail.com | consent  | consent |
    And I have a consent with the following details
      | client_id                            | client_name | scopes | description | client_type  | redirect_uri     | response_type | consent_id                           |
      | 48684fe2-43fa-46b8-ba6b-78cfc7196fb8 | ride        | openid | description | confidential | http://localhost | code          | 48684fe2-43fa-46b8-ba6b-78cfc7196fb8 |

  @success
  Scenario Outline: Valid Conset ID
    Given I have a consent with ID "<consent_id>"
    When I request consent Data
    Then I should get valid consent data
    Examples:
      | consent_id                           | 
      | 48684fe2-43fa-46b8-ba6b-78cfc7196fb8 |

