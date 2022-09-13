package main

import (
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "net/http"
  "html/template"
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
  http.HandleFunc("/create", create)
  http.HandleFunc("/save_product", save_product)

  fmt.Println("Server is listening...")
  http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
  t, err := template.ParseFiles("templates/index.html", "templates/header.html",
                                "templates/footer.html")
  if err != nil {
    fmt.Fprintf(w, err.Error())
  }
  t.ExecuteTemplate(w, "index", nil)
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

    http.Redirect(w, r, "/", http.StatusSeeOther)
  }
}
