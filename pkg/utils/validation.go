package utils

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) error {
    return validate.Struct(s)
}

type ValidationError struct {
    Field   string `json:"field"`
    Tag     string `json:"tag"`
    Param   string `json:"param,omitempty"`
    Message string `json:"message"`
}

func ValidateStructDetailed(s interface{}) ([]ValidationError, error) {
    err := validate.Struct(s)
    if err == nil {
        return nil, nil
    }
    verrs, ok := err.(validator.ValidationErrors)
    if !ok {
        return nil, err
    }
    out := make([]ValidationError, 0, len(verrs))
    for _, fe := range verrs {
        msg := buildMessage(fe)
        out = append(out, ValidationError{
            Field:   fe.Field(),
            Tag:     fe.Tag(),
            Param:   fe.Param(),
            Message: msg,
        })
    }
    return out, nil
}

func buildMessage(fe validator.FieldError) string {
    switch fe.Tag() {
    case "required":
        return fe.Field() + " is required"
    case "email":
        return fe.Field() + " must be a valid email"
    case "min":
        return fe.Field() + " must be at least " + fe.Param() + " characters"
    case "max":
        return fe.Field() + " must be at most " + fe.Param() + " characters"
    case "oneof":
        return fe.Field() + " must be one of: " + fe.Param()
    default:
        return fe.Field() + " is invalid"
    }
}


