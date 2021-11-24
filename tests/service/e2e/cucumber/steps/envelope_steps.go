package steps

import (
	"bytes"
	"context"
	json2 "encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-kit/kit/transport/http/jsonrpc"

	clientutils "github.com/consensys/orchestrate/pkg/toolkit/app/http/client-utils"

	"github.com/cenkalti/backoff/v4"
	pkgcryto "github.com/consensys/orchestrate/pkg/crypto/ethereum"
	"github.com/consensys/orchestrate/pkg/encoding/rlp"
	"github.com/consensys/orchestrate/pkg/types/api"
	utils4 "github.com/consensys/orchestrate/pkg/utils"
	utils3 "github.com/consensys/orchestrate/tests/utils"

	"github.com/Shopify/sarama"
	"github.com/consensys/orchestrate/pkg/encoding/json"
	encoding "github.com/consensys/orchestrate/pkg/encoding/sarama"
	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/ethereum/account"
	utils2 "github.com/consensys/orchestrate/pkg/toolkit/ethclient/utils"
	"github.com/consensys/orchestrate/pkg/types/tx"
	"github.com/consensys/orchestrate/tests/service/e2e/cucumber/alias"
	"github.com/consensys/orchestrate/tests/service/e2e/utils"
	"github.com/cucumber/godog"
	gherkin "github.com/cucumber/messages-go/v10"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const aliasHeaderValue = "alias"

var aliasRegex = regexp.MustCompile("{{([^}]*)}}")
var AddressPtrType = reflect.TypeOf(new(common.Address))

func (sc *ScenarioContext) sendEnvelope(topic string, e *tx.Envelope) error {
	// Prepare message to be sent
	msg := &sarama.ProducerMessage{
		Topic: viper.GetString(fmt.Sprintf("topic.%v", topic)),
		Key:   sarama.StringEncoder(e.PartitionKey()),
	}

	err := encoding.Marshal(e.TxEnvelopeAsRequest(), msg)
	if err != nil {
		return err
	}

	// Send message
	_, _, err = sc.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"id":            e.GetID(),
		"scenario.id":   sc.Pickle.Id,
		"scenario.name": sc.Pickle.Name,
	}).Debugf("scenario: envelope sent")

	return nil
}

func (sc *ScenarioContext) iSendEnvelopesToTopic(topic string, table *gherkin.PickleStepArgument_PickleTable) error {
	// Parse table
	if err := sc.replaceAliases(table); err != nil {
		return err
	}

	envelopes, err := utils.ParseEnvelope(table)
	if err != nil {
		return errors.DataError("could not parse tx request - got %v", err)
	}

	// Set trackers for each envelope
	sc.setTrackers(sc.newTrackers(envelopes))

	// Send envelopes
	for _, t := range sc.trackers {
		err := sc.sendEnvelope(topic, t.Current)
		if err != nil {
			return errors.InternalError("could not send tx request - got %v", err)
		}
	}

	return nil
}

func (sc *ScenarioContext) registerEnvelopeTracker(value string) error {
	envelopeID, ok := sc.aliases.Get(sc.Pickle.Id, value)
	if !ok {
		envelopeID, ok = sc.aliases.Get("global", value)
		if !ok {
			envelopeID = value
		}
	}

	evlp := tx.NewEnvelope()
	_ = evlp.SetID(envelopeID.(string)).
		SetContextLabelsValue("debug", "true").
		SetContextLabelsValue("scenario.id", sc.Pickle.Id)

	sc.setTrackers(append(sc.trackers, sc.newTracker(evlp)))

	return nil
}

func (sc *ScenarioContext) envelopeShouldBeInTopic(topic string) error {
	for i, t := range sc.trackers {
		err := t.Load(topic, viper.GetDuration(CucumberTimeoutViperKey))
		if err != nil {
			e := t.Load("tx.recover", time.Millisecond)
			if e != nil {
				return fmt.Errorf("%v: envelope n°%v neither in topic %q nor in %q", sc.Pickle.Id, i, topic, "tx.recover")
			}
			return fmt.Errorf("%v: envelope n°%v not in topic %q but found in %q - envelope.Errors %q", sc.Pickle.Id, i, topic, "tx.recover", t.Current.Error())
		}
	}

	// Waiting for job to be updated after notifying (Hacky and ugly)
	if topic == utils3.TxDecodedTopicKey || topic == utils3.TxRecoverTopicKey {
		time.Sleep(time.Second)
	}
	return nil
}

