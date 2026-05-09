package web

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"steamshopemulator/internal/domain"
	"steamshopemulator/internal/review"
	"steamshopemulator/internal/store"
)

type Server struct {
	store      *store.Store
	review     *review.Service
	templates  *template.Template
	projectDir string
	uploadDir  string
}

type pageData struct {
	Storefront domain.Storefront
	Review     domain.ReviewReport
}

type uploadResponse struct {
	Slot    string   `json:"slot"`
	Path    string   `json:"path,omitempty"`
	Gallery []string `json:"gallery,omitempty"`
}

func New(store *store.Store, review *review.Service, projectDir, uploadDir string) (*Server, error) {
	tpls, err := template.New("").Funcs(template.FuncMap{
		"contentHTML":   contentHTML,
		"joinLinks":     joinLinks,
		"mainHeroImage": mainHeroImage,
		"thumbAt":       thumbAt,
	}).ParseGlob(filepath.Join(projectDir, "web", "templates", "*.html"))
	if err != nil {
		return nil, err
	}

	return &Server{
		store:      store,
		review:     review,
		templates:  tpls,
		projectDir: projectDir,
		uploadDir:  uploadDir,
	}, nil
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleStorefront)
	mux.HandleFunc("/admin", s.handleAdmin)
	mux.HandleFunc("/admin/save", s.handleAdminSave)
	mux.HandleFunc("/admin/upload", s.handleUpload)
	mux.HandleFunc("/api/review", s.handleReviewAPI)
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(filepath.Join(s.projectDir, "web", "static")))))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(s.uploadDir))))
	return mux
}

func (s *Server) handleStorefront(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	page := s.store.Get()
	data := pageData{
		Storefront: page,
		Review:     s.review.Evaluate(page),
	}

	s.render(w, "storefront.html", data)
}

func (s *Server) handleAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	page := s.store.Get()
	data := pageData{
		Storefront: page,
		Review:     s.review.Evaluate(page),
	}

	s.render(w, "admin.html", data)
}

func (s *Server) handleAdminSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	current := s.store.Get()
	current.AppID = r.FormValue("app_id")
	current.Slug = r.FormValue("slug")
	current.Title = r.FormValue("title")
	current.Subtitle = r.FormValue("subtitle")
	current.ShortDescription = r.FormValue("short_description")
	current.LongDescription = r.FormValue("long_description")
	current.Developer = r.FormValue("developer")
	current.Publisher = r.FormValue("publisher")
	current.ReleaseDate = r.FormValue("release_date")
	current.Price = r.FormValue("price")
	current.DiscountText = r.FormValue("discount_text")
	current.ReviewSummary = r.FormValue("review_summary")
	current.ReviewCount = r.FormValue("review_count")
	current.Tags = splitCSV(r.FormValue("tags"))
	current.Features = splitCSV(r.FormValue("features"))
	current.ContentHTML = r.FormValue("content_html")
	current.CommunityLinks = parseLinks(r.FormValue("community_links"))
	current.DiscoveryLinks = parseLinks(r.FormValue("discovery_links"))
	current.AboutBlocks = []domain.AboutBlock{
		{Heading: "关于这款游戏", Body: current.LongDescription},
		{Heading: "核心卖点", Body: strings.Join(current.Features, " / ")},
	}

	if err := s.store.Save(current); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(16 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("asset")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := os.MkdirAll(s.uploadDir, 0o755); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name := strings.ReplaceAll(strings.ToLower(header.Filename), " ", "-")
	filename := time.Now().Format("20060102150405") + "-" + name
	targetPath := filepath.Join(s.uploadDir, filename)

	target, err := os.Create(targetPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer target.Close()

	if _, err := io.Copy(target, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	current := s.store.Get()
	publicPath := "/uploads/" + filename
	slot := r.FormValue("slot")
	switch slot {
	case "hero":
		current.HeroImage = publicPath
	case "hero_thumb_1", "hero_thumb_2", "hero_thumb_3", "hero_thumb_4":
		current.HeroThumbs = ensureThumbSlots(current.HeroThumbs, 4)
		index := int(slot[len(slot)-1] - '1')
		current.HeroThumbs[index] = publicPath
		if index == 0 {
			current.HeroImage = publicPath
		}
	case "capsule":
		current.CapsuleImage = publicPath
	case "content_image":
		// Stored on disk only; inserted into ContentHTML on the client before save.
	case "detail_image_1", "detail_image_2", "detail_image_3":
		current.DetailSections = ensureDetailSections(current.DetailSections, 3)
		index := int(slot[len(slot)-1] - '1')
		current.DetailSections[index].Image = publicPath
	default:
		current.Gallery = append(current.Gallery, publicPath)
	}

	if err := s.store.Save(current); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
		response := uploadResponse{Slot: slot}
		switch slot {
		case "hero":
			response.Path = current.HeroImage
		case "hero_thumb_1", "hero_thumb_2", "hero_thumb_3", "hero_thumb_4":
			index := int(slot[len(slot)-1] - '1')
			response.Path = current.HeroThumbs[index]
		case "capsule":
			response.Path = current.CapsuleImage
		case "content_image":
			response.Path = publicPath
		case "detail_image_1", "detail_image_2", "detail_image_3":
			index := int(slot[len(slot)-1] - '1')
			response.Path = current.DetailSections[index].Image
		default:
			response.Gallery = current.Gallery
		}
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (s *Server) handleReviewAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var page domain.Storefront
	if err := json.NewDecoder(r.Body).Decode(&page); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	report := s.review.Evaluate(page)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(report)
}

func (s *Server) render(w http.ResponseWriter, name string, data pageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.templates.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func splitCSV(value string) []string {
	raw := strings.Split(value, ",")
	items := make([]string, 0, len(raw))
	for _, item := range raw {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			items = append(items, trimmed)
		}
	}
	return items
}

func parseLinks(value string) []domain.ExternalLink {
	lines := strings.Split(value, "\n")
	links := make([]domain.ExternalLink, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 2)
		label := strings.TrimSpace(parts[0])
		url := "#"
		if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
			url = strings.TrimSpace(parts[1])
		}

		if label != "" {
			links = append(links, domain.ExternalLink{
				Label: label,
				URL:   url,
			})
		}
	}
	return links
}

