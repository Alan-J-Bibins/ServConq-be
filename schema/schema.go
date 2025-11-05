package schema

import (
	"reflect"
	"time"

	"github.com/nrednav/cuid2"
	"gorm.io/gorm"
)

type TeamMemberRole string

const (
	TeamMemberRoleOwner    TeamMemberRole = "OWNER"    // Make datacenters under their team and be able to delete datacenters
	TeamMemberRoleAdmin    TeamMemberRole = "ADMIN"    // Make Servers / Delete Servers but cannot delete datacenters
	TeamMemberRoleOperator TeamMemberRole = "OPERATOR" // Make Containers and Run shit
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
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	TeamMembers []TeamMember `gorm:"foreignKey:UserID" json:"teamMembers"`
}

type Team struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"unique;not null" json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	JoinToken   string    `json:"joinToken"`

	TeamMembers []TeamMember `gorm:"foreignKey:TeamID" json:"teamMembers"`
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

	Servers []Server `gorm:"foreignKey:DataCenterID"`

	Logs []Log `gorm:"foreignKey:DataCenterID"`
}

type Server struct {
	ID               string `gorm:"primaryKey"`
	DataCenterID     string
	Hostname         string `gorm:"not null"`
	ConnectionString string `gorm:"unique;not null"`
	CreatedAt        time.Time
}

type ContainerImage struct {
	ID          string `gorm:"primaryKey"`
	Name        string
	Version     string
	RegistryURL string
}

type Log struct {
	ID           string     `gorm:"primaryKey" json:"id"`
	DataCenterID string     `gorm:"not null" json:"dataCenterId"`
	TeamMemberID string     `gorm:"not null" json:"teamMemberId"`
	TeamMember   TeamMember `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"teamMember"`
	Message      string     `gorm:"not null" json:"message"`
	CreatedAt    time.Time  `json:"createdAt"`
}

// -------------------- Supporting Tables --------------------

type TeamMember struct {
	ID       string         `gorm:"primaryKey" json:"id"`
	UserID   string         `gorm:"not null" json:"userId"`
	User     User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	TeamID   string         `json:"teamId"`
	Team     Team           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"team"`
	Role     TeamMemberRole `json:"role"`
	JoinedAt time.Time      `json:"joinedAt"`
}