func (sc *ScenarioContext) envelopesShouldHaveTheFollowingValues(table *gherkin.PickleStepArgument_PickleTable) error {
	header := table.Rows[0]
	rows := table.Rows[1:]
	if len(rows) != len(sc.trackers) {
		return fmt.Errorf("expected as much rows as envelopes tracked")
	}

	for r, row := range rows {
		val := reflect.ValueOf(sc.trackers[r].Current).Elem()
		sEvlp, err := json.Marshal(*sc.trackers[r].Current)
		log.WithError(err).Debugf("Marshaled envelope: %s", utils4.ShortString(fmt.Sprint(sEvlp), 30))
		for c, col := range row.Cells {
			fieldName := header.Cells[c].Value
			field, err := utils.GetField(fieldName, val)
			if err != nil {
				return err
			}

			if err := utils.CmpField(field, col.Value); err != nil {
				return fmt.Errorf("(%d/%d) %v %v", r+1, len(rows), fieldName, err)
			}
		}
	}

	return nil
}

func (sc *ScenarioContext) iRegisterTheFollowingEnvelopeFields(table *gherkin.PickleStepArgument_PickleTable) (err error) {

	evlps := make(map[string]*tx.Envelope)
	for _, tracker := range sc.trackers {
		evlp := tracker.Current
		evlps[evlp.GetID()] = evlp
		evlps[evlp.GetContextLabelsValue("id")] = evlp
	}

	header := table.Rows[0]
	rows := table.Rows[1:]
	for idx, h := range []string{"id", "alias", "path"} {
		if header.Cells[idx].Value != h {
			return fmt.Errorf("invalid first column table header: expected '%s', found '%s'", h, header.Cells[idx].Value)
		}
	}

	for i, row := range rows {
		evlpID := row.Cells[0].Value
		if aliasRegex.MatchString(evlpID) {
			evlpID = aliasRegex.FindString(evlpID)
		}

		evlp, ok := evlps[evlpID]
		if !ok {
			return fmt.Errorf("envelope %s is not found: %q", evlpID, row)
		}

		a := row.Cells[1].Value
		bodyPath := table.Rows[i+1].Cells[2].Value
		val, err := utils.GetField(bodyPath, reflect.ValueOf(evlp))
		if err != nil {
			return err
		}

		switch val.Type() {
		case AddressPtrType:
			sc.aliases.Set(val.Interface().(*common.Address).Hex(), sc.Pickle.Id, a)
		default:
			sc.aliases.Set(val, sc.Pickle.Id, a)
		}
	}

	return nil
}

func (sc *ScenarioContext) tearDown(s *gherkin.Pickle, err error) {
	var wg sync.WaitGroup
	wg.Add(len(sc.TearDownFunc))

	for _, f := range sc.TearDownFunc {
		f := f
		go func() {
			defer wg.Done()
			f()
		}()
	}
	wg.Wait()
}

func (sc *ScenarioContext) iHaveTheFollowingTenant(table *gherkin.PickleStepArgument_PickleTable) error {
	headers := table.Rows[0]
	for _, row := range table.Rows[1:] {
		tenantMap := make(map[string]interface{})
		var a string
		var tenantID string
		var username string

		for i, cell := range row.Cells {
			switch v := headers.Cells[i].Value; {
			case v == aliasHeaderValue:
				a = cell.Value
			case v == "tenantID":
				tenantID = cell.Value
			case v == "username":
				username = cell.Value
			default:
				tenantMap[v] = cell.Value
			}
		}
		if a == "" {
			return errors.DataError("need an alias")
		}
		if tenantID == "" {
			tenantID = uuid.Must(uuid.NewV4()).String()
		}

		tenantMap["tenantID"] = tenantID
		if username != "" {
			tenantMap["username"] = username
		}
		sc.aliases.Set(tenantMap, sc.Pickle.Id, a)
	}

	return nil
}

