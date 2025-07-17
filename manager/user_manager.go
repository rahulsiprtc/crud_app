package manager

import (
	"crud-app/mongoDatabase"
	"crud-app/request"
	"crud-app/response"
)

type UserManager struct {
	mongodbQuery *mongoDatabase.Queryis
}

func (m *UserManager) CreateUser(req request.CreateUserRequest) (response.UserResponse, error) {
	return m.mongodbQuery.InsertUser(req)
}

func (m *UserManager) GetAllUsers(req request.PaginationRequest) (response.UserPaginationResponse, error) {
	return m.mongodbQuery.GetAllUsers(req)
}

func (m *UserManager) UpdateUser(id string, req request.UpdateUserRequest) (string, error) {
	return m.mongodbQuery.UpdateUser(id, req)
}

func (m *UserManager) DeleteUser(id string) error {
	return m.mongodbQuery.DeleteUser(id)
}
