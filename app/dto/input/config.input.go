package input

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
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

type DashboardStatsFilter struct {
	CustomerType *string `json:"customer_type,omitempty"`
	FromDate     *string `json:"from_date,omitempty"`
	ToDate       *string `json:"to_date,omitempty"`
}

type ServiceConfig struct {
	CustomerServiceURL string
}
