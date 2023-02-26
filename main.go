package main

import (
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "net/http"
  "html/template"
  "golang.org/x/crypto/bcrypt"
)

var database *sql.DB

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

type SignupForm struct {
  Email string `json:"email"`
  FirstName string `json:"first_name"`
  LastName string `json:"last_name"`
  Password string `json:"password"`
  CPassword string `json:"confirm_password"`
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

    sData := SignupForm{
      Email:  r.FormValue("email"),
      FirstName:  r.FormValue("first_name"),
      LastName: r.FormValue("last_name"),
      Password: r.FormValue("password"),
      CPassword:  r.FormValue("confirm_password"),
    }

    cHash, _ := HashPassword(sData.CPassword)

    _, err = database.Exec("INSERT INTO users (email, first_name, last_name, "+
                          "password, confirm_password) VALUES (?, ?, ?, ?, ?)",
                          sData.Email, sData.FirstName, sData.LastName, sData.Password, cHash)

    if err != nil {
      panic(err)
    }

    http.Redirect(w, r, "/", 301)
  }
}

func HashPassword(password string) (string, error) {
  bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
  return string(bytes), err
}

type LoginForm struct {
  Email string `json:"email"`
  Password string `json:"password"`
}

func login(w http.ResponseWriter, r *http.Request) {
  if r.Method == "POST" {
    err := r.ParseForm()
    if err != nil {
      panic(err)
    }

    login := LoginForm{
      Email:  r.FormValue("email"),
      Password: r.FormValue("password"),
    }

    var pwd string
    u := "SELECT confirm_password FROM users WHERE email = ?"
    row := database.QueryRow(u, login.Email)
    err = row.Scan(&pwd)
    if err != nil {
      http.Redirect(w, r, "/login.html", http.StatusSeeOther)
    }
    err = bcrypt.CompareHashAndPassword([]byte(pwd), []byte(login.Password))
    if err == nil {
      http.Redirect(w, r, "/main_page", http.StatusSeeOther)
    }
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
