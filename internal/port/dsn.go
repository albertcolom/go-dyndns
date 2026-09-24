//go:generate go tool mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

package port

var SchemeAliases = map[string]string{
	"sqlite":  "sqlite3",
	"sqlite3": "sqlite3",
	"mysql":   "mysql",
	"file":    "file",
}

type DSN struct {
	Driver     string
	DataSource string
	Raw        string
	Normalized string
}

type DSNService interface {
	ParseDSN(raw string) (*DSN, error)
}
