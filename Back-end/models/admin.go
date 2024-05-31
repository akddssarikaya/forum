package models

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func HandleAdmin(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	tmpl, ok := tmplCache["panel"]
	if !ok {
		http.Error(w, "Could not load panel template", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		userIdStr := r.FormValue("userId")
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		if err := deletetable(db, userId); err != nil {
			http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
			log.Println("Failed to delete user:", err)
			return
		}

		fmt.Fprintf(w, "User with ID %d deleted successfully", userId)
	} else {
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func deletetable(database *sql.DB, userId int) error {
	// Transaction başlat
	tx, err := database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Kullanıcıyı sil
	deleteStmt := "DELETE FROM users WHERE id = ?;"
	if _, err := tx.Exec(deleteStmt, userId); err != nil {
		return err
	}

	// Tabloda kalan satır sayısını kontrol et
	rowCount := 0
	countStmt := "SELECT COUNT(*) FROM users;"
	if err := tx.QueryRow(countStmt).Scan(&rowCount); err != nil {
		return err
	}

	// Eğer tablo boş ise, otomatik artan değeri sıfırla
	if rowCount == 0 {
		resetAIStmt := "DELETE FROM SQLITE_SEQUENCE WHERE NAME = 'users';"
		if _, err := tx.Exec(resetAIStmt); err != nil {
			return err
		}
	}

	// Transaction'ı commit et
	return tx.Commit()
}
