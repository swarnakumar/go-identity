package cookie

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/swarnakumar/go-identity/config"
	"net/http"
)

type Cookie struct {
	cookie *securecookie.SecureCookie
}

func New(cookieHashKey, cookieBlockKey string) (*Cookie, error) {
	if len(cookieHashKey) != 32 && len(cookieHashKey) != 64 {
		return nil, fmt.Errorf("CookieHashKey must be either 32 or 64 characters long. Passed Key Length: %d", len(cookieHashKey))
	}

	if len(cookieBlockKey) != 32 && len(cookieBlockKey) != 64 {
		return nil, fmt.Errorf("CookieBlockKey must be either 32 or 64 characters long. Passed Key Length: %d", len(cookieBlockKey))
	}

	cookie := securecookie.New([]byte(config.CookieHashKey), []byte(config.CookieHashKey))
	return &Cookie{cookie: cookie}, nil
}

func (s *Cookie) Set(w http.ResponseWriter, cookieName string, value string) error {
	encoded, err := s.cookie.Encode(cookieName, value)
	if err == nil {
		cookie := &http.Cookie{
			Name:  cookieName,
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)

	}

	return err
}

func (s *Cookie) Read(r *http.Request, cookieName string) (interface{}, error) {

	if cookie, err := r.Cookie(cookieName); err == nil {
		var value string
		err = s.cookie.Decode(cookieName, cookie.Value, &value)
		return value, err

	} else {
		return nil, err
	}
}

func (s *Cookie) Delete(w http.ResponseWriter, cookieName string) {
	encoded, _ := s.cookie.Encode(cookieName, "")
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    encoded,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   0,
	}
	http.SetCookie(w, cookie)
}
