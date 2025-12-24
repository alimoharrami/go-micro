package migrations

import (
	"auth/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"log"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&domain.User{},
		&domain.Role{},
		&domain.Permission{},
		&domain.RolePermission{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	// Seed default role/permission and assign to a user
	// Create or get 'admin' role
	var role domain.Role
	if err := db.Where("name = ?", "admin").First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			role = domain.Role{Name: "admin", Slug: "admin", Description: "Administrator"}
			if err := db.Create(&role).Error; err != nil {
				log.Printf("Failed to create admin role: %v", err)
			}
		} else {
			log.Printf("Error querying role: %v", err)
		}
	}

	// Create or get 'post:create' permission
	var perm domain.Permission
	if err := db.Where("key = ?", "post:create").First(&perm).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			perm = domain.Permission{Key: "post:create", Description: "Create posts"}
			if err := db.Create(&perm).Error; err != nil {
				log.Printf("Failed to create permission: %v", err)
			}
		} else {
			log.Printf("Error querying permission: %v", err)
		}
	}

	// Assign permission to role if not already assigned
	if role.ID != 0 && perm.ID != 0 {
		rp := domain.RolePermission{RoleID: role.ID, PermissionID: perm.ID}
		if err := db.First(&domain.RolePermission{}, "role_id = ? AND permission_id = ?", role.ID, perm.ID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&rp).Error; err != nil {
					log.Printf("Failed to assign permission to role: %v", err)
				}
			} else {
				log.Printf("Error checking role_permission: %v", err)
			}
		}
	}

	// Attach role to `admin@local` user (create if missing)
	var user domain.User
	if err := db.Where("email = ?", "admin@local").First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// create default admin user
			pw := "admin123"
			hash, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
			user = domain.User{Email: "admin@local", Password: string(hash), FirstName: "Admin", LastName: "User", Active: true, RoleID: role.ID}
			if err := db.Create(&user).Error; err != nil {
				log.Printf("Failed to create default admin user: %v", err)
			} else {
				log.Printf("Created default admin user: admin@local (password: %s)", pw)
			}
		} else {
			log.Printf("Error querying user by email: %v", err)
		}
	} else {
		// ensure user's role is set to admin
		if role.ID != 0 && user.RoleID != role.ID {
			user.RoleID = role.ID
			if err := db.Save(&user).Error; err != nil {
				log.Printf("Failed to assign role to existing admin user: %v", err)
			}
		}
	}
}
