package restfire

import "fmt"

type verifyerr struct {
	Error_ struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Errors  []struct {
			Message string `json:"message"`
			Domain  string `json:"domain"`
			Reason  string `json:"reason"`
		} `json:"errors"`
	} `json:"error"`
}

type verifyTokenResp struct {
	IDToken_      string `json:"idToken"`
	RefreshToken_ string `json:"refreshToken"`
}

func (p *firebaseAuth) VerifyCustomToken(customtoken string) (TokenInfo, error) {
	var data struct {
		Token       string `json:"token"`
		SecureToken bool   `json:"returnSecureToken"`
	}
	data.Token = customtoken
	data.SecureToken = true
	var out verifyTokenResp
	if err := requestdata("POST", fmt.Sprintf(GoogleCustomAuthUrl, p.key), data, &out, &verifyerr{}); err != nil {
		return nil, err
	}
	return &out, nil
}
func (x *verifyTokenResp) IDToken() string {
	return x.IDToken_
}
func (x *verifyTokenResp) RefreshToken() string {
	return x.RefreshToken_
}

func (msg *verifyerr) Error() string {
	return msg.Error_.Message
}
