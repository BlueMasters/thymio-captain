package main

type AuthRequest struct {
	Iss           string `json:"iss"`
	AtHash        string `json:"at_hash"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	EmailVerified bool   `json:"email_verified"`
	Azp           string `json:"azp"`
	Hd            string `json:"hd"`
	Email         string `json:"email"`
	Iat           int64  `json:"iat,string"`
	Exp           int64  `json:"exp,string"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
	Alg           string `json:"alg"`
	Kid           string `json:"kid"`
}

type AuthError struct {
	ErrorDescription string `json:"error_description"`
}
