package schema

import (
	"reflect"
	"time"

	"github.com/nrednav/cuid2"
	"gorm.io/gorm"
)

type TeamMemberRole string

const (
	TeamMemberRoleOwner    TeamMemberRole = "OWNER"
	TeamMemberRoleAdmin    TeamMemberRole = "ADMIN"
	TeamMemberRoleOperator TeamMemberRole = "OPERATOR"
	TeamMemberRoleViewer   TeamMemberRole = "VIEWER"
)

// NOTE: We need this function to initialize the ID field of tables with CUID
func RegisterCUIDCallback(db *gorm.DB) {
	db.Callback().Create().Before("gorm:before_create").Register("set_cuid_if_empty", func(tx *gorm.DB) {
		if tx.Statement.Schema != nil {
			idField := tx.Statement.Schema.LookUpField("ID")
			if idField != nil && idField.FieldType.Kind() == reflect.String {
				val, _ := idField.ValueOf(tx.Statement.Statement.Context, tx.Statement.ReflectValue)
				if str, ok := val.(string); ok && str == "" {
					tx.Statement.SetColumn("ID", cuid2.Generate())
				}
			}
		}
	})
}

// -------------------- Access Control --------------------

type User struct {
	ID        string `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	TeamMembers []TeamMember `gorm:"foreignKey:UserID"`
}

type Team struct {
	ID          string `gorm:"primaryKey"`
	Name        string `gorm:"unique;not null"`
	Description string
	CreatedAt   time.Time

	TeamMembers []TeamMember `gorm:"foreignKey:TeamID"`
}

// -------------------- Infrastructure --------------------

type DataCenter struct {
	ID          string `gorm:"primaryKey"`
	Name        string `gorm:"unique;not null"`
	Location    string
	Description string
	CreatedAt   time.Time
	TeamID      string
	Team        Team `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	Servers          []Server           `gorm:"foreignKey:DataCenterID"`
	NetworkingDevice []NetworkingDevice `gorm:"foreignKey:DataCenterID"`
	Logs             []Log              `gorm:"foreignKey:DataCenterID"`
}

type Server struct {
	ID                    string `gorm:"primaryKey"`
	DataCenterID          string
	Hostname              string `gorm:"unique;not null"`
	IPAddress             string `gorm:"unique;not null"`
	OS                    string
	StorageSystemID       string
	PowerInfrastructureID string
	AgentBinaryID         string
	CreatedAt             time.Time

	Containers []Container `gorm:"foreignKey:ServerID"`
	Events     []Event     `gorm:"foreignKey:ServerID"`
	Logs       []Log       `gorm:"foreignKey:ServerID"`
}

type Container struct {
	ID               string `gorm:"primaryKey"`
	ServerID         string
	Name             string
	ContainerImageID string
	StorageLimitGB   int
	CreatedAt        time.Time

	EnvVars []ContainerEnvVar `gorm:"foreignKey:ContainerID"`
	Ports   []ContainerPort   `gorm:"foreignKey:ContainerID"`
}

type ContainerImage struct {
	ID          string `gorm:"primaryKey"`
	Name        string
	Version     string
	RegistryURL string

	Containers []Container `gorm:"foreignKey:ContainerImageID"`
}

type ContainerEnvVar struct {
	ID          string `gorm:"primaryKey"`
	ContainerID string
	Key         string
	Value       string
}

type ContainerPort struct {
	ID          string `gorm:"primaryKey"`
	ContainerID string
	Port        int
	Protocol    string
}

type AgentBinary struct {
	ID          string `gorm:"primaryKey"`
	Version     string
	Checksum    string
	Description string

	Servers []Server `gorm:"foreignKey:AgentBinaryID"`
}

type StorageSystem struct {
	ID          string `gorm:"primaryKey"`
	TypeID      string
	CapacityGB  int
	Description string

	Servers []Server `gorm:"foreignKey:StorageSystemID"`
}

type PowerInfrastructure struct {
	ID          string `gorm:"primaryKey"`
	Type        string
	CapacityKW  int
	Description string

	Servers []Server `gorm:"foreignKey:PowerInfrastructureID"`
}

type Event struct {
	ID        string `gorm:"primaryKey"`
	ServerID  string
	EventType string
	Message   string
	Timestamp time.Time

	Logs []Log `gorm:"foreignKey:EventID"`
}

type NetworkingDevice struct {
	ID           string `gorm:"primaryKey"`
	DataCenterID string
	TypeID       string
	Manufacturer string
	Model        string
	IPAddress    string
	Description  string
}

type Log struct {
	ID           string `gorm:"primaryKey"`
	DataCenterID string
	ServerID     string
	EventID      string
	Message      string
	Level        string
	Timestamp    time.Time
}

// -------------------- Supporting Tables --------------------

type TeamMember struct {
	ID       string `gorm:"primaryKey"`
	UserID   string
	User     User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TeamID   string
	Team     Team           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Role     TeamMemberRole `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	JoinedAt time.Time
}

type NetworkingDeviceType struct {
	ID          string `gorm:"primaryKey"`
	Name        string
	Description string

	Devices []NetworkingDevice `gorm:"foreignKey:TypeID"`
}

type StorageSystemType struct {
	ID          string `gorm:"primaryKey"`
	Name        string
	Description string

	Systems []StorageSystem `gorm:"foreignKey:TypeID"`
}

// type TeamDataCenterAccess struct {
// 	ID           string `gorm:"primaryKey"`
// 	TeamID       string
// 	DataCenterID string
// }
