package handlers

import "github.com/muhammedshamil8/url-shortener/internal/models"

type FakeRepository struct {
	GetURLByCodeFunc   func(string) (string, error)
	DeleteURLFunc      func(int) error
	GetAllURLsFunc     func() ([]models.URL, error)
	HealthFunc         func() error
	CreateShortURLFunc func(string, string) (int64, error)
}

func (m *FakeRepository) Health() error {
	if m.HealthFunc != nil {
		return m.HealthFunc()
	}
	return nil
}

func (m *FakeRepository) CreateShortURL(shortCode, url string) (int64, error) {
	if m.CreateShortURLFunc != nil {
		return m.CreateShortURLFunc(shortCode, url)
	}
	return 1, nil
}

func (m *FakeRepository) GetURLByCode(code string) (string, error) {
	if m.GetURLByCodeFunc != nil {
		return m.GetURLByCodeFunc(code)
	}
	return "https://google.com", nil
}

func (m *FakeRepository) DeleteURL(id int) error {
	if m.DeleteURLFunc != nil {
		return m.DeleteURLFunc(id)
	}
	return nil
}

func (m *FakeRepository) GetAllURLs() ([]models.URL, error) {
	if m.GetAllURLsFunc != nil {
		return m.GetAllURLsFunc()
	}
	return []models.URL{}, nil
}
