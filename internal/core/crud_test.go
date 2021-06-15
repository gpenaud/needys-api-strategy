package core

import(
	fmt 		 "fmt"
	godog 	 "github.com/cucumber/godog"
	http 		 "net/http"
	httptest "net/http/httptest"
	mux 		 "github.com/gorilla/mux"
)

type apiFeature struct {
  resp *httptest.ResponseRecorder
	router *mux.Router
}

func (a *apiFeature) iSendRequestTo(method, endpoint string) error {
	req, err := http.NewRequest(method, endpoint, nil)

	if err != nil {
		return err
	}

	switch endpoint {
	case "/strategy":
		if (method == "GET") {
			GetStrategies(a.resp, req)
		} else if (method == "POST") {
			CreateStrategy(a.resp, req)
		} else {
			a.resp.WriteHeader(http.StatusMethodNotAllowed)
		}

	default:
		a.resp.WriteHeader(http.StatusNotFound)
	}

	return err
}

func (a *apiFeature) theResponseCodeShouldBe(code int) error {
	if code != a.resp.Code {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, a.resp.Code)
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	api := &apiFeature{}

  ctx.BeforeScenario(func(*godog.Scenario) {
		api.resp = httptest.NewRecorder()
		api.router = mux.NewRouter()
	})

	ctx.Step(`^I send "([^"]*)" request to "([^"]*)"$`, api.iSendRequestTo)
	ctx.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
}
