package api

import (
	"context"
	"fmt"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/types"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
)

func user(server mattrax.Server, builder *schemabuilder.Schema) {
	userObject := builder.Object("User", types.User{})
	userObject.Description = "A user."
	userObject.FieldFunc("password", func(args struct{ password string }) string {
		return ""
	})

	query := builder.Query()
	query.FieldFunc("users", server.UserService.GetAll)
	query.FieldFunc("user", func(ctx context.Context) (types.User, error) {
		fmt.Println(ctx)
		// TODO
		return types.User{}, nil
	})

	mutation := builder.Mutation()
	mutation.FieldFunc("createUser", func(req struct {
		DisplayName string
		Email       string
		Password    string
		Permissions types.Permissions `graphql:",optional"`
	}) types.User {
		hashedPassword, err := server.UserService.HashPassword([]byte(req.Password))
		if err != nil {
			panic(err) // TODO
		}
		user := types.User{
			DisplayName: req.DisplayName,
			Email:       req.Email,
			Password:    hashedPassword,
			Permissions: req.Permissions,
		}
		server.UserService.CreateOrEdit(user.Email, user)
		user.Password = []byte{}
		return user
	})
	// mutation.FieldFunc("editUser", func(args struct {
	// 	email string
	// 	user  types.User
	// }) types.User {
	// 	server.UserService.CreateOrEdit(args.email, args.user)
	// 	args.user.Email = args.email
	// 	return args.user
	// })
}
