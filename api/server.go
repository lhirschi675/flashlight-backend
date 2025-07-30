package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Student struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Server struct {
	*mux.Router
	studentList []Student
}

func NewServer() *Server {
	s := &Server{
		Router:      mux.NewRouter(),
		studentList: []Student{},
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.HandleFunc("/students", s.createStudent()).Methods("POST")
	s.HandleFunc("/students", s.listStudents()).Methods("GET")
	s.HandleFunc("/students/{id}", s.getStudent()).Methods("GET")
	s.HandleFunc("/students/{id}", s.deleteStudent()).Methods("DELETE")
}

func (s *Server) createStudent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		newStudent := Student{
			ID:   uuid.New(),
			Name: input.Name,
		}
		s.studentList = append(s.studentList, newStudent)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newStudent)
	}
}

func (s *Server) listStudents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(s.studentList)
	}
}

func (s *Server) getStudent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := mux.Vars(r)["id"]
		id, err := uuid.Parse(idParam)
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		for _, student := range s.studentList {
			if student.ID == id {
				json.NewEncoder(w).Encode(student)
				return
			}
		}
		http.NotFound(w, r)
	}
}

func (s *Server) deleteStudent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := mux.Vars(r)["id"]
		id, err := uuid.Parse(idParam)
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		for i, student := range s.studentList {
			if student.ID == id {
				// Delete from list
				s.studentList = append(s.studentList[:i], s.studentList[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		http.NotFound(w, r)
	}
}
