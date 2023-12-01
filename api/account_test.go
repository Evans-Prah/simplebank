package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/Evans-Prah/simplebank/db/mock"
	db "github.com/Evans-Prah/simplebank/db/sqlc"
	"github.com/Evans-Prah/simplebank/db/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T)  {
	account := randomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// build stubs
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)

	// start test server and send request
	server := newTestServer(t, store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/api/accounts/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	// check response
	require.Equal(t, http.StatusOK, recorder.Code)
	//requireBodyMatchAccount(t, recorder.Body, account)
}

func randomAccount() db.Account {
	return db.Account{
		ID: util.RandomInt(1, 1000),
		Owner: util.RandomOwner(),
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
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


