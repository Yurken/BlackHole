package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"main/models"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Init 初始化数据库
func Init() error {
	homeDir, _ := os.UserHomeDir()
	dbDir := filepath.Join(homeDir, ".blackhole")
	os.MkdirAll(dbDir, 0755)

	dbPath := filepath.Join(dbDir, "history.db")

	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	// 创建历史记录表
	createTable := `
	CREATE TABLE IF NOT EXISTS history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		original_path TEXT NOT NULL,
		original_name TEXT NOT NULL,
		new_path TEXT NOT NULL,
		new_name TEXT NOT NULL,
		rule_name TEXT,
		action TEXT NOT NULL,
		status TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_timestamp ON history(timestamp DESC);
	CREATE TABLE IF NOT EXISTS rules (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		icon TEXT,
		color TEXT,
		destination TEXT,
		action TEXT,
		keep_original INTEGER,
		file_types TEXT,
		custom_extensions TEXT,
		allow_all_files INTEGER,
		name_template TEXT,
		date_source TEXT,
		ai_enabled INTEGER,
		quick_access INTEGER,
		enabled INTEGER,
		created_at DATETIME,
		updated_at DATETIME
	);
	`

	_, err = DB.Exec(createTable)
	if err != nil {
		return err
	}

	log.Println("📊 Database initialized:", dbPath)
	return nil
}

// SaveHistory 保存历史记录
func SaveHistory(originalPath, originalName, newPath, newName, ruleName, action, status string) error {
	stmt, err := DB.Prepare(`
		INSERT INTO history (original_path, original_name, new_path, new_name, rule_name, action, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(originalPath, originalName, newPath, newName, ruleName, action, status)
	return err
}

// GetHistory 获取历史记录
func GetHistory() ([]models.HistoryRecord, error) {
	query := `
		SELECT id, original_path, original_name, new_path, new_name, rule_name, action, status, 
		       strftime('%Y-%m-%d %H:%M:%S', timestamp) as timestamp
		FROM history
		ORDER BY timestamp DESC
		LIMIT 100
	`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.HistoryRecord
	for rows.Next() {
		var record models.HistoryRecord
		err := rows.Scan(
			&record.ID,
			&record.OriginalPath,
			&record.OriginalName,
			&record.NewPath,
			&record.NewName,
			&record.RuleName,
			&record.Action,
			&record.Status,
			&record.Timestamp,
		)
		if err != nil {
			continue
		}
		records = append(records, record)
	}

	return records, nil
}

// ClearHistory 清除历史记录
func ClearHistory() error {
	_, err := DB.Exec("DELETE FROM history")
	return err
}
