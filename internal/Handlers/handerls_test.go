package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/GitEagleY/BookingsWebApp/internal/models"
)

// /////////////////////////////////////////////////////////////
// //////////////HANDLERS//////////////////////////////////////
// ///////////////////////////////////////////////////////////
var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{

	{"home", "/", "GET", 200},
	{"about", "/about", "GET", http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", http.StatusOK},
	{"search-availability", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
}

func TestHandlers(t *testing.T) {
	// Get the routes for testing
	routes := getRoutes()
	// Create a new TLS test server with the defined routes
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	// Loop through each test case
	for _, e := range theTests {
		if e.method == "GET" {
			// Send a GET request to the test server with the specified URL
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			// Compare received status code with expected status code
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}

// /////////////////////////////////////////////////////////////
// //////////////RESERVATION///////////////////////////////////
// ///////////////////////////////////////////////////////////
// TestRepository_Reservation tests the Reservation handler function from the repository.

type ReservationTestCase struct {
	Name           string
	Reservation    models.Reservation
	ExpectedStatus int
}

func TestRepository_Reservation(t *testing.T) {
	testCases := []ReservationTestCase{
		{
			Name: "ValidReservation",
			Reservation: models.Reservation{
				RoomID: 1,
				Room: models.Room{
					ID:       1,
					RoomName: "Generals Quarters",
				},
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:           "NoReservationRedirect",
			ExpectedStatus: http.StatusSeeOther,
		},
		{
			Name: "InvalidReservationRoomID",
			Reservation: models.Reservation{
				RoomID: 100,
				Room: models.Room{
					ID:       1,
					RoomName: "Generals Quarters",
				},
			},
			ExpectedStatus: http.StatusSeeOther,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			runReservationTest(t, tc.Reservation, tc.ExpectedStatus)
		})
	}
}

func runReservationTest(t *testing.T, reservation models.Reservation, expectedStatus int) {
	t.Helper()

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	if reservation.RoomID != 0 {
		session.Put(ctx, "reservation", reservation)
	}

	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != expectedStatus {
		t.Errorf("Reservation handler returned wrong response code for '%s': got %d, wanted %d", reservation.Room.RoomName, rr.Code, expectedStatus)
	}
}

// /////////////////////////////////////////////////////////////
// //////////////POST RESERVATION//////////////////////////////
// ///////////////////////////////////////////////////////////
func TestRepository_PostReservation(t *testing.T) {

	postedData := url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "john@smith.com")
	postedData.Add("room_id", "1")
	postedData.Add("phone", "1231231234")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Postreservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	//test for missing post body

	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post reservation handler returned wrong response code for missing post body : got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	//test for invalis startDate
	reqBody := "start_date=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1231231234")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post reservation handler returned wrong response code for invalid start date: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
	//test for invalis endDate
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1231231234")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post reservation handler returned wrong response code for end date: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
	//test for invalid room id
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1231231234")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post reservation handler returned wrong response code for room id: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
	//test for invalid data
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=J")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1231231234")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post reservation handler returned wrong response code for invalid data: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
	//test for failure for insert reservation into db
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1231231234")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post reservation failed when trying to fail inserting reservation: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
	//test for failure for insert restriction into db
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1000")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1231231234")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post reservation failed when trying to fail inserting reservation: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

// /////////////////////////////////////////////////////////////
// //////////////AVAILABILITY JSON/////////////////////////////
// ///////////////////////////////////////////////////////////
func TestRepository_AvailiabilityJSON(t *testing.T) {
	reqBody := "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-01")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	// create reques
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "x-www-form-urlencoded")
	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var j jsonResponce

	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json")
		/*****************************************
		// second case -- rooms not available
		*****************************************/
		// create our request body
		reqBody = "start=2040-01-01"
		reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2040-01-02")
		reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

		// create our request
		req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

		// get the context with session
		ctx = getCtx(req)
		req = req.WithContext(ctx)

		// set the request header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// create our response recorder, which satisfies the requirements
		// for http.ResponseWriter
		rr = httptest.NewRecorder()

		// make our handler a http.HandlerFunc
		handler = http.HandlerFunc(Repo.AvailabilityJSON)

		// make the request to our handler
		handler.ServeHTTP(rr, req)

		// this time we want to parse JSON and get the expected response
		err = json.Unmarshal([]byte(rr.Body.String()), &j)
		if err != nil {
			t.Error("failed to parse json!")
		}

		// since we specified a start date < 2049-12-31, we expect availability
		if !j.OK {
			t.Error("Got no availability when some was expected in AvailabilityJSON")
		}

		/*****************************************
		// third case -- no request body
		*****************************************/
		// create our request
		req, _ = http.NewRequest("POST", "/search-availability-json", nil)

		// get the context with session
		ctx = getCtx(req)
		req = req.WithContext(ctx)

		// set the request header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// create our response recorder, which satisfies the requirements
		// for http.ResponseWriter
		rr = httptest.NewRecorder()

		// make our handler a http.HandlerFunc
		handler = http.HandlerFunc(Repo.AvailabilityJSON)

		// make the request to our handler
		handler.ServeHTTP(rr, req)

		// this time we want to parse JSON and get the expected response
		err = json.Unmarshal([]byte(rr.Body.String()), &j)
		if err != nil {
			t.Error("failed to parse json!")
		}

		// since we specified a start date < 2049-12-31, we expect availability
		if j.OK || j.Message != "Internal server error" {
			t.Error("Got availability when request body was empty")
		}

		/*****************************************
		// fourth case -- database error
		*****************************************/
		// create our request body
		reqBody = "start=2060-01-01"
		reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2060-01-02")
		reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
		req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

		// get the context with session
		ctx = getCtx(req)
		req = req.WithContext(ctx)

		// set the request header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// create our response recorder, which satisfies the requirements
		// for http.ResponseWriter
		rr = httptest.NewRecorder()

		// make our handler a http.HandlerFunc
		handler = http.HandlerFunc(Repo.AvailabilityJSON)

		// make the request to our handler
		handler.ServeHTTP(rr, req)

		// this time we want to parse JSON and get the expected response
		err = json.Unmarshal([]byte(rr.Body.String()), &j)
		if err != nil {
			t.Error("failed to parse json!")
		}

		// since we specified a start date < 2049-12-31, we expect availability
		if j.OK || j.Message != "Error querying database" {
			t.Error("Got availability when simulating database error")
		}
	}
}

