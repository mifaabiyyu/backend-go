package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mifaabiyyu/backend-go/utils"
)

type Handler struct {
	Service Service
	*utils.AppWrapper
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	// Validasi dengan validator
	if err := utils.Validate.Struct(req); err != nil {
		var messages []string
		for _, e := range err.(validator.ValidationErrors) {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				messages = append(messages, field+" is required")
			case "email":
				messages = append(messages, "invalid email format")
			case "min":
				messages = append(messages, field+" must be at least "+e.Param()+" characters")
			case "max":
				messages = append(messages, field+" must be at most "+e.Param()+" characters")
			case "alphanum":
				messages = append(messages, field+" must contain only letters and numbers")
			default:
				messages = append(messages, "invalid "+field)
			}
		}
		h.BadRequestResponse(w, r, errors.New(strings.Join(messages, ", ")))
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	user, err := h.Service.Register(r.Context(), req)
	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	input.Password = strings.TrimSpace(input.Password)

	if err := utils.Validate.Struct(input); err != nil {
		var messages []string
		for _, e := range err.(validator.ValidationErrors) {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				messages = append(messages, field+" is required")
			case "email":
				messages = append(messages, "invalid email format")
			case "min":
				messages = append(messages, field+" must be at least "+e.Param()+" characters")
			case "max":
				messages = append(messages, field+" must be at most "+e.Param()+" characters")
			case "alphanum":
				messages = append(messages, field+" must contain only letters and numbers")
			default:
				messages = append(messages, "invalid "+field)
			}
		}
		h.BadRequestResponse(w, r, errors.New(strings.Join(messages, ", ")))
		return
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	token, err := h.Service.Login(r.Context(), input.Email, input.Password)
	if err != nil {
		h.UnauthorizedErrorResponse(w, r, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, LoginResponse{Token: token})
}
