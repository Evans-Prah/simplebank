package api

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	db "github.com/Evans-Prah/simplebank/db/sqlc"
	"github.com/Evans-Prah/simplebank/db/util"
	"github.com/stretchr/testify/require"
)

// func TestGetAccountAPI(t *testing.T)  {
// 	user := randomUser()
// 	account := randomAccount(user.Username)

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	store := mockdb.NewMockStore(ctrl)

// 	// build stubs
// 	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)

// 	// start test server and send request
// 	server := newTestServer(t, store)
// 	recorder := httptest.NewRecorder()

// 	url := fmt.Sprintf("/api/accounts/%d", account.ID)
// 	request, err := http.NewRequest(http.MethodGet, url, nil)
// 	require.NoError(t, err)

// 	server.router.ServeHTTP(recorder, request)
// 	// check response
// 	require.Equal(t, http.StatusOK, recorder.Code)
// 	//requireBodyMatchAccount(t, recorder.Body, account)
// }

func randomAccount(owner string) db.Account {
	return db.Account{
		ID: util.RandomInt(1, 1000),
		Owner: owner,
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func randomUser() db.User {
	return db.User{
		ID: util.RandomInt(1, 1000),
		FullName: util.RandomOwner(),
		Username: util.RandomString(8),
		Email: util.RandomEmail(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account)  {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var existingAccount db.Account
	err = json.Unmarshal(data, &existingAccount)
	require.NoError(t, err)
	require.Equal(t, account, existingAccount)
}


