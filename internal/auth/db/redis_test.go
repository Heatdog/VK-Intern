package authDb

import (
	"testing"

	"github.com/go-redis/redismock/v9"
)

func TestSetToken(t *testing.T) {
	db, mock := redismock.NewClientMock()

}
