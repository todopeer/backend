package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"
	"github.com/todopeer/backend/graph"
	"github.com/todopeer/backend/graph/resolver"
	"github.com/todopeer/backend/orm"
	"github.com/todopeer/backend/services/auth"

	_ "github.com/mattn/go-sqlite3"
)

const defaultPort = "8081"

func main() {
	// setups

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	secret := os.Getenv("APP_SECRET")
	if len(secret) > 0 {
		auth.SetJWTKey(secret)
	}

	db, err := gorm.Open("sqlite3", orm.DBPATH)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	userOrm := orm.NewUserORM(db)

	// Initialize the ORMs
	r := resolver.NewResolver(
		resolver.ResolverOptionWithUserOrm(userOrm),
		resolver.ResolverOptionWithTaskOrm(orm.NewTaskORM(db)),
		resolver.ResolverOptionWithEventOrm(orm.NewEventOrm(db)),
	)
	c := graph.Config{Resolvers: r}
	c.Directives.Auth = auth.AuthDirective

	var srv http.Handler = handler.NewDefaultServer(graph.NewExecutableSchema(c))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	middlewares := []middleware{cors.Default().Handler, auth.AuthMiddleware(userOrm), logRequestBody}

	for _, middleware := range middlewares {
		srv = middleware(srv)
	}
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

type middleware func(http.Handler) http.Handler

// LogRequestBody logs the request body for debugging purpose
func logRequestBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Log the request body
		log.Println("Request Body:", string(body))

		// Restore the request body for downstream handlers
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
