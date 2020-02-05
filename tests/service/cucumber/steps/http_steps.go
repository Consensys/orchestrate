package steps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
)

func (sc *ScenarioContext) resetResponse(interface{}) {
	sc.httpResponse = &http.Response{}
}

var r = regexp.MustCompile("{{([^}]*)}}")

func (sc *ScenarioContext) iSendRequestTo(method, endpoint string) error {
	sc.httpResponse = &http.Response{}

	endpoint, err := sc.replaceEndpointAliases(endpoint)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", viper.GetString(client.ChainRegistryURLViperKey), endpoint), nil)
	if err != nil {
		return err
	}

	resp, err := sc.httpClient.Do(req)
	if err != nil {
		return err
	} else if resp == nil {
		return fmt.Errorf("invalid response when sending %s request on %s", method, endpoint)
	}
	sc.httpResponse = resp
	return nil
}

func (sc *ScenarioContext) iSendRequestToWithJSON(method, endpoint string, body *gherkin.DocString) error {
	sc.httpResponse = &http.Response{}

	endpoint, err := sc.replaceEndpointAliases(endpoint)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", viper.GetString(client.ChainRegistryURLViperKey), endpoint), bytes.NewBuffer([]byte(body.Content)))
	if err != nil {
		return err
	}

	resp, err := sc.httpClient.Do(req)
	if err != nil {
		return err
	} else if resp == nil {
		return fmt.Errorf("invalid response when sending %s request on %s", method, endpoint)
	}
	sc.httpResponse = resp
	return nil
}

func (sc *ScenarioContext) theResponseCodeShouldBe(code int) error {
	if code != sc.httpResponse.StatusCode {
		body, err := ioutil.ReadAll(sc.httpResponse.Body)
		if err != nil {
			return fmt.Errorf("expected response code %d, but actual is %q also could not read body", code, sc.httpResponse.Status)
		}
		sc.httpResponse.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		return fmt.Errorf("expected response code %d, but actual is %q with error: %s", code, sc.httpResponse.Status, body)
	}
	return nil
}

func (sc *ScenarioContext) theResponseShouldMatchJSON(expectedBytes *gherkin.DocString) (err error) {
	var expected, body []byte
	var data interface{}
	if err = json.Unmarshal([]byte(expectedBytes.Content), &data); err != nil {
		return
	}
	if expected, err = json.Marshal(data); err != nil {
		return
	}
	defer func() {
		closeErr := sc.httpResponse.Body.Close()
		if closeErr != nil {
			log.Error("could not properly close response body")
		}
	}()

	body, err = ioutil.ReadAll(sc.httpResponse.Body)
	if err != nil {
		return
	}
	sc.httpResponse.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	if !bytes.Equal(body, expected) {
		err = fmt.Errorf("expected json, does not match actual: %s", string(body))
	}
	return
}

func (sc *ScenarioContext) iStoreTheUUIDAs(alias string) (err error) {
	defer func() {
		closeErr := sc.httpResponse.Body.Close()
		if closeErr != nil {
			log.Error("could not properly close response body")
		}
	}()
	if sc.httpResponse == nil {
		return fmt.Errorf("no http response stored, cannot retrieve ChainID")
	}

	bodyBytes, err := ioutil.ReadAll(sc.httpResponse.Body)
	if err != nil {
		return
	}

	var data struct {
		UUID string `json:"uuid"`
	}
	if err = json.Unmarshal(bodyBytes, &data); err != nil {
		return
	}
	sc.httpAliases.Set(sc.ID, alias, data.UUID)
	return
}

func (sc *ScenarioContext) replaceEndpointAliases(endpoint string) (string, error) {
	for _, alias := range r.FindAllStringSubmatch(endpoint, -1) {
		v, ok := sc.httpAliases.Get(sc.ID, alias[1])
		if !ok {
			return "", fmt.Errorf("could not replace alias %s", v)
		}
		endpoint = strings.Replace(endpoint, alias[0], v, 1)
	}
	return endpoint, nil
}

func initHTTP(s *godog.Suite, sc *ScenarioContext) {

	s.BeforeScenario(sc.resetResponse)

	s.Step(`^I send "(GET|POST|PATCH|PUT|DELETE)" request to "([^"]*)"$`, sc.iSendRequestTo)
	s.Step(`^I send "(GET|POST|PATCH|PUT|DELETE)" request to "([^"]*)" with json:$`, sc.iSendRequestToWithJSON)
	s.Step(`^I store the UUID as "([^"]*)"`, sc.iStoreTheUUIDAs)
	s.Step(`^the response code should be (\d+)$`, sc.theResponseCodeShouldBe)
	s.Step(`^the response should match json:$`, sc.theResponseShouldMatchJSON)
}
