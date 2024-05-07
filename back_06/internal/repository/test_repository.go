package repository

import (
	"errors"
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

	if id >= 100 {
		return models.Room{}, errors.New("wrong")
	}
	var room models.Room
	return room, nil

}

func (t *TestRepository) GetUserById(id int) (models.User, error) {
	var u models.User
	return u, nil
}
func (t *TestRepository) UpdateUser(m models.User) error {
	return nil
}
func (t *TestRepository) Authenticate(email, password string) (int, error) {
	return 1, nil
}
