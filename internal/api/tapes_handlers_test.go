package api

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gcaixeta/neo-answering-machine/tape"
	"github.com/google/uuid"
)

type fakeTapeRepo struct {
	saved   *tape.Tape
	saveErr error
}

func (f *fakeTapeRepo) Save(ctx context.Context, t *tape.Tape) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.saved = t
	return nil
}
func (f *fakeTapeRepo) FindByID(ctx context.Context, id uuid.UUID) (*tape.Tape, error) {
	return nil, tape.ErrNotFound
}
func (f *fakeTapeRepo) ListByMailboxID(ctx context.Context, mailboxID uuid.UUID) ([]*tape.Tape, error) {
	return nil, nil
}

func buildMultipart(t *testing.T, fields map[string]string, includeFile bool) (*bytes.Buffer, string) {
	t.Helper()
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	for k, v := range fields {
		if err := w.WriteField(k, v); err != nil {
			t.Fatal(err)
		}
	}

	if includeFile {
		fw, err := w.CreateFormFile("audio", "message.wav")
		if err != nil {
			t.Fatal(err)
		}
		fw.Write([]byte("fake audio bytes"))
	}

	w.Close()
	return body, w.FormDataContentType()
}

func TestUploadNewTape(t *testing.T) {
	t.Cleanup(func() { os.RemoveAll(tapesDir) })

	validFields := map[string]string{
		"mailboxId":  uuid.New().String(),
		"userId":     uuid.New().String(),
		"recordedAt": time.Now().UTC().Format(time.RFC3339),
	}

	t.Run("happy path", func(t *testing.T) {
		repo := &fakeTapeRepo{}
		h := &TapeHandler{repo: repo}

		body, ct := buildMultipart(t, validFields, true)
		req := httptest.NewRequest(http.MethodPost, "/tape/upload", body)
		req.Header.Set("Content-Type", ct)
		rec := httptest.NewRecorder()

		h.UploadNewTape(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
		}
		if repo.saved == nil {
			t.Fatal("expected repo.Save to be called")
		}
		if _, err := os.Stat(tapesDir + "/" + repo.saved.ID.String()); err != nil {
			t.Fatalf("expected file to be written: %v", err)
		}
	})

	t.Run("repo save fails", func(t *testing.T) {
		before, _ := os.ReadDir(tapesDir)

		repo := &fakeTapeRepo{saveErr: errors.New("fk violation")}
		h := &TapeHandler{repo: repo}

		body, ct := buildMultipart(t, validFields, true)
		req := httptest.NewRequest(http.MethodPost, "/tape/upload", body)
		req.Header.Set("Content-Type", ct)
		rec := httptest.NewRecorder()

		h.UploadNewTape(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want 500, body = %s", rec.Code, rec.Body.String())
		}

		after, err := os.ReadDir(tapesDir)
		if err != nil {
			t.Fatalf("reading tapesDir: %v", err)
		}
		if len(after) != len(before) {
			t.Fatalf("expected no orphaned files left behind, before=%d after=%d", len(before), len(after))
		}
	})

	t.Run("wrong method", func(t *testing.T) {
		h := &TapeHandler{repo: &fakeTapeRepo{}}
		req := httptest.NewRequest(http.MethodGet, "/tape/upload", nil)
		rec := httptest.NewRecorder()
		h.UploadNewTape(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("status = %d, want 405", rec.Code)
		}
	})

	t.Run("missing audio file", func(t *testing.T) {
		h := &TapeHandler{repo: &fakeTapeRepo{}}
		body, ct := buildMultipart(t, validFields, false)
		req := httptest.NewRequest(http.MethodPost, "/tape/upload", body)
		req.Header.Set("Content-Type", ct)
		rec := httptest.NewRecorder()
		h.UploadNewTape(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400, body = %s", rec.Code, rec.Body.String())
		}
	})

	missingFieldCases := []string{"mailboxId", "userId", "recordedAt"}
	for _, field := range missingFieldCases {
		t.Run("missing "+field, func(t *testing.T) {
			fields := map[string]string{}
			for k, v := range validFields {
				if k != field {
					fields[k] = v
				}
			}
			h := &TapeHandler{repo: &fakeTapeRepo{}}
			body, ct := buildMultipart(t, fields, true)
			req := httptest.NewRequest(http.MethodPost, "/tape/upload", body)
			req.Header.Set("Content-Type", ct)
			rec := httptest.NewRecorder()
			h.UploadNewTape(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400, body = %s", rec.Code, rec.Body.String())
			}
		})
	}

	t.Run("invalid recordedAt", func(t *testing.T) {
		fields := map[string]string{}
		for k, v := range validFields {
			fields[k] = v
		}
		fields["recordedAt"] = "not-a-date"
		h := &TapeHandler{repo: &fakeTapeRepo{}}
		body, ct := buildMultipart(t, fields, true)
		req := httptest.NewRequest(http.MethodPost, "/tape/upload", body)
		req.Header.Set("Content-Type", ct)
		rec := httptest.NewRecorder()
		h.UploadNewTape(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400, body = %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("invalid mailboxId", func(t *testing.T) {
		fields := map[string]string{}
		for k, v := range validFields {
			fields[k] = v
		}
		fields["mailboxId"] = "not-a-uuid"
		h := &TapeHandler{repo: &fakeTapeRepo{}}
		body, ct := buildMultipart(t, fields, true)
		req := httptest.NewRequest(http.MethodPost, "/tape/upload", body)
		req.Header.Set("Content-Type", ct)
		rec := httptest.NewRecorder()
		h.UploadNewTape(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400, body = %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("invalid userId", func(t *testing.T) {
		fields := map[string]string{}
		for k, v := range validFields {
			fields[k] = v
		}
		fields["userId"] = "not-a-uuid"
		h := &TapeHandler{repo: &fakeTapeRepo{}}
		body, ct := buildMultipart(t, fields, true)
		req := httptest.NewRequest(http.MethodPost, "/tape/upload", body)
		req.Header.Set("Content-Type", ct)
		rec := httptest.NewRecorder()
		h.UploadNewTape(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400, body = %s", rec.Code, rec.Body.String())
		}
	})
}