func (sc *ScenarioContext) iHaveTheFollowingAccount(table *gherkin.PickleStepArgument_PickleTable) error {
	headers := table.Rows[0]
	for _, row := range table.Rows[1:] {
		accountMap := make(map[string]interface{})
		var aliass string

		for i, cell := range row.Cells {
			switch v := headers.Cells[i].Value; {
			case v == aliasHeaderValue:
				aliass = cell.Value
			default:
				accountMap[v] = cell.Value
			}
		}

		if aliass == "" {
			return errors.DataError("need an alias")
		}

		w := account.NewAccount()
		_ = w.Generate()
		privBytes := crypto.FromECDSA(w.Priv())
		accountMap["address"] = w.Address().Hex()
		accountMap["private_key"] = hexutil.Encode(privBytes)
		sc.aliases.Set(accountMap, sc.Pickle.Id, aliass)
	}

	return nil
}

func (sc *ScenarioContext) iHaveCreatedTheFollowingAccounts(table *gherkin.PickleStepArgument_PickleTable) error {
	tenantCol := utils.ExtractColumns(table, []string{"Tenant"})
	apiKeyCol := utils.ExtractColumns(table, []string{"API-KEY"})
	accIDCol := utils.ExtractColumns(table, []string{"ID"})
	accChainCol := utils.ExtractColumns(table, []string{"ChainName"})
	aliasCol := utils.ExtractColumns(table, []string{aliasHeaderValue})
	if aliasCol == nil {
		return errors.DataError("alias column is mandatory")
	}

	for idx := range apiKeyCol.Rows[1:] {
		req := &api.CreateAccountRequest{}

		if accIDCol != nil {
			req.Alias = accIDCol.Rows[idx+1].Cells[0].Value
		}

		if accChainCol != nil {
			req.Chain = accChainCol.Rows[idx+1].Cells[0].Value
		}

		tenant := ""
		if tenantCol != nil {
			tenant = tenantCol.Rows[idx+1].Cells[0].Value
		}
		headers := utils.GetHeaders(apiKeyCol.Rows[idx+1].Cells[0].Value, tenant, "")
		accRes, err := sc.client.CreateAccount(context.WithValue(context.Background(), clientutils.RequestHeaderKey, headers), req)

		if err != nil {
			return err
		}

		sc.aliases.Set(accRes.Address, sc.Pickle.Id, aliasCol.Rows[idx+1].Cells[0].Value)
	}

	return nil
}

func (sc *ScenarioContext) iRegisterTheFollowingChains(table *gherkin.PickleStepArgument_PickleTable) error {
	aliasCol := utils.ExtractColumns(table, []string{"alias"})
	tenantCol := utils.ExtractColumns(table, []string{"Tenant"})
	apiKeyCol := utils.ExtractColumns(table, []string{"API-KEY"})

	interfaceSlices, err := utils.ParseTable(api.RegisterChainRequest{}, table)
	if err != nil {
		return err
	}

	onTearDown := func(uuid string, headers map[string]string) func() {
		return func() {
			_ = sc.client.DeleteChain(context.WithValue(context.Background(), clientutils.RequestHeaderKey, headers), uuid)
		}
	}

	ctx := context.Background()
	for i, chain := range interfaceSlices {
		apiKey := apiKeyCol.Rows[i+1].Cells[0].Value

		tenant := ""
		if tenantCol != nil {
			tenant = tenantCol.Rows[i+1].Cells[0].Value
		}

		headers := utils.GetHeaders(apiKey, tenant, "")
		res, err := sc.client.RegisterChain(context.WithValue(context.Background(), clientutils.RequestHeaderKey, headers), chain.(*api.RegisterChainRequest))
		if err != nil {
			return err
		}
		sc.TearDownFunc = append(sc.TearDownFunc, onTearDown(res.UUID, headers))

		apiURL, _ := sc.aliases.Get(alias.GlobalAka, "api")
		proxyURL := utils4.GetProxyURL(apiURL.(string), res.UUID)
		err = backoff.RetryNotify(
			func() error {
				_, err2 := sc.ec.Network(ctx, proxyURL)
				return err2
			},
			backoff.WithMaxRetries(backoff.NewConstantBackOff(time.Second), 5),
			func(err error, duration time.Duration) {
				log.WithFields(log.Fields{
					"chain_uuid": res.UUID,
				}).WithError(err).Debug("scenario: chain proxy is still not ready")
			},
		)

		if err != nil {
			return err
		}

		// set aliases
		sc.aliases.Set(res, sc.Pickle.Id, aliasCol.Rows[i+1].Cells[0].Value)
	}

	return nil
}

