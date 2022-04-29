package migrations

import (
	"github.com/go-pg/migrations/v7"
	log "github.com/sirupsen/logrus"
)

func upgradeTxRequestValues(db migrations.DB) error {
	log.Debug("Applying migration of tx request params...")
	_, err := db.Exec(`
UPDATE transaction_requests
	SET params = jsonb_set(params, '{gas}', to_jsonb(TO_NUMBER((params->'gas')::text, '9G999g999')), false)
	WHERE jsonb_typeof(params->'gas') = 'string';

UPDATE transaction_requests
	SET params = jsonb_set(params, '{nonce}', to_jsonb(TO_NUMBER((params->'nonce')::text, '9G999g999')), false)
	WHERE jsonb_typeof(params->'nonce') = 'string';

UPDATE jobs
	SET internal_data = jsonb_set(internal_data, '{chainID}', to_jsonb(TO_NUMBER((internal_data->'chainID')::text, '9G999g999')), false)
	WHERE jsonb_typeof(internal_data->'chainID') = 'string';
`)
	if err != nil {
		return err
	}
	log.Info("Applied migration of tx request params")

	return nil
}

func downgradeTxRequestValues(db migrations.DB) error {
	log.Debug("Undoing migration of tx request params")
	_, err := db.Exec(`
UPDATE transaction_requests
	SET params = jsonb_set(params, '{gas}', to_jsonb(CAST((params->'gas')::text as text)), false)
	WHERE jsonb_typeof(params->'gas') = 'number';

UPDATE transaction_requests
	SET params = jsonb_set(params, '{nonce}', to_jsonb(CAST((params->'nonce')::text as text)), false)
	WHERE jsonb_typeof(params->'nonce') = 'number';

UPDATE transaction_requests
	SET params = jsonb_set(internal_data, '{chainID}', to_jsonb(CAST((internal_data->'chainID')::text as text)), false)
	WHERE jsonb_typeof(internal_data->'chainID') = 'string';
`)
	if err != nil {
		return err
	}
	log.Info("Downgraded migration of tx request params")

	return nil
}

func init() {
	Collection.MustRegisterTx(upgradeTxRequestValues, downgradeTxRequestValues)
}
