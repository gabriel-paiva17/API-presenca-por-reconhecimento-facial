package user 

type User struct {
    ID       	 string `json:"id"`
    Username 	 string `json:"username"`
    Email 	 	 string	`json:"email"`
	Password 	 string `json:"password"`
	RegisteredAt string `json:"registeredAt"`
}

