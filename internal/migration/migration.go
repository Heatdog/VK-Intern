package migrations

import (
	"context"

	"github.com/Heater_dog/Vk_Intern/internal/user"
	"github.com/Heater_dog/Vk_Intern/pkg/client"
	cryptohash "github.com/Heater_dog/Vk_Intern/pkg/cryptoHash"
)

func InitDb(client client.Client) error {
	users := []struct {
		User     user.User
		Password string
	}{
		{
			User: user.User{
				Login: "Admin",
				Role:  "Admin",
			},
			Password: "Admin",
		},
		{
			User: user.User{
				Login: "VK Senior",
				Role:  "Admin",
			},
			Password: "678",
		},
		{
			User: user.User{
				Login: "John",
				Role:  "User",
			},
			Password: "123",
		},
		{
			User: user.User{
				Login: "David",
				Role:  "User",
			},
			Password: "123456",
		},
		{
			User: user.User{
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
