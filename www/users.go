package main

type User struct {
	Iss   string `json:"iss"`
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
}
