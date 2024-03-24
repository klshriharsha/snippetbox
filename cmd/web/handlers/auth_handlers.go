package handlers

import (
	"errors"
	"net/http"

	"github.com/klshriharsha/snippetbox/cmd/web/config"
	"github.com/klshriharsha/snippetbox/internal/models"
	"github.com/klshriharsha/snippetbox/internal/validator"
)

type signupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// SignupHandler renders a signup form
func SignupHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := app.NewTemplateData(r)
		data.Form = signupForm{}

		app.RenderPage(w, http.StatusOK, "signup.go.tmpl", data)
	}
}

// SignupPostHandler receives a POST request and creates a new request in the database
func SignupPostHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := signupForm{Validator: validator.Validator{FieldErrors: make(map[string]string)}}
		if err := app.DecodePostForm(r, &form); err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		// validate all the fields
		form.CheckField(validator.NotBlank(form.Name), "name", "Name cannot be empty")
		form.CheckField(validator.NotBlank(form.Email), "email", "Email cannot be empty")
		form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Email must be a valid email address")
		form.CheckField(validator.NotBlank(form.Password), "password", "Password cannot be empty")
		form.CheckField(validator.MinChars(form.Password, 8), "password", "Password must be at least 8 characters long")
		if !form.Valid() {
			// re-render the form with validation errors
			data := app.NewTemplateData(r)
			data.Form = form
			app.RenderPage(w, http.StatusUnprocessableEntity, "signup.go.tmpl", data)
			return
		}

		if err := app.Users.Insert(form.Name, form.Email, form.Password); err != nil {
			// re-render the form with duplicate email error
			if errors.Is(err, models.ErrDuplicateEmail) {
				data := app.NewTemplateData(r)
				form.AddFieldError("email", "Email address is already in use")
				data.Form = form
				app.RenderPage(w, http.StatusUnprocessableEntity, "signup.go.tmpl", data)
				return
			}

			app.ServerError(w, err)
			return
		}

		// set a flash message and redirect to login
		app.SessionManager.Put(r.Context(), "Flash", "Your signup was successful. Please login.")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}
}

type loginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func LoginHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := app.NewTemplateData(r)
		data.Form = loginForm{}

		app.RenderPage(w, http.StatusOK, "login.go.tmpl", data)
	}
}

func LoginPostHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := loginForm{}
		if err := app.DecodePostForm(r, &form); err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		// validate the login form data
		form.CheckField(validator.NotBlank(form.Email), "email", "Email cannot be empty")
		form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Email must be a valid email address")
		form.CheckField(validator.NotBlank(form.Password), "password", "Password cannot be empty")
		if !form.Valid() {
			data := app.NewTemplateData(r)
			data.Form = form

			app.RenderPage(w, http.StatusUnprocessableEntity, "login.go.tmpl", data)
			return
		}

		// attempt to authenticate with the given credentials
		id, err := app.Users.Authenticate(form.Email, form.Password)
		if err != nil {
			if errors.Is(err, models.ErrInvalidCredentials) {
				// re-render the page with an invalid credentials error
				form.AddNonFieldError("Email or Password is incorrect")
				data := app.NewTemplateData(r)
				data.Form = form

				app.RenderPage(w, http.StatusUnauthorized, "login.go.tmpl", data)
				return
			}

			app.ServerError(w, err)
			return
		}

		// upon successful authentication, renew the user's session token as a good security measure so that the
		// session id changes. Also, set the user's id in session context for identification in future requests
		if err := app.SessionManager.RenewToken(r.Context()); err != nil {
			app.ServerError(w, err)
			return
		}
		app.SessionManager.Put(r.Context(), "authenticatedUserID", id)

		http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
	}
}

func LogoutPostHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := app.SessionManager.RenewToken(r.Context()); err != nil {
			app.ServerError(w, err)
			return
		}

		app.SessionManager.Remove(r.Context(), "authenticatedUserID")
		app.SessionManager.Put(r.Context(), "flash", "You've been logged out successfully")

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
