package sqlstore_test

import (
	"dialogue/internal/models"
	"dialogue/internal/store/sqlstore"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")
	s := sqlstore.New(db)
	u := models.TestUser(t)
	db.AutoMigrate(&u)

	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_Get(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")
	s := sqlstore.New(db)
	u := models.TestUser(t)
	db.AutoMigrate(&u)

	_, err := s.User().FindByID(1)
	assert.Error(t, err)

	s.User().Create(u)
	_, err = s.User().FindByID(u.ID)
	assert.NoError(t, err)
}
