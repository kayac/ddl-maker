package example

import (
	"database/sql"
	"time"

	"github.com/kayac/ddl-maker/dialect"
	"github.com/kayac/ddl-maker/dialect/mysql"
)

type User struct {
	Id        uint64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u User) Table() string {
	return "player"
}

func (u User) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id")
}

type Entry struct {
	Id        int32   `ddl:"auto"`
	Title     string  `ddl:"size=100"`
	Public    bool    `ddl:"default=0"`
	Content   *string `ddl:"text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (e Entry) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id", "created_at")
}

func (e Entry) Indexes() dialect.Indexes {
	return dialect.Indexes{
		mysql.AddUniqueIndex("created_at_uniq_idx", "created_at"),
		mysql.AddIndex("title_idx", "title"),
		mysql.AddIndex("created_at_idx", "created_at"),
	}
}

type PlayerComment struct {
	Id        int32          `ddl:"auto,size=100" json:"id"`
	PlayerID  int32          `json:"player_id"`
	EntryID   int32          `json:"entry_id"`
	Comment   sql.NullString `json:"comment" ddl:"null,size=99"`
	CreatedAt time.Time      `ddl:"-" json:"created_at"`
	updatedAt time.Time
}

func (pc PlayerComment) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id")
}

func (pc PlayerComment) Indexes() dialect.Indexes {
	return dialect.Indexes{
		mysql.AddIndex("player_id_entry_id_idx", "player_id", "entry_id"),
	}
}

type Bookmark struct {
	Id        int32     `ddl:"size=100" json:"id"`
	UserId    int32     `json:"user_id"`
	EntryId   int32     `json:"entry_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (b Bookmark) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id")
}

func (b Bookmark) Indexes() dialect.Indexes {
	return dialect.Indexes{
		mysql.AddUniqueIndex("user_id_entry_id", "user_id", "entry_id"),
	}
}
