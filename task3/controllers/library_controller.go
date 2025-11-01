package controllers

import (
	"fmt"

	"github.com/hababisha/ABT/task3/models"
	"github.com/hababisha/ABT/task3/services"
)

func Start(library *services.Library) {
	for {
		fmt.Println("\n===== Library Management System =====")
		fmt.Println("1. Add Book")
		fmt.Println("2. Remove Book")
		fmt.Println("3. Borrow Book")
		fmt.Println("4. Return Book")
		fmt.Println("5. List Available Books")
		fmt.Println("6. List Borrowed Books by Member")
		fmt.Println("7. Exit")
		fmt.Print("Enter your choice: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			var id int
			var title, author string
			fmt.Print("Enter Book ID: ")
			fmt.Scan(&id)
			fmt.Print("Enter Book Title: ")
			fmt.Scan(&title)
			fmt.Print("Enter Book Author: ")
			fmt.Scan(&author)
			library.AddBook(models.Book{ID: id, Title: title, Author: author})
			fmt.Println("Book added successfully!")
		case 2:
			var id int
			fmt.Print("Enter Book ID to remove: ")
			fmt.Scan(&id)
			library.RemoveBook(id)
			fmt.Println("Book removed successfully!")
		case 3:
			var bookID, memberID int
			fmt.Print("Enter Book ID: ")
			fmt.Scan(&bookID)
			fmt.Print("Enter Member ID: ")
			fmt.Scan(&memberID)
			err := library.BorrowBook(bookID, memberID)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Book borrowed successfully!")
			}
		case 4:
			var bookID, memberID int
			fmt.Print("Enter Book ID: ")
			fmt.Scan(&bookID)
			fmt.Print("Enter Member ID: ")
			fmt.Scan(&memberID)
			err := library.ReturnBook(bookID, memberID)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Book returned successfully!")
			}
		case 5:
			fmt.Println("Available Books:")
			for _, book := range library.ListAvailableBooks() {
				fmt.Printf("ID: %d | %s by %s\n", book.ID, book.Title, book.Author)
			}
		case 6:
			var memberID int
			fmt.Print("Enter Member ID: ")
			fmt.Scan(&memberID)
			books := library.ListBorrowedBooks(memberID)
			if len(books) == 0 {
				fmt.Println("No borrowed books found for this member.")
			} else {
				fmt.Println("Borrowed Books:")
				for _, b := range books {
					fmt.Printf("ID: %d | %s by %s\n", b.ID, b.Title, b.Author)
				}
			}
		case 7:
			fmt.Println("Exiting... Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Try again.")
		}
	}
}