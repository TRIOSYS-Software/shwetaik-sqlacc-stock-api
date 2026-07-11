package repositories

import (
	"fmt"
	"shwetaik-sqlacc-stock-api/internal/domain/entities"
	"shwetaik-sqlacc-stock-api/internal/domain/repositories"
	"strings"

	"gorm.io/gorm"
)

type GLAccRepositoryImpl struct {
	db *gorm.DB
}

func NewGLAccRepository(db *gorm.DB) repositories.GLAccRepository {
	return &GLAccRepositoryImpl{db: db}
}

func (r *GLAccRepositoryImpl) GetAllGLAccs(filter map[string]any) ([]entities.GLAcc, error) {
	whereClauses := []string{}
	args := []interface{}{}

	if parent, ok := filter["parent"].(int); ok && parent > 0 {
		whereClauses = append(whereClauses, "PARENT = ?")
		args = append(args, parent)
	}
	if accType, ok := filter["acctype"].(string); ok && accType != "" {
		whereClauses = append(whereClauses, "ACCTYPE = ?")
		args = append(args, accType)
	}
	if description, ok := filter["description"].(string); ok && description != "" {
		whereClauses = append(whereClauses, "DESCRIPTION LIKE ?")
		args = append(args, "%"+description+"%")
	}
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

	query := fmt.Sprintf(`
		SELECT DOCKEY, PARENT, CODE, DESCRIPTION, DESCRIPTION2, ACCTYPE, SPECIALACCTYPE, TAX, CASHFLOWTYPE, SIC
		FROM GL_ACC
		%s
		ORDER BY CODE
		%s
	`, whereSQL, paginationSQL)

	var glAccs []entities.GLAcc
	if err := r.db.Raw(query, args...).Scan(&glAccs).Error; err != nil {
		return nil, err
	}
	return glAccs, nil
}

func (r *GLAccRepositoryImpl) GetGLAccByCode(code string) (*entities.GLAcc, error) {
	var glAcc entities.GLAcc
	err := r.db.Where("code = ?", code).First(&glAcc).Error
	return &glAcc, err
}
