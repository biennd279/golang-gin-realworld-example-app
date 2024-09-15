Feature: Check login api security
    As a user
    I want to check the security of the login api
    So that I can ensure that the api is secure

    Scenario: Check login api security
        Given I have a login api
        When I send a request to the login api with the following data
        Then I should get a response from the login api
        And the response should be a 200 status code
        And the response should contain a token
