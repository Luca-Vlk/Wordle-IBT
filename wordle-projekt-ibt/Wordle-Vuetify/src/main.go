package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type GameEndRequest struct {
	Name     string `json:"name"`
	Versuche string `json:"versuche"`
	Wort     string `json:"wort"`
}

func main() {
	fmt.Println("main.go wird ausgeführt")
	if checkDB() {
		http.HandleFunc("/api/insert", methodInsert)
		http.HandleFunc("/api/selectNames", methodSelectNames)
		http.HandleFunc("/api/selectWinCount", selectWinGameCount)
		http.HandleFunc("/api/selectGameCount", selectGameCount)
		http.HandleFunc("/api/selectLooseCount", selectLooseCount)
		http.HandleFunc("/api/selectTryCount", selectTryCount)
		http.HandleFunc("/api/selectAVGTry", selectAVGTry)
	}

	http.HandleFunc("/", setup)

	fmt.Println("Server läuft auf http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func checkDB() bool {
	db, err := sql.Open("mysql", "root:1234@/wordledb")
	connWork := true

	if err != nil {
		connWork = false
	}
	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		connWork = false
	}

	db.Close()
	return connWork
}

func setup(w http.ResponseWriter, r *http.Request) {
	// Oberfläche bauen
	fs := http.FileServer(http.Dir("./dist"))

	// api calls abfangen
	if strings.HasPrefix(r.URL.Path, "/api/") {
		http.NotFound(w, r)
		return
	}

	// Datei wird gesucht
	if _, err := http.Dir("./dist").Open(r.URL.Path); err != nil {
		// Datei nicht gefunden -> index.html (für SPA-Routing)
		http.ServeFile(w, r, "./dist/index.html")
		return
	}

	// wird geliefert
	fs.ServeHTTP(w, r)
}

func methodInsert(w http.ResponseWriter, r *http.Request) {
	//Verbingungsaufbau
	db, err := sql.Open("mysql", "root:1234@/wordledb")

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	var req GameEndRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Fatal(err)
	}

	lVersuche, err := strconv.Atoi(req.Versuche)
	if err != nil {
		log.Fatal(err)
	}

	//Dinge in DB hinzufügen
	result, err := db.ExecContext(ctx,
		"INSERT INTO data (name, versuche, wort) VALUES (?, ?, ?)",
		req.Name,
		lVersuche,
		req.Wort,
	)

	lastInsertId, err := result.LastInsertId()
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Last inserted ID: %d \n", lastInsertId)
	fmt.Printf("Rows affected: %d\n", rowsAffected)

	//schließen der verbindung
	db.Close()
	w.WriteHeader(http.StatusOK)
}

