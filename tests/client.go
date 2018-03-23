package tests

import (
	"github.com/phob0s-pl/perfchat/chat"
)

var (
	admin = &chat.User{
		AuthID: "admin",
		Name:   "admin",
		Role:   chat.AdminRole,
		Token:  "admin",
	}

	dummyuser = &chat.User{
		AuthID: "dummyID",
		Name:   "dummyName",
		Role:   chat.UserRole,
		Token:  "dummyToken",
	}

	useralpha = &chat.User{
		AuthID: "alphaID",
		Name:   "alphaName",
		Role:   chat.UserRole,
		Token:  "alphaToken",
	}

	userbeta = &chat.User{
		AuthID: "betaID",
		Name:   "betaName",
		Role:   chat.UserRole,
		Token:  "betaToken",
	}
)

var (
	roomAdmin = &chat.Room{
		Creator: "admin",
		Name:    "main",
	}

	roomDummy = &chat.Room{
		Creator: "dummyName",
		Name:    "dummyroom",
	}

	roomAlpha = &chat.Room{
		Creator: "alphaName",
		Name:    "alpharoom",
	}

	roomBeta = &chat.Room{
		Creator: "betaName",
		Name:    "betaroom",
	}
)
