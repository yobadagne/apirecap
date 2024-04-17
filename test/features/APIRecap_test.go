package features_test

import (
	"github.com/gin-gonic/gin"

	"github.com/cucumber/godog"
	"github.com/yobadagne/user_registration/model"
	"github.com/yobadagne/user_registration/service"
)

func aUserWithTheUsernameIsAlreadyRegistered(arg1 string) error {
	return godog.ErrPending
}

func givenIAmRegisteringWithAPasswordThatDoesNotMeetTheStrengthRequirements() error {
	return godog.ErrPending
}

func givenIAmRegisteringWithAUsernameLessThanCharactersLong(arg1 int) error {
	return godog.ErrPending
}

func iAmARegisteredUserWithValidCredentials() error {
	return godog.ErrPending
}

func iAmAttemptingToLogInWithAnInvalidPassword() error {
	return godog.ErrPending
}

func iAmAttemptingToLogInWithAnInvalidUsername() error {
	return godog.ErrPending
}

func iAmRegisteringWithAWeakPassword() error {
	return godog.ErrPending
}

func iAmRegisteringWithAnInvalidEmailFormat() error {
	return godog.ErrPending
}

func iAmRegisteringWithValidCredentials() error {
	usertoregister := model.User{
		Username: "eyobdagne",
		Email:    "yobadagne@gmail.com",
		Password: "abcABC123@",
	}
	err := service.NewServiceLayer().Register(usertoregister, &gin.Context{})
	return err
}

func iAttemptToRegisterWithTheSameUsername() error {
	return godog.ErrPending
}

func iLogInWithMyUsernameAndPassword() error {
	return godog.ErrPending
}

func iShouldBeSuccessfullyRegistered() error {
	return godog.ErrPending
}

func iSubmitTheLoginForm() error {
	return godog.ErrPending
}

// func iSubmitTheRegistrationForm() error {
// 	return godog.ErrPending
// }

func iSubmitTheRegistrationForm() error {
	return godog.ErrPending
}

func theSystemShouldGenerateAJWTTokenForAuthenticationAndIssueARefreshToken() error {
	return godog.ErrPending
}

func theSystemShouldReturnAnErrorMessageIndicatingThatTheEmailFormatIsInvalid() error {
	return godog.ErrPending
}

func theSystemShouldReturnAnErrorMessageIndicatingThatThePasswordIsIncorrect() error {
	return godog.ErrPending
}

func theSystemShouldReturnAnErrorMessageIndicatingThatThePasswordIsNotStrongEnough() error {
	return godog.ErrPending
}

func theSystemShouldReturnAnErrorMessageIndicatingThatTheUsernameAlreadyExists() error {
	return godog.ErrPending
}

func theSystemShouldReturnAnErrorMessageIndicatingThatTheUsernameIsNotRegistered() error {
	return godog.ErrPending
}

func theSystemShouldReturnAnErrorMessageIndicatingThatTheUsernameMustBeAtLeastCharactersLong(arg1 int) error {
	return godog.ErrPending
}

func theSystemShouldReturnAnErrorMessageIndicatingThePasswordRequirements() error {
	return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^a user with the username "([^"]*)" is already registered$`, aUserWithTheUsernameIsAlreadyRegistered)
	ctx.Step(`^Given I am registering with a password that does not meet the strength requirements$`, givenIAmRegisteringWithAPasswordThatDoesNotMeetTheStrengthRequirements)
	ctx.Step(`^Given I am registering with a username less than (\d+) characters long$`, givenIAmRegisteringWithAUsernameLessThanCharactersLong)
	ctx.Step(`^I am a registered user with valid credentials$`, iAmARegisteredUserWithValidCredentials)
	ctx.Step(`^I am attempting to log in with an invalid password$`, iAmAttemptingToLogInWithAnInvalidPassword)
	ctx.Step(`^I am attempting to log in with an invalid username$`, iAmAttemptingToLogInWithAnInvalidUsername)
	ctx.Step(`^I am registering with a weak password$`, iAmRegisteringWithAWeakPassword)
	ctx.Step(`^I am registering with an invalid email format,$`, iAmRegisteringWithAnInvalidEmailFormat)
	ctx.Step(`^I am registering with valid credentials,$`, iAmRegisteringWithValidCredentials)
	ctx.Step(`^I attempt to register with the same username$`, iAttemptToRegisterWithTheSameUsername)
	ctx.Step(`^I log in with my username and password$`, iLogInWithMyUsernameAndPassword)
	ctx.Step(`^I should be successfully registered$`, iShouldBeSuccessfullyRegistered)
	ctx.Step(`^I submit the login form$`, iSubmitTheLoginForm)
	ctx.Step(`^I submit the registration form$`, iSubmitTheRegistrationForm)
	ctx.Step(`^I submit the registration form,$`, iSubmitTheRegistrationForm)
	ctx.Step(`^the system should generate a JWT token for authentication and issue a refresh token$`, theSystemShouldGenerateAJWTTokenForAuthenticationAndIssueARefreshToken)
	ctx.Step(`^the system should return an error message indicating that the email format is invalid$`, theSystemShouldReturnAnErrorMessageIndicatingThatTheEmailFormatIsInvalid)
	ctx.Step(`^the system should return an error message indicating that the password is incorrect$`, theSystemShouldReturnAnErrorMessageIndicatingThatThePasswordIsIncorrect)
	ctx.Step(`^the system should return an error message indicating that the password is not strong enough$`, theSystemShouldReturnAnErrorMessageIndicatingThatThePasswordIsNotStrongEnough)
	ctx.Step(`^the system should return an error message indicating that the username already exists$`, theSystemShouldReturnAnErrorMessageIndicatingThatTheUsernameAlreadyExists)
	ctx.Step(`^the system should return an error message indicating that the username is not registered$`, theSystemShouldReturnAnErrorMessageIndicatingThatTheUsernameIsNotRegistered)
	ctx.Step(`^the system should return an error message indicating that the username must be at least (\d+) characters long$`, theSystemShouldReturnAnErrorMessageIndicatingThatTheUsernameMustBeAtLeastCharactersLong)
	ctx.Step(`^the system should return an error message indicating the password requirements$`, theSystemShouldReturnAnErrorMessageIndicatingThePasswordRequirements)
}