func methodSelectNames(w http.ResponseWriter, r *http.Request) {
	//Verbingungsaufbau
	db, err := sql.Open("mysql", "root:1234@/wordledb")

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	//Abrufen von DB dingen
	rows, err := db.QueryContext(ctx, "SELECT name FROM data GROUP BY name")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	namen := ""
	//Ausgabe der erhaltenen Rows
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		if namen == "" {
			namen = name
		} else {
			namen += "," + name
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// CORS-Header für lokale Entwicklung
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/plain")

	//Schreibt den string namen in die response
	fmt.Fprint(w, namen)

	//schließen der verbindung
	db.Close()
	w.WriteHeader(http.StatusOK)
}

func selectWinGameCount(w http.ResponseWriter, r *http.Request) {
	//Verbingungsaufbau
	db, err := sql.Open("mysql", "root:1234@/wordledb")

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	userName := "'" + r.URL.Query().Get("name") + "'"

	//Abrufen von DB dingen
	rows, err := db.QueryContext(ctx, "SELECT COUNT(*) AS count FROM data WHERE Versuche != 0 AND name="+userName)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var count string
	//Ausgabe der erhaltenen Rows
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Fatal(err)
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// CORS-Header für lokale Entwicklung
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/plain")

	//Schreibt den string namen in die response
	fmt.Fprint(w, count)

	//schließen der verbindung
	db.Close()
	w.WriteHeader(http.StatusOK)
}

func selectGameCount(w http.ResponseWriter, r *http.Request) {
	//Verbingungsaufbau
	db, err := sql.Open("mysql", "root:1234@/wordledb")

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	userName := "'" + r.URL.Query().Get("name") + "'"

	//Abrufen von DB dingen
	rows, err := db.QueryContext(ctx, "SELECT COUNT(*) AS count FROM data WHERE name="+userName)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var count string
	//Ausgabe der erhaltenen Rows
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Fatal(err)
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// CORS-Header für lokale Entwicklung
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/plain")

	//Schreibt den string namen in die response
	fmt.Fprint(w, count)

	//schließen der verbindung
	db.Close()
	w.WriteHeader(http.StatusOK)
}
func selectLooseCount(w http.ResponseWriter, r *http.Request) {
	//Verbingungsaufbau
	db, err := sql.Open("mysql", "root:1234@/wordledb")

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	userName := "'" + r.URL.Query().Get("name") + "'"

	//Abrufen von DB dingen
	rows, err := db.QueryContext(ctx, "SELECT COUNT(*) AS count FROM data WHERE Versuche = 0 AND name="+userName)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var count string
	//Ausgabe der erhaltenen Rows
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Fatal(err)
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// CORS-Header für lokale Entwicklung
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/plain")

	//Schreibt den string namen in die response
	fmt.Fprint(w, count)

	//schließen der verbindung
	db.Close()
	w.WriteHeader(http.StatusOK)
}
func selectAVGTry(w http.ResponseWriter, r *http.Request) {
	//Verbingungsaufbau
	db, err := sql.Open("mysql", "root:1234@/wordledb")

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	userName := "'" + r.URL.Query().Get("name") + "'"

	nameExists := false
	row, err := db.QueryContext(ctx, "SELECT name FROM data WHERE name = "+userName)
	for row.Next() {
		var name string
		if err := row.Scan(&name); err != nil {
			log.Fatal(err)
		}
		if name != "" {
			nameExists = true
		}
	}

	var avg string
	if nameExists {
		//Abrufen von DB dingen
		rows, err := db.QueryContext(ctx, "SELECT AVG(Versuche) AS avg FROM data WHERE Versuche != 0 AND name="+userName)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		//Ausgabe der erhaltenen Rows
		for rows.Next() {
			if err := rows.Scan(&avg); err != nil {
				log.Fatal(err)
			}
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
	} else {
		avg = ""
	}
	// CORS-Header für lokale Entwicklung
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/plain")

	//Schreibt den string namen in die response
	fmt.Fprint(w, avg)

	//schließen der verbindung
	db.Close()
	w.WriteHeader(http.StatusOK)
}
func selectTryCount(w http.ResponseWriter, r *http.Request) {
	//Verbingungsaufbau
	db, err := sql.Open("mysql", "root:1234@/wordledb")

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	userName := "'" + r.URL.Query().Get("name") + "'"

	//Abrufen von DB dingen
	rows, err := db.QueryContext(ctx, "SELECT Versuche, COUNT(*) AS count FROM data WHERE Versuche != 0 AND name="+userName+" GROUP BY Versuche ORDER BY Versuche")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var count string
	var versuche string
	resp := ""
	//Ausgabe der erhaltenen Rows
	for rows.Next() {
		if err := rows.Scan(&versuche, &count); err != nil {
			log.Fatal(err)
		}
		if resp == "" {
			resp = versuche + "," + count
		} else {
			resp += "," + versuche + "," + count
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// CORS-Header für lokale Entwicklung
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/plain")

	//Schreibt den string namen in die response
	fmt.Fprint(w, resp)

	//schließen der verbindung
	db.Close()
	w.WriteHeader(http.StatusOK)
}
