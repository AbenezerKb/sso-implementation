Feature: Get All Clients
  As an Admin,
  I want to get all registered clients
  So that I can view details of each client

  Background:
    Given The following clients are registered on the system
      | name        | client_type  | redirect_uris          | scopes          | logo_url                               | status   |
      | clientOne   | confidential | https://google.com     | profile email   | https://ww.google.com/error-image1.png | active   |
      | clientTwo   | public       | https://youtube.com    | profile balance | https://ww.google.com/error-image2.png | inactive |
      | clientThree | public       | https://facebook.com   | profile openid  | https://ww.google.com/error-image3.png | active   |
      | clientFour  | confidential | https://google.com     | profile         | https://ww.google.com/error-image4.png | active   |
      | clientFive  | confidential | https://2f-capital.com | openid          | https://ww.google.com/error-image5.png | inactive |
    And I am logged in as admin user
      | email           | password      | role       |
      | admin@gmail.com | adminPassword | super-user |

  @success
  Scenario Outline: I get all the clients
    When I request to get all the clients with the following preferences
      | page   | per_page   |
      | <page> | <per_page> |
    Then I should get the list of clients that pass my preferences
    Examples:
      | page | per_page |
      | 0    | 10       |
      | 0    | 3        |
      | 1    | 2        |
      | 1    | 5        |

  Scenario Outline: I fail to get all the clients due to invalid request
    When I request to get all the clients with the following preferences
      | page   | per_page   |
      | <page> | <per_page> |
    Then I should get error message "<message>"
    Examples:
      | page  | per_page | message               |
#      | hello | 10       | invalid filter params |
#      | 1     | hello    | invalid filter params |
