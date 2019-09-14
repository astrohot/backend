package like

// Like ...
type Like struct {
	ID      string `firestore:"-"`
	MainID  string `firestore:"mainID"`
	CrushID string `firestore:"crushID"`
}

// NewLike ...
type NewLike struct {
	MainID  string `firestore:"mainID"`
	CrushID string `firestore:"crushID"`
}
