package handler

import (
	"encoding/json"
	"fmt"
	"ggkit_learn_service/internals/app/models"
	"ggkit_learn_service/internals/app/processor"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CommentHandler struct {
	process *processor.CommentProcessor
}

func NewCommentHandler(process *processor.CommentProcessor) *CommentHandler {
	return &CommentHandler{process: process}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {

	w, r, err := UserIdentify(w, r)
	if err != nil {
		WrapErrorWithStatus(w, err, http.StatusUnauthorized)
		return
	}
	user_id := w.Header().Get(UserCtx)

	var comment models.Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		WrapErrorWithStatus(w, fmt.Errorf("failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	if comment.ThemeID == 0 {
		http.Error(w, "theme_id is required", http.StatusBadRequest)
		return
	}
	if comment.Content == "" {
		http.Error(w, "content is required", http.StatusBadRequest)
		return
	}

	comment.UserID = user_id

	err = h.process.CreateComment(comment)
	if err != nil {
		WrapErrorWithStatus(w, err, http.StatusInternalServerError)
		return
	}

	WrapOK(w, map[string]interface{}{
		"status":  "ok",
		"message": "Comment created successfully",
		"data":    comment,
	})
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	w, r, err := UserIdentify(w, r)
	if err != nil {
		WrapErrorWithStatus(w, err, http.StatusUnauthorized)
		return
	}
	user_id := w.Header().Get(UserCtx)
	user_role := w.Header().Get(UserRole)
	vars := mux.Vars(r)
	commentID, err := strconv.Atoi(vars["comment_id"])
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	err = h.process.DeleteComment(commentID, user_id,user_role)
	if err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CommentHandler) GetCommentsByThemeID(w http.ResponseWriter, r *http.Request) {
	w, r, err := UserIdentify(w, r)
	if err != nil {
		WrapErrorWithStatus(w, err, http.StatusUnauthorized)
		return
	}
	user_id := w.Header().Get(UserCtx)
	user_role := w.Header().Get(UserRole)
	vars := mux.Vars(r)
	themeID, err := strconv.Atoi(vars["theme_id"])
	if err != nil {
		http.Error(w, "Invalid theme ID", http.StatusBadRequest)
		return
	}

	comments, err := h.process.GetCommentsByThemeID(themeID, user_id, user_role)
	if err != nil {
		http.Error(w, "Failed to get comments", http.StatusInternalServerError)
		return
	}

	WrapOK(w, map[string]interface{}{
		"status": "ok",
		"data":   comments,
	})
}

func (h *CommentHandler) CreateReply(w http.ResponseWriter, r *http.Request) {
	w, r, err := UserIdentify(w, r)
	if err != nil {
		WrapErrorWithStatus(w, err, http.StatusUnauthorized)
		return
	}
	user_role := w.Header().Get(UserRole)
	user_id := w.Header().Get(UserCtx)
	if user_role != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var reply models.AdminReply
	err = json.NewDecoder(r.Body).Decode(&reply)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	reply.UserID = user_id
	err = h.process.CreateReply(reply)
	if err != nil {
		http.Error(w, "Failed to create reply", http.StatusInternalServerError)
		return
	}

	WrapOK(w, map[string]interface{}{
		"status": "ok",
		"data":   "ответ на комментарий отправлен",
	})
}

func (h *CommentHandler) GetRepliesByCommentID(w http.ResponseWriter, r *http.Request) {
	w, r, err := UserIdentify(w, r)
	if err != nil {
		WrapErrorWithStatus(w, err, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	commentID, err := strconv.Atoi(vars["comment_id"])
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}
	replies, err := h.process.GetRepliesByCommentID(commentID)
	if err != nil {
		http.Error(w, "Failed to get replies", http.StatusInternalServerError)
		return
	}
	WrapOK(w, map[string]interface{}{
		"status": "ok",
		"data":   replies,
	})
}
