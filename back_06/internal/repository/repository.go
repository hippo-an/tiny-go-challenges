package repository

import "github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/models"

type Repository interface {
	AllUser()
	InsertReservation(res models.Reservation) error
}
