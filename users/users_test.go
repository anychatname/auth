package users

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var (
	nickname        = "exampleuser"
	invalidNickname = "$**exampleNickname"
	password        = "1234"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("Given a valid user with a external sign When creating new user Then same user with User.CreatedAt field filled", func(t *testing.T) {
		user := User{
			Nickname:   "exampleuser",
			SignedWith: []ExternalSigned{{Platform: "examplePlatform"}},
		}
		gotUser, err := New(user)
		assert.NoError(t, err)

		if equalDateTime(t, time.Now(), gotUser.CreatedAt) {
			user.CreatedAt = time.Time{}
			gotUser.CreatedAt = time.Time{}
		}
		assert.Equal(t, user, gotUser)
	})
}

func TestErrorNew(t *testing.T) {
	t.Parallel()

	type args struct {
		userR User
	}
	tests := []struct {
		name     string
		args     args
		wantUser User
		wantErr  string
	}{
		{
			name: "Given a user with a invalid username When creating new user Then invalid nickname error",
			args: args{
				userR: User{
					Nickname: invalidNickname,
					Password: password,
				},
			},
			wantErr: fmt.Sprintf("invalid nickname: invalid nickname format of %s", invalidNickname),
		},
		{
			name: "Given a user with a empty password When creating new user Then invalid pasword error",
			args: args{
				userR: User{
					Nickname: nickname,
				},
			},
			wantErr: "invalid password: empty or nil value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := New(tt.args.userR)
			assert.EqualError(t, err, tt.wantErr)
		})
	}
}
func TestHashPassword(t *testing.T) {
	t.Parallel()

	t.Run("Given valid password When encrypting password provided Then success", func(t *testing.T) {
		pw := password
		assert.NoError(t, HashPassword(&pw))
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(pw), []byte(password)))
	})
}

func TestErrorHashPassword(t *testing.T) {
	t.Parallel()

	t.Run("Given a nil or empty password When validating password content Then invalid password error", func(t *testing.T) {
		err := HashPassword(nil)
		assert.EqualError(t, err, "invalid password: empty or nil value")

		emptyStr := ""
		err = HashPassword(&emptyStr)
		assert.EqualError(t, err, "invalid password: empty or nil value")
	})
}

func TestValidateNickname(t *testing.T) {
	t.Parallel()

	t.Run("Given a valid nickname When validating nickname Then success", func(t *testing.T) {
		assert.NoError(t, ValidateNickname(nickname))
	})
}

func TestErrorValidateNickname(t *testing.T) {
	t.Parallel()

	t.Run("Given invalid or empty nickname format When validating nickname Then invalid nickname format", func(t *testing.T) {
		err := ValidateNickname(invalidNickname)
		assert.EqualError(t, err, fmt.Sprintf("invalid nickname: invalid nickname format of %s", invalidNickname))
	})
}

func TestValidateEmail(t *testing.T) {
	t.Parallel()

	t.Run("Given valid email When validating email Then success", func(t *testing.T) {

		assert.NoError(t, ValidateEmail("example@gmail.com"))
	})
}

func TestErrorValidateEmail(t *testing.T) {
	t.Parallel()

	t.Run("Given invalid email format When validating email format Then invalid email format", func(t *testing.T) {
		t.Parallel()

		invalidEmail := "invalid.email.format"

		err := ValidateEmail(invalidEmail)
		assert.Contains(t, err.Error(), fmt.Sprintf("invalid email format: %s is not valid, cause", invalidEmail))
	})

	t.Run("Given invalid email host When validating email host Then invalid email host error", func(t *testing.T) {
		t.Parallel()

		invalidHost := "invalid.host"

		err := ValidateEmail("example@" + invalidHost)
		assert.EqualError(t, err, fmt.Sprintf("invalid email host: %s not exists", invalidHost))
	})
}

func equalDateTime(t *testing.T, expected, got time.Time) bool {
	t.Helper()

	return assert.True(
		t,
		expected.Format("200601021504") == got.Format("200601021504"),
	)
}