func (sc *ScenarioContext) iRegisterTheFollowingFaucets(table *gherkin.PickleStepArgument_PickleTable) error {
	tenantCol := utils.ExtractColumns(table, []string{"Tenant"})
	apiKeyCol := utils.ExtractColumns(table, []string{"API-KEY"})

	interfaceSlices, err := utils.ParseTable(api.RegisterFaucetRequest{}, table)
	if err != nil {
		return err
	}

	f := func(ctx context.Context, uuid string) func() {
		return func() { _ = sc.client.DeleteFaucet(ctx, uuid) }
	}

	for i, faucet := range interfaceSlices {
		apiKey := apiKeyCol.Rows[i+1].Cells[0].Value

		tenant := ""
		if tenantCol != nil {
			tenant = tenantCol.Rows[i+1].Cells[0].Value
		}

		ctx := context.WithValue(context.Background(), clientutils.RequestHeaderKey, utils.GetHeaders(apiKey, tenant, ""))
		res, err := sc.client.RegisterFaucet(ctx, faucet.(*api.RegisterFaucetRequest))
		if err != nil {
			return err
		}
		sc.TearDownFunc = append(sc.TearDownFunc, f(ctx, res.UUID))
	}

	return nil
}

func (sc *ScenarioContext) replace(s string) (string, error) {
	for _, matchedAlias := range aliasRegex.FindAllStringSubmatch(s, -1) {
		aka := []string{matchedAlias[1]}
		if strings.HasPrefix(matchedAlias[1], "random.") {
			random := strings.Split(matchedAlias[1], ".")
			switch random[1] {
			case "uuid":
				s = strings.Replace(s, matchedAlias[0], uuid.Must(uuid.NewV4()).String(), 1)
			case "account":
				w := account.NewAccount()
				_ = w.Generate()
				s = strings.Replace(s, matchedAlias[0], w.Address().Hex(), 1)
			case "private_key":
				w := account.NewAccount()
				_ = w.Generate()
				privBytes := crypto.FromECDSA(w.Priv())
				s = strings.Replace(s, matchedAlias[0], hexutil.Encode(privBytes)[2:], 1)
			case "int":
				s = strings.Replace(s, matchedAlias[0], fmt.Sprintf("%d", rand.Int()), 1)
			}
			continue
		}

		if !strings.HasPrefix(matchedAlias[1], "global.") && !strings.HasPrefix(matchedAlias[1], "chain.") {
			aka = append([]string{sc.Pickle.Id}, aka...)
		}
		v, ok := sc.aliases.Get(aka...)
		if !ok {
			return "", fmt.Errorf("could not replace alias '%s'", matchedAlias[1])
		}

		val := reflect.ValueOf(v)

		var str string
		switch val.Kind() {
		case reflect.Array, reflect.Slice:
			strb, _ := json.Marshal(v)
			str = string(strb)
		default:
			str = fmt.Sprintf("%v", v)
		}
		s = strings.Replace(s, matchedAlias[0], str, 1)
	}
	return s, nil
}

func (sc *ScenarioContext) replaceAliases(table *gherkin.PickleStepArgument_PickleTable) error {
	for _, row := range table.Rows {
		for _, r := range row.Cells {
			s, err := sc.replace(r.Value)
			if err != nil {
				return err
			}
			r.Value = s
		}
	}
	return nil
}

func (sc *ScenarioContext) iRegisterTheFollowingAliasAs(table *gherkin.PickleStepArgument_PickleTable) error {
	aliasTable := utils.ExtractColumns(table, []string{aliasHeaderValue})
	if aliasTable == nil {
		return errors.DataError("alias column is mandatory")
	}

	for i, row := range aliasTable.Rows[1:] {
		a := row.Cells[0].Value
		value := table.Rows[i+1].Cells[0].Value
		ok := sc.aliases.Set(value, sc.Pickle.Id, a)
		if !ok {
			return errors.DataError("could not register alias")
		}
	}
	return nil
}

