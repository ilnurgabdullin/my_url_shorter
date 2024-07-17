package storage 

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "time"
     "crypto/sha256"
    
    "encoding/hex"
    "math/rand"
    

    _ "github.com/lib/pq" // Импорт драйвера PostgreSQL
)



var DB *sql.DB


func InitDB() {
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

        var err error
    createTableQuery := `
            CREATE TABLE IF NOT EXISTS public.urls (
                longUrl VARCHAR(500),
                shortUrl VARCHAR(10),
                id SERIAL NOT NULL,
                PRIMARY KEY (id)
            );`
    maxRetries := 10
    for i := 0; i < maxRetries; i++ {
        DB, err = sql.Open("postgres", psqlInfo)
        if err != nil {
            log.Printf("Error connecting to the database: %v", err)
        } else {
            err = DB.Ping()
            if err == nil {
                break
            }
            log.Printf("Error pinging the database: %v", err)
        }
        log.Println("Retrying in 5 seconds...")
        time.Sleep(5 * time.Second)


        // Выполнение запроса
        _, err := DB.Exec(createTableQuery)
        if err != nil {
            log.Fatalf("Error creating table: %v", err)

        }
        fmt.Println("Table 'urls' created or already exists.")
        
    }

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


func checkHashExists(hash string) (bool, error) {
    query := "SELECT EXISTS(SELECT 1 FROM urls WHERE shortUrl = $1)"
    var exists bool
    err := DB.QueryRow(query, hash).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("Error executing query: %v", err)
    }
    return exists, nil
}

func GenerateUniqueShortHash(original string, length int) (string, error) {
    shortHash := GetShortHash(original, length)
    exists, err := checkHashExists(shortHash)
    if err != nil {
        return "", fmt.Errorf("Error checking hash existence: %v", err)
    }

    if exists {
        rand.Seed(time.Now().UnixNano())
        randomValue := rand.Intn(1000) // Случайное значение для добавления
        newInput := fmt.Sprintf("%s%d", original, randomValue)
        return GenerateUniqueShortHash(newInput, length)
    }

    return shortHash, nil
}

func GetShortHash(input string, length int) string {
    hash := sha256.New()
    hash.Write([]byte(input))
    hashBytes := hash.Sum(nil)

    hashString := hex.EncodeToString(hashBytes)

    if length > len(hashString) {
        length = len(hashString)
    }
    return hashString[:length]
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