func joinLinks(links []domain.ExternalLink) string {
	lines := make([]string, 0, len(links))
	for _, link := range links {
		line := strings.TrimSpace(link.Label)
		if strings.TrimSpace(link.URL) != "" {
			line += " | " + strings.TrimSpace(link.URL)
		}
		if line != "" {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n")
}

func ensureThumbSlots(values []string, size int) []string {
	if len(values) >= size {
		return values
	}

	next := make([]string, size)
	copy(next, values)
	return next
}

func thumbAt(values []string, index int) string {
	if index < 0 || index >= len(values) {
		return ""
	}
	return values[index]
}

func mainHeroImage(page domain.Storefront) string {
	if len(page.HeroThumbs) > 0 && strings.TrimSpace(page.HeroThumbs[0]) != "" {
		return page.HeroThumbs[0]
	}
	return page.HeroImage
}

func contentHTML(page domain.Storefront) template.HTML {
	return template.HTML(page.ContentHTML)
}

func buildDetailSections(headings, paragraphs, images []string) []domain.DetailSection {
	size := maxLen(len(headings), len(paragraphs), len(images))
	sections := make([]domain.DetailSection, 0, size)
	for i := 0; i < size; i++ {
		heading := valueAt(headings, i)
		paragraphText := valueAt(paragraphs, i)
		image := valueAt(images, i)
		parts := splitParagraphBlocks(paragraphText)

		if heading == "" && len(parts) == 0 && image == "" {
			continue
		}

		sections = append(sections, domain.DetailSection{
			Heading:    heading,
			Paragraphs: parts,
			Image:      image,
		})
	}
	return sections
}

func extractDetailImages(sections []domain.DetailSection) []string {
	images := make([]string, len(sections))
	for i, section := range sections {
		images[i] = section.Image
	}
	return images
}

func ensureDetailSections(sections []domain.DetailSection, size int) []domain.DetailSection {
	if len(sections) >= size {
		return sections
	}

	next := make([]domain.DetailSection, size)
	copy(next, sections)
	return next
}

func splitParagraphBlocks(value string) []string {
	normalized := strings.ReplaceAll(value, "\r\n", "\n")
	raw := strings.Split(normalized, "\n\n")
	items := make([]string, 0, len(raw))
	for _, item := range raw {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			items = append(items, trimmed)
		}
	}
	return items
}

func valueAt(values []string, index int) string {
	if index < 0 || index >= len(values) {
		return ""
	}
	return strings.TrimSpace(values[index])
}

func maxLen(values ...int) int {
	max := 0
	for _, value := range values {
		if value > max {
			max = value
		}
	}
	return max
}
