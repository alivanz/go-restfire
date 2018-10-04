package restfire

import "encoding/json"

type FirebaseAuthProvider interface {
	VerifyCustomToken(string) (TokenInfo, error)
	RefreshToken(string) (RefreshAuth, error)
}
type TokenInfo interface {
	IDToken() string
	RefreshToken() string
}
type RefreshAuth interface {
	TokenInfo
	AccessToken() string
	UserID() string
	ProjectID() string
}

type RealtimeDatabase interface {
	Get(string, interface{}) error
	Write(string, interface{}) error
	Update(string, interface{}) error
	Push(string, interface{}) (string, error)
	Delete(string) error
	Watch(string, RealtimeDatabaseListener) error
}
type RealtimeDatabaseListener interface {
	OnPut(string, json.RawMessage)
	OnPatch(string, json.RawMessage)
	OnDelete(string)
	OnCancel()
	OnAuthRevoked()
}
