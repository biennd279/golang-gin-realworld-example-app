

Feature: Check login api security
    As a user
    I want to check the security of the login api
    So that I can ensure that the api is secure

    Background:
        Given I have a valid email and password is "test@gmail.com" and "password"
        And I have a invalid username and password is "invalid" and "invalid"

    Scenario: Check login api security
        When I send a request to the login api with valid credentials
        Then the response should be a 200 status code
        And the response should contain a token
