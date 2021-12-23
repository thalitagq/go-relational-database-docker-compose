package main

import (
    "database/sql"
    "fmt"
		"time"
    "log"

    _ "github.com/go-sql-driver/mysql"
		"github.com/jinzhu/gorm"
)

type Album struct {
    ID     int64
    Title  string
    Artist string
    Price  float32
}

var db *gorm.DB

func main() {
	db = sqlConnect()
  defer db.Close()

	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
			log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	alb, err := albumByID(2)
	if err != nil {
			log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", alb)

	albID, err := addAlbum(Album{
    Title:  "The Modern Sound of Betty Carter",
    Artist: "Betty Carter",
    Price:  49.99,
	})
	if err != nil {
			log.Fatal(err)
	}
	fmt.Printf("ID of added album: %v\n", albID)

}

func sqlConnect() (database *gorm.DB) {
  DBMS := "mysql"
  USER := "user"
  PASS := "password"
  PROTOCOL := "tcp(db:3306)"
  DBNAME := "db"

  CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME 
  
  count := 0
  db, err := gorm.Open(DBMS, CONNECT)
  if err != nil {
    for {
      if err == nil {
        fmt.Println("")
        break
      }
      fmt.Print(".")
      time.Sleep(time.Second)
      count++
      if count > 180 {
        fmt.Println("")
        fmt.Println("DB connection failure")
        panic(err)
      }
      db, err = gorm.Open(DBMS, CONNECT)
    }
  }
  fmt.Println("DB connection successful")

  return db
}

func albumsByArtist(name string) ([]Album, error) {
    // An albums slice to hold data from returned rows.
    var albums []Album
		// Raw SQL
		rows, err := db.Raw("select * from album where artist = ?", name).Rows()
    if err != nil {
        return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
    }
    defer rows.Close()
    // Loop through rows, using Scan to assign column data to struct fields.
    for rows.Next() {
        var alb Album
        if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
            return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
        }
        albums = append(albums, alb)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
    }

    return albums, nil
}

func albumByID(id int64) (Album, error) {
    // An album to hold data from the returned row.
    var alb Album

    row := db.Raw("SELECT * FROM album WHERE id = ?", id).Row()
    if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
        if err == sql.ErrNoRows {
            return alb, fmt.Errorf("albumsById %d: no such album", id)
        }
        return alb, fmt.Errorf("albumsById %d: %v", id, err)
    }
    return alb, nil
}

func addAlbum(alb Album) (int64, error) {

	album := Album{Title: alb.Title, Artist: alb.Artist, Price: alb.Price}

	result := db.Table("album").Create(&album) // pass pointer of data to Create
	fmt.Println("NEW ID", album.ID)
	if result.Error  != nil {
		return 0, fmt.Errorf("addAlbum: %v", result.Error )
	}

	return album.ID, nil
}
