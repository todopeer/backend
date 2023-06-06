package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/flyfy1/diarier/graph"
	"github.com/flyfy1/diarier/graph/resolver"
	"github.com/flyfy1/diarier/orm"
	"github.com/flyfy1/diarier/services/auth"
	"github.com/jinzhu/gorm"

	_ "github.com/mattn/go-sqlite3"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
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
	)
	c := graph.Config{Resolvers: r}
	c.Directives.Auth = auth.AuthDirective

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(c))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", auth.AuthMiddleware(userOrm)(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
