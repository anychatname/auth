package auth

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/coffemanfp/chat/auth"
	"github.com/coffemanfp/chat/database"
	"github.com/coffemanfp/chat/users"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type authRepositoryImpl struct {
	m             sync.Mutex
	userSerial    int
	sessionSerial int
	users         map[string]users.User
	session       map[int]auth.Session
}

func newAuthRepositoryImpl() authRepositoryImpl {
	return authRepositoryImpl{
		m:       sync.Mutex{},
		users:   map[string]users.User{},
		session: map[int]auth.Session{},
	}
}

// This statement is to check if the authRepositoryImpl mock is doing well with the database.AuthRepository interface.
var _ database.AuthRepository = &authRepositoryImpl{}

func (a *authRepositoryImpl) SignUp(user users.User, session auth.Session) (id int, err error) {
	a.m.Lock()
	if _, ok := a.users[user.Nickname]; !ok {
		// Increase serial
		a.userSerial++
		user.ID = a.userSerial

		// Save user
		a.users[user.Nickname] = user
		a.session[session.ID] = session
		id = user.ID
	} else {
		err = errors.New("already exists nickname")
	}
	a.m.Unlock()
	return
}

func (a *authRepositoryImpl) GetPasswordHash(user users.User) (id int, pass string, err error) {
	a.m.Lock()
	if u, ok := a.users[user.Nickname]; ok {
		id = u.ID
		pass = u.Password
	} else {
		err = errors.New("not found: user don't exists")
	}
	a.m.Unlock()
	return
}

func (a *authRepositoryImpl) SaveSudo(sudo auth.Sudo) (err error) {
	// TODO
	return
}

func (a *authRepositoryImpl) UpsertSession(session auth.Session) (id int, err error) {
	a.m.Lock()
	if s, ok := a.session[session.ID]; ok {
		s.LastSeenAt = session.LastSeenAt
		a.session[session.ID] = s
	} else {
		a.sessionSerial++
		session.ID = a.sessionSerial
		a.session[session.ID] = session
	}
	id = session.ID
	a.m.Unlock()
	return
}

var (
	now  = time.Now()
	user = users.User{
		Nickname:  "example",
		Email:     "example@host.com",
		Password:  "1234",
		Picture:   "http://host.example/image.png",
		CreatedAt: now,
	}
	facebookSign = users.ExternalSigned{
		ID:        "laksdjakljsdf",
		Email:     "example@host.com",
		Picture:   "http://facebook.com/image.png",
		Platform:  "facebook",
		CreatedAt: now,
	}
	googleSign = users.ExternalSigned{
		ID:        "laksdjakljsdf",
		Email:     "example@host.com",
		Picture:   "http://google.com/image.png",
		Platform:  "google",
		CreatedAt: now,
	}
	session = auth.Session{
		LoggedAt:   now,
		LastSeenAt: now,
		LoggedWith: "system",
		Actived:    true,
	}
	sudo = auth.Sudo{
		DurationInSecs: 900,
		CreatedAt:      now,
	}
)

func TestAuthRepoSignUp(t *testing.T) {
	authRepo := newAuthRepositoryImpl()

	t.Run("Given complete user info When sign up Then success", func(t *testing.T) {
		userExp := newExpectedUser(t, user)

		id, err := authRepo.SignUp(userExp, session)
		assert.NoError(t, err)
		assert.Equal(t, authRepo.userSerial, id)

		userExp.ID = id

		userGot := authRepo.users[userExp.Nickname]
		assert.Equal(t, userExp, userGot)
	})
	t.Run("Given already exists user When sign up Then error", func(t *testing.T) {
		userExp := newExpectedUser(t, user)
		authRepo.users[userExp.Nickname] = userExp

		id, err := authRepo.SignUp(userExp, session)
		assert.Empty(t, id)
		assert.EqualError(t, err, "already exists nickname")
	})
}

func TestGetPasswordHash(t *testing.T) {
	authRepo := newAuthRepositoryImpl()

	t.Run("Given a existent user When getting password hash and id Then success", func(t *testing.T) {
		userExp := newExpectedUser(t, user)
		userExp.ID = authRepo.userSerial + 1
		authRepo.users[userExp.Nickname] = userExp

		id, pass, err := authRepo.GetPasswordHash(userExp)
		assert.NoError(t, err)
		assert.Equal(t, userExp.ID, id)
		assert.Equal(t, userExp.Password, pass)
	})

	t.Run("Given a non-existent user When getting password hash and id Then error", func(t *testing.T) {
		id, pass, err := authRepo.GetPasswordHash(user)
		assert.EqualError(t, err, "not found: user don't exists")
		assert.Empty(t, id)
		assert.Empty(t, pass)
	})
}

func TestUpsertSession(t *testing.T) {
	authRepo := newAuthRepositoryImpl()

	t.Run("Given a existent session When updating or inserting session Then session updated", func(t *testing.T) {
		sessionExp := session
		sessionExp.ID = authRepo.sessionSerial + 1
		authRepo.session[sessionExp.ID] = sessionExp

		id, err := authRepo.UpsertSession(sessionExp)
		assert.NoError(t, err)
		assert.Equal(t, sessionExp.ID, id)

		sessionGot := authRepo.session[sessionExp.ID]
		assert.Equal(t, sessionExp, sessionGot)
	})

	t.Run("Given a non-existent session When updateing or inserting a session Then session inserted", func(t *testing.T) {
		sessionExp := session
		prevSessionSerial := authRepo.sessionSerial

		id, err := authRepo.UpsertSession(sessionExp)
		assert.NoError(t, err)
		assert.Equal(t, prevSessionSerial+1, id)

		sessionExp.ID = id
		sessionGot := authRepo.session[id]

		assert.Equal(t, sessionExp, sessionGot)
	})
}

func newExpectedUser(t *testing.T, user users.User) (r users.User) {
	t.Helper()

	r = user
	r.Nickname += uuid.NewString()
	r.Email = uuid.NewString() + r.Email
	return
}
