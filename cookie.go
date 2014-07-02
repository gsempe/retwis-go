package main

import (
	"net/http"

	"github.com/gorilla/securecookie"
)

var (
	cookieHandler = securecookie.New(
		[]byte("in0y(>'@N+#N6A5*=iL%lM}[U`|AH#8ltj@02>e9gwsU&Wu'JNuhCRFPri7Z{*H1"),
		[]byte("zy-hA<!J(oXKp]sy]HoJkiuqK_R&8EnB"))
)

const (
	cookieName = "auth"
	cookieKey  = "auth"
)

func getAuth(r *http.Request) (auth string) {
	if cookie, err := r.Cookie(cookieName); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode(cookieName, cookie.Value, &cookieValue); err == nil {
			auth = cookieValue[cookieKey]
		}
	}
	return auth
}

func setSession(auth string, w http.ResponseWriter) {
	value := map[string]string{
		cookieKey: auth,
	}
	if encoded, err := cookieHandler.Encode(cookieName, value); err == nil {
		cookie := &http.Cookie{
			Name:  cookieName,
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func clearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   cookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
