package controller

import (
	"net/http"

	"github.com/geobe/gostip/go/model"

	"fmt"
	"html"
	"github.com/geobe/gostip/go/view"
	"github.com/justinas/nosurf"
)

func PasswordChange(w http.ResponseWriter, r *http.Request)  {
	login := r.FormValue("login")
	password := r.FormValue("password")
	db := model.Db()
	var user model.User

	db.First(&user, "login = ?",login)

	user.UpdatedBy = fmt.Sprintf("%s(%s)", user.Fullname, user.Login)
	user.Updater = user.ID
	user.Password = model.Encrypt(password+user.Login)
	db.Save(&user)

}

func UserFullnameChange(w http.ResponseWriter, r *http.Request)  {
	fullname := r.FormValue("fullname")
	login := r.FormValue("login")

	db := model.Db()
	var user model.User

	db.First(&user,"login = ?",login)

	user.UpdatedBy = fmt.Sprintf("%s(%s)", user.Fullname, user.Login)
	user.Updater = user.ID
	user.Fullname = fullname

	db.Save(&user)

	w.Write([]byte(fullname))
}

func FindUser(w http.ResponseWriter,r *http.Request)  {
	login := html.EscapeString(r.PostFormValue("lastname"))
	fullname := html.EscapeString(r.PostFormValue("firstname"))
	//action := html.EscapeString(r.PostFormValue("action"))
	//flag := html.EscapeString(r.PostFormValue("flag"))
	fmt.Println(login,fullname)
	db := model.Db()
	var users []model.User

	db.Where("users.deleted_at IS NULL").
		Where("users.login like ?", login + "%").
		Where("users.fullname like ?", fullname + "%").Find(&users)

	var res []map[string]string = make([]map[string]string, len(users))

	for i,value := range users{
		res[i] = map[string]string{
			"id":        fmt.Sprintf("%d", value.ID),
			"lastname":  value.Login,
			"firstname": value.Fullname,
		}
	}
	view.Views().ExecuteTemplate(w, "qresult", res)
}

func ShowUser(w http.ResponseWriter, r *http.Request)  {

	i18nlanguage := view.PreferedLanguages(r) [0]
	values := viewmodel{
		"csrftoken": nosurf.Token(r),
		"csrfid": "csrf_id_enrol",
		"language":     i18nlanguage,
	}
	var user model.User
	userid := atoint(html.EscapeString(r.PostFormValue("appid")))
	db := model.Db()
	db.First(&user,userid)
	setViewModelUser(user,values)
	view.Views().ExecuteTemplate(w, "work_edituser", values)
}

func setViewModelUser(user model.User, vmod viewmodel)  {
	vmod["userid"] = user.ID
	vmod["login"] = user.Login
	vmod["fullname"] = user.Fullname
	vmod["password"] = user.Password
	vmod["role"] = user.Role
}

func SubmitUser(w http.ResponseWriter, r *http.Request)  {
	login := r.FormValue("login")
	fullname := r.FormValue("fullname")
	password := r.FormValue("password")
	role := r.FormValue("role")

	var user model.User

	db := model.Db()
	db.First(&user,"login = ?",login)
	if user.ID>0 {
		user.Fullname = fullname
		user.Password = model.Encrypt(password + login)
		user.Authorize(atoint(role))
		user.Role = atoint(role)
	}else{
		user = *model.NewUser(login,fullname,password,atoint(role))
	}
	db.Save(&user)
}