package cfg

import (
	"database/sql"
	"log"
	"os"

	"github.com/AnhBigBrother/enlighten-backend/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	DBQueries       *database.Queries
	DbUri           string
	JwtSecret       string
	Port            string
	AccessTokenAge  int
	RefreshTokenAge int
	CookieAge       int

	GithubGetUserDataUrl    string
	GoogleGetUserDataUrl    string
	MicrosoftGetUserDataUrl string
	DiscordGetUserDataUrl   string
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	dbUri := os.Getenv("DB_URI")
	jwtSecret := os.Getenv("JWT_SECRET")
	port := os.Getenv("PORT")

	if dbUri == "" {
		log.Fatal("DB_URI is not found in the environment")
	}
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not found in the environment")
	}
	if port == "" {
		log.Fatal("PORT is not found in the environment")
	}

	conn, err := sql.Open("postgres", dbUri)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	log.Println("Connected to database")

	DBQueries = database.New(conn)
	DbUri = dbUri
	JwtSecret = jwtSecret
	Port = port
	AccessTokenAge = 30 * 60           // in second
	RefreshTokenAge = 7 * 24 * 60 * 60 // in second
	CookieAge = 7 * 24 * 60 * 60       // in second

	GithubGetUserDataUrl = "https://api.github.com/user"
	GoogleGetUserDataUrl = "https://www.googleapis.com/oauth2/v3/userinfo"
	MicrosoftGetUserDataUrl = "https://graph.microsoft.com/oidc/userinfo"
	DiscordGetUserDataUrl = "https://discord.com/api/v10/users/@me"
}
