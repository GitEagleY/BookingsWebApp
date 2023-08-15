package handlers

import (
	"encoding/json"
	"fmt"
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
			ExpectedStatus: http.StatusTemporaryRedirect,
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
			ExpectedStatus: http.StatusTemporaryRedirect,
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
	if rr.Code != http.StatusTemporaryRedirect {
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
	if rr.Code != http.StatusTemporaryRedirect {
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
	if rr.Code != http.StatusTemporaryRedirect {
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
	if rr.Code != http.StatusTemporaryRedirect {
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
	if rr.Code != http.StatusTemporaryRedirect {
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
	if rr.Code != http.StatusTemporaryRedirect {
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
			ExpectedStatus: http.StatusTemporaryRedirect,
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
			ExpectedStatus: http.StatusTemporaryRedirect,
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
			ExpectedStatus: http.StatusTemporaryRedirect,
		},
		{
			Name:           "InvalidStartDate",
			Start:          "invalid",
			End:            "2040-01-02",
			ExpectedStatus: http.StatusTemporaryRedirect,
		},
		{
			Name:           "InvalidEndDate",
			Start:          "2040-01-01",
			End:            "invalid",
			ExpectedStatus: http.StatusTemporaryRedirect,
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
