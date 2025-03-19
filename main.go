package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Marker struct {
	ID        int     `json:"id"`         // マーカーID
	UserID    int     `json:"user_id"`    // ユーザーID
	Lat       float64 `json:"lat"`        // 緯度
	Lng       float64 `json:"lng"`        // 経度
	CreatedAt string  `json:"created_at"` // 作成日時
	UpdatedAt string  `json:"updated_at"` // 更新日時
}

func main() {
	db, err := sql.Open("mysql", "root:fybrid@tcp(127.0.0.1:3306)/jimop")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/add-marker", withCORS(addMarkerHandler(db)))
	http.HandleFunc("/get-markers", withCORS(getMarkersHandler(db)))

	fs := http.FileServer(http.Dir("../"))
	http.Handle("/", fs)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func addMarkerHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var marker Marker
		err := json.NewDecoder(r.Body).Decode(&marker)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("INSERT INTO markers (user_id, lat, lng) VALUES (?, ?, ?)", marker.UserID, marker.Lat, marker.Lng)
		if err != nil {
			http.Error(w, "Failed to insert marker", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(marker); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func getMarkersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, user_id, lat, lng, created_at, updated_at FROM markers")
		if err != nil {
			http.Error(w, "Failed to fetch markers", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var markers []Marker
		for rows.Next() {
			var marker Marker
			if err := rows.Scan(&marker.ID, &marker.UserID, &marker.Lat, &marker.Lng, &marker.CreatedAt, &marker.UpdatedAt); err != nil {
				http.Error(w, "Failed to scan marker", http.StatusInternalServerError)
				return
			}
			markers = append(markers, marker)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Failed to iterate markers", http.StatusInternalServerError)
			return
		}

		// 取得したレコードをログに出力
		log.Printf("Fetched markers: %+v\n", markers)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(markers); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}
