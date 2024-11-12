package authenticationmaster

import (
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/hmcalister/AuthSSO/database"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/microcosm-cc/bluemonday"
)

const (
	passwordMaxLen = 1024

	// Tokens expire after this time. Effectively logs users out, prevents cookie stealing attacks.
	tokenExpirationDuration time.Duration = 6 * time.Hour

	// Maximum time a database query should take before failing.
	maximumDatabaseQueryDuration = 5 * time.Second
)

// Struct for holding state of authentication. Includes connection to database where credentials are held, and router to accept login / registration attempts.
type AuthenticationMaster struct {
	databaseConnection *database.DatabaseManager
	tokenSecretKey     []byte
	tokenAuth          *jwtauth.JWTAuth
	htmlSanitizer      *bluemonday.Policy
}

// Create a new authentication master.
//
// mainRouter is taken to mount the apiRouter to the endpoint "/api".
// db is a connection to the database holding user credentials.
// tokenSecretKey is a byte array holding the secret key for signing the JWT.
func NewAuthenticationMaster(db *database.DatabaseManager, tokenSecretKey []byte) *AuthenticationMaster {
	// Tokens must be signed using the secret key, be signed with the correct signing method.
	// Tokens must also be issued by this server, and have a subject claim (the subject is the UserID).
	tokenAuth := jwtauth.New(tokenSigningMethod.Alg(),
		tokenSecretKey,
		tokenSecretKey,
		jwt.WithIssuer(issuerString),
		jwt.WithRequiredClaim(jwt.SubjectKey),
	)

	authMaster := &AuthenticationMaster{
		databaseConnection: db,
		tokenSecretKey:     tokenSecretKey,
		tokenAuth:          tokenAuth,
		htmlSanitizer:      bluemonday.UGCPolicy(),
	}

	return authMaster
}
