package concurrency

import (
	"fmt"
	"sync"
	"time"

	"github.com/hababisha/ABT/task4/models"
	"github.com/hababisha/ABT/task4/services"
)

type ReservationRequest struct {
	BookID   int
	MemberID int
	Result   chan error 
}

type ReservationStatus struct {
	BookID    int
	MemberID  int
	ExpiresAt time.Time
}
type ReservationWorker struct {
	RequestQueue chan ReservationRequest

	ActiveReservations map[int]ReservationStatus

	LibraryService services.LibraryManager

	mu sync.Mutex
}

func NewReservationWorker(library services.LibraryManager, queueSize int) *ReservationWorker {
	worker := &ReservationWorker{
		RequestQueue:     make(chan ReservationRequest, queueSize),
		ActiveReservations: make(map[int]ReservationStatus),
		LibraryService:   library,
	}
	go worker.Start()
	go worker.StartAutoCancellation()
	return worker
}

func (w *ReservationWorker) Start() {
	fmt.Println("[CONCURRENT] Reservation Worker started...")
	for req := range w.RequestQueue {
		req.Result <- w.processReservation(req.BookID, req.MemberID)
	}
}

func (w *ReservationWorker) StartAutoCancellation() {
	fmt.Println("[CONCURRENT] Auto-Cancellation Goroutine started...")
	ticker := time.NewTicker(1 * time.Second) 
	defer ticker.Stop()

	for range ticker.C {
		w.cleanupExpiredReservations()
	}
}
func (w *ReservationWorker) cleanupExpiredReservations() {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now()
	for bookID, status := range w.ActiveReservations {
		if status.ExpiresAt.Before(now) {
			// Find the book, reset its status to 'Available'
			book, exists := w.LibraryService.FindBookByID(bookID)
			if exists {
				book.Status = models.Available
				// Note: We use an internal method on the service here, assuming it's available.
				// For now, we'll assume we can safely update the book's status via the library service.
				w.LibraryService.UpdateBookStatus(bookID, models.Available)
				fmt.Printf("[CONCURRENT] Reservation for Book ID %d expired and was cancelled.\n", bookID)
			}
			delete(w.ActiveReservations, bookID)
		}
	}
}

func (w *ReservationWorker) processReservation(bookID int, memberID int) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, ok := w.ActiveReservations[bookID]; ok {
		return fmt.Errorf("book with ID %d is already reserved", bookID)
	}

	book, exists := w.LibraryService.FindBookByID(bookID)
	if !exists {
		return fmt.Errorf("book with ID %d not found", bookID)
	}
	if book.Status != models.Available {
		return fmt.Errorf("book with ID %d is currently %s", bookID, book.Status)
	}

	if _, exists := w.LibraryService.FindMemberByID(memberID); !exists {
		return fmt.Errorf("member with ID %d not found", memberID)
	}

	w.LibraryService.UpdateBookStatus(bookID, models.Reserved) 

	w.ActiveReservations[bookID] = ReservationStatus{
		BookID:    bookID,
		MemberID:  memberID,
		ExpiresAt: time.Now().Add(5 * time.Second), // 5-second timeout
	}

	fmt.Printf("[CONCURRENT] Book ID %d successfully reserved by Member ID %d. Expires in 5s.\n", bookID, memberID)
	return nil
}

func (w *ReservationWorker) IsBookReserved(bookID int) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	_, ok := w.ActiveReservations[bookID]
	return ok
}

func (w *ReservationWorker) CancelReservation(bookID int) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, ok := w.ActiveReservations[bookID]; ok {
		delete(w.ActiveReservations, bookID)
		w.LibraryService.UpdateBookStatus(bookID, models.Available)
		fmt.Printf("[CONCURRENT] Reservation for Book ID %d cancelled manually.\n", bookID)
	}
}