package api

import (
	"encoding/json"
	"flashlight-backend/db"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/gorm"
)

type Student struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `json:"name"`
	Course    string    `json:"course"`
	Status    string    `json:"status"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Age       int       `json:"age"`
	Grade     int       `json:"grade"`
}

type Server struct {
	*mux.Router
	db *gorm.DB
}

func NewServer() *Server {
	s := &Server{
		Router: mux.NewRouter(),
		db:     db.DB,
	}

	s.routes()

	handler := cors.New(cors.Options{
		// Replace http://localhost:5173 with the path you want to allow
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}).Handler(s.Router)

	s.Router = mux.NewRouter()
	s.PathPrefix("/").Handler(handler)

	return s
}

func (s *Server) routes() {
	s.HandleFunc("/students", s.createStudent()).Methods("POST")
	s.HandleFunc("/students/{id}", s.getStudent()).Methods("GET")
	s.HandleFunc("/students", s.listStudents()).Methods("GET")
}

func (s *Server) createStudent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student Student
		if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		student.ID = uuid.New()
		student.CreatedAt = time.Now()
		student.UpdatedAt = time.Now()
		student.Grade = 0

		if err := s.db.Create(&student).Error; err != nil {
			http.Error(w, "Failed to create student", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(student)
	}
}

func (s *Server) getStudent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		var student Student
		if err := s.db.First(&student, "id = ?", id).Error; err != nil {
			http.NotFound(w, r)
			return
		}

		json.NewEncoder(w).Encode(student)
	}
}

func (s *Server) listStudents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var students []Student
		if err := s.db.Find(&students).Error; err != nil {
			http.Error(w, "Failed to list students", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(students)
	}
}
