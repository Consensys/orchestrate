package migrations

import (
	"github.com/go-pg/migrations"
	log "github.com/sirupsen/logrus"
)

func addTriggerOnEnvelopeUpdate(db migrations.DB) error {
	log.Debugf("Add trigger envelope_trig")

	// Remark: you will note that we consider that chain ID should be max a uint256
	_, err := db.Exec(
		`CREATE OR REPLACE FUNCTION envelope_updated() RETURNS TRIGGER AS 
	$$
	BEGIN
		NEW.status = 'stored';
		NEW.stored_at = (now() at time zone 'utc');
		NEW.error_at = NULL;
		NEW.sent_at = NULL;
		NEW.mined_at = NULL;
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

CREATE TRIGGER envelope_trig 
	BEFORE UPDATE OF envelope ON envelopes
	FOR EACH ROW 
	EXECUTE PROCEDURE envelope_updated();`,
	)

	if err != nil {
		return err
	}

	log.Infof("Added trigger envelope_trig")

	return nil
}

func removeTriggerOnEnvelopeUpdate(db migrations.DB) error {
	log.Debugf("Removing trigger envelope_trig")

	_, err := db.Exec(
		`DROP TRIGGER envelope_trig ON envelopes;
DROP FUNCTION envelope_updated();`,
	)
	if err != nil {
		return err
	}

	log.Infof("Removed trigger envelope_trig")

	return nil
}

func init() {
	Collection.MustRegisterTx(addTriggerOnEnvelopeUpdate, removeTriggerOnEnvelopeUpdate)
}
