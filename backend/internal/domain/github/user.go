package github

type User struct {
	Login      string
	Name       string
	AvatarURL  string
	ProfileURL string
	Bio        string
}

type Relationship struct {
	Username string
	Relation string
}