// /////////////////////////////////////////////////////////////
// //////////////CHOOSE ROOM///////////////////////////////////
// ///////////////////////////////////////////////////////////
type ChooseRoomTestCase struct {
	Name           string
	URL            string
	Reservation    models.Reservation
	ExpectedStatus int
}

func TestRepository_ChooseRoom(t *testing.T) {
	/*testCases := []ChooseRoomTestCase{
		{
			Name: "ReservationInSession",
			URL:  "/choose-room/1",
			Reservation: models.Reservation{
				RoomID: 1,
				Room: models.Room{
					ID:       1,
					RoomName: "General's Quarters",
				},
			},
			ExpectedStatus: http.StatusSeeOther,
		},
		{
			Name:           "ReservationNotInSession",
			URL:            "/choose-room/1",
			ExpectedStatus: http.StatusTemporaryRedirect,
		},
		{
			Name:           "MissingOrMalformedURLParameter",
			URL:            "/choose-room/fish",
			ExpectedStatus: http.StatusTemporaryRedirect,
		},
	}

		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				runChooseRoomTest(t, tc.URL, tc.Reservation, tc.ExpectedStatus)
			})
		}
	*/
}

func runChooseRoomTest(t *testing.T, url string, reservation models.Reservation, expectedStatus int) {
	t.Helper()

	req, _ := http.NewRequest("GET", url, nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = url

	rr := httptest.NewRecorder()
	if reservation.RoomID != 0 {
		session.Put(ctx, "reservation", reservation)
	}

	handler := http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != expectedStatus {
		t.Errorf("ChooseRoom handler returned wrong response code for '%s': got %d, wanted %d", url, rr.Code, expectedStatus)
	}
}

// //////////////////////////////////////////////////////////////////////////
// ///////////////////////BOOK ROOM/////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////
type BookRoomTestCase struct {
	Name           string
	URL            string
	Reservation    models.Reservation
	ExpectedStatus int
}

func TestRepository_BookRoom(t *testing.T) {
	testCases := []BookRoomTestCase{
		{
			Name: "DatabaseWorks",
			URL:  "/book-room?s=2050-01-01&e=2050-01-02&id=1",
			Reservation: models.Reservation{
				RoomID: 1,
				Room: models.Room{
					ID:       1,
					RoomName: "General's Quarters",
				},
			},
			ExpectedStatus: http.StatusSeeOther,
		},
		{
			Name:           "DatabaseFailed",
			URL:            "/book-room?s=2040-01-01&e=2040-01-02&id=4",
			ExpectedStatus: http.StatusSeeOther,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			runBookRoomTest(t, tc.URL, tc.Reservation, tc.ExpectedStatus)
		})
	}
}

