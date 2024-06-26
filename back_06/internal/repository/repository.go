package repository

import (
	"time"

	"github.com/hippo-an/tiny-go-challenges/back_06/internal/models"
)

type Repository interface {
	AllUser()
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(roomRestriction models.RoomRestriction) error
	SearchAvailabilityByDateByRoomId(start, end time.Time, roomId int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomById(id int) (models.Room, error)

	GetUserById(id int) (models.User, error)
	UpdateUser(m models.User) error
	Authenticate(email, password string) (int, error)
}
