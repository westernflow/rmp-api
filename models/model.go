package model

type TeacherSearchResults struct {
	Data struct {
		School struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"school"`
		Search struct {
			Teachers struct {
				Edges []struct {
					Cursor string `json:"cursor"`
					Node   struct {
						FirstName string `json:"firstName"`
						ID        string `json:"id"`
						LastName  string `json:"lastName"`
						Ratings   struct {
							Edges []struct {
								Node struct {
									ClarityRating    int    `json:"clarityRating"`
									Comment          string `json:"comment"`
									Date             string `json:"date"`
									DifficultyRating int    `json:"difficultyRating"`
									HelpfulRating    int    `json:"helpfulRating"`
									ID               string `json:"id"`
									QualityRating    int    `json:"qualityRating"`
								} `json:"node"`
							} `json:"edges"`
						} `json:"ratings"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"teachers"`
		} `json:"search"`
	} `json:"data"`
}
