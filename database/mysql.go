package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	//_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
	"platzi.com/go/rest-ws/models"
)

type MySqlRepository struct {
	db *sql.DB
}

func NewMySqlRepository(url string) (*MySqlRepository, error) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	} /*
		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)*/
	return &MySqlRepository{db}, nil
}

func (repo *MySqlRepository) InsertUser(ctx context.Context, user *models.User) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO users (id, email, password) VALUES('"+user.Id+"','"+user.Email+"','"+user.Password+"')")
	return err
}

func (repo *MySqlRepository) InsertPost(ctx context.Context, post *models.Post) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO posts (id, post_content, user_id) VALUES('"+post.Id+"','"+post.PostContent+"','"+post.UserId+"')")
	return err
}

func (repo *MySqlRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, email FROM users WHERE id='"+id+"'")

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var user = models.User{}

	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Email); err == nil {
			return &user, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *MySqlRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, email, password FROM users WHERE email LIKE '"+email+"'")

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var user = models.User{}
	numero := 0

	for rows.Next() {
		numero++

		if err = rows.Scan(&user.Id, &user.Email, &user.Password); err == nil {
			return &user, err
		}

	}

	if err = rows.Err(); err != nil {
		log.Fatal("ERROR 2")
		return nil, err
	}

	return &user, nil
}

func (repo *MySqlRepository) GetPostById(ctx context.Context, id string) (*models.Post, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, post_content, created_at, user_id FROM posts WHERE id='"+id+"'")

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var post = models.Post{}

	for rows.Next() {
		if err = rows.Scan(&post.Id, &post.PostContent, &post.CreatedAt, &post.UserId); err == nil {
			return &post, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &post, nil
}

func (repo *MySqlRepository) Close() error {
	return repo.db.Close()
}

func (repo *MySqlRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	_, err := repo.db.ExecContext(ctx, "UPDATE posts SET post_content='"+post.PostContent+"' WHERE id='"+post.Id+"' AND user_id='"+post.UserId+"'")

	return err
}

func (repo *MySqlRepository) DeletePost(ctx context.Context, id string, userId string) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM posts WHERE id='"+id+"' AND user_id='"+userId+"'")

	return err
}

func (repo *MySqlRepository) ListPost(ctx context.Context, page uint64) ([]*models.Post, error) {

	rows, err := repo.db.QueryContext(ctx, "SELECT id, post_content, user_id, created_at FROM posts LIMIT 2 OFFSET "+fmt.Sprint(page*2))

	if err != nil {
		fmt.Println("ERROR")
		return nil, err
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var posts []*models.Post

	for rows.Next() {

		var post = models.Post{}

		if err = rows.Scan(&post.Id, &post.PostContent, &post.UserId, &post.CreatedAt); err != nil {
			fmt.Println(2)
			posts = append(posts, &post)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
