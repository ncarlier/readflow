package postgres

import (
	"database/sql"
	"errors"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/ncarlier/readflow/pkg/model"
)

var ruleColumns = []string{
	"id",
	"user_id",
	"alias",
	"category_id",
	"priority",
	"rule",
	"created_at",
	"updated_at",
}

func mapRowToRule(row *sql.Row) (*model.Rule, error) {
	rule := &model.Rule{}

	err := row.Scan(
		&rule.ID,
		&rule.UserID,
		&rule.Alias,
		&rule.CategoryID,
		&rule.Priority,
		&rule.Rule,
		&rule.CreatedAt,
		&rule.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return rule, nil
}

func (pg *DB) createRule(rule model.Rule) (*model.Rule, error) {
	query, args, _ := pg.psql.Insert(
		"rules",
	).Columns(
		"user_id", "alias", "category_id", "priority", "rule",
	).Values(
		rule.UserID,
		rule.Alias,
		rule.CategoryID,
		rule.Priority,
		rule.Rule,
	).Suffix(
		"RETURNING " + strings.Join(ruleColumns, ","),
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	return mapRowToRule(row)
}

func (pg *DB) updateRule(rule model.Rule) (*model.Rule, error) {
	update := map[string]interface{}{
		"alias":       rule.Alias,
		"category_id": rule.CategoryID,
		"priority":    rule.Priority,
		"rule":        rule.Rule,
		"updated_at":  "NOW()",
	}
	query, args, _ := pg.psql.Update(
		"rules",
	).SetMap(update).Where(
		sq.Eq{"id": rule.ID},
	).Where(
		sq.Eq{"user_id": rule.UserID},
	).Suffix(
		"RETURNING " + strings.Join(ruleColumns, ","),
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	return mapRowToRule(row)
}

// CreateOrUpdateRule creates or updates a rule into the DB
func (pg *DB) CreateOrUpdateRule(rule model.Rule) (*model.Rule, error) {
	if rule.ID != nil {
		return pg.updateRule(rule)
	}
	return pg.createRule(rule)
}

// GetRuleByID get a rule from the DB
func (pg *DB) GetRuleByID(id uint) (*model.Rule, error) {
	query, args, _ := pg.psql.Select(ruleColumns...).From(
		"rules",
	).Where(
		sq.Eq{"id": id},
	).ToSql()
	row := pg.db.QueryRow(query, args...)
	return mapRowToRule(row)
}

// GetRuleByUserIDAndAlias get a rule from the DB.
func (pg *DB) GetRuleByUserIDAndAlias(uid uint, alias string) (*model.Rule, error) {
	query, args, _ := pg.psql.Select(ruleColumns...).From(
		"rules",
	).Where(
		sq.Eq{"user_id": uid},
	).Where(
		sq.Eq{"alias": alias},
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	return mapRowToRule(row)
}

// GetRulesByUserID returns rules of an user from DB
func (pg *DB) GetRulesByUserID(uid uint) ([]model.Rule, error) {
	query, args, _ := pg.psql.Select(ruleColumns...).From(
		"rules",
	).Where(
		sq.Eq{"user_id": uid},
	).OrderBy("priority DESC").ToSql()
	rows, err := pg.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Rule

	for rows.Next() {
		rule := model.Rule{}
		err = rows.Scan(
			&rule.ID,
			&rule.UserID,
			&rule.Alias,
			&rule.CategoryID,
			&rule.Priority,
			&rule.Rule,
			&rule.CreatedAt,
			&rule.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, rule)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteRule removes a rule from the DB
func (pg *DB) DeleteRule(rule model.Rule) error {
	query, args, _ := pg.psql.Delete("rules").Where(
		sq.Eq{"id": rule.ID},
	).ToSql()
	result, err := pg.db.Exec(query, args...)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("no rule has been removed")
	}

	return nil
}

// DeleteRules removes rules from the DB
func (pg *DB) DeleteRules(uid uint, ids []uint) (int64, error) {
	query, args, _ := pg.psql.Delete("rules").Where(
		sq.Eq{"user_id": uid},
	).Where(
		sq.Eq{"id": ids},
	).ToSql()
	result, err := pg.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
