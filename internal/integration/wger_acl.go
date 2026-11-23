package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
)

// WgerExercise representa el modelo canónico interno desacoplado de la estructura externa de wger.
type WgerExercise struct {
	ID          int
	Name        string
	MuscleGroup string
	Description string
}

// WgerACL define el contrato para consultar ejercicios desde la API externa.
type WgerACL interface {
	SearchExercises(ctx context.Context, query string, limit int) ([]WgerExercise, error)
	GetExerciseByID(ctx context.Context, id int) (*WgerExercise, error)
}

type wgerHTTPACL struct {
	baseURL    string
	httpClient *http.Client
	htmlTagRe  *regexp.Regexp
}

// NewWgerACL crea una nueva instancia de la capa anticorrupción para wger con timeout configurado.
func NewWgerACL(baseURL string, timeout time.Duration) WgerACL {
	if baseURL == "" {
		baseURL = "https://wger.de/api/v2"
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	return &wgerHTTPACL{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
		htmlTagRe: regexp.MustCompile(`<[^>]*>`),
	}
}

// Estructuras internas privadas para mapear el payload crudo de wger sin contaminar el dominio.
type wgerExerciseInfoResponse struct {
	Count   int                  `json:"count"`
	Results []wgerExerciseResult `json:"results"`
}

type wgerExerciseResult struct {
	ID           int                     `json:"id"`
	Category     wgerCategory            `json:"category"`
	Translations []wgerTranslationResult `json:"translations"`
}

type wgerCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type wgerTranslationResult struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Language    int    `json:"language"` // 4 = Español, 2 = Inglés
}

func (a *wgerHTTPACL) SearchExercises(ctx context.Context, query string, limit int) ([]WgerExercise, error) {
	cleanQuery := strings.TrimSpace(query)
	if len([]rune(cleanQuery)) < 3 {
		return nil, nil
	}

	if limit <= 0 {
		limit = 50
	}

	endpoint, err := url.Parse(fmt.Sprintf("%s/exerciseinfo/", a.baseURL))
	if err != nil {
		return nil, fmt.Errorf("parsing wger endpoint: %w", err)
	}

	params := endpoint.Query()
	params.Set("language__code", "es")
	params.Set("limit", fmt.Sprintf("%d", limit))
	endpoint.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating wger request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing wger request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from wger: %d", resp.StatusCode)
	}

	var rawResp wgerExerciseInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawResp); err != nil {
		return nil, fmt.Errorf("decoding wger response: %w", err)
	}

	normalized := a.translateAndFilter(rawResp.Results, cleanQuery)
	return normalized, nil
}

func (a *wgerHTTPACL) GetExerciseByID(ctx context.Context, id int) (*WgerExercise, error) {
	endpoint, err := url.Parse(fmt.Sprintf("%s/exerciseinfo/%d/", a.baseURL, id))
	if err != nil {
		return nil, fmt.Errorf("parsing wger endpoint: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating wger request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing wger request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from wger: %d", resp.StatusCode)
	}

	var rawResp wgerExerciseResult
	if err := json.NewDecoder(resp.Body).Decode(&rawResp); err != nil {
		return nil, fmt.Errorf("decoding wger response: %w", err)
	}

	name, desc := a.extractTranslation(rawResp.Translations)
	if name == "" {
		return nil, nil
	}

	return &WgerExercise{
		ID:          rawResp.ID,
		Name:        name,
		MuscleGroup: rawResp.Category.Name,
		Description: a.stripHTML(desc),
	}, nil
}

func (a *wgerHTTPACL) translateAndFilter(results []wgerExerciseResult, query string) []WgerExercise {
	lowerQuery := strings.ToLower(query)
	exercises := make([]WgerExercise, 0, len(results))

	for _, item := range results {
		name, desc := a.extractTranslation(item.Translations)
		if name == "" {
			continue
		}

		// Filtrado case-insensitive por nombre o grupo muscular
		lowerName := strings.ToLower(name)
		lowerGroup := strings.ToLower(item.Category.Name)
		if !strings.Contains(lowerName, lowerQuery) && !strings.Contains(lowerGroup, lowerQuery) {
			continue
		}

		exercises = append(exercises, WgerExercise{
			ID:          item.ID,
			Name:        name,
			MuscleGroup: item.Category.Name,
			Description: a.stripHTML(desc),
		})
	}

	// Orden ascendente alfabético por nombre de ejercicio
	sort.Slice(exercises, func(i, j int) bool {
		return strings.ToLower(exercises[i].Name) < strings.ToLower(exercises[j].Name)
	})

	return exercises
}

func (a *wgerHTTPACL) extractTranslation(translations []wgerTranslationResult) (string, string) {
	var fallbackName, fallbackDesc string

	for _, t := range translations {
		trimmedName := strings.TrimSpace(t.Name)
		if trimmedName == "" {
			continue
		}
		// Prioridad al español (language == 4 en wger)
		if t.Language == 4 {
			return trimmedName, t.Description
		}
		if fallbackName == "" {
			fallbackName = trimmedName
			fallbackDesc = t.Description
		}
	}

	return fallbackName, fallbackDesc
}

func (a *wgerHTTPACL) stripHTML(raw string) string {
	cleaned := a.htmlTagRe.ReplaceAllString(raw, " ")
	return strings.Join(strings.Fields(cleaned), " ")
}