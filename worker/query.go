package worker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	model "rmpParser/models"
	"strings"
)

func GetKProfessorAtCursor(k int, cursor string) model.TeacherSearchResults {
	url := "https://www.ratemyprofessors.com/graphql"
	method := "POST"

	payload := GetAllProfsQuery(k, cursor)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(string(payload)))
	if err != nil {
		fmt.Println(err)
		panic(err)

	}
	req.Header.Add("Authorization", "Basic dGVzdDp0ZXN0")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		panic(err)

	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		panic(err)

	}

	var response model.TeacherSearchResults
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return response
}
