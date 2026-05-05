package graph

type Node struct {
	ID         string `json:"id"`
	Label      string `json:"label"`
	AvatarURL  string `json:"avatarUrl,omitempty"`
	ProfileURL string `json:"profileUrl,omitempty"`
}

type Edge struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Label  string `json:"label,omitempty"`
}

type Step struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Relation string `json:"relation"`
}

type SourceMeta struct {
	Mode   string `json:"mode"`
	Reason string `json:"reason,omitempty"`
}

type Path struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Degree int    `json:"degree"`
	Meta   SourceMeta `json:"meta"`
	Nodes  []Node `json:"nodes"`
	Edges  []Edge `json:"edges"`
	Steps  []Step `json:"steps"`
}
