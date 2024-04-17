Feature: User Registration and Login

  Scenario: User Registration with Valid Credentials
    Given I am registering with valid credentials,
    When I submit the registration form,
    Then I should be successfully registered

  Scenario: Duplicate Username Handling
    Given a user with the username "testuser" is already registered
    When I attempt to register with the same username
    Then the system should return an error message indicating that the username already exists

  Scenario: Invalid Email Format
    Given I am registering with an invalid email format,
    When I submit the registration form,
    Then the system should return an error message indicating that the email format is invalid

  Scenario: Weak Password Handling
    Given I am registering with a weak password
    When I submit the registration form
    Then the system should return an error message indicating that the password is not strong enough

  Scenario: Username Length Requirement
    Given Given I am registering with a username less than 5 characters long
    When I submit the registration form
    Then the system should return an error message indicating that the username must be at least 5 characters long

  Scenario: Password Strength Requirement
    Given Given I am registering with a password that does not meet the strength requirements
    When I submit the registration form
    Then the system should return an error message indicating the password requirements

# Feature: User Login

  Scenario: Login with Valid Credentials
    Given I am a registered user with valid credentials
    When I log in with my username and password
    Then the system should generate a JWT token for authentication and issue a refresh token

  Scenario: Login with Invalid Username
    Given I am attempting to log in with an invalid username
    When I submit the login form
    Then the system should return an error message indicating that the username is not registered

  Scenario: Login with Invalid Password
    Given I am attempting to log in with an invalid password
    When I submit the login form
    Then the system should return an error message indicating that the password is incorrect
