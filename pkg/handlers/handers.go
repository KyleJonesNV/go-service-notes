package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/KyleJonesNV/go-service-notes/pkg/notes"
)

// @title           Notes API

const (
	ErrIDNotFound = "id not found"
	ErrInvalidPayload = "invalid payload"
)

type Response struct {
	StatusCode int
	Body       any
}

type ErrorBody struct {
	ErrorMsg string `json:"error,omitempty"`
}

type User struct {
	ID string `json:"id,omitempty"`
}

type Note struct {
	Title string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type Topic struct {
	Title string `json:"title,omitempty"`
	Notes []Note `json:"notes,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type InsertTopicRequest struct {
	UserID string `json:"userId,omitempty"`
	Title string `json:"title,omitempty"`
}

type DeleteTopicRequest struct {
	UserID string `json:"userId,omitempty"`
	Title string `json:"title,omitempty"`
}

type InsertNoteRequest struct {
	UserID string `json:"userId,omitempty"`
	Title string `json:"title,omitempty"`
	Note Note `json:"note,omitempty"`
}

type DeleteNoteRequest struct {
	UserID string `json:"userId,omitempty"`
	Title string `json:"title,omitempty"`
	NoteTitle string `json:"noteTitle,omitempty"`
}

type GetAllNotesRequest struct {
	UserID string `json:"userId,omitempty"`
	Title string `json:"title,omitempty"`
}


// getAll godoc
// @Summary      Get all movies
// @Description  Get all movies
// @Tags         movies
// @Accept       json
// @Produce      json
// @Success      200  {object}  []movies.Movie
// @Failure      400  {object}  ErrorBody
// @Failure      404  {object}  ErrorBody
// @Failure      500  {object}  ErrorBody
// @Router       /getAll [get]
func GetAllForUser(req *http.Request) Response {
	var user = User{}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return Response{
			http.StatusBadRequest,
			ErrorBody{err.Error()},
		}
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return Response{
			http.StatusBadRequest,
			ErrorBody{ErrInvalidPayload},
		}
	}	

	

	topics, err := notes.GetAllForUser(req.Context(), user.ID)
	if err != nil {
		return Response{http.StatusInternalServerError, ErrorBody{err.Error()}}
	}
	return Response{http.StatusOK, topics}
}

func InsertTopic(req *http.Request) Response {
	var insertTopicRequest = InsertTopicRequest{}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return Response{
			http.StatusBadRequest,
			ErrorBody{err.Error()},
		}
	}

	err = json.Unmarshal(body, &insertTopicRequest)
	if err != nil {
		return Response{
			http.StatusBadRequest,
			ErrorBody{ErrInvalidPayload},
		}
	}	

	err = notes.InsertTopic(req.Context(), insertTopicRequest.UserID, insertTopicRequest.Title)
	if err != nil {
		return Response{
			http.StatusInternalServerError,
			ErrorBody{fmt.Sprintf("insert, %s", err)},
		}
	}

	return Response{
		http.StatusOK,
		nil,
	}
}

func DeleteTopic(req *http.Request) Response {
	var deleteTopicRequest = DeleteTopicRequest{}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return Response{
			http.StatusBadRequest,
			ErrorBody{err.Error()},
		}
	}

	err = json.Unmarshal(body, &deleteTopicRequest)
	if err != nil {
		return Response{
			http.StatusBadRequest,
			ErrorBody{err.Error()},
		}
	}

	err = notes.DeleteTopic(req.Context(), deleteTopicRequest.UserID, deleteTopicRequest.Title)
	if err != nil {
		return Response{
			http.StatusInternalServerError,
			ErrorBody{err.Error()},
		}
	}

	return Response{
		http.StatusOK,
		nil,
	}
}

func InsertNote(req *http.Request) Response {
	var insertNoteRequest = InsertNoteRequest{}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return Response{
			http.StatusBadRequest,
			ErrorBody{err.Error()},
		}
	}

	err = json.Unmarshal(body, &insertNoteRequest)
	if err != nil {
		return Response{
			http.StatusBadRequest,
			ErrorBody{ErrInvalidPayload},
		}
	}	

	dbNote := notes.Note{
		Title: insertNoteRequest.Note.Title,
		Content: insertNoteRequest.Note.Content,
	}

	err = notes.InsertNote(req.Context(), insertNoteRequest.UserID, insertNoteRequest.Title, dbNote)
	if err != nil {
		return Response{
			http.StatusInternalServerError,
			ErrorBody{fmt.Sprintf("insert, %s", err)},
		}
	}

	return Response{
		http.StatusOK,
		nil,
	}
}

func DeleteNote(req *http.Request) Response {
	var deleteNoteRequest = DeleteNoteRequest{}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return Response{
			http.StatusBadRequest,
			ErrorBody{err.Error()},
		}
	}

	err = json.Unmarshal(body, &deleteNoteRequest)
	if err != nil {
		return Response{
			http.StatusBadRequest,
			ErrorBody{ErrInvalidPayload},
		}
	}	


	err = notes.DeleteNote(req.Context(), deleteNoteRequest.UserID, deleteNoteRequest.Title, deleteNoteRequest.NoteTitle)
	if err != nil {
		return Response{
			http.StatusInternalServerError,
			ErrorBody{fmt.Sprintf("delete, %s", err)},
		}
	}

	return Response{
		http.StatusOK,
		nil,
	}
}

func GetAllNotes(req *http.Request) Response {
	var getAllNotesRequest = GetAllNotesRequest{}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return Response{
			http.StatusBadRequest,
			ErrorBody{err.Error()},
		}
	}

	err = json.Unmarshal(body, &getAllNotesRequest)
	if err != nil {
		return Response{
			http.StatusBadRequest,
			ErrorBody{ErrInvalidPayload},
		}
	}	

	topics, err := notes.GetUserTopicByTitle(req.Context(), getAllNotesRequest.UserID, getAllNotesRequest.Title)
	if err != nil {
		return Response{http.StatusInternalServerError, ErrorBody{err.Error()}}
	}
	return Response{http.StatusOK, topics}
}
