package users

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	nickname := "exampleuser"
	now := time.Now()
	exampleSignedWith := []ExternalSigned{{Platform: "examplePlatform"}}
	type args struct {
		userR User
	}
	tests := []struct {
		name     string
		args     args
		wantUser User
		wantErr  error
	}{
		{
			name: "Given a valid user with a external sign When creating new user Then same user with User.CreatedAt field filled",
			args: args{
				userR: User{
					Nickname:   nickname,
					SignedWith: exampleSignedWith,
				},
			},
			wantUser: User{
				Nickname:   nickname,
				SignedWith: exampleSignedWith,
				CreatedAt:  now,
			},
		},
		{
			name: "Given a user with a invalid username When creating new user Then invalid nickname error",
			args: args{
				userR: User{
					Nickname: "$**exampleNickname",
				},
			},
			wantErr: errors.New("invalid nickname: invalid nickname format of $**exampleNickname"),
		},
		{
			name: "Given a user with a empty password When creating new user Then invalid pasword error",
			args: args{
				userR: User{
					Nickname: nickname,
				},
			},
			wantErr: errors.New("invalid password: empty value"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, err := New(tt.args.userR)
			if err != nil && tt.wantErr == nil {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil && tt.wantErr != nil {
				assert.EqualError(t, tt.wantErr, err.Error())
				return
			}

			if equalDateTime(t, tt.wantUser.CreatedAt, gotUser.CreatedAt) {
				tt.wantUser.CreatedAt = time.Time{}
				gotUser.CreatedAt = time.Time{}
			}
			assert.Equal(t, tt.wantUser, gotUser)
		})
	}
}

func equalDateTime(t *testing.T, expected, got time.Time) bool {
	t.Helper()

	return assert.True(
		t,
		expected.Format("200601021504") == got.Format("200601021504"),
	)
}
