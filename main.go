package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
)

var (
	templateRender = render.New(render.Options{
		Layout:        "layout",
		IsDevelopment: true,
	})
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	templateParams := map[string]interface{}{}
	u, err := isLogin(getAuth(r))
	templateParams["user"] = u
	if nil != err {
		templateRender.HTML(w, http.StatusOK, "welcome", templateParams)
	} else {
		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

func Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	templateParams := map[string]interface{}{}
	u, err := isLogin(getAuth(r))
	if nil != err {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	templateParams["user"] = u

	var start int64
	if "" == r.FormValue("start") {
		start = int64(0)
	} else {
		start, err = strconv.ParseInt(r.FormValue("start"), 10, 64)
		if err != nil {
			start = int64(0)
		}
	}
	posts, rest, err := getUserPosts("posts:"+u.Id, start, 10)
	if err == nil {
		if start > 0 {
			templateParams["prev"] = start - 10
		}
		templateParams["posts"] = posts
		if rest > 0 {
			templateParams["next"] = start + 10
		}
	}
	templateRender.HTML(w, http.StatusOK, "home", templateParams)
}

func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	password2 := r.PostFormValue("password2")
	if username == "" || password == "" || password2 == "" {
		GoBack(w, r, errors.New("Every field of the registration form is needed!"))
		return
	}
	if password != password2 {
		GoBack(w, r, errors.New("The two password fileds don't match!"))
		return
	}
	auth, err := register(username, password)
	if err != nil {
		GoBack(w, r, err)
		return
	}
	setSession(auth, w)
	templateParams := map[string]interface{}{}
	templateParams["username"] = username
	templateRender.HTML(w, http.StatusOK, "register", templateParams)
}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	if username == "" || password == "" {
		GoBack(w, r, errors.New("You need to enter both username and password to login."))
		return
	}
	auth, err := login(username, password)
	if err != nil {
		GoBack(w, r, err)
		return
	}
	setSession(auth, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u, err := isLogin(getAuth(r))
	if nil != err {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	logout(u)
	http.Redirect(w, r, "/", http.StatusFound)
}

func Publish(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u, err := isLogin(getAuth(r))
	if nil != err {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	status := r.PostFormValue("status")
	if status == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	err = post(u, status)
	http.Redirect(w, r, "/", http.StatusFound)
}

func Timeline(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	templateParams := map[string]interface{}{}
	users, err := getLastUsers()
	if err != nil {
		log.Println(err)
	} else {
		templateParams["users"] = users
	}
	posts, _, err := getUserPosts("timeline", 0, 50)
	if err != nil {
		log.Println(err)
	} else {
		templateParams["posts"] = posts
	}
	templateRender.HTML(w, http.StatusOK, "timeline", templateParams)
}

func Profile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	templateParams := map[string]interface{}{}
	// get username
	username := r.FormValue("u")
	if username == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// Get profile
	p, err := profileByUsername(username)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	templateParams["profile"] = p
	// Get logged in user
	u, err := isLogin(getAuth(r))
	if nil == err {
		templateParams["user"] = u
	}

	var start int64
	if "" == r.FormValue("start") {
		start = int64(0)
	} else {
		start, err = strconv.ParseInt(r.FormValue("start"), 10, 64)
		if err != nil {
			start = int64(0)
		}
	}
	posts, rest, err := getUserPosts("posts:"+p.Id, start, 10)
	if err == nil {
		if start > 0 {
			templateParams["prev"] = start - 10
		}
		templateParams["posts"] = posts
		if rest > 0 {
			templateParams["next"] = start + 10
		}
	}
	templateRender.HTML(w, http.StatusOK, "profile", templateParams)
}

func Follow(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// get the user id
	userId := r.FormValue("uid")
	if userId == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// Get the action to do
	doFollow := false
	switch r.FormValue("f") {
	case "1":
		doFollow = true
	case "0":
		doFollow = false
	default:
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	u, err := isLogin(getAuth(r))
	if nil != err {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if userId == u.Id {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if doFollow {
		u.Follow(&User{Id: userId})
	} else {
		u.Unfollow(&User{Id: userId})
	}
	p, err := profileByUserId(userId)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/profile?u="+p.Username, http.StatusFound)
}

func GoBack(w http.ResponseWriter, r *http.Request, err error) {
	templateParams := map[string]interface{}{}
	templateParams["error"] = err
	templateRender.HTML(w, http.StatusOK, "error", templateParams)
}

func main() {

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/home", Home)
	router.POST("/register", Register)
	router.POST("/login", Login)
	router.GET("/logout", Logout)
	router.POST("/post", Publish)
	router.GET("/timeline", Timeline)
	router.GET("/profile", Profile)
	router.GET("/follow", Follow)

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":8080")
}
