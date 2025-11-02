package services

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/hababisha/ABT/task4/models"
)

type ReservationCanceller interface {
	CancelReservation(bookID int)
	IsBookReserved(bookID int) bool
}

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
	SearchBooks(query string) []models.Book
	AddMember(member models.Member)
	RemoveMember(memberID int) error
	UpdateBookStatus(bookID int, status models.Status)
	FindBookByID(bookID int) (models.Book, bool)
	FindMemberByID(memberID int) (models.Member, bool)
}

type Library struct {
	Books map[int]models.Book
	Members map[int]models.Member
	mu sync.Mutex
	ResCanceller ReservationCanceller
}

func NewLibrary() *Library {
	return &Library{
		Books:   make(map[int]models.Book),
		Members: make(map[int]models.Member),
	}
}

func (l *Library) SetReservationHandler(canceller ReservationCanceller) {
	l.ResCanceller = canceller
}

func (l *Library) FindBookByID(bookID int) (models.Book, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	book, exists := l.Books[bookID]
	return book, exists
}

func (l *Library) FindMemberByID(memberID int) (models.Member, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	member, exists := l.Members[memberID]
	return member, exists
}

func (l *Library) UpdateBookStatus(bookID int, status models.Status) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if book, exists := l.Books[bookID]; exists {
		book.Status = status
		l.Books[bookID] = book
	}
}

func (l *Library) AddBook(book models.Book) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Books[book.ID] = book
}

func (l *Library) RemoveBook(bookID int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.Books, bookID)
}

func (l *Library) AddMember(member models.Member) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Members[member.ID] = member
}

func (l *Library) RemoveMember(memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	member, exists := l.Members[memberID]
	if !exists {
		return fmt.Errorf("member with ID %d not found", memberID)
	}
	if len(member.BorrowedBooks) > 0 {
		return fmt.Errorf("cannot remove member %s (ID: %d) as they have %d borrowed book(s)", member.Name, memberID, len(member.BorrowedBooks))
	}
	delete(l.Members, memberID)
	return nil
}

func (l *Library) BorrowBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, exist := l.Books[bookID]
	if !exist {
		return fmt.Errorf("book not found")
	}

	if book.Status == models.Borrowed {
		return errors.New("book is currently borrowed")
	}

	member, ok := l.Members[memberID]
	if !ok {
		return fmt.Errorf("member not found")
	}

	if book.Status == models.Reserved {
		if l.ResCanceller != nil && l.ResCanceller.IsBookReserved(bookID) {
			l.ResCanceller.CancelReservation(bookID)
		} else {
			return fmt.Errorf("book is reserved and not available for general borrowing")
		}
	} else if book.Status != models.Available {
		return fmt.Errorf("book status is %s, not available for borrowing", book.Status)
	}

	book.Status = models.Borrowed
	l.Books[bookID] = book
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.Members[memberID] = member
	return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, exist := l.Books[bookID]
	if !exist {
		return fmt.Errorf("book not found")
	}
	if book.Status != models.Borrowed {
		return errors.New("book is not currently borrowed")
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

	book.Status = models.Available
	l.Books[bookID] = book
	l.Members[memberID] = member
	return nil
}

func (l *Library) ListAvailableBooks() []models.Book {
	l.mu.Lock()
	defer l.mu.Unlock()
	availableBooks := []models.Book{}
	for _, book := range l.Books {
		if book.Status == models.Available {
			availableBooks = append(availableBooks, book)
		}
	}
	return availableBooks
}

func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	l.mu.Lock()
	defer l.mu.Unlock()
	member, ok := l.Members[memberID]
	if !ok {
		return []models.Book{}
	}
	return member.BorrowedBooks
}

func (l *Library) SearchBooks(query string) []models.Book {
	l.mu.Lock()
	defer l.mu.Unlock()
	results := []models.Book{}
	lowerQuery := strings.ToLower(query)

	for _, book := range l.Books {
		lowerTitle := strings.ToLower(book.Title)
		lowerAuthor := strings.ToLower(book.Author)

		if strings.Contains(lowerTitle, lowerQuery) || strings.Contains(lowerAuthor, lowerQuery) {
			results = append(results, book)
		}
	}
	return results
}