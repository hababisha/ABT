package services

import (
	"errors"
	"fmt"

	"github.com/hababisha/ABT/task3/models"
)

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
}

type Library struct{
	Books map[int]models.Book
	Members map[int]models.Member
}

func NewLibrary() *Library{
	return &Library{
		Books: make(map[int]models.Book),
		Members: make(map[int]models.Member),
	}
}

func (l *Library) AddBook(book models.Book) {
	l.Books[book.ID] = book
}
func (l *Library) RemoveBook(bookID int){
	delete(l.Books, bookID)
}

func (l *Library) BorrowBook(bookID int, memberID int) error{
	book, exist := l.Books[bookID]
	if !exist{
		return fmt.Errorf("book not found")
	}
	if book.Status =="Borrowed"{
		return errors.New("book is borrowed")
	}

	member, ok := l.Members[memberID]
	if !ok{
		return fmt.Errorf("member not found")
	}

	book.Status = "Borrowed"
	l.Books[bookID] = book
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.Members[memberID] = member

	return nil
}


func (l *Library) ReturnBook(bookID int, memberID int) error {
	book, exist := l.Books[bookID]
	if !exist {
		return fmt.Errorf("book not found")
	}
	if book.Status == "available" {
		return errors.New("book is already returned")
	}
	member, ok := l.Members[memberID]
	if !ok {
		return fmt.Errorf("member not found")
	}
	found := false
	for i, b := range member.BorrowedBooks {
		if b.ID == bookID {
			member.BorrowedBooks = append(member.BorrowedBooks[:i], member.BorrowedBooks[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("this member did not borrow the book")
	}

	book.Status = "available"
	l.Books[bookID] = book
	l.Members[memberID] = member

	return nil
}

func (l *Library) ListAvailableBooks() []models.Book {
	availableBooks := []models.Book{}
	for _, book := range l.Books {
		if book.Status == "available" {
			availableBooks = append(availableBooks, book)
		}
	}
	fmt.Println(availableBooks)
	return availableBooks
}

func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	member, ok := l.Members[memberID]
	if !ok {
		return []models.Book{}
	}
	return member.BorrowedBooks
}

