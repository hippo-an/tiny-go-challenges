package repository

import (
	"context"
	"fmt"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/driver"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/models"
	"time"
)

type PostgresRepository struct {
	db *driver.DB
}

func NewPostgresRepository(db *driver.DB) Repository {
	return &PostgresRepository{
		db: db,
	}
}

func (p *PostgresRepository) AllUser() {
	fmt.Println("hello this is sehyeong")
}

func (p *PostgresRepository) InsertReservation(res models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into reservations(first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
			`

	now := time.Now()
	_, err := p.db.SQL.ExecContext(
		ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		now,
		now,
	)

	if err != nil {
		return err
	}
	return nil
}
