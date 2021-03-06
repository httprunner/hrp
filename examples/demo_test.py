# NOTE: Generated By HttpRunner v3.1.6
# FROM: hrp/examples/demo.json


from httprunner import HttpRunner, Config, Step, RunRequest, RunTestCase


class TestCaseDemo(HttpRunner):

    config = (
        Config("demo with complex mechanisms")
        .variables(
            **{
                "a": 12.3,
                "b": 3.45,
                "n": 5,
                "varFoo1": "${gen_random_string($n)}",
                "varFoo2": "${max($a, $b)}",
            }
        )
        .base_url("https://postman-echo.com")
    )

    teststeps = [
        Step(
            RunRequest("get with params")
            .with_variables(**{"b": 34.5, "n": 3, "varFoo2": "${max($a, $b)}"})
            .get("/get")
            .with_params(**{"foo1": "$varFoo1", "foo2": "$varFoo2"})
            .with_headers(**{"User-Agent": "HttpRunnerPlus"})
            .extract()
            .with_jmespath("body.args.foo1", "varFoo1")
            .validate()
            .assert_equal("status_code", 200)
            .assert_equal('headers."Content-Type"', "application/json")
            .assert_equal("body.args.foo1", 5)
            .assert_equal("$varFoo1", 5)
            .assert_equal("body.args.foo2", "34.5")
        ),
        Step(
            RunRequest("post json data")
            .post("/post")
            .validate()
            .assert_equal("status_code", 200)
            .assert_equal("body.json.foo1", 5)
            .assert_equal("body.json.foo2", 12.3)
        ),
        Step(
            RunRequest("post form data")
            .post("/post")
            .with_headers(
                **{"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8"}
            )
            .validate()
            .assert_equal("status_code", 200)
            .assert_equal("body.form.foo1", 5)
            .assert_equal("body.form.foo2", "12.3")
        ),
    ]


if __name__ == "__main__":
    TestCaseDemo().test_start()
