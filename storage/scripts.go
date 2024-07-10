package storage 

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq" // Импорт драйвера PostgreSQL
)

const (
    host = "127.0.0.1"
    port = 5432
    user = "postgres"
    password = "123"
    dbname = "gobase"
)



var DB *sql.DB


func InitDB() {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    var err error
    DB, err = sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatalf("Error opening database: %v", err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }

    fmt.Println("Successfully connected to the database!")
}

func GetAllRecords() ([]Record, error) {
    query := "SELECT * FROM urls"
    rows, err := DB.Query(query)
    if err != nil {
        return nil, fmt.Errorf("Error executing query: %v", err)
    }
    defer rows.Close()

    var records []Record
    for rows.Next() {
        var record Record
        err := rows.Scan(&record.Long, &record.Short, &record.Id)
        fmt.Println(record.Long, record.Short, record.Id)

        if err != nil {
            return nil, fmt.Errorf("Error scanning row: %v", err)
        }
        records = append(records, record)
    }

    err = rows.Err()
    if err != nil {
        return nil, fmt.Errorf("Error during iteration: %v", err)
    }

    return records, nil
}


func GetRecordByHash(hash string) (Record, error) {
    query := "SELECT * FROM urls WHERE shorturl = $1"
    var record Record

    // Выполнение запроса к базе данных
    row := DB.QueryRow(query, hash)
    err := row.Scan(&record.Long, &record.Short, &record.Id)

    if err != nil {
        if err == sql.ErrNoRows {
            return Record{}, fmt.Errorf("Record not found for id: %s", hash)
        }
        return Record{}, fmt.Errorf("Error scanning row: %v", err)
    }

    return record, nil
}


type Record struct {
    Long string
    Short string
    Id int
}

func InsertRecord(long string, short string) error {
    query := "INSERT INTO urls VALUES ($1, $2)"
    var id int
    err := DB.QueryRow(query, long, short).Scan(&id)
    if err != nil {
        return fmt.Errorf("Error inserting new record: %v", err)
    }

    fmt.Printf("Inserted new record with ID: %d\n", id)
    return nil
}