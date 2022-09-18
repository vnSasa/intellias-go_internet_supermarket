package main

import (
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "net/http"
  "html/template"
  "strings"
)

var database *sql.DB

type SignupForm struct {
  Email string `json : "email"`
  FirstName string `json : "first_name"`
  LastName string `json : "last_name"`
  Password string `json : "password"`
  CPassword string `json : "confirm_password"`
}

type LoginForm struct {
  Email string `json : "email"`
  Password string `json : "password"`
}

func main()  {
  handleFunc()
}

func handleFunc()  {
  db, err := sql.Open("mysql", "sasa:110513@tcp(127.0.0.1:3306)/go-periodicals")
  if err != nil {
    panic(err)
  }
  database = db
  defer db.Close()

  http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static/"))))
  http.HandleFunc("/", index)
  http.HandleFunc("/login", login)
  http.HandleFunc("/signup", signup)
  http.HandleFunc("/new_signup", new_signup)
  http.HandleFunc("/main_page", main_page)
  http.HandleFunc("/create", create)
  http.HandleFunc("/save_product", save_product)

  fmt.Println("Server is listening...")
  http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
  t, err := template.ParseFiles("templates/index.html")
  if err != nil {
    fmt.Fprintf(w, err.Error())
  }
  t.ExecuteTemplate(w, "index", nil)
}

func main_page(w http.ResponseWriter, r *http.Request) {
  t, err := template.ParseFiles("templates/main_page.html", "templates/header.html",
                                "templates/footer.html")
  if err != nil {
    fmt.Fprintf(w, err.Error())
  }
  t.ExecuteTemplate(w, "main_page", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
  if r.Method == "POST" {
    err := r.ParseForm()
    if err != nil {
      panic(err)
    }

    email := r.FormValue("email")
    password := r.FormValue("password")
    var pwd string
    u := "SELECT password FROM users WHERE email = ?"
    row := database.QueryRow(u, email)
    err = row.Scan(&pwd)
    if err != nil {
      http.Redirect(w, r, "/login.html", http.StatusSeeOther)
    }
    if strings.Compare(pwd, password) == 0 {
      http.Redirect(w, r, "/main_page", http.StatusSeeOther)
    }
  }
}

func signup(w http.ResponseWriter, r *http.Request) {
  t, err := template.ParseFiles("templates/signup.html", "templates/header.html",
                                "templates/footer.html")
  if err != nil {
    fmt.Fprintf(w, err.Error())
  }
  t.ExecuteTemplate(w, "signup", nil)
}

func new_signup(w http.ResponseWriter, r *http.Request) {
  if r.Method == "POST" {
    err := r.ParseForm()
    if err != nil {
      panic(err)
    }

    email := r.FormValue("email")
    firstName := r.FormValue("first_name")
    lastName := r.FormValue("last_name")
    password := r.FormValue("password")
    cPassword := r.FormValue("confirm_password")

    _, err = database.Exec("INSERT INTO users (email, first_name, last_name, "+
                          "password, confirm_password) VALUES (?, ?, ?, ?, ?)",
                          email, firstName, lastName, password, cPassword)

    if err != nil {
      panic(err)
    }

    http.Redirect(w, r, "/", 301)
  }
}

func create(w http.ResponseWriter, r *http.Request) {
  t, err := template.ParseFiles("templates/create.html", "templates/header.html",
                                "templates/footer.html")
    if err != nil {
      fmt.Fprintf(w, err.Error())
    }
    t.ExecuteTemplate(w, "create", nil)
}

func save_product(w http.ResponseWriter, r *http.Request) {

  if r.Method == "POST" {
    err := r.ParseForm()
    if err != nil {
      panic(err)
    }

    description := r.FormValue("description")
    name_product := r.FormValue("name_product")
    price := r.FormValue("price")

    _, err = database.Exec("INSERT INTO products (description, name_product, price) " +
                        "VALUES (?, ?, ?)", description, name_product, price)
    if err != nil {
      panic(err)
    }

    http.Redirect(w, r, "/main_page", http.StatusSeeOther)
  }
}
