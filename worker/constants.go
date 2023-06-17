package worker

// this URL  defines the context for the page to be scraped
const westernProfessorListURL = "https://www.ratemyprofessors.com/search/professors/1491?q=*"

// Selectors for goquery to select the specific nodes based on CSS classes
const teacherCardSelector = "a.TeacherCard__StyledTeacherCard-syjs0d-0.dLJIlx"
const cardNameSelector = "div.CardName__StyledCardName-sc-1gyrgim-0.cJdVEK"
const cardDifficultySelector = "div.CardFeedback__CardFeedbackNumber-lq6nix-2.hroXqf"
const cardDepartmentSelector = "div.CardSchool__Department-sc-19lmz2k-0.haUIRO"
const cardRatingSectionSelector = "div.CardNumRating__StyledCardNumRating-sc-17t4b9u-0.eWZmyX"
const buttonSelector = "gjQZal"

const ProfQuery = `
query NewSearchTeachersQuery($text: String!, $schoolID: ID!)
{
  newSearch {
    teachers(query: {text: $text, schoolID: $schoolID}) {
      edges {
        cursor
        node {
          id
          firstName
          lastName
          school {
            name
            id
          }
        }
      }
    }
  }
}
`
const HomePageQuery = `query TeacherSearchResultsPageQuery(
	$query: TeacherSearchQuery!
	$schoolID: ID
  ) {
	search: newSearch {
	  ...TeacherSearchPagination_search_1ZLmLD
	}
	school: node(id: $schoolID) {
	  __typename
	  ... on School {
		name
	  }
	  id
	}
  }

  fragment TeacherSearchPagination_search_1ZLmLD on newSearch {
	teachers(query: $query, first: 8, after: "") {
	  didFallback
	  edges {
		cursor
		node {
		  ...TeacherCard_teacher
		  id
		  __typename
		}
	  }
	  pageInfo {
		hasNextPage
		endCursor
	  }
	  resultCount
	  filters {
		field
		options {
		  value
		  id
		}
	  }
	}
  }

  fragment TeacherCard_teacher on Teacher {
	id
	legacyId
	avgRating
	numRatings
	...CardFeedback_teacher
	...CardSchool_teacher
	...CardName_teacher
	...TeacherBookmark_teacher
  }

  fragment CardFeedback_teacher on Teacher {
	wouldTakeAgainPercent
	avgDifficulty
  }

  fragment CardSchool_teacher on Teacher {
	department
	school {
	  name
	  id
	}
  }

  fragment CardName_teacher on Teacher {
	firstName
	lastName
  }

  fragment TeacherBookmark_teacher on Teacher {
	id
	isSaved
  }`
const DepartmentQuery = "query TeacherSearchResultsPageQuery(\n  $query: TeacherSearchQuery!\n  $schoolID: ID\n) {\n  search: newSearch {\n    ...TeacherSearchPagination_search_1ZLmLD\n  }\n  school: node(id: $schoolID) {\n    __typename\n    ... on School {\n      name\n    }\n    id\n  }\n}\n\nfragment TeacherSearchPagination_search_1ZLmLD on newSearch {\n  teachers(query: $query, first: 3000, after: \"\") {\n    didFallback\n    edges {\n      cursor\n      node {\n        ...TeacherCard_teacher\n        id\n        __typename\n      }\n    }\n    pageInfo {\n      hasNextPage\n      endCursor\n    }\n    resultCount\n    filters {\n      field\n      options {\n        value\n        id\n      }\n    }\n  }\n}\n\nfragment TeacherCard_teacher on Teacher {\n  id\n  legacyId\n  avgRating\n  numRatings\n  ...CardFeedback_teacher\n  ...CardSchool_teacher\n  ...CardName_teacher\n  ...TeacherBookmark_teacher\n}\n\nfragment CardFeedback_teacher on Teacher {\n  wouldTakeAgainPercent\n  avgDifficulty\n}\n\nfragment CardSchool_teacher on Teacher {\n  department\n  school {\n    name\n    id\n  }\n}\n\nfragment CardName_teacher on Teacher {\n  firstName\n  lastName\n}\n\nfragment TeacherBookmark_teacher on Teacher {\n  id\n  isSaved\n}\n"
const AllProfsAndReviewsQuery = `{"query":"query TeacherSearchResultsPageQuery(\r\n\t$query: TeacherSearchQuery!\r\n\t$schoolID: ID\r\n  ) {\r\n\tsearch: newSearch {\r\n\t  ...TeacherSearchPagination_search_1ZLmLD\r\n\t}\r\n\tschool: node(id: $schoolID) {\r\n\t#   __typename\r\n\t  ... on School {\r\n\t\tname\r\n\t  }\r\n\t  id\r\n\t}\r\n  }\r\n  \r\n  fragment TeacherSearchPagination_search_1ZLmLD on newSearch {\r\n\tteachers(query: $query, first: 1, after: \"\") {\r\n\t  edges {\r\n\t\tnode {\r\n\t\t  ...TeacherCard_teacher\r\n\t\t  id\r\n        ratings {\r\n            edges{\r\n                node{\r\n                    qualityRating\r\n                    difficultyRating\r\n                    date\r\n                    comment\r\n                    helpfulRating\r\n                    clarityRating\r\n                    id\r\n                }\r\n            }\r\n        }\r\n\t\t}\r\n\t  }\r\n\t}\r\n  }\r\n  \r\n  fragment TeacherCard_teacher on Teacher {\r\n\tid\r\n\t...CardName_teacher\r\n  }\r\n  \r\n  \r\n  fragment CardName_teacher on Teacher {\r\n\tfirstName\r\n\tlastName\r\n  }","variables":{"query":{"text":"","schoolID":"U2Nob29sLTE0OTE=","fallback":false},"schoolID":"U2Nob29sLTE0OTE="}}`
const WesternID = "U2Nob29sLTE0OTE="

const AuthID = "dGVzdDp0ZXN0"
