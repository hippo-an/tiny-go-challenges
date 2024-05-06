package repository

import (
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/models"
	"time"
)

type Repository interface {
	AllUser()
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(roomRestriction models.RoomRestriction) error
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomById(id int) (models.Room, error)
}
