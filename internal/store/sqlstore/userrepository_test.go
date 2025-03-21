package sqlstore_test

import (
	"dialogue/internal/models"
	"dialogue/internal/store/sqlstore"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Insert(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")
	s := sqlstore.New(db)
	u := models.TestUser(t)
	db.AutoMigrate(&u)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}
