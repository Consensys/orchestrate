package steps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"github.com/cucumber/godog"
	gherkin "github.com/cucumber/messages-go/v10"
	merror "github.com/hashicorp/go-multierror"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/e2e/utils"
)

func (sc *ScenarioContext) resetResponse(*gherkin.Pickle) {
	sc.httpResponse = &http.Response{}
}

func (sc *ScenarioContext) iSetTheHeaders(table *gherkin.PickleStepArgument_PickleTable) error {
	headers := make(map[string]string)
	for _, v := range table.Rows[1:] {
		if len(v.Cells) != 2 {
			return errors.DataError("headers should be a 2 column table with key/value only")
		}
		headers[v.Cells[0].GetValue()] = v.Cells[1].GetValue()
	}
	sc.aliases.Set(headers, sc.Pickle.Id, "HTTP.Headers")

	return nil
}

func (sc *ScenarioContext) iInjectHeaders(req *http.Request) error {
	headers, ok := sc.aliases.Get(sc.Pickle.Id, "HTTP.Headers")
	if !ok {
		return nil
	}

	for k, v := range headers.(map[string]string) {
		req.Header.Add(k, v)
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

	endpoint, err := sc.replace(endpoint)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return err
	}

	err = sc.iInjectHeaders(req)
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

	endpoint, err := sc.replace(endpoint)
	if err != nil {
		return err
	}

	reqBody, err := sc.replace(body.Content)
	if err != nil {
		return err
	}

	sc.logger.Debugf("Request %s %s with body %v", method, endpoint, reqBody)

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		return err
	}

	err = sc.iInjectHeaders(req)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

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
	var resp interface{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return err
	}

	for c, col := range rowResponse.Cells {
		fieldName := header.Cells[c].Value
		field, err := utils.GetField(fieldName, reflect.ValueOf(resp))
		if err != nil {
			return err
		}

		if err := utils.CmpField(field, col.Value); err != nil {
			return fmt.Errorf("(%d/%d) %v %v", c+1, len(rowResponse.Cells), fieldName, err)
		}
	}
	return nil
}

func (sc *ScenarioContext) responseShouldHaveHeaders(table *gherkin.PickleStepArgument_PickleTable) (err error) {
	header := table.Rows[0]
	rowResponse := table.Rows[1]

	headers := sc.httpResponse.Header
	for c, col := range rowResponse.Cells {
		headerName := header.Cells[c].Value
		field := headers.Get(headerName)
		if err := utils.CmpField(reflect.ValueOf(field), col.Value); err != nil {
			return fmt.Errorf("(%d/%d) %v %v", c+1, len(rowResponse.Cells), headerName, err)
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

	sc.aliases.Set(data.UUID, sc.Pickle.Id, alias)
	return
}

func (sc *ScenarioContext) iRegisterTheFollowingResponseFields(table *gherkin.PickleStepArgument_PickleTable) (err error) {
	body, err := ioutil.ReadAll(sc.httpResponse.Body)
	if err != nil {
		return fmt.Errorf("expected response code body to math field but it errored with %s", err.Error())
	}
	sc.httpResponse.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	var resp interface{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return err
	}

	aliasTable := utils.ExtractColumns(table, []string{aliasHeaderValue})
	for i, row := range aliasTable.Rows[1:] {

		alias := row.Cells[0].Value
		bodyPath := table.Rows[i+1].Cells[0].Value
		val, err := utils.GetField(bodyPath, reflect.ValueOf(resp))
		if err != nil {
			return err
		}
		sc.aliases.Set(val, sc.Pickle.Id, alias)
	}

	return nil
}

func (sc *ScenarioContext) iSleep(s string) error {
	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	time.Sleep(d)

	return nil
}

func initHTTP(s *godog.ScenarioContext, sc *ScenarioContext) {

	s.BeforeScenario(sc.resetResponse)

	s.Step(`^I send "(GET|POST|PATCH|PUT|DELETE)" request to "([^"]*)"$`, sc.iSendRequestTo)
	s.Step(`^I send "(GET|POST|PATCH|PUT|DELETE)" request to "([^"]*)" with json:$`, sc.iSendRequestToWithJSON)
	s.Step(`^I store the UUID as "([^"]*)"`, sc.iStoreTheUUIDAs)
	s.Step(`^I register the following response fields$`, sc.preProcessTableStep(sc.iRegisterTheFollowingResponseFields))
	s.Step(`^the response code should be (\d+)$`, sc.theResponseCodeShouldBe)
	s.Step(`^the response should match json:$`, sc.theResponseShouldMatchJSON)
	s.Step(`^Response should have the following fields$`, sc.preProcessTableStep(sc.responseShouldHaveFields))
	s.Step(`^Response should have the following headers$`, sc.preProcessTableStep(sc.responseShouldHaveHeaders))
	s.Step(`^I set the headers$`, sc.preProcessTableStep(sc.iSetTheHeaders))
	s.Step(`^I sleep "([^"]*)"$`, sc.iSleep)
}
