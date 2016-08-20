package authentication

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthSimplePlain(t *testing.T) {

	users = []User{}

	// No authorized users
	user, err := simple("admin", "test")
	assert.Equal(t, &User{}, user)
	assert.NotNil(t, err)

	admin := User{Username: "admin", Password: "test"}
	users = append(users, admin)

	// Successful login
	user, err = simple("admin", "test")
	assert.Equal(t, "admin", user.Username)
	assert.Nil(t, err)

	// Wrong password for existing user
	user, err = simple("admin", "testwrong")
	assert.Equal(t, &User{}, user)
	assert.NotNil(t, err)

	// APR MD5 hash
	users[0].Password = "{crypt}$apr1$YhYWYmA/$QE2UAxx9.tLWGZiLt9nPF."
	user, err = simple("admin", "testapr")
	assert.Equal(t, "admin", user.Username)
	assert.Nil(t, err)

	// MD5-crypt hash
	users[0].Password = "{crypt}$1$3B039HF6$8/5NgCnTm/WzUenUJtBQn1"
	user, err = simple("admin", "testmd5")
	assert.Equal(t, "admin", user.Username)
	assert.Nil(t, err)

	// SHA256-crypt hash
	users[0].Password = "{crypt}$5$Y9fqjXnx.y2jzcJ/$8muV4OL12axrvjG.uPvId0sVvkjJ5t9M2SuV2FTbrm/"
	user, err = simple("admin", "testsha256")
	assert.Equal(t, "admin", user.Username)
	assert.Nil(t, err)

	// SHA512-crypt hash
	users[0].Password = "{crypt}$6$KYg/Ceo.HQW6R7D0$4d3vWTQRhoyC27nNusOo/fnYh6wEmOm1YoQsTz4K9mYlmoUGL.LWO2ez1g5QqZFamWgo.VSeOEHLwe0EeKAC3/"
	user, err = simple("admin", "testsha512")
	assert.Equal(t, "admin", user.Username)
	assert.Nil(t, err)

	// Wrong password for existing user with SHA512-crypt hash
	user, err = simple("admin", "testwrong")
	assert.NotNil(t, err)

}
