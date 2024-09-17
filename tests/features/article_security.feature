Feature: Article security
  As a user
  I want to check the security of the article api
  So that I can ensure that the api is secure

  Background:
    Given I have a valid email and password is "test@gmail.com" and "password"

  Scenario Template: Check article api security
    Given I am unauthenticated with invalid token
    When I send a request to "<action>" the article api
    Then the response status code should be a <status_code>
#    And the response should contain an error message

    Examples:
      | action | status_code |
      | get    | 200         |
      | create | 401         |
      | update | 401         |
      | delete | 401         |

  Scenario Template: Check article api security
    Given I am authenticated with valid token
    When I send a request to "<action>" the article api
    Then the response status code should be a <status_code>
#    And the response should contain an error message

    Examples:
      | action | status_code |
      | get    | 200         |
#      | create | 201         |
#      | update | 200         |
#      | delete | 200         |

