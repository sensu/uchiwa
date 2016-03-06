package auth;


import (
    "testing"
    "github.com/stretchr/testify/assert"
)


func TestAuthSimplePlain(t *testing.T) {

    users = []User{};

    user, err := simple("admin", "test");
    assert.Equal(t, &User{}, user);
    assert.NotNil(t, err);

    admin := User{Username: "admin", Password: "test"};
    users = append(users, admin);

    user, err = simple("admin", "test");
    assert.Equal(t, "admin", user.Username);
    assert.Nil(t, err);

    user, err = simple("admin", "testwrong");
    assert.Equal(t, &User{}, user);
    assert.NotNil(t, err);

    admin.Password = "$6$rounds=1000$$vDKCc9rOGoJbVpjMvQImrBCbdha0O.xOzDOISi93TtOhw50y5pfOawbWUBl/.bvAQ9GYV3/rTXJemXg429BHy/";
    user, err = simple("admin", "test");
    assert.Equal(t, "admin", user.Username);
    assert.Nil(t, err);

    user, err = simple("admin", "testwrong");
    assert.Equal(t, &User{}, user);
    assert.NotNil(t, err);

}
