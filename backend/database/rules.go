package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"main/models"
)

func CreateRule(rule models.Rule) (models.Rule, error) {
	if rule.ID == "" {
		rule.ID = fmt.Sprintf("rule_%d", time.Now().UnixNano())
	}
	now := time.Now().Format(time.RFC3339)
	rule.CreatedAt = now
	rule.UpdatedAt = now

	_, err := DB.Exec(`
		INSERT INTO rules (
			id, name, icon, color, destination, action, keep_original, file_types,
			custom_extensions, allow_all_files, name_template, date_source,
			ai_enabled, quick_access, enabled, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		rule.ID,
		rule.Name,
		rule.Icon,
		rule.Color,
		rule.Destination,
		rule.Action,
		boolToInt(rule.KeepOriginal),
		marshalStringSlice(rule.FileTypes),
		marshalStringSlice(rule.CustomExtensions),
		boolToInt(rule.AllowAllFiles),
		marshalStringSlice(rule.NameTemplate),
		rule.DateSource,
		boolToInt(rule.AIEnabled),
		boolToInt(rule.QuickAccess),
		boolToInt(rule.Enabled),
		rule.CreatedAt,
		rule.UpdatedAt,
	)
	if err != nil {
		return models.Rule{}, err
	}

	return rule, nil
}

func UpdateRule(rule models.Rule) (models.Rule, error) {
	rule.UpdatedAt = time.Now().Format(time.RFC3339)

	result, err := DB.Exec(`
		UPDATE rules SET
			name = ?,
			icon = ?,
			color = ?,
			destination = ?,
			action = ?,
			keep_original = ?,
			file_types = ?,
			custom_extensions = ?,
			allow_all_files = ?,
			name_template = ?,
			date_source = ?,
			ai_enabled = ?,
			quick_access = ?,
			enabled = ?,
			updated_at = ?
		WHERE id = ?
	`,
		rule.Name,
		rule.Icon,
		rule.Color,
		rule.Destination,
		rule.Action,
		boolToInt(rule.KeepOriginal),
		marshalStringSlice(rule.FileTypes),
		marshalStringSlice(rule.CustomExtensions),
		boolToInt(rule.AllowAllFiles),
		marshalStringSlice(rule.NameTemplate),
		rule.DateSource,
		boolToInt(rule.AIEnabled),
		boolToInt(rule.QuickAccess),
		boolToInt(rule.Enabled),
		rule.UpdatedAt,
		rule.ID,
	)
	if err != nil {
		return models.Rule{}, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return models.Rule{}, err
	}
	if rows == 0 {
		return models.Rule{}, sql.ErrNoRows
	}

	return rule, nil
}

func DeleteRule(id string) error {
	result, err := DB.Exec(`DELETE FROM rules WHERE id = ?`, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func GetRule(id string) (models.Rule, error) {
	row := DB.QueryRow(`
		SELECT id, name, icon, color, destination, action, keep_original, file_types,
		       custom_extensions, allow_all_files, name_template, date_source,
		       ai_enabled, quick_access, enabled, created_at, updated_at
		FROM rules
		WHERE id = ?
	`, id)

	return scanRule(row)
}

func GetRules() ([]models.Rule, error) {
	rows, err := DB.Query(`
		SELECT id, name, icon, color, destination, action, keep_original, file_types,
		       custom_extensions, allow_all_files, name_template, date_source,
		       ai_enabled, quick_access, enabled, created_at, updated_at
		FROM rules
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []models.Rule
	for rows.Next() {
		rule, err := scanRule(rows)
		if err != nil {
			continue
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

func scanRule(scanner interface {
	Scan(dest ...interface{}) error
}) (models.Rule, error) {
	var rule models.Rule
	var keepOriginal int
	var allowAllFiles int
	var aiEnabled int
	var quickAccess int
	var enabled int
	var fileTypes string
	var customExtensions string
	var nameTemplate string

	err := scanner.Scan(
		&rule.ID,
		&rule.Name,
		&rule.Icon,
		&rule.Color,
		&rule.Destination,
		&rule.Action,
		&keepOriginal,
		&fileTypes,
		&customExtensions,
		&allowAllFiles,
		&nameTemplate,
		&rule.DateSource,
		&aiEnabled,
		&quickAccess,
		&enabled,
		&rule.CreatedAt,
		&rule.UpdatedAt,
	)
	if err != nil {
		return models.Rule{}, err
	}

	rule.KeepOriginal = keepOriginal == 1
	rule.AllowAllFiles = allowAllFiles == 1
	rule.AIEnabled = aiEnabled == 1
	rule.QuickAccess = quickAccess == 1
	rule.Enabled = enabled == 1
	rule.FileTypes = unmarshalStringSlice(fileTypes)
	rule.CustomExtensions = unmarshalStringSlice(customExtensions)
	rule.NameTemplate = unmarshalStringSlice(nameTemplate)

	return rule, nil
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func marshalStringSlice(items []string) string {
	if len(items) == 0 {
		return "[]"
	}
	data, err := json.Marshal(items)
	if err != nil {
		return "[]"
	}
	return string(data)
}

func unmarshalStringSlice(data string) []string {
	if data == "" {
		return []string{}
	}
	var items []string
	if err := json.Unmarshal([]byte(data), &items); err != nil {
		return []string{}
	}
	return items
}
