package steps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
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

	endpoint, err := sc.replaceAllMatchesAliases(endpoint)
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

	endpoint, err := sc.replaceAllMatchesAliases(endpoint)
	if err != nil {
		return err
	}

	reqBody, err := sc.replaceAllMatchesAliases(body.Content)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer([]byte(reqBody)))
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

func (sc *ScenarioContext) responseShouldHaveFields(table *gherkin.PickleStepArgument_PickleTable) (err error) {
	header := table.Rows[0]
	rowResponse := table.Rows[1]

	body, err := ioutil.ReadAll(sc.httpResponse.Body)
	if err != nil {
		return fmt.Errorf("expected response code body to math field but it errored with %s", err.Error())
	}
	sc.httpResponse.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	for c, col := range rowResponse.Cells {
		fieldName := header.Cells[c].Value
		respVal, err := navJSONResponse(fieldName, body)
		if err != nil {
			return err
		}
		if respVal == nil {
			continue
		}

		field := reflect.ValueOf(respVal)
		if col.Value == "~" {
			if isEqual("", field) {
				return fmt.Errorf("response did not expected %s to be empty", fmt.Sprintf("%v", fieldName))
			}
			continue
		}

		var aliasRE = regexp.MustCompile(`{{(.*)}}`)
		if aliasRE.MatchString(col.Value) {
			alias := aliasRE.FindStringSubmatch(col.Value)[1]
			val, _ := sc.aliases.Get(sc.Pickle.Id, alias)
			if !isEqual(val, field) {
				return fmt.Errorf("response %s expected %s but got %s", fieldName, val, fmt.Sprintf("%v", field))
			}

			continue
		}

		if !isEqual(col.Value, field) {
			return fmt.Errorf("response %s expected %s but got %s", fieldName, col.Value, fmt.Sprintf("%v", field))
		}
	}
	return nil
}

func (sc *ScenarioContext) iStoreTheUUIDAs(alias string) (err error) {
	body, err := ioutil.ReadAll(sc.httpResponse.Body)
	if err != nil {
		return fmt.Errorf("expected response code body to math field but it errored with %s", err.Error())
	}
	sc.httpResponse.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var data struct {
		UUID string `json:"uuid"`
	}
	if err = json.Unmarshal(body, &data); err != nil {
		return
	}

	sc.aliases.Set(sc.Pickle.Id, alias, data.UUID)
	return
}

func (sc *ScenarioContext) iStoreResponseFieldAs(navigation, alias string) (err error) {
	body, err := ioutil.ReadAll(sc.httpResponse.Body)
	if err != nil {
		return fmt.Errorf("expected response code body to math field but it errored with %s", err.Error())
	}
	sc.httpResponse.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var val interface{}
	val, err = navJSONResponse(navigation, body)
	if err != nil || val == nil {
		return
	}

	sc.aliases.Set(sc.Pickle.Id, alias, val.(string))
	return
}

var r = regexp.MustCompile("{{([^}]*)}}")

func (sc *ScenarioContext) replaceAllMatchesAliases(foo string) (string, error) {
	for _, alias := range r.FindAllStringSubmatch(foo, -1) {
		v, ok := sc.aliases.Get(sc.Pickle.Id, alias[1])
		if !ok {
			v, ok = sc.aliases.Get(GenericNamespace, alias[1])
			if !ok {
				return "", fmt.Errorf("could not replace alias %s", v)
			}
		}
		foo = strings.Replace(foo, alias[0], v, 1)
	}
	return foo, nil
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
		if resp == nil {
			return "", fmt.Errorf("could not find response field '%s'", nav)
		}
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
	s.Step(`^Response should have the following fields:$`, sc.responseShouldHaveFields)
}
