package database

import (
	"database/sql"
	"fmt"
	"log"
	"tstkbot/private"
	_ "github.com/lib/pq"
)

// Controller represents controller for database
type Controller struct {
	DataBase *sql.DB
}

// InitDatabase represents database initialization
func Init() *Controller {
        dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", private.DB_USER, private.DB_PASSWORD, private.DB_NAME)
        db, err := sql.Open("postgres", dbinfo)
        checkErr("Open db failed", err)

	err = db.Ping()
	checkErr("Ping db failed", err)

	return &Controller{DataBase: db}
}

func (c *Controller) AddUser(userId int) {
	exists := c.CheckUser(userId)
	if !exists {
		stmt, err := c.DataBase.Prepare("INSERT INTO " + DB_NAME + ".user(telegram_id) VALUES($1)")
		checkErr("Prepare insert user failed", err)
		_, err = stmt.Exec(userId)
		checkErr("Exec insert user failed", err)
	}
}

func (c *Controller) GetUser(userId int) int {
	var id int
	var telegramId int
	query := "SELECT * FROM " + DB_NAME + ".user where telegram_id=$1"
	err := c.DataBase.QueryRow(query, userId).Scan(&id, &telegramId)
	if err != nil && err != sql.ErrNoRows {
        	log.Fatalf("getting user failed %s %v", userId, err)
	}
	return id
}

func (c *Controller) CheckUser(userId int) bool {
	var exists bool
	query := "SELECT exists (SELECT * FROM " + DB_NAME + ".user where telegram_id=$1)"
	err := c.DataBase.QueryRow(query, userId).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
        	log.Fatalf("checking if row exists failed %s %v", userId, err)
	}
	return exists
}

func (c *Controller) AddJudge(userId int, phrase string) {
	exists := c.CheckUser(userId)
	if !exists {
		stmt, err := c.DataBase.Prepare("INSERT INTO " + DB_NAME + ".user(telegram_id) VALUES($1)")
		checkErr("Prepare insert user failed", err)
		_, err = stmt.Exec(userId)
		checkErr("Exec insert user failed", err)
	}
	id := c.GetUser(userId)
	stmt, err := c.DataBase.Prepare("INSERT INTO " + DB_NAME + ".judge(phrase, author_id) VALUES($1, $2)")
	checkErr("Prepare insert judge failed", err)
	_, err = stmt.Exec(phrase, id)
	checkErr("Exec insert judge failed", err)
}

func (c *Controller) GetJudge(phrase string) int {
	var id int
        var judge string
	var authorId int
	query := "SELECT * FROM " + DB_NAME + ".judge where phrase=$1"
	err := c.DataBase.QueryRow(query, phrase).Scan(&id, &judge, &authorId)
	if err != nil && err != sql.ErrNoRows {
        	log.Fatalf("getting judge failed %s %v", phrase, err)
	}
	return id
}

func (c *Controller) CheckJudge(phrase string) bool {
	var exists bool
	query := "SELECT exists (SELECT * FROM " + DB_NAME + ".judge where phrase=$1)"
	err := c.DataBase.QueryRow(query, phrase).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
        	log.Fatalf("checking if row exists failed %s %v", phrase, err)
	}
	return exists
}

func (c *Controller) AddVote(userId int, phrase string) {
	exists := c.CheckUser(userId)
	if !exists {
		stmt, err := c.DataBase.Prepare("INSERT INTO " + DB_NAME + ".user(telegram_id) VALUES($1)")
		checkErr("Prepare insert user failed", err)
		_, err = stmt.Exec(userId)
		checkErr("Exec insert user failed", err)
	}

	judgeExists := c.CheckJudge(phrase)
	id := c.GetUser(userId)
	if !judgeExists {
		stmt, err := c.DataBase.Prepare("INSERT INTO " + DB_NAME + ".judge(phrase, author_id) VALUES($1, $2)")
		checkErr("Prepare insert judge failed", err)
		_, err = stmt.Exec(phrase, id)
		checkErr("Exec insert judge failed", err)
	}

	judgeId := c.GetJudge(phrase)
	stmt, err := c.DataBase.Prepare("INSERT INTO " + DB_NAME + ".vote(judge_id, user_id) VALUES($1, $2)")
	checkErr("Prepare insert vote failed", err)
	_, err = stmt.Exec(judgeId, id)
	checkErr("Exec insert vote failed", err)
}

func checkErr(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)        
	}
}
