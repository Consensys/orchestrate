package steps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cucumber/godog"
	gherkin "github.com/cucumber/messages-go/v10"
	merror "github.com/hashicorp/go-multierror"
	log "github.com/sirupsen/logrus"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/utils"
)

func (sc *ScenarioContext) resetResponse(*gherkin.Pickle) {
	sc.httpResponse = &http.Response{}
}

func (sc *ScenarioContext) resetAuth(*gherkin.Pickle) {
	sc.authSetup = AuthSetup{}
}

func (sc *ScenarioContext) iSetAuth(authMethod, value string) error {
	switch authMethod {
	case "API-Key":
		log.Tracef("API Key Set")
		sc.authSetup.authMethod = authutils.APIKeyHeader
		sc.authSetup.authData = value
		return nil
	case "JWT":
		log.Tracef("JWT Set")
		sc.authSetup.authMethod = authutils.AuthorizationHeader
		sc.authSetup.authData = value
		return nil
	default:
		sc.logger.Trace("No authentication set for the request")
		return nil
	}
}

func (sc *ScenarioContext) iInjectAuth(req *http.Request) error {
	switch sc.authSetup.authMethod {
	case authutils.APIKeyHeader:
		log.Tracef("API Key Set")
		req.Header.Add(authutils.APIKeyHeader, sc.authSetup.authData)
	case authutils.AuthorizationHeader:
		log.Tracef("JWT")
		authorization, err := sc.parser.JWTGenerator.GenerateAccessTokenWithTenantID(sc.authSetup.authData, 24*time.Hour)
		if err != nil {
			return err
		}
		log.Tracef("Auth JWT token: %s", authorization)
		req.Header.Add(authutils.AuthorizationHeader, "Bearer "+authorization)
	default:
		req.Header.Del(authutils.AuthorizationHeader)
		req.Header.Del(authutils.APIKeyHeader)
	}
	return nil
}

func (sc *ScenarioContext) iManageResponse(req *http.Request) error {
	resp, err := sc.httpClient.Do(req)
	if err != nil {
		return err
	} else if resp == nil {
		return fmt.Errorf("invalid response: ")
	}
	sc.httpResponse = resp
	return nil
}

func (sc *ScenarioContext) iSendRequestTo(method, endpoint string) error {
	sc.resetResponse(nil)

	endpoint, err := sc.replaceEndpointAliases(endpoint)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return err
	}

	err = sc.iInjectAuth(req)
	if err != nil {
		return err
	}

	err = sc.iManageResponse(req)
	if err != nil {
		err = merror.Append(err, fmt.Errorf("when sending %s request on %s", method, endpoint))
		return err
	}

	return nil
}

func (sc *ScenarioContext) iSendRequestToWithJSON(method, endpoint string, body *gherkin.PickleStepArgument_PickleDocString) error {
	sc.resetResponse(nil)

	endpoint, err := sc.replaceEndpointAliases(endpoint)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer([]byte(body.Content)))
	if err != nil {
		return err
	}

	err = sc.iInjectAuth(req)
	if err != nil {
		return err
	}

	err = sc.iManageResponse(req)
	if err != nil {
		err = merror.Append(err, fmt.Errorf("when sending %s request on %s", method, endpoint))
		return err
	}

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

func (sc *ScenarioContext) theResponseShouldMatchJSON(expectedBytes *gherkin.PickleStepArgument_PickleDocString) (err error) {
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

	sc.aliases.Set(sc.Pickle.Id, alias, data.UUID)
	return
}

func (sc *ScenarioContext) iStoreResponseFieldAs(navigation, alias string) (err error) {
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

	var val interface{}
	val, err = navJSONResponse(navigation, bodyBytes)
	if err != nil || val == nil {
		return
	}

	sc.aliases.Set(sc.Pickle.Id, alias, val.(string))
	return
}

var r = regexp.MustCompile("{{([^}]*)}}")

func (sc *ScenarioContext) replaceEndpointAliases(endpoint string) (string, error) {
	for _, alias := range r.FindAllStringSubmatch(endpoint, -1) {
		v, ok := sc.aliases.Get(sc.Pickle.Id, alias[1])
		if !ok {
			v, ok = sc.aliases.Get(GenericNamespace, alias[1])
			if !ok {
				return "", fmt.Errorf("could not replace alias %s", v)
			}
		}
		endpoint = strings.Replace(endpoint, alias[0], v, 1)
	}
	return endpoint, nil
}

func navJSONResponse(nav string, bodyBytes []byte) (interface{}, error) {
	var resp interface{}
	if err := json.Unmarshal(bodyBytes, &resp); err != nil {
		return "", err
	}

	var result interface{}
	// Navigate throw the json response
	navigation := strings.Split(nav, ".")
	for _, navStep := range navigation {
		if jdx, err := strconv.Atoi(navStep); err == nil {
			respAcum := resp.([]interface{})
			result = respAcum[jdx]
			resp = result
		} else {
			respAcum := resp.(map[string]interface{})
			result = respAcum[navStep]
			resp = result
		}
	}

	return result, nil
}

func initHTTP(s *godog.Suite, sc *ScenarioContext) {

	s.BeforeScenario(sc.resetResponse)
	s.BeforeScenario(sc.resetAuth)

	s.Step(`^I set authentication method "(API-Key|JWT)" with "([^"]*)"$`, sc.iSetAuth)
	s.Step(`^I send "(GET|POST|PATCH|PUT|DELETE)" request to "([^"]*)"$`, sc.iSendRequestTo)
	s.Step(`^I send "(GET|POST|PATCH|PUT|DELETE)" request to "([^"]*)" with json:$`, sc.iSendRequestToWithJSON)
	s.Step(`^I store the UUID as "([^"]*)"`, sc.iStoreTheUUIDAs)
	s.Step(`^I store response field "([^"]*)" as "([^"]*)"`, sc.iStoreResponseFieldAs)
	s.Step(`^the response code should be (\d+)$`, sc.theResponseCodeShouldBe)
	s.Step(`^the response should match json:$`, sc.theResponseShouldMatchJSON)
}
