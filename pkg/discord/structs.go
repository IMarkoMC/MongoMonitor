package discord

type Message struct {
	Username  *string  `json:"username,omitempty"`
	AvatarUrl *string  `json:"avatar_url,omitempty"`
	Content   *string  `json:"content,omitempty"`
	Embeds    *[]Embed `json:"embeds,omitempty"`
}

type Embed struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Color       int64    `json:"color"`
	Fields      []Fields `json:"fields"`
	Timestamp   int64    `json:"timestamp"`
	Footer      Footer   `json:"footer"`
}

type Fields struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type Footer struct {
	Text      string `json:"text"`
	IconURI   string `json:"icon_url"`
	ProxyIcon string `json:"proxy_icon_url"`
}
