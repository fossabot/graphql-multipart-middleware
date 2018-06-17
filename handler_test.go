package graphqlmultipart_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	graphqlmultipart "github.com/lucassabreu/graphql-multipart-middleware"
	"github.com/lucassabreu/graphql-multipart-middleware/testutil"

	"github.com/stretchr/testify/require"
)

func TestMultipartMiddleware_ForwardsRequestsOtherThanMultipart(t *testing.T) {
	newRequest := func(contentType string) *http.Request {
		r, _ := http.NewRequest("GET", "/graphql", strings.NewReader(""))
		r.Header.Set("Content-type", contentType)
		return r
	}

	cases := map[string]*http.Request{
		"application/json":                  newRequest("application/json"),
		"application/graphql":               newRequest("application/graphql"),
		"application/x-www-form-urlencoded": newRequest("application/x-www-form-urlencoded"),
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("reached me"))
	})

	for n, r := range cases {
		t.Run(n, func(t *testing.T) {
			resp := httptest.NewRecorder()

			mh := graphqlmultipart.NewHandler(
				&testutil.Schema,
				1*1024,
				h,
			)

			mh.ServeHTTP(resp, r)

			body, _ := ioutil.ReadAll(resp.Result().Body)

			require.Equal(t, string(body), "reached me")
		})
	}
}
