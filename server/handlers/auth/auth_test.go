package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/coffemanfp/chat/auth"
	"github.com/coffemanfp/chat/config"
	"github.com/coffemanfp/chat/database"
	"github.com/coffemanfp/chat/server/handlers"
	"github.com/coffemanfp/chat/users"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

type authRepositoryImpl struct {
	m         sync.Mutex
	idCounter int
	users     map[string]users.User
	session   map[string]auth.Session
}

var authRepository database.AuthRepository = &authRepositoryImpl{
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

var conf config.ConfigInfo = config.ConfigInfo{
	OAuth: config.OAuth{
		Google: config.OAuthProperties{
			Endpoint:     google.Endpoint,
			RedirectURIS: make([]string, 1),
		},
		Facebook: config.OAuthProperties{
			Endpoint:     facebook.Endpoint,
			RedirectURIS: make([]string, 1),
		},
	},
}

var fbHandler = newFacebookHandler(conf)
var gHandler = newGoogleHandler(conf)

func TestAuthHandler_handleExternalSign(t *testing.T) {
	type fields struct {
		config               config.ConfigInfo
		repository           database.AuthRepository
		writer               handlers.ResponseWriter
		reader               handlers.RequestReader
		userReaders          map[handlerName]userReader
		externalSignHandlers map[handlerName]externalSignUpHandler
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantCode        int
		wantURLRedirect string
	}{
		{
			name: "Given When Then",
			fields: fields{
				config:     conf,
				repository: authRepository,
				writer:     handlers.GetResponseWriterImpl(),
				reader:     handlers.GetRequestReaderImpl(),
				userReaders: map[handlerName]userReader{
					systemHandlerName: systemUserReader{
						reader: handlers.GetRequestReaderImpl(),
						writer: handlers.GetResponseWriterImpl(),
					},
					googleHandlerName:   gHandler,
					facebookHandlerName: fbHandler,
				},
				externalSignHandlers: map[handlerName]externalSignUpHandler{
					googleHandlerName:   gHandler,
					facebookHandlerName: fbHandler,
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/", nil),
			},
			wantCode: http.StatusTemporaryRedirect,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AuthHandler{
				config:               tt.fields.config,
				repository:           tt.fields.repository,
				writer:               tt.fields.writer,
				reader:               tt.fields.reader,
				userReaders:          tt.fields.userReaders,
				externalSignHandlers: tt.fields.externalSignHandlers,
			}
			a.handleExternalSign(tt.args.w, tt.args.r)
			rec := tt.args.w.(*httptest.ResponseRecorder)
			fmt.Println(rec.Code)
		})
	}
}
