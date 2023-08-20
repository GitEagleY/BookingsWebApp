package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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

// NewRepo creates a new repository and returns a pointer to it.
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,                        // Assign the provided AppConfig to the repository.
		DB:  dbrepo.NewTestingRepo(a), // Create a new PostgreSQL repository using the provided DB connection and AppConfig.
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Get the room details based on the room ID stored in the reservation.
	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		// If there's an error fetching the room details, set an error message in the session and redirect.
		m.App.Session.Put(r.Context(), "error", "can't find room!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	// 2020-01-01 -- 01/02 03:04:05PM '06 -0700

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get parse end date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
		Room:      room,
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

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

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation into database!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert room restriction!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//send email notification to guest
	htmlMessage := fmt.Sprintf(`
	<strong>Reservation Confiramtion</strong><br>
	 %s, this is confirm your reservation from %s to %s.
	`, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))
	msg := models.MailData{
		To:       reservation.Email,
		From:     "me@here.com",
		Subject:  "Reservation Confiramtion",
		Content:  htmlMessage,
		Template: "basic.html",
	}
	m.App.MailChan <- msg

	//send email notification to owner
	htmlMessage = fmt.Sprintf(`
	<strong>Reservation Confiramtion</strong><br>
	%s, this is confirm your reservation for %s from %s to %s.
	`, reservation.FirstName, reservation.Room.RoomName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))
	msg = models.MailData{
		To:      "property@owner.com",
		From:    "me@here.com",
		Subject: "Reservation Confiramtion",
		Content: htmlMessage,
	}
	m.App.MailChan <- msg

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// PostAvailability handles the search availability form submission
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	// Parse the form data from the request.
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Parse the end date string into a time.Time object.
	endDate, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Search for room availability for the specified dates.
	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get availability for rooms")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
	err := r.ParseForm()
	if err != nil {
		resp := jsonResponce{
			OK:      false,
			Message: "internal server eror",
		}
		out, _ := json.MarshalIndent(resp, "", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	availiable, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		resp := jsonResponce{
			OK:      false,
			Message: "Error connection to db",
		}
		out, _ := json.MarshalIndent(resp, "", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

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
		http.Redirect(w, r, "/", http.StatusSeeOther)                                 // Redirect to the home page.
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Parse the end date string into a time.Time object.
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	var res models.Reservation
	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		// If there's an error fetching the room details, set an error message in the session and redirect.
		m.App.Session.Put(r.Context(), "error", "can't find room!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	res.Room.RoomName = room.RoomName
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})

}
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {

		log.Println(err)
	}
	var email string
	var password string

	email = r.Form.Get("email")
	password = r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{Form: form})
		return
	}
	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "invalid login credantials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "logged in successfullly")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {

	_ = m.App.Session.Destroy(r.Context())

	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}

func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["reservations"] = reservations
	render.Template(w, r, "admin-new-reservations.page.tmpl", &models.TemplateData{Data: data})
}

func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["reservations"] = reservations
	render.Template(w, r, "admin-all-reservations.page.tmpl", &models.TemplateData{Data: data})
}

func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	spltd := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(spltd[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	src := spltd[3]

	strMap := make(map[string]string)
	strMap["src"] = src

	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return

	}

	data := make(map[string]interface{})
	data["reservation"] = res
	render.Template(w, r, "admin-reservations-show.page.tmpl", &models.TemplateData{
		StringMap: strMap,
		Data:      data})
}
func (m *Repository) AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	spltd := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(spltd[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	src := spltd[3]

	strMap := make(map[string]string)
	strMap["src"] = src

	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")

	err = m.DB.UpdateReservation(res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "flash", "Changes saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)

}
func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")

	_ = m.DB.UpdateProcessedForReservation(id, 1)

	m.App.Session.Put(r.Context(), "flash", "Reservation marked as processed")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")

	_ = m.DB.DeleteReservation(id)

	m.App.Session.Put(r.Context(), "flash", "Reservation deleted")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		mounth, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(mounth), 1, 0, 0, 0, 0, time.UTC)
	}
	data := make(map[string]interface{})
	data["now"] = now
	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMounth := next.Format("01")
	nextMounthYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	stringMap := make(map[string]string)
	stringMap["next_mounth"] = nextMounth
	stringMap["next_mounth_year"] = nextMounthYear
	stringMap["last_mounth"] = lastMonth
	stringMap["last_mounth_year"] = lastMonthYear

	stringMap["this_mounth"] = now.Format("01")
	stringMap["this_mounth_year"] = now.Format("2006")

	currentYear, currentMounth, _ := now.Date()

	currentLocation := now.Location()
	firstOfMouth := time.Date(currentYear, currentMounth, 1, 0, 0, 0, 0, currentLocation)

	lastOfTheMouth := firstOfMouth.AddDate(0, 1, -1)

	intMap := make(map[string]int)
	intMap["days_in_mounth"] = lastOfTheMouth.Day()

	rooms, err := m.DB.AllRooms()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["rooms"] = rooms

	for _, x := range rooms {
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstOfMouth; d.After(lastOfTheMouth) == false; d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0
		}

		restrictions, err := m.DB.GetRestrictionsForRoomByDate(x.ID, firstOfMouth, lastOfTheMouth)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		for _, y := range restrictions {
			if y.ReservationID > 0 {
				for d := y.StartDate; d.After(y.EndDate) == false; d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2006-01-2")] = y.ReservationID
				}
			} else {
				blockMap[y.StartDate.Format("2006-01-2")] = y.RestrictionID

			}

			data[fmt.Sprintf("reservation_map_%d", x.ID)] = reservationMap
			data[fmt.Sprintf("block_map_%d", x.ID)] = blockMap
			m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", x.ID), blockMap)
		}
	}

	render.Template(w, r, "admin-reservations-calendar.page.tmpl", &models.TemplateData{StringMap: stringMap,
		Data: data})
}
