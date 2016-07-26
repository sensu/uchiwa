package authentication

import "fmt"

func (a *Config) login(user, pass string) (*User, error) {
	// Authenticate the user with the authentication driver
	u, err := a.DriverFn(user, pass)
	if err != nil {
		return nil, fmt.Errorf("Authentication failed: %s", err)
	}

	// Obfuscate the user's salt & hash
	u.PasswordHash = ""
	u.PasswordSalt = ""

	token, err := GetToken(&u.Role, user)
	if err != nil {
		return nil, fmt.Errorf("Authentication failed, could not create the token: %s", err)
	}

	// Add token to the user struct
	u.Token = token

	return u, nil
}