func (sc *ScenarioContext) iTrackTheFollowingEnvelopes(table *gherkin.PickleStepArgument_PickleTable) error {
	if len(table.Rows[0].Cells) != 1 {
		return errors.DataError("invalid table")
	}

	var envelopes []*tx.Envelope
	for _, r := range table.Rows[1:] {
		if r.Cells[0].Value != "" {
			envelopes = append(envelopes, tx.NewEnvelope().SetID(r.Cells[0].Value))
		}
	}
	sc.setTrackers(sc.newTrackers(envelopes))

	return nil
}

func (sc *ScenarioContext) iSignTheFollowingTransactions(table *gherkin.PickleStepArgument_PickleTable) error {
	tenantCol := utils.ExtractColumns(table, []string{"Tenant"})
	apiKeyCol := utils.ExtractColumns(table, []string{"API-KEY"})

	helpersColumns := []string{aliasHeaderValue, "privateKey"}
	helpersTable := utils.ExtractColumns(table, helpersColumns)
	if helpersTable == nil {
		return errors.DataError("One of the following columns is missing %q", helpersColumns)
	}

	envelopes, err := utils.ParseEnvelope(table)
	if err != nil {
		return err
	}

	// Sign tx for each envelope
	ctx := utils2.RetryConnectionError(context.Background(), true)
	for i, e := range envelopes {
		apiKey := apiKeyCol.Rows[i+1].Cells[0].Value

		tenant := ""
		if tenantCol != nil {
			tenant = tenantCol.Rows[i+1].Cells[0].Value
		}

		headers := utils.GetHeaders(apiKey, tenant, "")
		err = sc.craftAndSignEnvelope(
			context.WithValue(ctx, clientutils.RequestHeaderKey, headers),
			e,
			helpersTable.Rows[i+1].Cells[1].Value,
		)
		if err != nil {
			return err
		}
		sc.aliases.Set(e, sc.Pickle.Id, helpersTable.Rows[i+1].Cells[0].Value)
	}

	return nil
}

