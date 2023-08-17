package dbrepo

import (
	"errors"
	"time"

	"github.com/GitEagleY/BookingsWebApp/internal/models"
)

// AllUsers returns a boolean value (always true in this case).
func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservations inserts a new reservation into the database and returns the new ID.
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.RoomID == 2 {
		return 0, errors.New("some error")
	}

	return 1, nil
}

// InsertRoomRestrictions inserts room restrictions into the database.
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 1000 {
		return errors.New("some error")
	}
	return nil
}

// SearchAvailabilityDatesByRoomsID checks if a room is available within a specified time range.
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	return false, nil
}

// SearchAvailabilityForAllRooms searches for available rooms within a specified time range.
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil
}

// GetRoomByID gets a room by id
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room

	if id > 2 {
		return room, errors.New("some erorr")
	}

	return room, nil
}
func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User
	return u, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}
