package handlers

import "github.com/muhammedshamil8/url-shortener/internal/models"

type MockRepository struct {
	GetURLByCodeFunc   func(string) (string, error)
	DeleteURLFunc      func(int) error
	GetAllURLsFunc     func() ([]models.URL, error)
}

func (m *MockRepository) CreateShortURL(shortCode, url string) (int64, error) {
	return 1, nil
}

func (m *MockRepository) GetURLByCode(code string) (string, error) {
	if m.GetURLByCodeFunc != nil {
		return m.GetURLByCodeFunc(code)
	}
	return "https://google.com", nil
}

func (m *MockRepository) DeleteURL(id int) error {
	if m.DeleteURLFunc != nil {
		return m.DeleteURLFunc(id)
	}
	return nil
}

func (m *MockRepository) GetAllURLs() ([]models.URL, error) {
	if m.GetAllURLsFunc != nil {
		return m.GetAllURLsFunc()
	}
	return []models.URL{}, nil
}