package delete_test

import (
	"bytes"
	"encoding/json"

	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	urlDel "url-shortener/internal/http-server/handlers/url/delete"
	"url-shortener/internal/http-server/handlers/url/delete/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"

	"github.com/stretchr/testify/require"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
		},
		{
			name:      "Empty alias",
			alias:     "",
			respError: "field Alias is a required field",
		},
		{
			name:      "DeleteURL Error",
			alias:     "test_alias",
			respError: "Failed to delete URL",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlDeleterMock := mocks.NewURLDeleter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlDeleterMock.On("DeleteURL", tc.alias).
					Return(tc.mockError).Once()
			}

			handler := urlDel.New(slogdiscard.NewDiscardLogger(), urlDeleterMock)

			input := fmt.Sprintf(`{"alias": "%s"}`, tc.alias)

			req, err := http.NewRequest(http.MethodPost, "/url/delete", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp struct {
				Error string `json:"error"`
			}

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
