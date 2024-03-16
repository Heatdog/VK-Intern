package migrations

import (
	"context"

	user_model "github.com/Heater_dog/Vk_Intern/internal/models/user"
	"github.com/Heater_dog/Vk_Intern/pkg/client"
	cryptohash "github.com/Heater_dog/Vk_Intern/pkg/cryptoHash"
)

func InitDb(client client.Client) error {
	users := []struct {
		User     user_model.User
		Password string
	}{
		{
			User: user_model.User{
				Login: "Admin",
				Role:  "Admin",
			},
			Password: "Admin",
		},
		{
			User: user_model.User{
				Login: "VK Senior",
				Role:  "Admin",
			},
			Password: "678",
		},
		{
			User: user_model.User{
				Login: "John",
				Role:  "User",
			},
			Password: "123",
		},
		{
			User: user_model.User{
				Login: "David",
				Role:  "User",
			},
			Password: "123456",
		},
		{
			User: user_model.User{
				Login: "Peter",
				Role:  "User",
			},
			Password: "890",
		},
	}

	q := `
			INSERT INTO Users
				(login, password, role) 
			VALUES 
				($1, $2, $3)
	`

	for _, el := range users {
		pswd, err := cryptohash.Hash(el.Password)
		if err != nil {
			return err
		}
		client.QueryRow(context.Background(), q, el.User.Login, string(pswd), el.User.Role)
	}
	return nil

}
