package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"ontopsolutions.net/gasperlf/social/internal/mailer"
	"ontopsolutions.net/gasperlf/social/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,max=255"`
	Password string `json:"passsword" validate:"required,min=3,max=72"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

// RegisterUser godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var request RegisterUserPayload

	if err := readJSON(w, r, &request); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(request); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &store.User{
		Username: request.Username,
		Email:    request.Email,
	}

	//hash password
	if err := user.Password.Set(request.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	ctx := r.Context()
	token := uuid.NewString()

	//store
	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])
	err := app.store.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.exp)

	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestResponse(w, r, err)
		case store.ErrDuplicateUsername:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	userWithToken := UserWithToken{
		User:  user,
		Token: token,
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, userWithToken.Token)
	isProdEnv := app.config.env == "prod"
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}
	//send mail
	status, err := app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		app.logger.Errorw("error sending welcome email", "error", err.Error())
		//rollback user creation if email fails
		if err := app.store.Users.Delete(ctx, user.ID); err != nil {
			app.logger.Errorw("Error deleting user", "error", err)
		}
		app.internalServerError(w, r, err)
		return
	}
	app.logger.Infow("Email sent with status: ", status)

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
