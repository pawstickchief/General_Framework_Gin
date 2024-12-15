package mysql

import (
	"General_Framework_Gin/schemas/business"
	"bufio"
	"fmt"
	"gorm.io/gorm"
	"os"
	"strings"
)

func UpdatePoliciesFromFile(db *gorm.DB, filename, operator string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "p,") {
			// Parse the policy line (e.g. p, admin, /users, GET)
			parts := strings.Split(line, ",")
			if len(parts) != 4 {
				continue // Skip invalid lines
			}
			role := strings.TrimSpace(parts[1])
			resource := strings.TrimSpace(parts[2])
			action := strings.TrimSpace(parts[3])

			// Check if the policy already exists
			var existingPolicy business.Policy
			db.Where("role = ? AND resource = ? AND action = ?", role, resource, action).First(&existingPolicy)

			if existingPolicy.ID == 0 {
				// Insert new policy
				newPolicy := business.Policy{
					Role:     role,
					Resource: resource,
					Action:   action,
					Operator: operator,
				}
				if err := db.Create(&newPolicy).Error; err != nil {
					return fmt.Errorf("failed to insert policy: %v", err)
				}
				fmt.Printf("Inserted new policy: %s %s %s\n", role, resource, action)
			} else {
				// Update existing policy
				existingPolicy.Remark = "Updated from file" // Example remark
				existingPolicy.Operator = operator
				if err := db.Save(&existingPolicy).Error; err != nil {
					return fmt.Errorf("failed to update policy: %v", err)
				}
				fmt.Printf("Updated policy: %s %s %s\n", role, resource, action)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	return nil
}
