package input

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	// Read Replica configuration (optional)
	ReadReplicaHost     string
	ReadReplicaPort     int
	ReadReplicaUser     string
	ReadReplicaPassword string
}

type ServerConfig struct {
	Host string
	Port int
}

type AppConfig struct {
	Environment string
	JWTSecret   string
	ServerPort  string
	// Comma-separated list of allowed origins for CORS (e.g. "http://localhost:3000,http://127.0.0.1:3000")
	AllowedOrigins string
}

type GCSConfig struct {
	BucketName    string
	ProjectID     string
	PublicBaseURL string
}


type OAuthConfig struct {
	GoogleClientID        string
	GoogleIOSClientID     string
	GoogleAndroidClientID string
	FirebaseProjectID     string
}

type CustomerServiceConfig struct {
	BaseURL string
}

// DashboardStatsFilter represents filters for dashboard statistics
type DashboardStatsFilter struct {
	CustomerType *string `json:"customer_type,omitempty"` // mobile_user, partner, vendor, admin, superadmin
	FromDate     *string `json:"from_date,omitempty"`     // Format: YYYY-MM-DD
	ToDate       *string `json:"to_date,omitempty"`       // Format: YYYY-MM-DD
}

type ServiceConfig struct {
	CustomerServiceURL string
}
