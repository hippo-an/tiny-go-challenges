package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	mockdb "github.com/hippo-an/tiny-go-challenges/mombank/db/mock"
	db "github.com/hippo-an/tiny-go-challenges/mombank/db/sqlc"
	"github.com/hippo-an/tiny-go-challenges/mombank/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func CreateUserEq(arg db.CreateUserParams, password string) gomock.Matcher {
	return createUserMatcher{
		arg:      arg,
		password: password,
	}
}

type createUserMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e createUserMatcher) Matches(x any) bool {
	c, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, c.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = c.HashedPassword

	return reflect.DeepEqual(e.arg, c)
}

func (e createUserMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func TestCreateUser(t *testing.T) {
	user := createRandomUser()
	password := "secret"

	testCases := []struct {
		Name      string
		Endpoint  string
		Req       createUserRequest
		BuildStub func(*mockdb.MockStore)
		CheckRes  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			Name:     "OK",
			Endpoint: "/users",
			Req: createUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			BuildStub: func(ms *mockdb.MockStore) {

				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}

				ms.EXPECT().
					CreateUser(gomock.Any(), CreateUserEq(arg, password)).
					Times(1).
					Return(user, nil)

			},
			CheckRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requiredBodyMatchUser(t, recorder.Body, user)
			},
		},

		{
			Name:     "UsernameBindingError",
			Endpoint: "/users",
			Req: createUserRequest{
				Username: "",
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			BuildStub: func(ms *mockdb.MockStore) {

				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}

				ms.EXPECT().
					CreateUser(gomock.Any(), CreateUserEq(arg, password)).
					Times(0)

			},
			CheckRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			Name:     "FullNameBindingError",
			Endpoint: "/users",
			Req: createUserRequest{
				Username: user.Username,
				Password: password,
				FullName: "",
				Email:    user.Email,
			},
			BuildStub: func(ms *mockdb.MockStore) {

				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}

				ms.EXPECT().
					CreateUser(gomock.Any(), CreateUserEq(arg, password)).
					Times(0)

			},
			CheckRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			Name:     "FullNameBindingError",
			Endpoint: "/users",
			Req: createUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    "",
			},
			BuildStub: func(ms *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}

				ms.EXPECT().
					CreateUser(gomock.Any(), CreateUserEq(arg, password)).
					Times(0)

			},
			CheckRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			Name:     "PasswordBindingError",
			Endpoint: "/users",
			Req: createUserRequest{
				Username: user.Username,
				Password: "1234",
				FullName: user.FullName,
				Email:    user.Email,
			},
			BuildStub: func(ms *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}

				ms.EXPECT().
					CreateUser(gomock.Any(), CreateUserEq(arg, password)).
					Times(0)

			},
			CheckRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			Name:     "CreateUserInternalServerError",
			Endpoint: "/users",
			Req: createUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			BuildStub: func(ms *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}

				ms.EXPECT().
					CreateUser(gomock.Any(), CreateUserEq(arg, password)).
					Times(1).
					Return(user, sql.ErrConnDone)

			},
			CheckRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.Name, func(t *testing.T) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()

			mockStore := mockdb.NewMockStore(ctr)

			tc.BuildStub(mockStore)
			server := NewTestServer(t, mockStore)
			recorder := httptest.NewRecorder()

			b, err := json.Marshal(tc.Req)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, tc.Endpoint, bytes.NewReader(b))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.CheckRes(t, recorder)
		})
	}
}

func createRandomUser() db.User {
	return db.User{
		ID:        util.RandomInt(1, 10000),
		Username:  util.RandomString(6),
		FullName:  util.RandomOwner(),
		Email:     util.RandomEmail(),
		CreatedAt: time.Now(),
	}
}

func requiredBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.WithinDuration(t, user.CreatedAt, gotUser.CreatedAt, time.Second)
}
