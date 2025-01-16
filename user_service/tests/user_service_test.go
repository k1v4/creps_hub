package tests

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	userv1 "github.com/k1v4/protos/gen/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"user_service/tests/suite"
)

func TestUserService_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	t.Logf("testing Add User Service")
	respAdd, err := st.UserClient.AddUser(ctx, &userv1.AddUserRequest{
		Name:     gofakeit.FirstName(),
		Surname:  gofakeit.LastName(),
		Username: gofakeit.Username(),
		UserId:   101,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respAdd.GetUserId())
	userId := respAdd.GetUserId()

	t.Logf("testing Get User Service")
	respGet, err := st.UserClient.GetUser(ctx, &userv1.GetUserRequest{
		UserId: userId,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respGet.GetUser())

	t.Logf("testing Update User Service")
	respUpdate, err := st.UserClient.UpdateUser(ctx, &userv1.UpdateUserRequest{
		UserId:   userId,
		Name:     gofakeit.FirstName(),
		Surname:  gofakeit.LastName(),
		Username: gofakeit.Username(),
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respUpdate.GetUser())

	t.Logf("testing Delete User Service")
	respDelete, err := st.UserClient.DeleteUser(ctx, &userv1.DeleteUserRequest{
		UserId: userId,
	})
	require.NoError(t, err)
	assert.Equal(t, true, respDelete.GetIsSuccessfully())
}

func TestUserServiceAdd_FailPath(t *testing.T) {
	ctx, st := suite.New(t)

	//TODO дополнить тесты
	tests := []struct {
		id          int64
		name        string
		surname     string
		username    string
		expectedErr string
	}{
		{
			id:          1,
			name:        gofakeit.FirstName(),
			surname:     gofakeit.LastName(),
			username:    "",
			expectedErr: "username is required",
		},
		{
			id:          2,
			name:        gofakeit.FirstName(),
			surname:     "",
			username:    gofakeit.Username(),
			expectedErr: "surname is required",
		},
		{
			id:          3,
			name:        "",
			surname:     gofakeit.LastName(),
			username:    gofakeit.Username(),
			expectedErr: "name is required",
		},
		{
			id:          -1,
			name:        gofakeit.FirstName(),
			surname:     gofakeit.LastName(),
			username:    gofakeit.Username(),
			expectedErr: "userId is wrong",
		},
	}

	t.Logf("testing Add User Service errors")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.UserClient.AddUser(ctx, &userv1.AddUserRequest{
				UserId:   tt.id,
				Name:     tt.name,
				Surname:  tt.surname,
				Username: tt.username,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)

		})
	}
}

func TestUserServiceGet_FailPath(t *testing.T) {
	ctx, st := suite.New(t)

	//TODO дополнить тесты
	tests := []struct {
		id          int64
		expectedErr string
	}{
		{
			id:          100001, // слишком большой несуществующий id
			expectedErr: "user not found",
		},
		{
			id:          -1,
			expectedErr: "userId is wrong",
		},
	}

	t.Logf("testing Get User Service errors")
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.id), func(t *testing.T) {
			_, err := st.UserClient.GetUser(ctx, &userv1.GetUserRequest{
				UserId: tt.id,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestUserServiceDelete_FailPath(t *testing.T) {
	ctx, st := suite.New(t)

	//TODO дополнить тесты
	tests := []struct {
		id          int64
		expectedErr string
	}{
		{
			id:          -1,
			expectedErr: "userId is wrong",
		},
	}

	t.Logf("testing Delete User Service errors")
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.id), func(t *testing.T) {
			_, err := st.UserClient.DeleteUser(ctx, &userv1.DeleteUserRequest{
				UserId: tt.id,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}
