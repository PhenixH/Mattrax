package mattrax

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/mattrax/Mattrax/internal/api"
	"github.com/mattrax/Mattrax/internal/authentication"
	"github.com/mattrax/Mattrax/internal/certificates"
	"github.com/mattrax/Mattrax/internal/db"
	"github.com/mattrax/Mattrax/internal/settings"
	"github.com/patrickmn/go-cache"
)

// Version contains the builds version. It is injected at build time
const Version = "v0.0.0-dev"

// VersionCommit contains the Git commit it of the build. It is injected at build time
const VersionCommit = "xxxxxx"

// VersionDate contains the time of the build. It is injected at build time
const VersionDate = "2020.11.12 00:00:00AM"

// Server contains the global server state
type Server struct {
	Args         Arguments
	GlobalRouter *mux.Router
	Router       *mux.Router // Subrouter which is only accessible via secure origins (configured by admin)

	API    *api.Service
	DB     *db.Queries
	DBConn *sql.DB
	Cache  *cache.Cache

	Cert     *certificates.Service
	Auth     *authentication.Service
	Settings *settings.Service
}

// Arguments are the command line flags
type Arguments struct {
	Domain  string `placeholder:"\"mdm.example.com\"" arg:"env:DOMAIN" help:"The domain your server is accessible from"`
	DB      string `placeholder:"\"postgres://localhost/Mattrax\"" arg:"env:DB" help:"The Postgres database connection url"`
	Addr    string `default:":443" placeholder:"\":443\"" arg:"env:ADDR" help:"The listen address of the https server"`
	TLSCert string `default:"./certs/tls.crt" placeholder:"\"./certs/tls.crt\"" harg:"env:TLS_CERT" elp:"The path for the tls certificate"`
	TLSKey  string `default:"./certs/tls.key" placeholder:"\"./certs/tls.key\"" arg:"env:TLS_KEY" help:"The path for the tls certificates key"`
	Zipkin  string `default:"" placeholder:"\"http://localhost:9411/api/v2/spans\"" arg:"env:ZIPKIN" help:"The url of the Zipkin server. This feature is optional"`

	GoogleServiceAccountPath string `default:"" placeholder:"\"./certs/serviceaccount.json\"" arg:"env:GOOGLE_APPLICATION_CREDENTIALS" help:"The path of the service account. required for the Android management"`
	AzureADClientID          string `default:"" placeholder:"" arg:"env:AZUREAD_CLIENT_ID" help:"The application (client) ID. required for the AzureAD integration"`
	AzureADClientSecret      string `default:"" placeholder:"" arg:"env:AZUREAD_CLIENT_SECRET" help:"The application client secret. required for the AzureAD integration"`

	MattraxCloud              bool   `default:"false" arg:"env:MATTRAX_CLOUD" help:"Enable extra Saas Mattrax API endpoints. DO NOT use unless you know what you are doing!"`
	MattraxCloudAuth          string `default:"" arg:"env:MATTRAX_CLOUD_AUTH" help:"Authentication key for Mattrax Cloud API endpoints. DO NOT use unless you know what you are doing!"`
	MattraxCloudAuthIP        string `default:"127.0.0.1" arg:"env:MATTRAX_CLOUD_AUTH_IP" help:"Allowed IP for Mattrax Cloud API endpoints. DO NOT use unless you know what you are doing!"`
	MattraxCloudVendorEmail   string `default:"" arg:"env:MATTRAX_CLOUD_VENDOR_EMAIL" help:"The support email for the Mattrax cloud vendor!"`
	MattraxCloudVendorWebsite string `default:"" arg:"env:MATTRAX_CLOUD_VENDOR_WEBSITE" help:"The support website for the Mattrax cloud vendor!"`
	MattraxCloudVendorPhone   string `default:"" arg:"env:MATTRAX_CLOUD_VENDOR_PHONE" help:"The support phone number for the Mattrax cloud vendor!"`

	Debug bool `help:"Enabled development mode. PLEASE DO NOT USE IN PRODUCTION!"`
}

// Description is for alexflint/go-args
func (Arguments) Description() string {
	return "Mattrax MDM Server. Created by Oscar Beaumont!"
}

// Version is for alexflint/go-args
func (Arguments) Version() string {
	return "Version: " + Version
}
