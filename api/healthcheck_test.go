// mqttbridge
// https://github.com/topfreegames/mqttbridge
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package api

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestHealthcheckHandler(t *testing.T) {
	g := Goblin(t)

	// special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Healthcheck Handler", func() {
		g.It("Should respond with default WORKING string", func() {
			a := GetDefaultTestApp()
			res := Get(a, "/healthcheck", t)
			fmt.Printf("%v\n", res.Raw())

			g.Assert(res.Raw().StatusCode).Equal(http.StatusOK)
			res.Body().Equal("WORKING")
		})

		g.It("Should respond with customized WORKING string", func() {
			a := GetDefaultTestApp()
			a.Config.SetDefault("healthcheck.workingText", "OTHERWORKING")
			res := Get(a, "/healthcheck", t)

			g.Assert(res.Raw().StatusCode).Equal(http.StatusOK)
			res.Body().Equal("OTHERWORKING")
		})
	})
}
