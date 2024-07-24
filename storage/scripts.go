package storage 

import (
    "database/sql"
    "fmt"
    "log"
    "net"
    "os"
    "time"
     "crypto/sha256"
    
    "encoding/hex"
    "math/rand"
    

    _ "github.com/go-sql-driver/mysql" // Импорт драйвера PostgreSQL
)



var DB *sql.DB


func GetLocalIP() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddress := conn.LocalAddr().(*net.UDPAddr)

    return localAddress.IP
}

func InitDB() {
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")
    psqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", //host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
     user, password, host, port, dbname)
    //psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
     //   host, port, user, password, dbname)

        var err error
    createTableQuery := `
            CREATE TABLE IF NOT EXISTS GO_DATA.urls (
                longUrl VARCHAR(500),
                shortUrl VARCHAR(10),
                id SERIAL NOT NULL,
                PRIMARY KEY (id)
            );`
    maxRetries := 10
    for i := 0; i < maxRetries; i++ {
        DB, err = sql.Open("mysql", psqlInfo)
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
        
        
    }

    if err != nil {
        log.Fatalf("Error opening database: %v", err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }

    
    _, err = DB.Exec(createTableQuery)
    if err != nil {
        log.Fatalf("Error creating table: %v", err)

        }
        fmt.Println("Table 'urls' created or already exists.")
    fmt.Println("Successfully connected to the database!")
}



func checkHashExists(hash string) (bool, error) {
    query := "SELECT EXISTS(SELECT 1 FROM urls WHERE shortUrl = ?)"
    var exists int
    err := DB.QueryRow(query, hash).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("Ошибка выполнения запроса: %v", err)
    }
    return exists == 1, nil
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
    query := "SELECT * FROM urls WHERE shorturl = ?"
    var record Record

    // Выполнение запроса к базе данных
    row := DB.QueryRow(query, hash)
    err := row.Scan(&record.Long, &record.Short, &record.Id)

    if err != nil {
        if err == sql.ErrNoRows {
            return Record{}, fmt.Errorf("Запись не найдена для id: %s", hash)
        }
        return Record{}, fmt.Errorf("Ошибка сканирования строки: %v", err)
    }

    return record, nil
}


type Record struct {
    Long string
    Short string
    Id int
}

func InsertRecord(long string, short string) error {
    query := "INSERT INTO urls (longUrl, shortUrl) VALUES (?, ?)"
    
    // Выполнение запроса к базе данных
    result, err := DB.Exec(query, long, short)
    if err != nil {
        return fmt.Errorf("Ошибка вставки новой записи: %v", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("Ошибка получения ID последней вставки: %v", err)
    }

    fmt.Printf("Вставлена новая запись с ID: %d\n", id)
    return nil
}
