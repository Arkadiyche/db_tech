package repository

import (
	"github.com/Arkadiyche/bd_techpark/internal/pkg/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
)

type ForumRepository struct {
	db *pgx.ConnPool
}

func NewForumRepository(db *pgx.ConnPool) *ForumRepository {
	return &ForumRepository{
		db: db,
	}
}

func (r *ForumRepository) Insert(forum models.Forum) (f *models.Forum, error *models.Error)  {
	err := r.db.QueryRow("INSERT INTO forums VALUES ($1, $2, $3, $4, $5) RETURNING nickname",
		forum.Slug,
		forum.Title,
		forum.User,
		forum.Posts,
		forum.Threads).Scan(&forum.User)
	if err != nil {
		//fmt.Println("111", err.(pgx.PgError).Code)
		switch err.(pgx.PgError).Code {
		case pgerrcode.ForeignKeyViolation:
			return nil, &models.Error{Message: models.NotExist.Error()}
		case pgerrcode.UniqueViolation:
			err := r.db.QueryRow("SELECT slug, title, nickname, posts, threads FROM forums WHERE slug = $1", forum.Slug).
				Scan(&forum.Slug,
					&forum.Title,
					&forum.User,
					&forum.Posts,
					&forum.Threads)
			if err != nil {
				return nil, &models.Error{Message: err.Error()}
			}
			return &forum, &models.Error{Message: models.Exist.Error()}
		default:
			return nil, &models.Error{Message: err.Error()}
		}
	}
	return &forum, nil
}

func (r *ForumRepository) GetForum(slug string) (f *models.Forum, error *models.Error) {
	forum := models.Forum{}
	err := r.db.QueryRow("SELECT slug, title, nickname, posts, threads FROM forums WHERE slug = $1", slug).
		Scan(&forum.Slug,
			&forum.Title,
			&forum.User,
			&forum.Posts,
			&forum.Threads)
	if err != nil {
		return nil, &models.Error{Message: models.NotExist.Error()}
	}
	return &forum, nil
}

func (r *ForumRepository) GetForumUsers(slug string, desc bool, since string, limit int) (us *models.Users, e *models.Error)  {
	users := models.Users{}
	query := "SELECT users.about, users.email, users.fullname, users.nickname " +
		"FROM forum_users " +
		"JOIN users on users.nickname = forum_users.author " +
		"WHERE slug = $1 "
	rows := &pgx.Rows{}
	var err error
	if limit > 0 && since != "" {
		if desc {
			query += "AND lower(users.nickname) < lower($2::text) ORDER BY users.nickname COLLATE \"C\" DESC LIMIT $3"
		} else {
			query += "AND lower(users.nickname)  > lower($2::text) ORDER BY users.nickname COLLATE \"C\" ASC LIMIT $3"
		}
		rows, err = r.db.Query(query, &slug, &since, &limit)
	} else {
		if limit > 0 {
			if desc {
				query += "ORDER BY users.nickname COLLATE \"C\" DESC LIMIT $2"
			} else {
				query += "ORDER BY users.nickname COLLATE \"C\" ASC LIMIT $2"
			}
			rows, err = r.db.Query(query, &slug, &limit)
		} else if since != "" {
			if desc {
				query += "AND lower(users.nickname) < lower($2::text) ORDER BY users.nickname COLLATE \"C\" DESC "
			} else {
				query += "AND lower(users.nickname) > lower($2::text) ORDER BY users.nickname COLLATE \"C\" ASC "
			}
			rows, err = r.db.Query(query, &slug, &since)
		} else {
			if desc {
				query += "ORDER BY users.nickname COLLATE \"C\" DESC "
			} else {
				query += "ORDER BY users.nickname COLLATE \"C\" ASC "
			}
			rows, err = r.db.Query(query, &slug)
		}
	}
	//fmt.Println(query)
	if err != nil {
		return nil,  &models.Error{Message: err.Error()}
	}
	for rows.Next() {
		user := models.User{}

		err := rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
		if err != nil {
			return nil, &models.Error{Message: err.Error()}
		}
		users = append(users, user)
	}
	return &users, nil
}
