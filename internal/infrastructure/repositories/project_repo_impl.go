package repositories

import (
	"fmt"
	"shwetaik-sqlacc-stock-api/internal/domain/entities"
	"shwetaik-sqlacc-stock-api/internal/domain/repositories"
	"strings"

	"gorm.io/gorm"
)

type ProjectRepositoryImpl struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) repositories.ProjectRepository {
	return &ProjectRepositoryImpl{db: db}
}

func (r *ProjectRepositoryImpl) GetAllProjects(filter map[string]any) ([]entities.Project, error) {
	whereClauses := []string{}
	args := []interface{}{}

	// Keyset pagination: the caller passes the CODE of the last item from
	// the previous page instead of an offset.
	if after, ok := filter["after"].(string); ok && after != "" {
		whereClauses = append(whereClauses, "CODE > ?")
		args = append(args, after)
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	limit, hasLimit := filter["limit"].(int)
	if !hasLimit || limit <= 0 {
		limit = 0
	}

	paginationSQL := ""
	if limit > 0 {
		paginationSQL = fmt.Sprintf("ROWS %d", limit)
	}

	// ATTACHMENTS is a BLOB column and isn't needed here, so it's left out
	// of the SELECT list.
	query := fmt.Sprintf(`
		SELECT CODE, DESCRIPTION, DESCRIPTION2, PROJECTVALUE, PROJECTCOST, ISACTIVE
		FROM PROJECT
		%s
		ORDER BY CODE
		%s
	`, whereSQL, paginationSQL)

	var projects []entities.Project
	if err := r.db.Raw(query, args...).Scan(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *ProjectRepositoryImpl) GetProjectByCode(code string) (*entities.Project, error) {
	var project entities.Project
	err := r.db.Where("code = ?", code).First(&project).Error
	return &project, err
}
