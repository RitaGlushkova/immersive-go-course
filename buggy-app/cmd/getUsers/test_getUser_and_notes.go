package main

import (
	"context"
	"fmt"
	// "io"
	"log"
	// "net/http"
	"os"
	"os/signal"
	// "sync"
	// "time"

	"github.com/CodeYourFuture/immersive-go-course/buggy-app/util"
	"github.com/jackc/pgx/v5"
)

type User struct {
	id       string
	password string
}

func main() {

	// Set up a default POSTGRES_PASSWORD_FILE because we know where it's likely to be...
	if os.Getenv("POSTGRES_PASSWORD_FILE") == "" {
		os.Setenv("POSTGRES_PASSWORD_FILE", "volumes/secrets/postgres-passwd")
	}
	// ... and the read it. $POSTGRES_USER will still take precedence.
	dbPasswd, err := util.ReadPasswd()
	if err != nil {
		log.Fatal(err)
	}

	// The NotifyContext will signal Done when these signals are sent, allowing others to shut down safely
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	// Connect to the database
	connString := fmt.Sprintf("postgres://postgres:%s@%s/%s?sslmode=disable", dbPasswd, "localhost:5432", "app")
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer conn.Close(ctx)
	// Get all users
	users, err := GetAllActiveUsers(ctx, conn)

	if err != nil {
		log.Fatalf("error getting users: %v", err)
	}
	fmt.Println(users)

	// var count int
	// // get note for each user
	// var wg sync.WaitGroup

	// for _, user := range users {
	// 	wg.Add(1)
	// 	go func(u User) {
	// 		res, err := GetNoteForUser(ctx, conn, u)
	// 		if err != nil {
	// 			log.Fatalf("error getting note for user: %v", err)
	// 		}
	// 		fmt.Println("user: ", u)
	// 		fmt.Println("response: ", string(res))
	// 		count += 1
	// 		wg.Done()
	// 	}(user)
	// }
	// wg.Wait()
}

func GetAllActiveUsers(ctx context.Context, conn *pgx.Conn) ([]User, error) {
	// return all active users
	var users []User
	rows, err := conn.Query(ctx, "SELECT id, password FROM public.user WHERE status = $1", "active")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user User
		err = rows.Scan(&user.id, &user.password)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	return users, nil
}

// func GetNoteForUser(ctx context.Context, conn *pgx.Conn, user User) ([]byte, error) {
// 	client := http.Client{Timeout: 5 * time.Second}

// 	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8090/1/my/notes.json", http.NoBody)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	req.SetBasicAuth(user.id, "K:rocks")

// 	res, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer res.Body.Close()

// 	resBody, err := io.ReadAll(res.Body)
// 	if res.StatusCode < 200 || res.StatusCode >= 300 {
// 		return nil, err
// 	}
// 	return resBody, nil
// }
