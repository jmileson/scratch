package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandle(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(register())
	defer server.Close()

	resp, err := http.Get(server.URL)
	assert.NoError(err)
	defer resp.Body.Close()

	assert.Equal(http.StatusOK, resp.StatusCode)

	b, err := ioutil.ReadAll(resp.Body)
	assert.NoError(err)
	assert.Equal("{\"result\":true}\n", string(b))
}
