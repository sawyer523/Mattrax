package api

import (
	"context"
	"errors"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/types"
	"github.com/rs/zerolog/log"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
)

// User contains the API endpoints and types for Users
func User(server *mattrax.Server, builder *schemabuilder.Schema) {
	userObject := builder.Object("User", types.User{})
	userObject.Description = "An end user or administrator of the Mattrax system."
	userObject.FieldFunc("password", func(args struct{ password string }) string {
		return ""
	})

	query := builder.Query()
	query.FieldFunc("users", server.UserService.GetAll)
	query.FieldFunc("user", func(ctx context.Context, req struct{ Email string }) (types.User, error) {
		if user, err := server.UserService.Get(req.Email); err != nil && err == types.ErrUserNotFound {
			return types.User{}, nil
		} else if err != nil {
			log.Error().Str("email", req.Email).Err(err).Msg("Error retieving user")
			return types.User{}, errors.New("error retrieving user")
		} else {
			return user, nil
		}
	})
	// TODO: after auth is added
	// query.FieldFunc("me", func(ctx context.Context) (types.User, error) {
	// 	return types.User{}, nil
	// })

	mutation := builder.Mutation()
	mutation.FieldFunc("createUser", func(user types.User) (types.User, error) {
		if err := user.Verify(); err != nil {
			return types.User{}, err
		}

		var err error
		if user.Password, err = server.UserService.HashPassword([]byte(user.Password)); err != nil {
			return types.User{}, err
		}

		if err = server.UserService.CreateOrEdit(user.Email, user); err != nil {
			return types.User{}, errors.New("error creating user")
		}

		user.Password = ""
		return user, nil
	})
	mutation.FieldFunc("editUser", func(req struct {
		NewEmail string
		User     types.User
	}) (types.User, error) {
		if err := req.User.Verify(); err != nil {
			return types.User{}, err
		}

		var err error
		if req.User.Password, err = server.UserService.HashPassword([]byte(req.User.Password)); err != nil {
			return types.User{}, err
		}

		if err = server.UserService.CreateOrEdit(req.NewEmail, req.User); err != nil {
			return types.User{}, errors.New("error updating user")
		}

		req.User.Password = ""
		return req.User, nil
	})
}
