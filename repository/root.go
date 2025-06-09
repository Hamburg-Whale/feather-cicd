package repository

import (
	"database/sql"
	"feather/config"
	"feather/types"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Repository struct {
	config *config.Config
	db     *sql.DB
}

const (
	users     = "feather.user"
	basecamps = "feather.basecamp"
	projects  = "feather.project"
)

func NewRepository(c *config.Config) (*Repository, error) {
	r := &Repository{config: c}
	var err error

	if r.db, err = sql.Open(c.DB.Database, c.DB.URL); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Repository) CreateUser(email string, password string, nickname string) error {
	if _, err := r.db.Exec("INSERT INTO feather.user(email, password, nickname) VALUES(?, ?, ?)",
		email, password, nickname); err != nil {
		return err
	}
	log.Println("CreateUser Query run successfully!")
	return nil
}

func (r *Repository) User(userId int64) (*types.User, error) {
	u := new(types.User)
	qs := query([]string{"SELECT * FROM", users, "WHERE id = ?"})
	if err := r.db.QueryRow(qs, userId).Scan(&u.ID, &u.Email, &u.Nickname); err != nil {
		if err := noResult(err); err != nil {
			return nil, err
		}
	}

	log.Println("User Query run successfully!")
	return u, nil
}

func (r *Repository) CreateBaseCamp(name string, url string, token string, userId int64) error {
	if _, err := r.db.Exec("INSERT INTO feather.basecamp(name, url, token, user_id) VALUES(?, ?, ?, ?)",
		name, url, token, userId); err != nil {
		return err
	}
	log.Println("CreateBaseCamp Query run successfully!")
	return nil
}

func (r *Repository) BaseCampsByUserId(userId int64) ([]*types.BaseCamp, error) {
	qs := query([]string{"SELECT * FROM", basecamps, "WHERE user_id = ?"})
	rows, err := r.db.Query(qs, userId)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var baseCamps []*types.BaseCamp

	for rows.Next() {
		b := new(types.BaseCamp)
		if err := rows.Scan(&b.ID, &b.Name, &b.URL, &b.Token, &b.User_ID); err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		baseCamps = append(baseCamps, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	log.Printf("Found %d basecamps for user_id=%d\n", len(baseCamps), userId)
	return baseCamps, nil
}

func (r *Repository) BaseCamp(baseCampId int64) (*types.BaseCamp, error) {
	b := new(types.BaseCamp)
	qs := query([]string{"SELECT * FROM", basecamps, "WHERE id = ?"})

	if err := r.db.QueryRow(qs, baseCampId).Scan(&b.ID, &b.Name, &b.URL, &b.Token, &b.User_ID); err != nil {
		if err := noResult(err); err != nil {
			return nil, err
		}
	}

	log.Println("BaseCamp Query run successfully!")
	return b, nil
}

func (r *Repository) CreateProject(name string, url string, owner string, private bool, baseCampId int64) error {
	if _, err := r.db.Exec("INSERT INTO feather.project(name, url, owner, private, basecamp_id) VALUES(?, ?, ?, ?, ?)",
		name, url, owner, private, baseCampId); err != nil {
		return err
	}
	log.Println("CreateProject Query run successfully!")
	return nil
}

func (r *Repository) ProjectsByBaseCampId(baseCampId int64) ([]*types.Project, error) {
	qs := query([]string{"SELECT * FROM", projects, "WHERE basecamp_id = ?"})
	rows, err := r.db.Query(qs, baseCampId)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var projects []*types.Project

	for rows.Next() {
		p := new(types.Project)
		if err := rows.Scan(&p.ID, &p.Name, &p.URL, &p.Owner, &p.Private, &p.BaseCamp_ID); err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	log.Printf("Found %d projects for basecamp_id=%d\n", len(projects), baseCampId)
	return projects, nil
}

func (r *Repository) Project(projectId int64) (*types.Project, error) {
	p := new(types.Project)
	qs := query([]string{"SELECT * FROM", projects, "WHERE id = ?"})

	if err := r.db.QueryRow(qs, projectId).Scan(&p.ID, &p.Name, &p.URL, &p.Owner, &p.Private, &p.BaseCamp_ID); err != nil {
		if err := noResult(err); err != nil {
			return nil, err
		}
	}

	log.Println("Project Query run successfully!")
	return p, nil
}

func query(qs []string) string {
	return strings.Join(qs, " ") + ";"
}

func noResult(err error) error {
	if strings.Contains(err.Error(), "sql: no rows in result set") {
		return nil
	} else {
		return err
	}
}
