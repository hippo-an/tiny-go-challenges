package repository

import (
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/models"
	"time"
)

type TestRepository struct {
}

func NewTestRepository() Repository {
	return &TestRepository{}
}

func (t *TestRepository) AllUser() {
}

func (t *TestRepository) InsertReservation(res models.Reservation) (int, error) {
	return 1, nil
}

func (t *TestRepository) InsertRoomRestriction(roomRestriction models.RoomRestriction) error {
	return nil
}

func (t *TestRepository) SearchAvailabilityByDateByRoomId(start, end time.Time, roomId int) (bool, error) {
	return false, nil
}

func (t *TestRepository) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

func (t *TestRepository) GetRoomById(id int) (models.Room, error) {
	var room models.Room
	return room, nil

}
