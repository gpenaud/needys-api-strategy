Feature: get all registered strategies using an API
  As a unauthenticated user
  I want to be able to get all registered strategies as a json array through an API endpoint
  So that I can remind that a lot of suitables strategies exist, so no need to despair.

  Scenario: doing a wrong query with GET method to get all strategies
    When I send "GET" request to "/strategies"
    Then the response code should be 200

  Scenario: doing a wrong query with PUT method to get all strategies
    When I send "POST" request to "/strategy"
    Then the response code should be 201
