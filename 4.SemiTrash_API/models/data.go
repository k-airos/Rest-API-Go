package models

var DB []Book

type Book struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Author        Author `json:"author"`
	YearPublished int    `json:"yearPublished"`
}

type Author struct {
	Name     string `json:"name"`
	LastName string `json:"lastName"`
	BornYear int    `json:"bornYear"`
}

func init() {
	book1 := Book{
		ID:    1,
		Title: "Lord of the rings",
		Author: Author{
			Name:     "J.R",
			LastName: "Tolkin",
			BornYear: 1892,
		},
	}
	DB = append(DB, book1)
}

func FindBookById(id int) (Book, bool) {
	var book Book
	var found bool
	for _, b := range DB {
		if b.ID == id {
			book = b
			found = true
		}
	}
	return book, found

}
