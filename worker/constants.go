package worker

import "fmt"

// this URL  defines the context for the page to be scraped
const westernProfessorListURL = "https://www.ratemyprofessors.com/search/professors/1491?q=*"

// Selectors for goquery to select the specific nodes based on CSS classes
const teacherCardSelector = "a.TeacherCard__StyledTeacherCard-syjs0d-0.dLJIlx"
const cardNameSelector = "div.CardName__StyledCardName-sc-1gyrgim-0.cJdVEK"
const cardDifficultySelector = "div.CardFeedback__CardFeedbackNumber-lq6nix-2.hroXqf"
const cardDepartmentSelector = "div.CardSchool__Department-sc-19lmz2k-0.haUIRO"
const cardRatingSectionSelector = "div.CardNumRating__StyledCardNumRating-sc-17t4b9u-0.eWZmyX"
const buttonSelector = "gjQZal"

func GetAllProfsQuery(first int, after string) string {
	payload := fmt.Sprintf(`{"query":"query TeacherSearchResultsPageQuery($query: TeacherSearchQuery! $schoolID: ID $first: Int $after: String) { search: newSearch { ...TeacherSearchPagination_search_1ZLmLD } school: node(id: $schoolID) { ... on School { name } id } } fragment TeacherSearchPagination_search_1ZLmLD on newSearch { teachers(query: $query, first: $first, after: $after) { edges { cursor node { ...TeacherCard_teacher id ratings { edges { node { qualityRating difficultyRating date comment helpfulRating clarityRating id } } } } } } } fragment TeacherCard_teacher on Teacher { id ...CardName_teacher } fragment CardName_teacher on Teacher { firstName lastName }","variables":{"query":{"text":"","schoolID":"U2Nob29sLTE0OTE=","fallback":false},"schoolID":"U2Nob29sLTE0OTE=","first":%d,"after":"%s"}}`, first, after)
	return payload
}

const AllProfsAndReviewsQuery = `"{\"query\":\"query TeacherSearchResultsPageQuery(\\r\\n\\t$query: TeacherSearchQuery!\\r\\n\\t$schoolID: ID\\r\\n    $first: Int\\r\\n    $after: String\\r\\n  ) {\\r\\n\\tsearch: newSearch {\\r\\n\\t  ...TeacherSearchPagination_search_1ZLmLD\\r\\n\\t}\\r\\n\\tschool: node(id: $schoolID) {\\r\\n\\t#   __typename\\r\\n\\t  ... on School {\\r\\n\\t\\tname\\r\\n\\t  }\\r\\n\\t  id\\r\\n\\t}\\r\\n  }\\r\\n  \\r\\n  fragment TeacherSearchPagination_search_1ZLmLD on newSearch {\\r\\n\\tteachers(query: $query, first: $first, after: $after) {\\r\\n\\t  edges {\\r\\n          cursor\\r\\n\\t\\tnode {\\r\\n\\t\\t  ...TeacherCard_teacher\\r\\n\\t\\t  id\\r\\n        ratings {\\r\\n            edges{\\r\\n                node{\\r\\n                    qualityRating\\r\\n                    difficultyRating\\r\\n                    date\\r\\n                    comment\\r\\n                    helpfulRating\\r\\n                    clarityRating\\r\\n                    id\\r\\n                }\\r\\n            }\\r\\n        }\\r\\n\\t\\t}\\r\\n\\t  }\\r\\n\\t}\\r\\n  }\\r\\n  \\r\\n  fragment TeacherCard_teacher on Teacher {\\r\\n\\tid\\r\\n\\t...CardName_teacher\\r\\n  }\\r\\n  \\r\\n  \\r\\n  fragment CardName_teacher on Teacher {\\r\\n\\tfirstName\\r\\n\\tlastName\\r\\n  }\",\"variables\":{\"query\":{\"text\":\"\",\"schoolID\":\"U2Nob29sLTE0OTE=\",\"fallback\":false},\"schoolID\":\"U2Nob29sLTE0OTE=\",\"first\":0,\"after\":\"\"}}"`
const WesternID = "U2Nob29sLTE0OTE="

const AuthID = "dGVzdDp0ZXN0"
