package worker

import (
	"fmt"
	"math/rand"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// Trim removes all irrelevant whitespace
func Trim(str string) string {
	// Trim the ends of whitespace
	str = strings.TrimSpace(str)

	// Trim the inner text of whitespace
	space := regexp.MustCompile(`\s+`)
	str = space.ReplaceAllString(str, " ")

	return str
}

// SleepRandom sleeps for a random amount of time in seconds between min and max
func SleepRandom(min int, max int) {
	time.Sleep(time.Duration(rand.Intn(max-min+1)+min) * time.Second)
}

// CreateData is a helper to create string map for form data
func CreateData(subject string) map[string][]string {
	data := url.Values{}
	data.Set("subject", subject)
	data.Set("command", "search")

	return data
}

// ParseEnjoymentAndDifficulty parses the scraped string that contains the enjoyment and difficulty
func parseEnjoymentAndDifficulty(str string) (enjoyment string, difficulty string, err error) {
	// the str is in the form {enjoyment}{difficulty}...
	// enjoyment can either be N/A or a percentage. E.g. N/A, 100%, 50%, ...
	// difficulty is a float. E.g. 1.0, 2.0, 3.0, 4.0, 5.0

	// check if str is less than 3 characters
	if len(str) < 3 {
		return "", "", fmt.Errorf("string is too short to be a valid enjoyment and difficulty string")
	}

	// if first three characters are N/A, then enjoyment is N/A and difficulty is the rest of the string
	if str[:3] == "N/A" {
		return str[:3], str[3:], nil
	}

	// if first three characters are not N/A, then split the string into two parts based on the first % sign
	splitStr := strings.Split(str, "%")

	// if the split string is not of length 2, then return an error
	if len(splitStr) != 2 {
		return "", "", fmt.Errorf("string is not in the correct format")
	}

	return splitStr[0] + "%", splitStr[1], nil
}
