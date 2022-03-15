package auth

import (
	"errors"
	"sync"

	"github.com/coffemanfp/chat/auth"
	"github.com/coffemanfp/chat/database"
	"github.com/coffemanfp/chat/users"
)

type authRepositoryImpl struct {
	m         sync.Mutex
	idCounter int
	users     map[string]users.User
	session   map[string]auth.Session
}

var _ database.AuthRepository = &authRepositoryImpl{
	m:       sync.Mutex{},
	users:   make(map[string]users.User),
	session: make(map[string]auth.Session),
}

func (a *authRepositoryImpl) SignUp(user users.User, session auth.Session) (id int, err error) {
	a.m.Lock()
	a.idCounter++
	if _, ok := a.users[user.Nickname]; !ok {
		a.users[user.Nickname] = user
		a.session[session.ID] = session
	} else {
		err = errors.New("already exists nickname")
	}
	a.m.Unlock()
	return
}

func (a *authRepositoryImpl) MatchCredentials(user users.User) (id int, err error) {
	a.m.Lock()
	if u, ok := a.users[user.Nickname]; ok {
		if u.Password == user.Password {
			id = u.ID
			return
		}
	} else {
		err = errors.New("not found: user don't exists")
	}
	a.m.Unlock()
	return
}

func (a *authRepositoryImpl) UpsertSession(session auth.Session) (err error) {
	a.m.Lock()
	if s, ok := a.session[session.ID]; ok {
		s.LastSeenAt = session.LastSeenAt
		s.TmpID = ""
		a.session[session.ID] = s
	} else {
		a.session[session.ID] = session
	}
	a.m.Unlock()
	return
}
