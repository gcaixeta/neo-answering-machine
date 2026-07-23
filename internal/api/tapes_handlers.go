package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gcaixeta/neo-answering-machine/tape"
	"github.com/google/uuid"
)

type TapeHandler struct {
	repo tape.Repository
}

const tapesDir = "uploads/tapes"

func (h *TapeHandler) UploadNewTape(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 20<<20)

	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, "error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("audio")
	if err != nil {
		http.Error(w, "error reading file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Campos de texto — vêm de FormValue, não de FormFile
	mailboxID := r.FormValue("mailboxId")
	userID := r.FormValue("userId")
	recordedAtRaw := r.FormValue("recordedAt")

	if mailboxID == "" || userID == "" || recordedAtRaw == "" {
		http.Error(w, "mailboxId, userId e recordedAt são obrigatórios", http.StatusBadRequest)
		return
	}

	recordedAt, err := time.Parse(time.RFC3339, recordedAtRaw)
	if err != nil {
		http.Error(w, "recordedAt inválido, use formato ISO 8601", http.StatusBadRequest)
		return
	}

	mailboxUUID, err := uuid.Parse(mailboxID)
	if err != nil {
		http.Error(w, "error parsing mailboxId:"+err.Error(), http.StatusBadRequest)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "error parsing userId:"+err.Error(), http.StatusBadRequest)
		return
	}

	t, err := tape.NewTape(mailboxUUID, userUUID, recordedAt)
	if err != nil {
		http.Error(w, "error creating new tape", http.StatusInternalServerError)
		return
	}

	os.MkdirAll(tapesDir, 0755)
	dstPath := filepath.Join(tapesDir, t.ID.String())
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "error creating file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "error copying audio to file", http.StatusInternalServerError)
		return
	}

	// Save new tape into database
	if err := h.repo.Save(r.Context(), t); err != nil {
		dst.Close()
		os.Remove(dstPath)
		http.Error(w, "error saving tape: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":   "ok",
		"filename": header.Filename,
		"tapeId":   t.ID.String(),
	})
}
