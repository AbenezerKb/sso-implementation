Feature: Refresh Token
  as a user
  I want my  access token to be refreshed
  so that I do not have to authenticate every time my access token expires

  Scenario: Refresh Token is not expired
    Given There is a registered user on the system:
      | first_name | middle_name | last_name | phone      | email            | password |
      | testuser1  | testuser1   | testuser1 | 0925252595 | test11@gmail.com | 1234567  |
    And I am logged in to the system and have a refresh token:
      | refresh_token   | expires_at                          |
      | myR1fr35h70k3n | 2023-09-26T09:06:36.525293389+03:00 |
    When  I refresh my access token using my refresh token
      | refresh_token  |
      | myR1fr35h70k3n |
    Then I should get a new access token

