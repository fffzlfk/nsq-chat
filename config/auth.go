package config

import (
	"math"
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/gplus"
)

var Store *sessions.FilesystemStore

func init() {
	Store = sessions.NewFilesystemStore(os.TempDir(), []byte(os.Getenv("SECRET_KEY")))
	Store.MaxLength(math.MaxInt64)
	gothic.Store = Store
}

var providers = make(map[string]*gplus.Provider)

func CreateProvider(redirectURI string) {

	if providers[redirectURI] == nil {
		providers[redirectURI] = gplus.New(
			"272503840736-s26k4gcvua9kamc03lg8vbp4siimh8mu.apps.googleusercontent.com",
			"W19V9lyrvOxxRQe4YI71-FDK",
			redirectURI,
		)

		goth.UseProviders(
			providers[redirectURI],
		)
	}
}