func (sc *ScenarioContext) iHaveTheFollowingJWTTokens(table *gherkin.PickleStepArgument_PickleTable) error {
	headers := table.Rows[0]
	for _, row := range table.Rows[1:] {
		tenantMap := make(map[string]interface{})
		var a string
		var audience string

		for i, cell := range row.Cells {
			switch v := headers.Cells[i].Value; {
			case v == aliasHeaderValue:
				a = cell.Value
			case v == "audience":
				audience = cell.Value
			default:
				tenantMap[v] = cell.Value
			}
		}
		if a == "" {
			return errors.DataError("need an alias")
		}
		if audience == "" {
			return errors.DataError("need an audience")
		}

		jwtToken, err := sc.getJWT(audience)
		if err != nil {
			return err
		}

		tenantMap["token"] = fmt.Sprintf("Bearer %s", jwtToken)
		sc.aliases.Set(tenantMap, sc.Pickle.Id, a)
	}

	return nil
}

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (sc *ScenarioContext) getJWT(audience string) (string, error) {
	clientID, ok := sc.aliases.Get(alias.GlobalAka, "oidc.clientID")
	if !ok {
		return "", errors.DataError("Could not find oidc client ID")
	}

	clientSecret, ok := sc.aliases.Get(alias.GlobalAka, "oidc.clientSecret")
	if !ok {
		return "", errors.DataError("Could not find oidc client secret")
	}

	idpURL, ok := sc.aliases.Get(alias.GlobalAka, "oidc.tokenURL")
	if !ok {
		return "", errors.DataError("Could not find token url")
	}

	body := new(bytes.Buffer)
	_ = json2.NewEncoder(body).Encode(map[string]interface{}{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"audience":      audience,
		"grant_type":    "client_credentials",
	})

	resp, err := http.DefaultClient.Post(idpURL.(string), jsonrpc.ContentType, body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	acessToken := &accessTokenResponse{}
	if resp.StatusCode == http.StatusOK {
		err = json2.NewDecoder(resp.Body).Decode(acessToken)
		if err != nil {
			return "", err
		}

		return acessToken.AccessToken, nil
	}

	// Read body
	respMsg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return "", fmt.Errorf(string(respMsg))
}

func (sc *ScenarioContext) craftAndSignEnvelope(ctx context.Context, e *tx.Envelope, privKey string) error {
	chainRegistry, ok := sc.aliases.Get(alias.GlobalAka, "api")
	if !ok {
		return errors.DataError("Could not find the api endpoint")
	}
	endpoint := utils4.GetProxyURL(chainRegistry.(string), e.GetChainUUID())
	if e.GetChainID() == nil && e.GetChainUUID() != "" {
		chainID, errNetwork := sc.ec.Network(utils2.RetryConnectionError(ctx, true), endpoint)
		if errNetwork != nil {
			log.WithError(errNetwork).Error("failed to get chain ID")
			return errNetwork
		}
		_ = e.SetChainID(chainID)
	}

	signer := pkgcryto.GetEIP155Signer(e.GetChainIDString())
	acc, err := crypto.HexToECDSA(privKey)
	if err != nil {
		log.WithError(err).WithField("private_key", privKey).Error("failed to create account using private key")
		return err
	}

	_ = e.SetFrom(crypto.PubkeyToAddress(acc.PublicKey))
	if e.GetGasPrice() == nil {
		gasPrice, errGasPrice := sc.ec.SuggestGasPrice(ctx, endpoint)
		if errGasPrice != nil {
			log.WithError(errGasPrice).Error("failed to suggest gas price")
			return errGasPrice
		}
		_ = e.SetGasPrice(gasPrice)
	}

	transaction, err := e.GetTransaction()
	if err != nil {
		log.WithError(err).Error("failed to get transaction from envelope")
		return err
	}

	signature, err := pkgcryto.SignTransaction(transaction, acc, signer)
	if err != nil {
		log.WithError(err).Error("failed to sign transaction")
		return err
	}

	signedTx, err := transaction.WithSignature(signer, signature)
	if err != nil {
		log.WithError(err).Error("failed to set signature in transaction")
		return err
	}

	signedRaw, err := rlp.Encode(signedTx)
	if err != nil {
		log.WithError(err).Error("failed to RLP encode signed transaction")
		return err
	}

	_ = e.SetRaw(signedRaw).SetTxHash(signedTx.Hash())
	return nil
}

func initEnvelopeSteps(s *godog.ScenarioContext, sc *ScenarioContext) {
	s.Step(`^I register the following chains$`, sc.preProcessTableStep(sc.iRegisterTheFollowingChains))
	s.Step(`^I register the following faucets$`, sc.preProcessTableStep(sc.iRegisterTheFollowingFaucets))
	s.Step(`^I have the following tenants$`, sc.preProcessTableStep(sc.iHaveTheFollowingTenant))
	s.Step(`^I have the following account`, sc.preProcessTableStep(sc.iHaveTheFollowingAccount))
	s.Step(`^I register the following alias$`, sc.preProcessTableStep(sc.iRegisterTheFollowingAliasAs))
	s.Step(`^I have created the following accounts$`, sc.preProcessTableStep(sc.iHaveCreatedTheFollowingAccounts))
	s.Step(`^I track the following envelopes$`, sc.preProcessTableStep(sc.iTrackTheFollowingEnvelopes))
	s.Step(`^I send envelopes to topic "([^"]*)"$`, sc.iSendEnvelopesToTopic)
	s.Step(`^Register new envelope tracker "([^"]*)"$`, sc.registerEnvelopeTracker)
	s.Step(`^Envelopes should be in topic "([^"]*)"$`, sc.envelopeShouldBeInTopic)
	s.Step(`^Envelopes should have the following fields$`, sc.preProcessTableStep(sc.envelopesShouldHaveTheFollowingValues))
	s.Step(`^I register the following envelope fields$`, sc.preProcessTableStep(sc.iRegisterTheFollowingEnvelopeFields))
	s.Step(`^I sign the following transactions$`, sc.preProcessTableStep(sc.iSignTheFollowingTransactions))
	s.Step(`^I have the following jwt tokens$`, sc.preProcessTableStep(sc.iHaveTheFollowingJWTTokens))
}
