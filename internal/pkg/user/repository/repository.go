package repository

import (
	"fmt"
	"github.com/Arkadiyche/bd_techpark/internal/pkg/models"
	"github.com/jackc/pgx"
)

type UserRepository struct {
	db *pgx.ConnPool
}

func NewUserRepository(db *pgx.ConnPool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Insert(user models.User) (error error)  {
	_, err := r.db.Exec("INSERT INTO users VALUES ($1, $2, $3, $4)",
		user.Nickname,
		user.Email,
		user.About,
		user.Fullname)
	if err != nil {
		err = models.Exist
	}
	return err
}

func (r *UserRepository) SelectByNicknameOrEmail(nickname, email string) (us *models.Users, error error) {
	user := models.User{}
	users := models.Users{}
	query, err := r.db.Query("SELECT nickname, fullname, about, email FROM users WHERE nickname = $1 or email = $2", nickname, email)
	if err != nil {
		return nil, err
	}
	defer query.Close()
	for query.Next(){
		query.Scan(&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Email)
		users = append(users, user)
	}
	return &users, nil
}

func (r *UserRepository) SelectByNickname(nickname string) (u *models.User, error *models.Error) {
	user := models.User{}
	err := r.db.QueryRow("SELECT nickname, fullname, about, email FROM users WHERE nickname = $1", nickname).
		Scan(&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Email)
	if err != nil {
		return nil, &models.Error{Message: models.NotExist.Error()}
	}
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) (error *models.Error) {
	fmt.Println(user)
	query := "UPDATE users SET  "
	if user.About != "" {
		query += fmt.Sprintf(" about = '%s', ", user.About)
	}
	if user.Email != "" {
		query += fmt.Sprintf(" email = '%s', ", user.Email)
	}
	if user.Fullname != "" {
		query += fmt.Sprintf(" fullname = '%s', ", user.Fullname)
	}
	query = query[:len(query)-2]
	query += "WHERE nickname = $1 RETURNING about, email, fullname, nickname"
	if user.About == "" && user.Email == "" && user.Fullname == "" {
		query = "SELECT about, email, fullname, nickname FROM users WHERE nickname = $1"
	}
	err := r.db.QueryRow(query,
		&user.Nickname).
		Scan(&user.About,
			&user.Email,
			&user.Fullname,
			&user.Nickname)
	fmt.Println(query)
	fmt.Println(err)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &models.Error{Message: models.NotExist.Error()}
		} else {
			return &models.Error{Message: models.Exist.Error()}
		}
	}
	return nil
}
