package migration

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type MigrationRecord struct {
	ID        int    `db:"id"`
	Filename  string `db:"filename"`
	AppliedAt string `db:"applied_at"`
}

type Migrator struct {
	DB            *sqlx.DB
	MigrationPath string
}

func NewMigrator(db *sqlx.DB, migrationPath string) *Migrator {
	return &Migrator{
		DB:            db,
		MigrationPath: migrationPath,
	}
}

func (m *Migrator) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			filename VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`

	_, err := m.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	return nil
}

func (m *Migrator) getAppliedMigrations() (map[string]bool, error) {
	applied := make(map[string]bool)

	var records []MigrationRecord
	err := m.DB.Select(&records, "SELECT filename FROM migrations ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("failed to get applied migrations: %w", err)
	}

	for _, record := range records {
		applied[record.Filename] = true
	}

	return applied, nil
}

func (m *Migrator) getSQLFiles() ([]string, error) {
	var files []string

	err := filepath.WalkDir(m.MigrationPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".sql") {
			files = append(files, d.Name())
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to read migration files: %w", err)
	}

	sort.Slice(files, func(i, j int) bool {
		numI := extractNumber(files[i])
		numJ := extractNumber(files[j])
		return numI < numJ
	})

	return files, nil
}

func extractNumber(filename string) int {

	name := strings.TrimSuffix(filename, ".sql")

	var numStr strings.Builder
	for _, char := range name {
		if char >= '0' && char <= '9' {
			numStr.WriteRune(char)
		} else {
			break
		}
	}

	if numStr.Len() == 0 {
		return 0
	}

	num, err := strconv.Atoi(numStr.String())
	if err != nil {
		return 0
	}

	return num
}

func (m *Migrator) executeSQLFile(filename string) error {
	filePath := filepath.Join(m.MigrationPath, filename)

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", filename, err)
	}

	statements := strings.Split(string(content), ";")

	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		_, err := m.DB.Exec(statement)
		if err != nil {
			return fmt.Errorf("failed to execute statement in %s: %w\nStatement: %s", filename, err, statement)
		}
	}

	_, err = m.DB.Exec("INSERT INTO migrations (filename) VALUES ($1)", filename)
	if err != nil {
		return fmt.Errorf("failed to record migration %s: %w", filename, err)
	}

	return nil
}

func (m *Migrator) RunMigrations() error {
	log.Println("Starting database migrations...")

	if err := m.createMigrationsTable(); err != nil {
		return err
	}

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return err
	}

	files, err := m.getSQLFiles()
	if err != nil {
		return err
	}

	if len(files) == 0 {
		log.Println("No migration files found")
		return nil
	}

	migrationsRun := 0
	for _, filename := range files {
		if applied[filename] {
			log.Printf("Migration %s already applied, skipping", filename)
			continue
		}

		log.Printf("Running migration: %s", filename)
		if err := m.executeSQLFile(filename); err != nil {
			return fmt.Errorf("migration %s failed: %w", filename, err)
		}

		log.Printf("Migration %s completed successfully", filename)
		migrationsRun++
	}

	if migrationsRun == 0 {
		log.Println("All migrations are up to date")
	} else {
		log.Printf("Successfully applied %d migrations", migrationsRun)
	}

	return nil
}

func AutoMigrate(db *sqlx.DB, migrationPath string) error {

	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		log.Printf("Migration path %s does not exist, skipping migrations", migrationPath)
		return nil
	}

	migrator := NewMigrator(db, migrationPath)
	return migrator.RunMigrations()
}
