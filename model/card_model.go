package model

type CardResponse struct {
	Data []CardInfo `json:"data"`
}

type CardInfo struct {
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Desc       string      `json:"desc"`
	Atk        int         `json:"atk"`
	Def        int         `json:"def"`
	Level      int         `json:"level"`
	Race       string      `json:"race"`
	Attribute  string      `json:"attribute"`
	CardImages []CardImage `json:"card_images"`
}

type CardImage struct {
	ID            int    `json:"id"`
	ImageURL      string `json:"image_url"`
	ImageURLSmall string `json:"image_url_small"`
}
