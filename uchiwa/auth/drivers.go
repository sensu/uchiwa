package auth

import (
    "fmt"
    "strings"
    "github.com/kless/osutil/user/crypt"
)

import _ "github.com/kless/osutil/user/crypt/apr1_crypt"
import _ "github.com/kless/osutil/user/crypt/md5_crypt"
import _ "github.com/kless/osutil/user/crypt/sha256_crypt"
import _ "github.com/kless/osutil/user/crypt/sha512_crypt"

func none(u, p string) (*User, error) {
	return &User{}, nil
}

func simple(u, p string) (*User, error) {
	for _, user := range users {
        if u != user.Username {
            continue
        }
        if strings.HasPrefix(user.Password, "{crypt}") {
            password := user.Password;
            password = strings.Replace(password, "{crypt}", "", 1);
            return &user, crypt.NewFromHash(password).Verify(password, []byte(p));
        } else {
            if p == user.Password {
                return &user, nil;
            }
        }
	}
	return &User{}, fmt.Errorf("invalid user '%s' or invalid password", u)
}
