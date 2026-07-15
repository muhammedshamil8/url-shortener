package handlers

import (
	"errors"

	"github.com/muhammedshamil8/url-shortener/internal/models"
)

type FakeRepository struct {
	GetURLByCodeFunc   func(string) (string, error)
	GetCodeByIDFunc    func(int) (string, error)
	DeleteURLFunc      func(int) error
	UpdateURLFunc      func(int, string) error
	GetAllURLsFunc     func(models.ListOptions) ([]models.URL, error)
	HealthFunc         func() error
	CreateShortURLFunc func(string, string, *int) (int64, error)

	CreateUserFunc            func(string, string, string) (int64, error)
	GetUserByEmailFunc        func(string) (*models.User, error)
	DeleteUserURLFunc         func(int) error
	UpdateUserURLFunc         func(int, string, string) error
	GetAllURLsByUserEmailFunc func(string) ([]models.URL, error)
	GetAllUsersFunc           func() ([]models.User, error)
	DeleteUserFunc            func(int) error
}

// CreateUser implements [Repository].
func (m *FakeRepository) CreateUser(username string, email string, password string) (int64, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(username, email, password)
	}
	return 1, nil
}

// GetUserByEmail implements [Repository].
func (m *FakeRepository) GetUserByEmail(email string) (*models.User, error) {
	if m.GetUserByEmailFunc != nil {
		return m.GetUserByEmailFunc(email)
	}
	return nil, errors.New("user not found")
}

func (m *FakeRepository) Health() error {
	if m.HealthFunc != nil {
		return m.HealthFunc()
	}
	return nil
}

func (m *FakeRepository) CreateShortURL(shortCode, url string, userID *int) (int64, error) {
	if m.CreateShortURLFunc != nil {
		return m.CreateShortURLFunc(shortCode, url, userID)
	}
	return 1, nil
}

func (m *FakeRepository) GetURLByCode(code string) (string, error) {
	if m.GetURLByCodeFunc != nil {
		return m.GetURLByCodeFunc(code)
	}
	return "https://google.com", nil
}

func (m *FakeRepository) GetCodeByID(id int) (string, error) {
	if m.GetCodeByIDFunc != nil {
		return m.GetCodeByIDFunc(id)
	}
	return "", nil
}

func (m *FakeRepository) DeleteURL(id int) error {
	if m.DeleteURLFunc != nil {
		return m.DeleteURLFunc(id)
	}
	return nil
}

func (m *FakeRepository) UpdateURL(id int, newURL string) error {
	if m.UpdateURLFunc != nil {
		return m.UpdateURLFunc(id, newURL)
	}
	return nil
}

func (m *FakeRepository) GetAllURLs(opts models.ListOptions) ([]models.URL, error) {
	if m.GetAllURLsFunc != nil {
		return m.GetAllURLsFunc(opts)
	}
	return []models.URL{}, nil
}

func (m *FakeRepository) DeleteUserURL(id int) error {
	if m.DeleteUserURLFunc != nil {
		return m.DeleteUserURLFunc(id)
	}
	return nil
}

func (m *FakeRepository) UpdateUserURL(id int, email string, newURL string) error {
	if m.UpdateUserURLFunc != nil {
		return m.UpdateUserURLFunc(id, email, newURL)
	}
	return nil
}

func (m *FakeRepository) GetAllURLsByUserEmail(email string) ([]models.URL, error) {
	if m.GetAllURLsByUserEmailFunc != nil {
		return m.GetAllURLsByUserEmailFunc(email)
	}
	return []models.URL{}, nil
}

func (m *FakeRepository) GetAllUsers() ([]models.User, error) {
	if m.GetAllUsersFunc != nil {
		return m.GetAllUsersFunc()
	}
	return []models.User{}, nil
}

func (m *FakeRepository) DeleteUser(id int) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(id)
	}
	return nil
}
