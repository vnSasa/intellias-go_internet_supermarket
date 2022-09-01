package main

import (
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

type User struct {
  Email string `json : "email"`
  FirstName string `json : "first_name"`
  LastName string `json : "last_name"`
  Password string `json : "password"`
  CPassword string `json : "confirm_password"`
}

func main()  {
  db, err := sql.Open("mysql", "sasa:110513@tcp(127.0.0.1:3306)/go-periodicals")
  if err != nil {
    panic(err)
  }
  defer db.Close()

  // insert, err := db.Query("INSERT INTO users (email, first_name, last_name, password, confirm_password)" +
  //                 " VALUES('sasa.yar@gmail.com', 'Olexandr', 'Yaremechko', '110513', '110513')")
  // if err != nil {
  //   panic(err)
  // }
  // defer insert.Close()

  res, err := db.Query("SELECT email, first_name, last_name, password, confirm_password FROM users")
  if err != nil {
    panic(err)
  }

  for res.Next() {
    var user User
    err = res.Scan(&user.Email, &user.FirstName, &user.LastName, &user.Password, &user.CPassword)
    if err != nil {
      panic(err)
    }

    fmt.Printf("User: email - %s, first_name - %s, last_name - %s, password - %s, confirm_password - %s\n",
                user.Email, user.FirstName, user.LastName, user.Password, user.CPassword)
  }

  fmt.Println("Work!!!")
}
