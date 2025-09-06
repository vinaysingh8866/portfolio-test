package repo

import (
	"bufio"
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ApplyMigrations runs all .sql files in the given directory in lexical order.
func ApplyMigrations(db *sql.DB, migrationsDir string) error {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, filepath.Join(migrationsDir, e.Name()))
		}
	}
	sort.Strings(files)
	for _, f := range files {
		if err := runSQLFile(db, f); err != nil {
			return err
		}
	}
	return nil
}

func runSQLFile(db *sql.DB, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	// Very simple splitter on ';' at line ends to keep it robust across files.
	var stmt strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		stmt.WriteString(line)
		stmt.WriteString("\n")
		if strings.HasSuffix(strings.TrimSpace(line), ";") {
			if _, err := db.Exec(stmt.String()); err != nil {
				return err
			}
			stmt.Reset()
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if strings.TrimSpace(stmt.String()) != "" {
		_, err := db.Exec(stmt.String())
		if err != nil {
			return err
		}
	}
	return nil
}
