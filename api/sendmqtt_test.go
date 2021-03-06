// arkadiko
// https://github.com/topfreegames/arkadiko
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package api

import (
	"encoding/json"
	"net/http"
	"testing"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestSendMqtt(t *testing.T) {
	g := Goblin(t)

	// special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Send Mqtt", func() {
		g.It("Should respond with 200 for a valid message", func() {
			a := GetDefaultTestApp()
			var jsonPayload JSON
			testJSON := `{"message": "hello"}`
			response := `{"topic": "test", "payload": {"message":"hello"}}`
			json.Unmarshal([]byte(testJSON), &jsonPayload)
			res := PostJSON(a, "/sendmqtt/test", t, jsonPayload)

			g.Assert(res.Raw().StatusCode).Equal(http.StatusOK)
			res.Body().Equal(response)
		})

		g.It("Should respond with 200 for a valid message with hierarchical topic", func() {
			a := GetDefaultTestApp()
			var jsonPayload JSON
			testJSON := `{"message": "hello"}`
			response := `{"topic": "test/topic", "payload": {"message":"hello"}}`
			json.Unmarshal([]byte(testJSON), &jsonPayload)
			url := "/sendmqtt/test/topic"
			res := PostJSON(a, url, t, jsonPayload)

			g.Assert(res.Raw().StatusCode).Equal(http.StatusOK)
			res.Body().Equal(response)
		})

		g.It("Should respond with 400 if malformed JSON", func() {
			a := GetDefaultTestApp()
			var jsonPayload JSON
			testJSON := `{"message": "hello"}}`
			json.Unmarshal([]byte(testJSON), &jsonPayload)
			res := PostJSON(a, "/sendmqtt/test/topic", t, jsonPayload)

			g.Assert(res.Raw().StatusCode).Equal(400)
		})
	})
}
