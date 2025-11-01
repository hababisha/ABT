# 📚 Console-Based Library Management System Documentation

This document outlines the architecture, data structures, and core functionalities of the simple console-based library management system implemented in Go.

---

## 1. System Overview

The system is a simple application designed to manage books and members within a library context. It uses Go's features like **structs**, **interfaces**, **maps**, and **slices** to create a clear and maintainable codebase following the **Service/Controller/Model** pattern.

## 2. Project Structure

The project follows a standard, organized folder structure:

| Folder/File | Description |
| :--- | :--- |
| `main.go` | The entry point of the application. Initializes the library and runs the main console menu. |
| `controllers/` | Contains logic for handling user input/output (console interaction) and invoking service layer methods. |
| `models/` | Contains the data structures (structs) used throughout the application. |
| `services/` | Contains the core business logic, data manipulation, and adherence to the `LibraryManager` interface. |
| `docs/` | Contains system documentation. |
| `go.mod` | Defines the module path and dependencies. |

---

## 3. Data Models (`models/`)

### 3.1. Book Struct (`models/book.go`)

Represents a book entity in the library.

| Field | Type | Description |
| :--- | :--- | :--- |
| `ID` | `int` | Unique identifier for the book. |
| `Title` | `string` | The title of the book. |
| `Author` | `string` | The author's name. |
| `Status` | `string` | Indicates the book's availability. Values are `"Available"` or `"Borrowed"`. **Default is "Available"**. |

### 3.2. Member Struct (`models/member.go`)

Represents a library member.

| Field | Type | Description |
| :--- | :--- | :--- |
| `ID` | `int` | Unique identifier for the member. |
| `Name` | `string` | The member's name. |
| `BorrowedBooks` | `[]Book` | A slice of `Book` structs currently borrowed by the member. |

---

## 4. Service Layer and Interface (`services/library_service.go`)

### 4.1. LibraryManager Interface

The interface defines the contract for any type that manages library operations.

```go
type LibraryManager interface {
    AddBook(book models.Book)
    RemoveBook(bookID int)
    BorrowBook(bookID int, memberID int) error
    ReturnBook(bookID int, memberID int) error
    ListAvailableBooks() []models.Book
    ListBorrowedBooks(memberID int) []models.Book
}