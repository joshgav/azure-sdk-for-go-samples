package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/subosito/gotenv"
	"golang.org/x/net/html"
)

const (
	envvarName  = "COOKIE_KEY"
	sessionName = "auth_sample"

	authenticatedKey = "authenticated"
	stateKey         = "state"
	nameKey          = "name"
	emailKey         = "email"
)

var store sessions.Store

func init() {
	gotenv.Load()
	key := os.Getenv(envvarName)
	if len(key) == 0 {
		log.Printf("Session (init): use envvar %v for cookie key\n", envvarName)
		key = "makemerandom"
	}

	log.Printf("Session (init): creating new cookie store\n")
	store = sessions.NewCookieStore([]byte(key))
}

// WithSession decorates an http.Handler to create a session and populates
// context.Context with info about the session.
// To access this info in a later handler:
//   `authenticated, ok := r.Context().Value(authenticatedKey).(bool)`
//   `state, ok := r.Context().Value(stateKey).(string)`
func WithSession(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("WithSession: getting session %v\n", sessionName)
		s, err := store.Get(r, sessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if _, ok := s.Values[stateKey].(string); ok == false {
			log.Printf("WithSession: no state in session, adding it\n")

			state := "notrandomatthemoment"
			s.Values[stateKey] = state

			log.Printf("WithSession: state set to [%s]\n", state)
		}

		if _, ok := s.Values[authenticatedKey].(bool); ok == false {
			log.Printf("WithSession: user not previously authenticated\n")
			s.Values[authenticatedKey] = false
		}

		log.Printf("WithSession: saving session [%v]\n", s)
		_ = s.Save(r, w)

		log.Printf("WithSession: adding session data to context for ensuing modules\n")
		var ctx context.Context
		ctx = context.WithValue(ctx, authenticatedKey, s.Values[authenticatedKey])
		ctx = context.WithValue(ctx, stateKey, s.Values[stateKey])
		ctx = context.WithValue(ctx, nameKey, s.Values[nameKey])
		ctx = context.WithValue(ctx, emailKey, s.Values[emailKey])

		log.Printf("WithSession: done, calling next with updated context\n", ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// saveSession saves info into the current session.
func SaveSession(info map[string]string, w http.ResponseWriter, r *http.Request) (*http.Request, error) {
	log.Printf("SetSession: getting session %v\n", sessionName)
	s, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, "failed to get session", http.StatusInternalServerError)
		return nil, err
	}

	s.Values[authenticatedKey] = true
	for k, v := range info {
		log.Printf("SaveSession: saving key [%s] value [%s]\n", k, v)
		s.Values[k] = v
		r = r.WithContext(context.WithValue(r.Context(), k, v))
	}
	err = s.Save(r, w)
	return r, err
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("UserInfoHandler: request received [%+v]\n", r)
	name := r.Context().Value("name").(string)
	fmt.Fprintf(w, "Hello %s!", html.EscapeString(name))
}
