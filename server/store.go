package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

type Stat struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Profile struct {
	Name      string `json:"name"`
	Tagline   string `json:"tagline"`
	Available bool   `json:"available"`
	Stats     []Stat `json:"stats"`
	HeroImage string `json:"heroImage"`
	ResumePDF string `json:"resumePdf"`
}

type Experience struct {
	Period  string `json:"period"`
	Role    string `json:"role"`
	Company string `json:"company"`
	Desc    string `json:"desc"`
}

type Tool struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
}

type Work struct {
	ID    string `json:"id"`
	Cat   string `json:"cat"`
	Type  string `json:"type"`
	Title string `json:"title"`
	Area  string `json:"area"`
	Year  string `json:"year"`
	Scope string `json:"scope"`
	Cover string `json:"cover"`
	PDF   string `json:"pdf"`
	Size  string `json:"size"`
}

type Contacts struct {
	Email    string `json:"email"`
	Telegram string `json:"telegram"`
	Phone    string `json:"phone"`
}

type Content struct {
	Profile    Profile      `json:"profile"`
	Experience []Experience `json:"experience"`
	Tools      []Tool       `json:"tools"`
	Skills     []string     `json:"skills"`
	Works      []Work       `json:"works"`
	Contacts   Contacts     `json:"contacts"`
}

type Store struct {
	mu   sync.RWMutex
	path string
	data Content
}

func OpenStore(path string) (*Store, error) {
	s := &Store{path: path}
	b, err := os.ReadFile(path)
	switch {
	case errors.Is(err, os.ErrNotExist):
		s.data = seedContent()
		if err := s.saveLocked(); err != nil {
			return nil, err
		}
	case err != nil:
		return nil, err
	default:
		if err := json.Unmarshal(b, &s.data); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *Store) Get() Content {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return clone(s.data)
}

// Update applies fn to a copy of the content and persists the result.
func (s *Store) Update(fn func(c *Content) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	c := clone(s.data)
	if err := fn(&c); err != nil {
		return err
	}
	prev := s.data
	s.data = c
	if err := s.saveLocked(); err != nil {
		s.data = prev
		return err
	}
	return nil
}

func (s *Store) saveLocked() error {
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

func clone(c Content) Content {
	b, _ := json.Marshal(c)
	var out Content
	_ = json.Unmarshal(b, &out)
	return out
}

func unsplash(id string) string {
	return "https://images.unsplash.com/photo-" + id + "?w=1200&q=70&auto=format&fit=crop"
}

func seedContent() Content {
	return Content{
		Profile: Profile{
			Name:      "Фаима Жолмухамедова",
			Tagline:   "Дизайнер‑проектировщик интерьеров. 4 года: обмеры → планировка → рабочая документация → авторский надзор.",
			Available: true,
			Stats:     []Stat{{"Стаж", "4 года"}, {"Объекты", "30+"}, {"Площадь", "5 000 м²"}},
			HeroImage: unsplash("1601993957728-1e56ab70c5a8"),
		},
		Experience: []Experience{
			{"2024 — сейчас", "Дизайнер-проектировщик", "Студия интерьеров «Линия»", "Веду жилые объекты от обмерного плана до авторского надзора. Планировки, рабочая документация, спецификации для закупки."},
			{"2022 — 2024", "Младший дизайнер-проектировщик", "Бюро «Форма»", "Рабочие чертежи, развёртки стен, схемы освещения и сантехники. Визуализации в 3ds Max + Corona."},
			{"2018 — 2022", "Бакалавр, дизайн архитектурной среды", "Архитектурный университет", "Диплом — интерьер общественного пространства в историческом здании."},
		},
		Tools:  []Tool{{"AutoCAD", 5}, {"ArchiCAD", 4}, {"3ds Max + Corona", 4}, {"Revit", 3}, {"SketchUp", 4}, {"Photoshop", 4}, {"Planoplan", 3}, {"Figma", 2}},
		Skills: []string{"Обмерные планы", "Планировочные решения", "Рабочая документация", "Развёртки стен", "Схемы освещения", "Электрика и сантехника", "Спецификации", "3D-визуализация", "Подбор материалов", "Авторский надзор"},
		Works: []Work{
			{"w1", "home", "Квартира", "Квартира у Патриарших", "142 м²", "2025", "Полный комплект: планировка, развёртки, освещение, спецификация мебели.", unsplash("1583847268964-b28dc8f51f92"), "", ""},
			{"w2", "house", "Дом", "Загородный дом в Истре", "280 м²", "2025", "Интерьер двух этажей, кухня-гостиная с камином, 72 листа документации.", unsplash("1705321963943-de94bb3f0dd3"), "", ""},
			{"w3", "home", "Квартира", "Студия для молодой пары", "46 м²", "2024", "Перепланировка с согласованием, встроенное хранение, спальная ниша.", unsplash("1605774337664-7a846e9cdf17"), "", ""},
			{"w4", "house", "Дом", "Лестничный холл", "38 м²", "2024", "Монолитная лестница, скрытая подсветка ступеней, узлы примыканий.", unsplash("1601993957728-1e56ab70c5a8"), "", ""},
			{"w5", "commercial", "Коммерческий", "Шоурум мебели", "210 м²", "2024", "Экспозиционные зоны, навигация, сценарии освещения.", unsplash("1583329550487-0fa300a4cd1a"), "", ""},
			{"w6", "home", "Квартира", "Апартаменты в Сити", "78 м²", "2023", "Светлая гостиная, подбор материалов, визуализации для заказчика.", unsplash("1609081144289-eacc3108cd03"), "", ""},
		},
		Contacts: Contacts{Email: "faima@design.ru", Telegram: "https://t.me/", Phone: "+7 000 000-00-00"},
	}
}
