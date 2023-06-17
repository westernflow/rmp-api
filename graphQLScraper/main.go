package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	model "rmpParser/models"
	"rmpParser/worker"
	"strings"
)

func main() {
	url := "https://www.ratemyprofessors.com/graphql"
	method := "POST"

	payload := worker.AllProfsAndReviewsQuery

	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(payload))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Basic dGVzdDp0ZXN0")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var response model.TeacherSearchResults
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(response.Data.Search.Teachers.Edges[0].Node.Ratings.Edges[0].Node)
}
