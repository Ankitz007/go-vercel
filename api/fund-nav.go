package mutualfunds

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Define the base URL as a constant
const baseURL = "https://api.mfapi.in/mf/"

// Define a Fund struct to match the API response structure
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

// Define a Response struct for the API response
type Response struct {
	Meta struct {
		FundHouse      string `json:"fund_house"`
		SchemeType     string `json:"scheme_type"`
		SchemeCategory string `json:"scheme_category"`
		SchemeCode     int    `json:"scheme_code"`
		SchemeName     string `json:"scheme_name"`
	} `json:"meta"`
	Period string      `json:"period,omitempty"`
	Data   []DataPoint `json:"data"`
}

// Define a DataPoint struct for individual data points
type DataPoint struct {
	Date string `json:"date"`
	Nav  string `json:"nav"`
}

// HTTP handler function to process the request
func Handler(w http.ResponseWriter, r *http.Request) {
	// Fetch query parameters
	mutualFundID := r.URL.Query().Get("mutualFundID")
	startDate := r.URL.Query().Get("start")
	endDate := r.URL.Query().Get("end")

	// Check if mutualFundID is provided and is a valid integer
	if mutualFundID == "" {
		createErrorResponse(w, http.StatusBadRequest, "mutualFundID query parameter is required")
		return
	}
	if !isValidInteger(mutualFundID) {
		createErrorResponse(w, http.StatusBadRequest, "mutualFundID must be an integer")
		return
	}

	// Validate and parse dates
	start, end, err := validateAndParseDates(startDate, endDate)
	if err != nil {
		createErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Fetch fund data from API
	fund, err := fetchFundData(mutualFundID)
	if err != nil {
		createErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check if the meta field is empty, indicating an invalid mutualFundID
	if fund.Meta == (struct {
		FundHouse      string `json:"fund_house"`
		SchemeType     string `json:"scheme_type"`
		SchemeCategory string `json:"scheme_category"`
		SchemeCode     int    `json:"scheme_code"`
		SchemeName     string `json:"scheme_name"`
	}{}) {
		createErrorResponse(w, http.StatusBadRequest, "Invalid mutualFundID")
		return
	}

	// Filter data based on date range
	filteredData := filterData(fund.Data, start, end)

	// Create a success response
	response := createSuccessResponse(fund.Meta, filteredData, start, end)

	// Marshal the response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		createErrorResponse(w, http.StatusInternalServerError, "error creating JSON response")
		return
	}

	// Write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// isValidInteger checks if a string can be parsed as an integer.
func isValidInteger(value string) bool {
	_, err := strconv.Atoi(value)
	return err == nil
}

// fetchFundData fetches the fund data from the API using the mutualFundID.
func fetchFundData(mutualFundID string) (Fund, error) {
	var fund Fund

	url := fmt.Sprintf("%s%s", baseURL, mutualFundID)
	resp, err := http.Get(url)
	if err != nil {
		return fund, fmt.Errorf("error fetching data from API: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&fund); err != nil {
		return fund, fmt.Errorf("error decoding API response: %w", err)
	}

	return fund, nil
}

// validateAndParseDates validates and parses the date strings from the query parameters.
func validateAndParseDates(startDate, endDate string) (time.Time, time.Time, error) {
	var start, end time.Time
	var err error

	if startDate != "" && endDate != "" {
		start, end, err = parseDates(startDate, endDate)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}

		if end.After(time.Now()) {
			return time.Time{}, time.Time{}, fmt.Errorf("end date cannot be in the future")
		}

		if start.After(end) {
			return time.Time{}, time.Time{}, fmt.Errorf("start date cannot be after end date")
		}
	} else if startDate == "" && endDate == "" {
		// No dates provided, return all data
		start, end = time.Time{}, time.Time{}
	} else {
		// Only one of the dates provided
		return time.Time{}, time.Time{}, fmt.Errorf("both start and end dates are required in the format dd-mm-yyyy")
	}

	return start, end, nil
}

// parseDates parses the start and end date strings into time.Time objects.
func parseDates(startDate, endDate string) (time.Time, time.Time, error) {
	start, err := time.Parse("02-01-2006", startDate)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start date format. use dd-mm-yyyy")
	}

	end, err := time.Parse("02-01-2006", endDate)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end date format. use dd-mm-yyyy")
	}

	return start, end, nil
}

// filterData filters the data based on the provided date range.
func filterData(data []struct {
	Date string `json:"date"`
	Nav  string `json:"nav"`
}, start, end time.Time) []DataPoint {
	var filteredData []DataPoint

	for _, item := range data {
		date, err := time.Parse("02-01-2006", item.Date)
		if err != nil {
			continue
		}
		if (start.IsZero() && end.IsZero()) || (date.Equal(start) || date.After(start)) && (date.Equal(end) || date.Before(end)) {
			filteredData = append(filteredData, DataPoint{Date: item.Date, Nav: item.Nav})
		}
	}

	return filteredData
}

// createSuccessResponse creates a successful response with the given data and period.
func createSuccessResponse(meta struct {
	FundHouse      string `json:"fund_house"`
	SchemeType     string `json:"scheme_type"`
	SchemeCategory string `json:"scheme_category"`
	SchemeCode     int    `json:"scheme_code"`
	SchemeName     string `json:"scheme_name"`
}, data []DataPoint, start, end time.Time) Response {
	response := Response{
		Meta: meta,
		Data: data,
	}

	if !start.IsZero() && !end.IsZero() {
		response.Period = fmt.Sprintf("%s to %s", start.Format("02-01-2006"), end.Format("02-01-2006"))
	}

	return response
}

// createErrorResponse creates an error response with the given status code and message.
func createErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
