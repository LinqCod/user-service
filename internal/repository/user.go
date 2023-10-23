package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/linqcod/user-service/internal/domain"
	"log"
	"strings"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u userRepository) Create(ctx context.Context, user *domain.User) error {
	if err := u.db.QueryRowContext(
		ctx,
		`
			INSERT INTO users (name, surname, patronymic, age, gender, nationality) 
			VALUES ($1, $2, $3, $4, $5, $6);
		`,
		user.Name,
		user.Surname,
		user.Patronymic,
		user.Age,
		user.Gender,
		user.Nationality,
	).Err(); err != nil {
		return err
	}

	return nil
}

func (u userRepository) Delete(ctx context.Context, id int64) error {
	if err := u.db.QueryRowContext(
		ctx,
		` DELETE FROM users WHERE id=$1;`,
		id,
	).Err(); err != nil {
		return err
	}

	return nil
}

func (u userRepository) Update(ctx context.Context, user *domain.User) error {
	if err := u.db.QueryRowContext(
		ctx,
		`
			UPDATE users
			SET name=coalesce($1, name), 
			    surname=coalesce($2, surname), 
			    patronymic=coalesce($3, patronymic), 
			    age=coalesce($4, age), 
			    gender=coalesce($5, gender), 
			    nationality=coalesce($6, nationality)
			WHERE id=$7;
		`,
		user.Name,
		user.Surname,
		user.Patronymic,
		user.Age,
		user.Gender,
		user.Nationality,
		user.Id,
	).Err(); err != nil {
		return err
	}

	return nil
}

func (u userRepository) GetById(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User

	if err := u.db.QueryRowContext(
		ctx,
		`
			SELECT id, name, surname, patronymic, age, gender, nationality
 			FROM users 
 			WHERE id=$1;
		`,
		id,
	).Scan(
		&user.Id,
		&user.Name,
		&user.Surname,
		&user.Patronymic,
		&user.Age,
		&user.Gender,
		&user.Nationality,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u userRepository) GetFilteredUsers(ctx context.Context, count string, nationality string, minAge string, maxAge string, gender string) ([]*domain.User, error) {
	var users []*domain.User

	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT id, name, surname, patronymic, age, gender, nationality FROM users")

	if nationality != "" || minAge != "" || maxAge != "" || gender != "" {
		queryBuilder.WriteString(" WHERE ")

		conditions := make([]string, 0)
		if nationality != "" {
			conditions = append(conditions, fmt.Sprintf("nationality='%s'", nationality))
		}
		if minAge != "" {
			conditions = append(conditions, fmt.Sprintf("age>=%s", minAge))
		}
		if maxAge != "" {
			conditions = append(conditions, fmt.Sprintf("age<=%s", maxAge))
		}
		if gender != "" {
			conditions = append(conditions, fmt.Sprintf("gender='%s'", gender))
		}

		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	if count != "" {
		queryBuilder.WriteString(fmt.Sprintf(" LIMIT %s", count))
	}

	log.Println(queryBuilder.String())

	rows, err := u.db.QueryContext(
		ctx,
		queryBuilder.String(),
	)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Surname,
			&user.Patronymic,
			&user.Age,
			&user.Gender,
			&user.Nationality,
		); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}
