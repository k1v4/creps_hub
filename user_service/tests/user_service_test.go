package tests

import (
	"github.com/brianvoe/gofakeit/v6"
	userv1 "github.com/k1v4/protos/gen/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"user_service/tests/suite"
)

func TestUserServiceAdd_HappyPath(t *testing.T) {
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

func TestUserServiceAdd_FailPath(t *testing.T) {}
