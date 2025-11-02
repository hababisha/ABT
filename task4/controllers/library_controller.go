package controllers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hababisha/ABT/task4/concurrency"
	"github.com/hababisha/ABT/task4/models"
	"github.com/hababisha/ABT/task4/services"
)

func Start(library *services.Library, worker *concurrency.ReservationWorker) {
	for {
		fmt.Println("\n===== Library Management System =====")
		fmt.Println("1. Add Book")
		fmt.Println("2. Remove Book")
		fmt.Println("3. Borrow Book")
		fmt.Println("4. Return Book")
		fmt.Println("5. List Available Books")
		fmt.Println("6. List Borrowed Books by Member")
		fmt.Println("7. Search Books")
		fmt.Println("8. Add Member")
		fmt.Println("9. Remove Member")
		fmt.Println("10. Reserve Book (Concurrent)")
		fmt.Println("11. Exit")
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
			fmt.Scanln()
			title, _ = bufio.NewReader(os.Stdin).ReadString('\n')
			title = strings.TrimSpace(title)
			fmt.Print("Enter Book Author: ")
			author, _ = bufio.NewReader(os.Stdin).ReadString('\n')
			author = strings.TrimSpace(author)

			library.AddBook(models.Book{ID: id, Title: title, Author: author, Status: models.Available})
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
				fmt.Printf("ID: %d | %s by %s | Status: %s\n", book.ID, book.Title, book.Author, book.Status)
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
			var query string
			fmt.Print("Enter search query (Title or Author): ")
			fmt.Scanln()
			query, _ = bufio.NewReader(os.Stdin).ReadString('\n')
			query = strings.TrimSpace(query)

			results := library.SearchBooks(query)
			if len(results) == 0 {
				fmt.Println("No books found matching the query.")
			} else {
				fmt.Println("--- Search Results ---")
				for _, book := range results {
					fmt.Printf("ID: %d | %s by %s | Status: %s\n", book.ID, book.Title, book.Author, book.Status)
				}
			}

		case 8:
			var id int
			var name string
			fmt.Print("Enter Member ID: ")
			fmt.Scan(&id)
			fmt.Print("Enter Member Name: ")
			fmt.Scanln()
			name, _ = bufio.NewReader(os.Stdin).ReadString('\n')
			name = strings.TrimSpace(name)

			library.AddMember(models.Member{ID: id, Name: name})
			fmt.Printf("Member %s added successfully!\n", name)

		case 9:
			var id int
			fmt.Print("Enter Member ID to remove: ")
			fmt.Scan(&id)

			err := library.RemoveMember(id)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Member removed successfully!")
			}

		case 10:
			var bookID, memberID int
			fmt.Print("Enter Book ID to reserve: ")
			fmt.Scan(&bookID)
			fmt.Print("Enter Member ID: ")
			fmt.Scan(&memberID)

			resultChan := make(chan error)

			worker.RequestQueue <- concurrency.ReservationRequest{
				BookID:   bookID,
				MemberID: memberID,
				Result:   resultChan,
			}

			select {
			case err := <-resultChan:
				if err != nil {
					fmt.Println("Reservation Failed:", err)
				} else {
					fmt.Println("Reservation Request **Sent and Processed** successfully!")
				}
			case <-time.After(5 * time.Second):
				fmt.Println("Reservation request timed out waiting for worker response.")
			}

		case 11:
			fmt.Println("Exiting... Goodbye!")
			return

		default:
			fmt.Println("Invalid choice. Try again.")
		}
	}
}