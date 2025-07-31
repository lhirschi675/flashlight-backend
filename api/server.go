package api

import (
	"encoding/json"
	"flashlight-backend/db"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
	StuID     string    `json:"stu_id"`
	Age       int       `json:"age"`
	Progress  int       `json:"progress"`
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

	s.Use(mux.CORSMethodMiddleware(s.Router))
	s.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*") // or specify frontend origin
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	s.routes()
	return s
}

func (s *Server) routes() {
	s.HandleFunc("/students", s.createStudent()).Methods("POST")
	s.HandleFunc("/students", s.listStudents()).Methods("GET")
	s.HandleFunc("/students/{id}", s.getStudent()).Methods("GET")
	s.HandleFunc("/students/{id}", s.updateStudent()).Methods("PUT")
	s.HandleFunc("/students/{id}", s.deleteStudent()).Methods("DELETE")
}

func (s *Server) createStudent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student Student
		if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		student.ID = uuid.New()
		student.StuID = "STU" + student.ID.String()
		student.CreatedAt = time.Now()
		student.UpdatedAt = time.Now()
		student.Progress = 0

		if err := s.db.Create(&student).Error; err != nil {
			http.Error(w, "Failed to create student", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
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

func (s *Server) updateStudent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		var input Student
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		var student Student
		if err := s.db.First(&student, "id = ?", id).Error; err != nil {
			http.NotFound(w, r)
			return
		}

		student.Name = input.Name
		student.Course = input.Course
		student.Status = input.Status
		student.Email = input.Email
		student.Phone = input.Phone
		student.Age = input.Age
		student.UpdatedAt = time.Now()

		if err := s.db.Save(&student).Error; err != nil {
			http.Error(w, "Failed to update student", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(student)
	}
}

func (s *Server) deleteStudent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		if err := s.db.Delete(&Student{}, "id = ?", id).Error; err != nil {
			http.Error(w, "Failed to delete student", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
