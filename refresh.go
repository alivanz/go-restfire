package restfire

import "fmt"

type refreshTokenResp struct {
	AccessToken_  string `json:"access_token"`
	TokenType     string `json:"token_type"`
	RefreshToken_ string `json:"refresh_token"`
	IDToken_      string `json:"id_token"`
	UserID_       string `json:"user_id"`
	ProjectID_    string `json:"project_id"`
}

func (p *firebaseAuth) RefreshToken(refreshtoken string) (RefreshAuth, error) {
	var data struct {
		RefreshToken string `json:"refresh_token"`
		GrantType    string `json:"grant_type"`
	}
	data.RefreshToken = refreshtoken
	data.GrantType = "refresh_token"
	var out refreshTokenResp
	if err := requestdata("POST", fmt.Sprintf(GoogleRefreshAuth, p.key), data, &out, &verifyerr{}); err != nil {
		return nil, err
	}
	return &out, nil
}
func (x *refreshTokenResp) AccessToken() string {
	return x.AccessToken_
}
func (x *refreshTokenResp) IDToken() string {
	return x.IDToken_
}
func (x *refreshTokenResp) ProjectID() string {
	return x.ProjectID_
}
func (x *refreshTokenResp) RefreshToken() string {
	return x.RefreshToken_
}
func (x *refreshTokenResp) UserID() string {
	return x.UserID_
}
