package model

// Fund struct to match the API response structure
type Fund struct {
	Meta struct {
		FundHouse      string `json:"fund_house"`
		SchemeType     string `json:"scheme_type"`
		SchemeCategory string `json:"scheme_category"`
		SchemeCode     int    `json:"scheme_code"`
		SchemeName     string `json:"scheme_name"`
	} `json:"meta"`
	Data []struct {
		Date string `json:"date"`
		Nav  string `json:"nav"`
	} `json:"data"`
}

// Response struct for the API response
type Response struct {
	Meta struct {
		FundHouse      string `json:"fund_house"`
		SchemeType     string `json:"scheme_type"`
		SchemeCategory string `json:"scheme_category"`
		SchemeCode     int    `json:"scheme_code"`
		SchemeName     string `json:"scheme_name"`
	} `json:"meta"`
	Period string    `json:"period,omitempty"`
	Data   []NAVData `json:"data"`
}

// Define a NAVData struct for individual data points
type NAVData struct {
	Date string `json:"date"`
	Nav  string `json:"nav"`
}
