package main

// import (
// 	"database/sql"
// 	"log"

// )

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// echoのサーバーのインスタンス？を作成
	e := echo.New()
	// エンドポイント
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello Echo")
	})
	e.POST("/signup", signUp)
	e.Logger.Fatal(e.Start(":1323"))
}

func signUp(c echo.Context) error {
	db, err := sql.Open("sqlite3", "./mydb.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// name, passwordを受け取る
	// TODO nameに値が入っているか検証の必要ある
	name := c.FormValue("name")
	// TODO 少なすぎる文字数の場合エラーにしたい
	password := c.FormValue("password")

	// TODO 同じnameのuserがいないか検証
	// rows, err := db.Query("SELECT name FROM users WHERE name = ?", name)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// passwordをハッシュ値にする
	// bcrypt.DefaultCost costとは
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	// uuidを生成
	uuidWithHyphen := uuid.New()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)

	// DBに保存
	res, err := db.Exec("INSERT INTO users (id, name, password) VALUES (?, ?, ?)", uuid, name, hash)
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, uuid)
}
