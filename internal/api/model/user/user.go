package user

// User ...
type User struct {
	ID    string   `firestore:"-"`
	Name  string   `firestore:"name"`
	Email string   `firestore:"email"`
	Likes []string `firestore:"likes"`
}

// NewUser ...
type NewUser struct {
	Name     string `firestore:"name"`
	Email    string `firestore:"email"`
	Password string `firestore:"password"`
}
