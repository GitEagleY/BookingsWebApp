package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	render "github.com/GitEagleY/BookingsWebApp/internal/Render"
	"github.com/GitEagleY/BookingsWebApp/internal/config"
	"github.com/GitEagleY/BookingsWebApp/internal/driver"
	"github.com/GitEagleY/BookingsWebApp/internal/forms"
	"github.com/GitEagleY/BookingsWebApp/internal/helpers"
	"github.com/GitEagleY/BookingsWebApp/internal/models"
	"github.com/GitEagleY/BookingsWebApp/internal/repository"
	"github.com/GitEagleY/BookingsWebApp/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository and returns a pointer to it.
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,                                 // Assign the provided AppConfig to the repository.
		DB:  dbrepo.NewPostgresRepo(db.SQL, a), // Create a new PostgreSQL repository using the provided DB connection and AppConfig.
	}
}

// NewHandlers sets the repository for the handlers.
func NewHandlers(r *Repository) {
	Repo = r // Set the global Repo variable to the provided repository.
}

// Home is the handler for the home page.
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{}) // Render the home page template.
}

// About is the handler for the about page.
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{}) // Render the about page template.
}

// Generals renders the "generals" room page.
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the "majors" room page.
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the "search availability" page.
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// Reservation renders the "make a reservation" page and displays the reservation form.
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	// Retrieve the reservation from the session.
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		// If unable to retrieve the reservation, set an error message in the session and redirect to the home page.
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Get the room details based on the room ID stored in the reservation.
	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		// If there's an error fetching the room details, set an error message in the session and redirect.
		m.App.Session.Put(r.Context(), "error", "can't find room!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Attach the room name to the reservation.
	res.Room.RoomName = room.RoomName

	// Update the reservation in the session to include room details.
	m.App.Session.Put(r.Context(), "reservation", res)

	// Format the start and end dates of the reservation for display purposes.
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	// Create a string map to hold formatted dates.
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	// Prepare data for the template rendering.
	data := make(map[string]interface{})
	data["reservation"] = res

	// Render the "make a reservation" page template with the reservation data and formatted dates.
	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil), // Initialize an empty form for form rendering.
		Data:      data,           // Reservation data.
		StringMap: stringMap,      // Formatted date strings.
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {

	// Retrieve the reservation from the session.
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		// If unable to retrieve the reservation, set an error message in the session and redirect to the home page.
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Parse the form data from the request.
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	/*
		// Retrieve start and end date strings from the form.
		sd := r.Form.Get("start_date")
		ed := r.Form.Get("end_date")

		// Define the layout for date parsing.
		layout := "2006-01-02"

		// Parse the start date string into a time.Time object.
		startDate, err := time.Parse(layout, sd)
		if err != nil {
			m.App.Session.Put(r.Context(), "error", "can't parse start date")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		// Parse the end date string into a time.Time object.
		endDate, err := time.Parse(layout, ed)
		if err != nil {
			m.App.Session.Put(r.Context(), "error", "can't parse end date")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		// Convert room ID from string to integer.
		roomID, err := strconv.Atoi(r.Form.Get("room_id"))
		if err != nil {
			m.App.Session.Put(r.Context(), "error", "invalid data!")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}*/

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")

	/*// Create a reservation object with form data.
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}*/

	// Create a new form and apply validation rules.
	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	// If the form is not valid, display error message and re-render the reservation form.
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		http.Error(w, "my own error message", http.StatusSeeOther)
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// Insert the reservation into the database and retrieve the new reservation ID.
	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation into database!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	// Create a room restriction based on the reservation.
	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	// Insert the room restriction into the database.
	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {

		m.App.Session.Put(r.Context(), "error", "can't insert room restriction!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Store the reservation in the session and redirect to the reservation summary page.
	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// PostAvailability handles the search availability form submission
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	// Parse the form data from the request.
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Retrieve start and end date strings from the form.
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	// Define the layout for date parsing.
	layout := "2006-01-02"

	// Parse the start date string into a time.Time object.
	startDate, err := time.Parse(layout, start)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Parse the end date string into a time.Time object.
	endDate, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Search for room availability for the specified dates.
	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get availability for rooms")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// If no rooms are available, display an error message and redirect to the search availability page.
	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	// Store the reservation data in the session and render the "choose-room" page.
	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// jsonResponce is a custom JSON response structure.
type jsonResponce struct {
	OK        bool   `json:"ok"`      // Field for indicating success.
	Message   string `json:"message"` // Field for adding a message.
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJSON handles the availability request and responds with a JSON response.
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	availiable, _ := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	resp := jsonResponce{
		OK:        availiable, // Set OK to true to indicate success.
		Message:   "",         // Provide a message indicating availability.
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	out, err := json.MarshalIndent(resp, "", "     ") // Marshal the response into JSON format.
	if err != nil {
		helpers.ServerError(w, err) // Handle and log server error if JSON marshaling fails.
		return
	}
	//log.Println(string(out)) // Log the JSON response.

	w.Header().Set("Content-type", "application/json") // Set the response content type to JSON.
	w.Write(out)                                       // Write the JSON response to the response writer.
}

// ReservationSummary displays the reservation summary page.
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	// Retrieve the reservation from the session.
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get reservation from session")
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session") // Store an error message in the session.
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)                        // Redirect to the home page.
		return
	}

	m.App.Session.Remove(r.Context(), "reservation") // Remove the reservation from the session.

	// Prepare the data for rendering the template.
	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed
	// Render the reservation summary page template.
	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	// Retrieve the room ID from the URL parameter.
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err) // Handle and log server error if room ID parsing fails.
		return
	}

	// Retrieve the reservation from the session.
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err) // Handle and log server error if reservation retrieval fails.
		return
	}

	// Assign the selected room ID to the reservation.
	res.RoomID = roomID

	// Store the updated reservation back in the session.
	m.App.Session.Put(r.Context(), "reservation", res)

	// Redirect to the make-reservation page to complete the reservation.
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// Book Room takes url parameters, builds a sessional variable and takes user to res screen
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {

	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")
	// Define the layout for date parsing.
	layout := "2006-01-02"

	// Parse the start date string into a time.Time object.
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Parse the end date string into a time.Time object.
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	var res models.Reservation
	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		// If there's an error fetching the room details, set an error message in the session and redirect.
		m.App.Session.Put(r.Context(), "error", "can't find room!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res.Room.RoomName = room.RoomName
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
