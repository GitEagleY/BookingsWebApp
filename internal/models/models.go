package models

import "time"

// User model
type User struct {
	ID         int
	FirstName  string
	LastName   string
	Email      string
	Password   string
	AcessLevel int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Room model
type Room struct {
	ID        int
	RoomName  string
	Rooms     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restriction model
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Reservation model
type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	RoomID    int
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room
}

// Room Restriction model
type RoomRestriction struct {
	ID            int
	RoomID        int
	ReservationID int
	RestrictionID int
	StartDate     time.Time
	EndDate       time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Room
	Reservations  Reservation
	Restriction   Restriction
}