func runBookRoomTest(t *testing.T, url string, reservation models.Reservation, expectedStatus int) {
	t.Helper()

	req, _ := http.NewRequest("GET", url, nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	if reservation.RoomID != 0 {
		session.Put(ctx, "reservation", reservation)
	}

	handler := http.HandlerFunc(Repo.BookRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != expectedStatus {
		t.Errorf("BookRoom handler returned wrong response code for '%s': got %d, wanted %d", url, rr.Code, expectedStatus)
	}
}

// //////////////////////////////////////////////////////////////////////////
// ///////////////////////RESERVATION SUMMARY///////////////////////////////
// ////////////////////////////////////////////////////////////////////////

type ReservationSummaryTestCase struct {
	Name           string
	Reservation    models.Reservation
	ExpectedStatus int
}

func TestRepository_ReservationSummary(t *testing.T) {
	testCases := []ReservationSummaryTestCase{
		{
			Name: "ReservationInSession",
			Reservation: models.Reservation{
				RoomID: 1,
				Room: models.Room{
					ID:       1,
					RoomName: "General's Quarters",
				},
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:           "ReservationNotInSession",
			ExpectedStatus: http.StatusSeeOther,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			runReservationSummaryTest(t, tc.Reservation, tc.ExpectedStatus)
		})
	}
}

func runReservationSummaryTest(t *testing.T, reservation models.Reservation, expectedStatus int) {
	t.Helper()

	req, _ := http.NewRequest("GET", "/reservation-summary", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	if reservation.RoomID != 0 {
		session.Put(ctx, "reservation", reservation)
	}

	handler := http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, req)

	if rr.Code != expectedStatus {
		t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", rr.Code, expectedStatus)
	}
}

// //////////////////////////////////////////////////////////////////////////
// ///////////////////////POST AVAILABILITY/////////////////////////////////
// ////////////////////////////////////////////////////////////////////////
type AvailabilityTestCase struct {
	Name           string
	Start          string
	End            string
	ExpectedStatus int
}

func TestRepository_PostAvailability(t *testing.T) {
	testCases := []AvailabilityTestCase{
		{
			Name:           "NoRoomsAvailable",
			Start:          "2050-01-01",
			End:            "2050-01-02",
			ExpectedStatus: http.StatusSeeOther,
		},
		{
			Name:           "RoomsAvailable",
			Start:          "2040-01-01",
			End:            "2040-01-02",
			ExpectedStatus: http.StatusSeeOther,
		},
		{
			Name:           "EmptyRequestBody",
			Start:          "",
			End:            "",
			ExpectedStatus: http.StatusSeeOther,
		},
		{
			Name:           "InvalidStartDate",
			Start:          "invalid",
			End:            "2040-01-02",
			ExpectedStatus: http.StatusSeeOther,
		},
		{
			Name:           "InvalidEndDate",
			Start:          "2040-01-01",
			End:            "invalid",
			ExpectedStatus: http.StatusSeeOther,
		},
		{
			Name:           "DatabaseQueryFails",
			Start:          "2060-01-01",
			End:            "2060-01-02",
			ExpectedStatus: http.StatusSeeOther,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			runAvailabilityTest(t, tc.Start, tc.End, tc.ExpectedStatus)
		})
	}
}

func runAvailabilityTest(t *testing.T, start, end string, expectedStatus int) {
	t.Helper()

	reqBody := fmt.Sprintf("start=%s&end=%s", start, end)
	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != expectedStatus {
		t.Errorf("Post availability with start=%s and end=%s gave wrong status code: got %d, wanted %d", start, end, rr.Code, expectedStatus)
	}
}

// reservationSummaryTests is the data to test ReservationSummary handler
var reservationSummaryTests = []struct {
	name               string
	reservation        models.Reservation
	url                string
	expectedStatusCode int
	expectedLocation   string
}{
	{
		name: "res-in-session",
		reservation: models.Reservation{
			RoomID: 1,
			Room: models.Room{
				ID:       1,
				RoomName: "General's Quarters",
			},
		},
		url:                "/reservation-summary",
		expectedStatusCode: http.StatusOK,
		expectedLocation:   "",
	},
	{
		name:               "res-not-in-session",
		reservation:        models.Reservation{},
		url:                "/reservation-summary",
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/",
	},
}

// TestReservationSummary tests the ReservationSummaryHandler
func TestReservationSummary(t *testing.T) {
	for _, e := range reservationSummaryTests {
		req, _ := http.NewRequest("GET", e.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		if e.reservation.RoomID > 0 {
			session.Put(ctx, "reservation", e.reservation)
		}

		handler := http.HandlerFunc(Repo.ReservationSummary)

		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}

		if e.expectedLocation != "" {
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}
	}
}

// chooseRoomTests is the data for ChooseRoom handler tests, /choose-room/{id}
var chooseRoomTests = []struct {
	name               string
	reservation        models.Reservation
	url                string
	expectedStatusCode int
	expectedLocation   string
}{
	{
		name: "reservation-in-session",
		reservation: models.Reservation{
			RoomID: 1,
			Room: models.Room{
				ID:       1,
				RoomName: "General's Quarters",
			},
		},
		url:                "/choose-room/1",
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/make-reservation",
	},
	{
		name:               "reservation-not-in-session",
		reservation:        models.Reservation{},
		url:                "/choose-room/1",
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/",
	},
	{
		name:               "malformed-url",
		reservation:        models.Reservation{},
		url:                "/choose-room/fish",
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/",
	},
}

// TestChooseRoom tests the ChooseRoom handler
/*
func TestChooseRoom(t *testing.T) {
	for _, e := range chooseRoomTests {
		req, _ := http.NewRequest("GET", e.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		// set the RequestURI on the request so that we can grab the ID from the URL
		req.RequestURI = e.url

		rr := httptest.NewRecorder()
		if e.reservation.RoomID > 0 {
			session.Put(ctx, "reservation", e.reservation)
		}

		handler := http.HandlerFunc(Repo.ChooseRoom)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}

		if e.expectedLocation != "" {
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}
	}
}*/

// bookRoomTests is the data for the BookRoom handler tests
var bookRoomTests = []struct {
	name               string
	url                string
	expectedStatusCode int
}{
	{
		name:               "database-works",
		url:                "/book-room?s=2050-01-01&e=2050-01-02&id=1",
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "database-fails",
		url:                "/book-room?s=2040-01-01&e=2040-01-02&id=4",
		expectedStatusCode: http.StatusSeeOther,
	},
}

// TestBookRoom tests the BookRoom handler
func TestBookRoom(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	for _, e := range bookRoomTests {
		req, _ := http.NewRequest("GET", e.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		session.Put(ctx, "reservation", reservation)

		handler := http.HandlerFunc(Repo.BookRoom)

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusSeeOther {
			t.Errorf("%s failed: returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}
	}
}

// loginTests is the data for the Login handler tests
var loginTests = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"valid-credentials",
		"me@here.ca",
		http.StatusSeeOther,
		"",
		"/",
	},
	{
		"invalid-credentials",
		"jack@nimble.com",
		http.StatusSeeOther,
		"",
		"/user/login",
	},
	{
		"invalid-data",
		"j",
		http.StatusOK,
		`action="/user/login"`,
		"",
	},
}

/*
	func TestLogin(t *testing.T) {
		// range through all tests
		for _, e := range loginTests {
			postedData := url.Values{}
			postedData.Add("email", e.email)
			postedData.Add("password", "password")

			// create request
			req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))
			ctx := getCtx(req)
			req = req.WithContext(ctx)

			// set the header
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()

			// call the handler
			handler := http.HandlerFunc(Repo.PostShowLogin)
			handler.ServeHTTP(rr, req)

			if rr.Code != e.expectedStatusCode {
				t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
			}

			if e.expectedLocation != "" {
				// get the URL from test
				actualLoc, _ := rr.Result().Location()
				if actualLoc.String() != e.expectedLocation {
					t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
				}
			}

			// checking for expected values in HTML
			if e.expectedHTML != "" {
				// read the response body into a string
				html := rr.Body.String()
				if !strings.Contains(html, e.expectedHTML) {
					t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
				}
			}
		}
	}

	var adminPostShowReservationTests = []struct {
		name                 string
		url                  string
		postedData           url.Values
		expectedResponseCode int
		expectedLocation     string
		expectedHTML         string
	}{

		{
			name: "valid-data-from-new",
			url:  "/admin/reservations/new/1/show",
			postedData: url.Values{
				"first_name": {"John"},
				"last_name":  {"Smith"},
				"email":      {"john@smith.com"},
				"phone":      {"555-555-5555"},
			},
			expectedResponseCode: http.StatusSeeOther,
			expectedLocation:     "/admin/reservations-new",
			expectedHTML:         "",
		},
		{
			name: "valid-data-from-all",
			url:  "/admin/reservations/all/1/show",
			postedData: url.Values{
				"first_name": {"John"},
				"last_name":  {"Smith"},
				"email":      {"john@smith.com"},
				"phone":      {"555-555-5555"},
			},
			expectedResponseCode: http.StatusSeeOther,
			expectedLocation:     "/admin/reservations-all",
			expectedHTML:         "",
		},
		{
			name: "valid-data-from-cal",
			url:  "/admin/reservations/cal/1/show",
			postedData: url.Values{
				"first_name": {"John"},
				"last_name":  {"Smith"},
				"email":      {"john@smith.com"},
				"phone":      {"555-555-5555"},
				"year":       {"2022"},
				"month":      {"01"},
			},
			expectedResponseCode: http.StatusSeeOther,
			expectedLocation:     "/admin/reservations-calendar?y=2022&m=01",
			expectedHTML:         "",
		},
	}

// TestAdminPostShowReservation tests the AdminPostReservation handler

	func TestAdminPostShowReservation(t *testing.T) {
		for _, e := range adminPostShowReservationTests {
			var req *http.Request
			if e.postedData != nil {
				req, _ = http.NewRequest("POST", "/user/login", strings.NewReader(e.postedData.Encode()))
			} else {
				req, _ = http.NewRequest("POST", "/user/login", nil)
			}
			ctx := getCtx(req)
			req = req.WithContext(ctx)
			req.RequestURI = e.url

			// set the header
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()

			// call the handler
			handler := http.HandlerFunc(Repo.AdminPostShowReservation)
			handler.ServeHTTP(rr, req)

			if rr.Code != e.expectedResponseCode {
				t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
			}

			if e.expectedLocation != "" {
				// get the URL from test
				actualLoc, _ := rr.Result().Location()
				if actualLoc.String() != e.expectedLocation {
					t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
				}
			}

			// checking for expected values in HTML
			if e.expectedHTML != "" {
				// read the response body into a string
				html := rr.Body.String()
				if !strings.Contains(html, e.expectedHTML) {
					t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
				}
			}
		}
	}
*/
var adminProcessReservationTests = []struct {
	name                 string
	queryParams          string
	expectedResponseCode int
	expectedLocation     string
}{
	{
		name:                 "process-reservation",
		queryParams:          "",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
	{
		name:                 "process-reservation-back-to-cal",
		queryParams:          "?y=2021&m=12",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
}

func TestAdminProcessReservation(t *testing.T) {
	for _, e := range adminProcessReservationTests {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/admin/process-reservation/cal/1/do%s", e.queryParams), nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AdminProcessReservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusSeeOther {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
		}
	}
}

var adminDeleteReservationTests = []struct {
	name                 string
	queryParams          string
	expectedResponseCode int
	expectedLocation     string
}{
	{
		name:                 "delete-reservation",
		queryParams:          "",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
	{
		name:                 "delete-reservation-back-to-cal",
		queryParams:          "?y=2021&m=12",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
}

func TestAdminDeleteReservation(t *testing.T) {
	for _, e := range adminDeleteReservationTests {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/admin/process-reservation/cal/1/do%s", e.queryParams), nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AdminDeleteReservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusSeeOther {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
		}
	}
}

// gets the context
func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
