package manager

import (
	"crud-app/models"
	"crud-app/mongoDatabase"
	"crud-app/request"
)

type UserManager struct {
	mongodbQuery *mongoDatabase.Queryis
}

func (m *UserManager) CreateUser(req request.CreateUserRequest) (*models.User, error) {
	return m.mongodbQuery.InsertUser(req)
}

func (m *UserManager) GetAllUsers(page, limit, minAge int64, nameContains string) ([]models.User, int64, error) {
	return m.mongodbQuery.GetAllUsers(page, limit, minAge, nameContains)
}

func (m *UserManager) GetUserByID(id string) (*models.User, error) {
	return m.mongodbQuery.GetUserByID(id)
}

func (m *UserManager) UpdateUser(id string, req request.UpdateUserRequest) (string, error) {
	return m.mongodbQuery.UpdateUser(id, req)
}

func (m *UserManager) DeleteUser(id string) error {
	return mongoDatabase.SoftDeleteUser(id)
}
