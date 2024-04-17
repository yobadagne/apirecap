package features_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cucumber/godog"
	//"github.com/gin-gonic/gin"
	"github.com/yobadagne/user_registration/model"
	"github.com/yobadagne/user_registration/service"
	"github.com/yobadagne/user_registration/util"
)

var usertoregister model.User
var NewServiceLayer = service.NewServiceLayer()

// Scenario 1 User Registration with Valid Credentials
func iAmRegisteringWithValidCredentials() error {
	usertoregister = model.User{
		Username: "eyobdagne",
		Email:    util.RandomEmail(),
		Password: "abcABC123@",
	}
	return nil
}
func iSubmitTheRegistrationForm() error {
	//_, err:= SendHTTPtoRegisterUser(usertoregister)
	return nil
}

func iShouldBeSuccessfullyRegistered() error {
	//we try to read here
	//c := gin.Context{}
	_, err := SendHTTPtoRegisterUser(usertoregister)
	if err != nil{
		return err
	}
	// _, err = NewServiceLayer.GetRegisteredUser(&c, "eyobdagne")
	// return err
	// err := iSubmitTheRegistrationForm(usertoregister)
	// return err
return nil
}

// Scenario 2 Duplicate Username Handling

func aUserWithTheUsernameIsAlreadyRegistered(arg1 string) error {
	usertoregister = model.User{
		Username: arg1,
		Email:    util.RandomEmail(),
		Password: "abcABC123@",
	}
	_ ,err := SendHTTPtoRegisterUser(usertoregister)
	if err != nil {
		return err
	}
	return nil
}

func iAttemptToRegisterWithTheSameUsername() error {
	usertoregister = model.User{
		Username: "testuser",
		Email:    util.RandomEmail(),
		Password: "abcABC123@",
	}
	return nil
}

func theSystemShouldReturnAnErrorMessageIndicatingThatTheUsernameAlreadyExists() error {
	_, err := SendHTTPtoRegisterUser(usertoregister)
	if err != nil{
		return nil
	}
	return err

}

// Scenario 3: Invalid Email Format
func iAmRegisteringWithAnInvalidEmailFormat() error {
	usertoregister = model.User{
		Username: util.RandomUsername(),
		Email: "1234",
		Password: "abcABC123@",
	}
	return nil
}

// func iSubmitTheRegistrationForm() error {
// 	return nil
// }

func theSystemShouldReturnAnErrorMessageIndicatingThatTheEmailFormatIsInvalid() error {
	_, err := SendHTTPtoRegisterUser(usertoregister)
	if err != nil {
		return nil
	}
	return err
}

//Scenario 4: Weak Password Handling
func iAmRegisteringWithAWeakPassword() error {
	usertoregister = model.User{
		Username: util.RandomUsername(),
		Email: util.RandomEmail(),
		Password: "abc",
	}
	return nil
} 
func theSystemShouldReturnAnErrorMessageIndicatingThatThePasswordIsNotStrongEnough() error {
	_, err := SendHTTPtoRegisterUser(usertoregister)
	if err != nil {
		return nil
	}
	return err
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





func iLogInWithMyUsernameAndPassword() error {
	return godog.ErrPending
}

func iSubmitTheLoginForm() error {
	return godog.ErrPending
}

func theSystemShouldGenerateAJWTTokenForAuthenticationAndIssueARefreshToken() error {
	return godog.ErrPending
}



func theSystemShouldReturnAnErrorMessageIndicatingThatThePasswordIsIncorrect() error {
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
	ctx.Step(`^I am registering with an invalid email format$`, iAmRegisteringWithAnInvalidEmailFormat)
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

func SendHTTPtoRegisterUser(usertoregister model.User) (int, error) {
	// Convert registrationForm struct to JSON
	jsonData, err := json.Marshal(usertoregister)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to marshal registration form data: %v", err)
	}
	// Create an HTTP client
	client := &http.Client{}

	// Prepare the request
	req, err := http.NewRequest("POST", "http://localhost:8080/register", bytes.NewBuffer(jsonData))
	if err != nil {
		return http.StatusInternalServerError,fmt.Errorf("failed to create request: %v", err)
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError,fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return http.StatusInternalServerError,fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.StatusCode, nil
}
