package restfire

type refresher struct {
	provider FirebaseAuthProvider
	token    TokenInfo
}

func NewAuthRefresher(provider FirebaseAuthProvider, token TokenInfo) AuthRefresher {
	return &refresher{provider, token}
}

func (r *refresher) AuthRefresh() error {
	token, err := r.provider.RefreshToken(r.token.RefreshToken())
	if err != nil {
		return err
	}
	r.token = token
	return nil
}
func (r *refresher) Token() TokenInfo {
	return r.token
}
