package main

// import (
// 	"database/sql"
// 	"log"

// )

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string
	Name     string
	Password string
}

func main() {
	// echoのサーバーのインスタンス？を作成
	e := echo.New()
	// ! secretはランダムな値に
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	// エンドポイント
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello Echo")
	})
	e.POST("/signup", signUp)
	e.POST("/signin", signIn)

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
	_, err = db.Exec("INSERT INTO users (id, name, password) VALUES (?, ?, ?)", uuid, name, hash)
	if err != nil {
		log.Fatalf("同じ名前のユーザーが既に存在します: %v", err)
	}

	// sessionの作成
	session, _ := session.Get("session", c)
	session.Values["username"] = name
	session.Save(c.Request(), c.Response())

	return c.String(http.StatusOK, "ユーザー登録が完了しました")
}

func signIn(c echo.Context) error {
	// name, passwordを受け取る
	name := c.FormValue("name")
	password := c.FormValue("password")

	db, err := sql.Open("sqlite3", "./mydb.db")
	if err != nil {
		log.Fatalf("DBに接続できませんでした err:%v", err)
	}
	defer db.Close()

	// DBにnameと同じユーザーがいるか見る
	rows, err := db.Query("SELECT id, name, password FROM users WHERE name = ?", name)
	if err != nil {
		log.Fatalf("DBに同じnameを持つユーザーがいませんでした err:%v", err)
	}

	// ? errがなければDBに存在していると思う
	u := &User{}
	for rows.Next() {
		if err := rows.Scan(&u.ID, &u.Name, &u.Password); err != nil {
			log.Fatalf("getRows rows.Scan error err:%v", err)
		}
		fmt.Println(u)
	}

	fmt.Println(u.ID)
	fmt.Println(u.Name)
	fmt.Println(u.Password)

	// passwordがあっているか検証
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		log.Fatalf("passwordが間違っている err:%v", err)
	}

	// session
	session, _ := session.Get("session", c)
	session.Values["username"] = name
	session.Save(c.Request(), c.Response())

	return c.String(http.StatusOK, "ログインしました")
}
