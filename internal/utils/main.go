package utils

import (
	"catalog-be/internal/domain"
	"math/rand"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Utils interface {
	GenerateRandomCode(length int) string
	Slugify(text string) (string, error)
	HashPassword(password string) (*string, *domain.Error)
}

type utils struct{}

// HashPassword implements Utils.
func (u *utils) HashPassword(password string) (*string, *domain.Error) {
	if password == "" {
		password = u.GenerateRandomCode(20)
	}

	hash, bcryptErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if bcryptErr != nil {
		return nil, &domain.Error{
			Code: 500,
			Err:  bcryptErr,
		}
	}

	stringHash := string(hash)
	return &stringHash, nil
}

// GenerateRandomCode implements Utils.
func (u *utils) GenerateRandomCode(length int) string {
	var letterRune = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, length)

	for i := range b {
		b[i] = letterRune[rand.Intn(len(letterRune))]
	}

	return string(b)
}

// Slugify implements Utils.
func (u *utils) Slugify(text string) (string, error) {
	// Remove special characters
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return "", err
	}
	processedString := reg.ReplaceAllString(text, " ")

	// Remove leading and trailing spaces
	processedString = strings.TrimSpace(processedString)

	// Replace spaces with dashes
	slug := strings.ReplaceAll(processedString, " ", "-")

	// Convert to lowercase
	slug = strings.ToLower(slug)

	return slug, nil
}

func NewUtils() Utils {
	return &utils{}
}
