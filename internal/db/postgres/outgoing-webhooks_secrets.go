package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/pkg/secret"
)

// ManageOutgoingWebhookSecrets manage protection of outgoing webhook secrets
func (pg *DB) ManageOutgoingWebhookSecrets(ctx context.Context, provider secret.EngineProvider, action secret.Action) (uint, error) {
	columns := []string{"id", "user_id", "secrets"}
	query, args, _ := pg.psql.Select(columns...).From(
		"outgoing_webhooks",
	).ToSql()
	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var count uint
	for rows.Next() {
		var id, uid uint
		var secrets secret.Secrets
		if err = rows.Scan(&id, &uid, &secrets); err != nil {
			break
		}
		if err = provider.Apply(action, &secrets); err != nil {
			break
		}
		if err = pg.updateOutgoingWebhookSecrets(&model.OutgoingWebhook{ID: &id, UserID: &uid}, &secrets); err != nil {
			break
		}
		count++
	}

	return count, err
}

func (pg *DB) updateOutgoingWebhookSecrets(webhook *model.OutgoingWebhook, secrets *secret.Secrets) error {
	for k, v := range *secrets {
		// if secret's value is not provided use the current value
		if current, ok := webhook.Secrets[k]; ok && v == "" {
			(*secrets)[k] = current
		}
	}
	query, args, _ := pg.psql.Update(
		"outgoing_webhooks",
	).Set("secrets", *secrets).Where(
		sq.Eq{"id": webhook.ID},
	).Where(
		sq.Eq{"user_id": webhook.UserID},
	).ToSql()

	_, err := pg.db.Exec(query, args...)
	webhook.Secrets = *secrets
	return err
}
