package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hippo-an/tiny-go-challenges/back_06/internal/driver"
	"github.com/hippo-an/tiny-go-challenges/back_06/internal/models"
	"golang.org/x/crypto/bcrypt"
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

func (p *PostgresRepository) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into reservations(first_name, last_name, email, phone, start_date, end_date, room_id)
			values ($1, $2, $3, $4, $5, $6, $7) RETURNING id
			`

	var reservationId int
	err := p.db.SQL.QueryRowContext(
		ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
	).Scan(&reservationId)

	if err != nil {
		return 0, err
	}
	return reservationId, nil
}

func (p *PostgresRepository) InsertRoomRestriction(roomRestriction models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions(start_date, end_date, room_id, reservation_id, restriction_id)
			values ($1, $2, $3, $4, $5)
			`

	_, err := p.db.SQL.ExecContext(
		ctx, stmt,
		roomRestriction.StartDate,
		roomRestriction.EndDate,
		roomRestriction.RoomID,
		roomRestriction.ReservationID,
		roomRestriction.RestrictionID,
	)

	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresRepository) SearchAvailabilityByDateByRoomId(start, end time.Time, roomId int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `select count(id)
				from room_restrictions
				where room_id = $3 and $1 < end_date and $2 > start_date;
			`

	var numRows int

	err := p.db.SQL.QueryRowContext(
		ctx, stmt,
		start, end, roomId,
	).Scan(&numRows)

	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

func (p *PostgresRepository) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
			select r.id, r.room_name
			from rooms r 
			where r.id not in (
			    select rr.room_id
			    from room_restrictions rr
			    where $1 < rr.end_date and $2 > rr.start_date
			)
			`

	var rooms []models.Room

	rows, err := p.db.SQL.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var room models.Room

		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			continue
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil
}

func (p *PostgresRepository) GetRoomById(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `
		select id, room_name, created_at, updated_at
		from rooms
		where id = $1
	`

	row := p.db.SQL.QueryRowContext(ctx, query, id)
	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}

	return room, nil

}

func (p *PostgresRepository) GetUserById(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u models.User
	query := `
				select id, first_name, last_name, email, password, access_level, created_at, updated_at
				from users
				where id = $1
			`
	row := p.db.SQL.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}

	return u, nil
}

func (p *PostgresRepository) UpdateUser(m models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
				update users set first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5
				where id = $6
			`

	_, err := p.db.SQL.ExecContext(ctx, query, m.FirstName, m.LastName, m.Email, m.AccessLevel, time.Now(), m.ID)

	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresRepository) Authenticate(email, password string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var id int
	var hashedPassword string

	row := p.db.SQL.QueryRowContext(ctx, `select id, password from users where email = $1`, email)

	err := row.Scan(&id, &hashedPassword)

	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return 0, errors.New("incorrect password")
	} else if err != nil {
		return 0, err
	}

	return id, nil
}
