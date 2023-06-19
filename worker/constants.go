package worker

import "fmt"

func GetAllProfsQuery(first int, after string) string {
	payload := fmt.Sprintf(`{"query":"query TeacherSearchResultsPageQuery($query: TeacherSearchQuery! $schoolID: ID $first: Int $after: String) { search: newSearch { ...TeacherSearchPagination_search_1ZLmLD } school: node(id: $schoolID) { ... on School { name } id } } fragment TeacherSearchPagination_search_1ZLmLD on newSearch { teachers(query: $query, first: $first, after: $after) { edges { cursor node { ...TeacherCard_teacher id ratings { edges { node { qualityRating difficultyRating date comment helpfulRating clarityRating id } } } } } } } fragment TeacherCard_teacher on Teacher { id ...CardName_teacher } fragment CardName_teacher on Teacher { firstName lastName }","variables":{"query":{"text":"","schoolID":"U2Nob29sLTE0OTE=","fallback":false},"schoolID":"U2Nob29sLTE0OTE=","first":%d,"after":"%s"}}`, first, after)
	return payload
}

const WesternID = "U2Nob29sLTE0OTE="
const AuthID = "dGVzdDp0ZXN0"
