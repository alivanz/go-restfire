package restfire

const (
	GoogleRefreshAuth   = "https://securetoken.googleapis.com/v1/token?key=%s"
	GoogleCustomAuthUrl = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyCustomToken?key=%s"
	// TODO
	GoogleGetUser                = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/getAccountInfo?key=%s"
	GoogleIdentityUrl            = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyAssertion?key=%s"
	GoogleSignUpUrl              = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/signupNewUser?key=%s"
	GooglePasswordUrl            = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyPassword?key=%s"
	GoogleDeleteUserUrl          = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/deleteAccount?key=%s"
	GoogleGetConfirmationCodeUrl = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/getOobConfirmationCode?key=%s"
	GoogleSetAccountUrl          = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/setAccountInfo?key=%s"
	GoogleCreateAuthUrl          = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/createAuthUri?key=%s"
)
