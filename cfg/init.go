package cfg

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/AnhBigBrother/enlighten-backend/internal/database"
	"github.com/joho/godotenv"
)

type contextKey string

type contextKeys struct {
	User contextKey
}

var (
	DBConnection    *sql.DB
	DBQueries       *database.Queries
	DbUri           string
	FrontendUrl     string
	BackendUrl      string
	JwtSecret       string
	Port            string
	AccessTokenAge  int
	RefreshTokenAge int
	CookieAge       int
	CtxKeys         contextKeys

	GithubClientId       string
	GithubClientSecret   string
	GithubCallbackUrl    string
	GithubGetUserDataUrl string
	GithubGetTokenUrl    string

	GoogleClientId       string
	GoogleClientSecret   string
	GoogleCallbackUrl    string
	GoogleGetUserDataUrl string
	GoogleGetTokenUrl    string

	MicrosoftClientId       string
	MicrosoftClientSecret   string
	MicrosoftCallbackUrl    string
	MicrosoftGetUserDataUrl string
	MicrosoftGetTokenUrl    string

	DiscordClientId       string
	DiscordClientSecret   string
	DiscordCallbackUrl    string
	DiscordGetUserDataUrl string
	DiscordGetTokenUrl    string
)

func init() {
	// err := godotenv.Load(".env.production")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	dbUri := os.Getenv("DB_URI")
	jwtSecret := os.Getenv("JWT_SECRET")
	port := os.Getenv("PORT")
	frontendUrl := os.Getenv("FRONTEND_URL")
	backendUrl := os.Getenv("BACKEND_URL")

	if dbUri == "" || jwtSecret == "" || port == "" || frontendUrl == "" || backendUrl == "" {
		log.Fatal("some of variables is not found in the environment: DB_URI, JWT_SECRET, PORT, FRONTEND_URL, BACKEND_URL")
	}

	githubClientId := os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	githubCallbackUrl := os.Getenv("GITHUB_CALLBACK_URL")

	if githubClientId == "" || githubClientSecret == "" || githubCallbackUrl == "" {
		log.Fatal("some of variables is not found in the environment: GITHUB_CLIENT_ID, GITHUB_CLIENT_SECRET, GITHUB_CALLBACK_URL")
	}

	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleCallbackUrl := os.Getenv("GOOGLE_CALLBACK_URL")

	if googleClientId == "" || googleClientSecret == "" || googleCallbackUrl == "" {
		log.Fatal("some of variables is not found in the environment: GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, GOOGLE_CALLBACK_URL")
	}

	microsoftClientId := os.Getenv("MICROSOFT_CLIENT_ID")
	microsoftClientSecret := os.Getenv("MICROSOFT_CLIENT_SECRET")
	microsoftCallbackUrl := os.Getenv("MICROSOFT_CALLBACK_URL")

	if microsoftClientId == "" || microsoftClientSecret == "" || microsoftCallbackUrl == "" {
		log.Fatal("some of variables is not found in the environment: MICROSOFT_CLIENT_ID, MICROSOFT_CLIENT_SECRET, MICROSOFT_CALLBACK_URL")
	}

	discordClientId := os.Getenv("DISCORD_CLIENT_ID")
	discordClientSecret := os.Getenv("DISCORD_CLIENT_SECRET")
	discordCallbackUrl := os.Getenv("DISCORD_CALLBACK_URL")

	if discordClientId == "" || discordClientSecret == "" || discordCallbackUrl == "" {
		log.Fatal("some of variables is not found in the environment: DISCORD_CLIENT_ID, DISCORD_CLIENT_SECRET, DISCORD_CALLBACK_URL")
	}

	conn, err := sql.Open("postgres", dbUri)
	if err != nil {
		log.Fatal("Cannot connect to database")
		log.Println(err)
	}
	log.Println("Connected to database")

	DBConnection = conn
	DBQueries = database.New(conn)
	DbUri = dbUri
	FrontendUrl = frontendUrl
	BackendUrl = backendUrl
	JwtSecret = jwtSecret
	Port = port
	CtxKeys = contextKeys{
		User: "user",
	}
	AccessTokenAge = 30 * 60           // in second
	RefreshTokenAge = 7 * 24 * 60 * 60 // in second
	CookieAge = 7 * 24 * 60 * 60       // in second

	GithubClientId = githubClientId
	GithubClientSecret = githubClientSecret
	GithubCallbackUrl = githubCallbackUrl
	GithubGetUserDataUrl = "https://api.github.com/user"
	GithubGetTokenUrl = "https://github.com/login/oauth/access_token"

	GoogleClientId = googleClientId
	GoogleClientSecret = googleClientSecret
	GoogleCallbackUrl = googleCallbackUrl
	GoogleGetUserDataUrl = "https://www.googleapis.com/oauth2/v3/userinfo"
	GoogleGetTokenUrl = "https://oauth2.googleapis.com/token"

	MicrosoftClientId = microsoftClientId
	MicrosoftClientSecret = microsoftClientSecret
	MicrosoftCallbackUrl = microsoftCallbackUrl
	MicrosoftGetUserDataUrl = "https://graph.microsoft.com/oidc/userinfo"
	MicrosoftGetTokenUrl = "https://login.microsoftonline.com/common/oauth2/v2.0/token"

	DiscordClientId = discordClientId
	DiscordClientSecret = discordClientSecret
	DiscordCallbackUrl = discordCallbackUrl
	DiscordGetUserDataUrl = "https://discord.com/api/v10/users/@me"
	DiscordGetTokenUrl = "https://discord.com/api/v10/oauth2/token"
}
