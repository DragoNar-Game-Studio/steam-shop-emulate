package domain

type Storefront struct {
	AppID              string             `json:"app_id"`
	Slug               string             `json:"slug"`
	Title              string             `json:"title"`
	Subtitle           string             `json:"subtitle"`
	ShortDescription   string             `json:"short_description"`
	LongDescription    string             `json:"long_description"`
	Developer          string             `json:"developer"`
	Publisher          string             `json:"publisher"`
	ReleaseDate        string             `json:"release_date"`
	Price              string             `json:"price"`
	DiscountText       string             `json:"discount_text"`
	ReviewSummary      string             `json:"review_summary"`
	ReviewCount        string             `json:"review_count"`
	Tags               []string           `json:"tags"`
	Features           []string           `json:"features"`
	ContentHTML        string             `json:"content_html"`
	HeroImage          string             `json:"hero_image"`
	HeroThumbs         []string           `json:"hero_thumbs"`
	CapsuleImage       string             `json:"capsule_image"`
	LogoImage          string             `json:"logo_image"`
	Gallery            []string           `json:"gallery"`
	AboutBlocks        []AboutBlock       `json:"about_blocks"`
	DetailSections     []DetailSection    `json:"detail_sections"`
	ReadMoreLabel      string             `json:"read_more_label"`
	SystemRequirements SystemRequirements `json:"system_requirements"`
	SidebarHighlights  []IconLabel        `json:"sidebar_highlights"`
	ControllerSupports []IconLabel        `json:"controller_supports"`
	LanguageSupport    []LanguageSupport  `json:"language_support"`
	SteamDeck          SteamDeckInfo      `json:"steam_deck"`
	AchievementCount   string             `json:"achievement_count"`
	AchievementBadges  []AchievementBadge `json:"achievement_badges"`
	MetadataEntries    []KeyValue         `json:"metadata_entries"`
	CommunityLinks     []ExternalLink     `json:"community_links"`
	DiscoveryLinks     []ExternalLink     `json:"discovery_links"`
}

type AboutBlock struct {
	Heading string `json:"heading"`
	Body    string `json:"body"`
}

type DetailSection struct {
	Heading    string   `json:"heading"`
	Paragraphs []string `json:"paragraphs"`
	Image      string   `json:"image"`
}

type SystemRequirements struct {
	Platforms []PlatformRequirements `json:"platforms"`
}

type PlatformRequirements struct {
	Name        string   `json:"name"`
	IsActive    bool     `json:"is_active"`
	Minimum     []string `json:"minimum"`
	Recommended []string `json:"recommended"`
}

type IconLabel struct {
	Icon  string `json:"icon"`
	Label string `json:"label"`
}

type LanguageSupport struct {
	Name      string `json:"name"`
	Interface bool   `json:"interface"`
	FullAudio bool   `json:"full_audio"`
	Subtitles bool   `json:"subtitles"`
}

type SteamDeckInfo struct {
	Status string `json:"status"`
	CTA    string `json:"cta"`
}

type AchievementBadge struct {
	Label string `json:"label"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ExternalLink struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

type ReviewReport struct {
	Score       int      `json:"score"`
	Status      string   `json:"status"`
	Highlights  []string `json:"highlights"`
	Warnings    []string `json:"warnings"`
	Suggestions []string `json:"suggestions"`
}
