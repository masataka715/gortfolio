package flash

import (
	"encoding/base64"
	"net/http"
	"time"
)

func Set(w http.ResponseWriter, name string, value []byte) {
	encodedValue := base64.URLEncoding.EncodeToString(value)
	c := &http.Cookie{
		Name:  name,
		Value: encodedValue,
		Path:  "/",
	}
	http.SetCookie(w, c)
}

func Get(w http.ResponseWriter, r *http.Request, name string) (string, error) {
	c, err := r.Cookie(name)
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return "", nil
		default:
			return "", err
		}
	}
	value, err := base64.URLEncoding.DecodeString(c.Value)
	if err != nil {
		return "", err
	}
	dc := &http.Cookie{
		Name:    name,
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Unix(1, 0),
	}
	http.SetCookie(w, dc)
	return string(value), nil
}
