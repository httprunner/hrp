package examples

import (
	"testing"

	"github.com/httprunner/hrp"
)

var rendezvousTestcase = &hrp.TestCase{
	Config: hrp.NewConfig("run request with functions").
		SetBaseURL("https://postman-echo.com").
		WithVariables(map[string]interface{}{
			"n": 5,
			"a": 12.3,
			"b": 3.45,
		}),
	TestSteps: []hrp.IStep{
		hrp.NewStep("waiting for all users").
			InsertRendezvousByPercent("rend0", 1, 0),
		hrp.NewStep("get with params").
			GET("/get").
			WithParams(map[string]interface{}{"foo1": "${gen_random_string($n)}", "foo2": "${max($a, $b)}"}).
			WithHeaders(map[string]string{"User-Agent": "HttpRunnerPlus"}).
			Extract().
			WithJmesPath("body.args.foo1", "varFoo1").
			Validate().
			AssertEqual("status_code", 200, "check status code"),
		hrp.NewStep("rendezvous1").
			InsertRendezvousByNumber("rend1", 400, 3000),
		hrp.NewStep("post json data with functions").
			POST("/post").
			WithHeaders(map[string]string{"User-Agent": "HttpRunnerPlus"}).
			WithBody(map[string]interface{}{"foo1": "${gen_random_string($n)}", "foo2": "${max($a, $b)}"}).
			Validate().
			AssertEqual("status_code", 200, "check status code").
			AssertLengthEqual("body.json.foo1", 5, "check args foo1").
			AssertEqual("body.json.foo2", 12.3, "check args foo2"),
		hrp.NewStep("rendezvous2").
			InsertRendezvousByNumber("rend2", 200, 2000),
	},
}

func TestRendezvous(t *testing.T) {
	err := hrp.NewRunner(t).Run(rendezvousTestcase)
	if err != nil {
		t.Fatalf("run testcase error: %v", err)
	}
}

func TestRendezvousDump2JSON(t *testing.T) {
	tCase, err := rendezvousTestcase.ToTCase()
	if err != nil {
		t.Fatalf("ToTCase error: %v", err)
	}
	err = tCase.Dump2JSON("rendezvous_test.json")
	if err != nil {
		t.Fatalf("dump to json error: %v", err)
	}
}
