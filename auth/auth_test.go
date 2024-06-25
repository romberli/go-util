package auth

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/common"
)

const (
	testSecretKey = "MIIBOwIBAAJBAMlVXSjfNKcUgmeqzbHON8ASYDZE0zyEsa8l8j0r2NPhggMC4VJEazJEOqrCq0kuERW+zK0AwpfbYBHth6r3lN0CAwEAAQJBALtM3/sLE6ewK9UXkH6usyzLq5gxFTcC125y5dXEudX6GDkQ7+c9WCMutDBF40D9xCvYfSVlNInBAGZVcC33WcECIQDRcFBwXIdXzj0lecjhkepkJHdC7+3zcDKx3lvj6rKxzQIhAPYXwhHL27AhvJ931dXL5tJGsajx5/xANAGZAn14+59RAiBbjfaL99buamjOfhtziB7nog1EhLAHcC+pE6Ql0Q5GrQIgUKSQcAyBvUIQ8aDvbdQXm6iW52n+P2c6o5tkeYF/00ECIQCyeOPbrbD8QMDkZzrvgKBMIG6ZW/hBTNXoTet0y3GB+Q=="
)

var testAuth *Auth

func init() {
	testAuth = NewAuth(common.StringToBytes(testSecretKey))
}

func TestAuth_Sign(t *testing.T) {
	asst := assert.New(t)

	jwtToken, err := testAuth.Sign()
	asst.Nil(err, common.CombineMessageWithError("test Sign() failed", err))
	t.Log(jwtToken)
}

func TestAuth_SignWithClaims(t *testing.T) {
	asst := assert.New(t)

	claims := jwt.MapClaims{
		"username": "test",
		"password": "test",
	}
	jwtToken, err := testAuth.SignWithMethodAndClaims(DefaultSignMethod, claims, nil)
	asst.Nil(err, common.CombineMessageWithError("test SignWithMethodAndClaims() failed", err))
	t.Log(jwtToken)
}

func TestAuth_ParseUnverified(t *testing.T) {
	asst := assert.New(t)

	claims := jwt.MapClaims{
		"username": "test",
		"password": "test",
	}
	token, err := testAuth.SignWithMethodAndClaims(DefaultSignMethod, claims, nil)
	asst.Nil(err, common.CombineMessageWithError("test ParseUnverified() failed", err))

	u := &struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	err = testAuth.ParseUnverified(token, &u)
	asst.Nil(err, common.CombineMessageWithError("test ParseUnverified() failed", err))
	t.Logf("%+v", u)
}
