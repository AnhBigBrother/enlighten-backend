package cfg

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/AnhBigBrother/enlighten-backend/internal/database"
)

type contextKeys struct {
	User string
}

var (
	DBConnection    *pgxpool.Pool
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

	// GithubClientId       string
	// GithubClientSecret   string
	// GithubCallbackUrl    string
	// GithubGetTokenUrl    string
	GithubGetUserDataUrl string

	// GoogleClientId       string
	// GoogleClientSecret   string
	// GoogleCallbackUrl    string
	// GoogleGetTokenUrl    string
	GoogleGetUserDataUrl string

	// MicrosoftClientId       string
	// MicrosoftClientSecret   string
	// MicrosoftCallbackUrl    string
	// MicrosoftGetTokenUrl    string
	MicrosoftGetUserDataUrl string

	// DiscordClientId       string
	// DiscordClientSecret   string
	// DiscordCallbackUrl    string
	// DiscordGetTokenUrl    string
	DiscordGetUserDataUrl string
)

func init() {
	env := flag.String("env", "developement", "Set environment")

	flag.Parse()

	load(*env)
}

func load(env string) {
	log.Println("Environment:", env)
	if env == "production" {
		err := godotenv.Load(".env.production")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
	}

	dbUri := os.Getenv("DB_URI")
	jwtSecret := os.Getenv("JWT_SECRET")
	port := os.Getenv("PORT")
	frontendUrl := os.Getenv("FRONTEND_URL")
	backendUrl := os.Getenv("BACKEND_URL")

	if dbUri == "" || jwtSecret == "" || port == "" || frontendUrl == "" || backendUrl == "" {
		log.Fatal("some of variables is not found in the environment: DB_URI, JWT_SECRET, PORT, FRONTEND_URL, BACKEND_URL")
	}

	// githubClientId := os.Getenv("GITHUB_CLIENT_ID")
	// githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	// githubCallbackUrl := os.Getenv("GITHUB_CALLBACK_URL")
	// if githubClientId == "" || githubClientSecret == "" || githubCallbackUrl == "" {
	// 	log.Fatal("some of variables is not found in the environment: GITHUB_CLIENT_ID, GITHUB_CLIENT_SECRET, GITHUB_CALLBACK_URL")
	// }

	// googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	// googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	// googleCallbackUrl := os.Getenv("GOOGLE_CALLBACK_URL")
	// if googleClientId == "" || googleClientSecret == "" || googleCallbackUrl == "" {
	// 	log.Fatal("some of variables is not found in the environment: GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, GOOGLE_CALLBACK_URL")
	// }

	// microsoftClientId := os.Getenv("MICROSOFT_CLIENT_ID")
	// microsoftClientSecret := os.Getenv("MICROSOFT_CLIENT_SECRET")
	// microsoftCallbackUrl := os.Getenv("MICROSOFT_CALLBACK_URL")
	// if microsoftClientId == "" || microsoftClientSecret == "" || microsoftCallbackUrl == "" {
	// 	log.Fatal("some of variables is not found in the environment: MICROSOFT_CLIENT_ID, MICROSOFT_CLIENT_SECRET, MICROSOFT_CALLBACK_URL")
	// }

	// discordClientId := os.Getenv("DISCORD_CLIENT_ID")
	// discordClientSecret := os.Getenv("DISCORD_CLIENT_SECRET")
	// discordCallbackUrl := os.Getenv("DISCORD_CALLBACK_URL")
	// if discordClientId == "" || discordClientSecret == "" || discordCallbackUrl == "" {
	// 	log.Fatal("some of variables is not found in the environment: DISCORD_CLIENT_ID, DISCORD_CLIENT_SECRET, DISCORD_CALLBACK_URL")
	// }

	conn, err := pgxpool.New(context.Background(), dbUri)
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

	// GithubClientId = githubClientId
	// GithubClientSecret = githubClientSecret
	// GithubCallbackUrl = githubCallbackUrl
	// GithubGetTokenUrl = "https://github.com/login/oauth/access_token"
	GithubGetUserDataUrl = "https://api.github.com/user"

	// GoogleClientId = googleClientId
	// GoogleClientSecret = googleClientSecret
	// GoogleCallbackUrl = googleCallbackUrl
	// GoogleGetTokenUrl = "https://oauth2.googleapis.com/token"
	GoogleGetUserDataUrl = "https://www.googleapis.com/oauth2/v3/userinfo"

	// MicrosoftClientId = microsoftClientId
	// MicrosoftClientSecret = microsoftClientSecret
	// MicrosoftCallbackUrl = microsoftCallbackUrl
	// MicrosoftGetTokenUrl = "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	MicrosoftGetUserDataUrl = "https://graph.microsoft.com/oidc/userinfo"

	// DiscordClientId = discordClientId
	// DiscordClientSecret = discordClientSecret
	// DiscordCallbackUrl = discordCallbackUrl
	// DiscordGetTokenUrl = "https://discord.com/api/v10/oauth2/token"
	DiscordGetUserDataUrl = "https://discord.com/api/v10/users/@me"
}
