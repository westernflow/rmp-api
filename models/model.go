package model

type Review struct {
	Professor  string `json:"professor"`
	Quality 	float64    `json:"quality"`
	Difficulty float64    `json:"difficulty"`
	Date       string `json:"date"`
	ReviewText string `json:"reviewText"`
	Course     Course `json:"course"`
	Helpful 	 float64 `json:"helpful"` // quality = helpful+clarity/2	
	Clarity 	 float64 `json:"clarity"`
}

type Course struct {
	Department string `json:"department"`
	Number     int    `json:"number"`
}

// create a model for a professor with the fields: name, rating, numRatings, department, level of difficulty, and reviews
type Professor struct {
	Name        string  `json:"name"`
	RMPId 		 string     `json:"rmpId"`
	Rating      float64 `json:"rating"`
	Department string `json:"department"`
	Difficulty  float64 `json:"difficulty"`
	Reviews     []Review `json:"reviews"`
	Courses []Course `json:"courseCodes"`
}

type Request struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type ProfessorData struct {
	// Typename      string  `json:"__typename"`
	AvgDifficulty float64 `json:"avgDifficulty"`
	AvgRating     float64 `json:"avgRating"`
	// CourseCodes   []struct {
		// CourseCount int    `json:"courseCount"`
		// CourseName  string `json:"courseName"`
	// } `json:"courseCodes"`
	Department        string `json:"department"`
	DepartmentID      string `json:"departmentId"`
	FirstName         string `json:"firstName"`
	ID                string `json:"id"`
	// IsProfCurrentUser bool   `json:"isProfCurrentUser"`
	// IsSaved           bool   `json:"isSaved"`
	LastName          string `json:"lastName"`
	// LegacyID          int    `json:"legacyId"`
	// LockStatus        string `json:"lockStatus"`
	NumRatings        int    `json:"numRatings"`
	Ratings           struct {
		Edges []struct {
			Cursor string `json:"cursor"`
			Node   struct {
				Typename            string        `json:"__typename"`
				// AdminReviewedAt     string        `json:"adminReviewedAt"`
				// AttendanceMandatory string        `json:"attendanceMandatory"`
				ClarityRating       float64           `json:"clarityRating"`
				Class               string        `json:"class"`
				Comment             string        `json:"comment"`
				CreatedByUser       bool          `json:"createdByUser"`
				Date                string        `json:"date"`
				DifficultyRating    float64           `json:"difficultyRating"`
				// FlagStatus          string        `json:"flagStatus"`
				Grade               string        `json:"grade"`
				HelpfulRating       float64           `json:"helpfulRating"`
				ID                  string        `json:"id"`
				// IsForCredit         bool          `json:"isForCredit"`
				// IsForOnlineClass    bool          `json:"isForOnlineClass"`
				// LegacyID            int           `json:"legacyId"`
				RatingTags          string        `json:"ratingTags"`
				// TeacherNote         interface{}   `json:"teacherNote"`
				// TextbookUse         int           `json:"textbookUse"`
				// Thumbs              []interface{} `json:"thumbs"`
				// ThumbsDownTotal     int           `json:"thumbsDownTotal"`
				// ThumbsUpTotal       int           `json:"thumbsUpTotal"`
				// WouldTakeAgain      interface{}   `json:"wouldTakeAgain"`
			} `json:"node"`
		} `json:"edges"`
		PageInfo struct {
			EndCursor   string `json:"endCursor"`
			HasNextPage bool   `json:"hasNextPage"`
		} `json:"pageInfo"`
	} `json:"ratings"`
	// RatingsDistribution struct {
		// R1    int `json:"r1"`
		// R2    int `json:"r2"`
		// R3    int `json:"r3"`
		// R4    int `json:"r4"`
		// R5    int `json:"r5"`
		// Total int `json:"total"`
	// } `json:"ratingsDistribution"`
	// RelatedTeachers []struct {
		// AvgRating float32    `json:"avgRating"`
		// FirstName string `json:"firstName"`
		// ID        string `json:"id"`
		// LastName  string `json:"lastName"`
		// LegacyID  int    `json:"legacyId"`
	// } `json:"relatedTeachers"`
	// School struct {
		// AvgRating  float32    `json:"avgRating"`
		// City       string `json:"city"`
		// ID         string `json:"id"`
		// LegacyID   int    `json:"legacyId"`
		// Name       string `json:"name"`
		// NumRatings int    `json:"numRatings"`
		// State      string `json:"state"`
	// } `json:"school"`
	// TeacherRatingTags     []interface{} `json:"teacherRatingTags"`
	// WouldTakeAgainPercent float64           `json:"wouldTakeAgainPercent"`
}

type Response struct {
	Data struct {
		ProfessorData ProfessorData `json:"node"`
	} `json:"data"`
}
		