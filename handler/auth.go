package handler

import (
	"io/ioutil"
	"net/http"
)

func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			username := r.Form.Get("username")
			token := r.Form.Get("token")
			if len(username) < 3 || !IsTokenValid(token) {
				data, err := ioutil.ReadFile("./static/view/signin.html")
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Write(data)
				return
			}
			h(w, r)
		})
}
