package service

import (
	"database/sql"
	"fmt"
	"log"
	"project/dto"
	"project/model"
)

type ClassService struct {
	DB *sql.DB
}

func (s *ClassService) CreateClass(req dto.CreateClassRequest) (*model.Class, error) {
	query := `INSERT INTO classes (name, jadwal_kelas, teacher, class_code) VALUES ($1, $2, $3, $4) RETURNING id, name, jadwal_kelas, created_at, teacher, class_code`
	var class model.Class
	err := s.DB.QueryRow(query, req.Name, req.JadwalKelas, req.Teacher, req.ClassCode).Scan(&class.ID, &class.Name, &class.JadwalKelas, &class.CreatedAt, &class.Teacher, &class.ClassCode)
	if err != nil {
		return nil, err
	}
	return &class, nil
}

func (s *ClassService) DeleteClass(id string) error {
	query := `DELETE FROM classes WHERE id = $1`
	result, err := s.DB.Exec(query, id)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to delete class: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *ClassService) GetClasses() ([]model.Class, error) {
	query := `SELECT id, name, jadwal_kelas, created_at, teacher, class_code FROM classes`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query classes: %v", err)
	}
	defer rows.Close()

	var classes []model.Class
	for rows.Next() {
		var class model.Class
		if err := rows.Scan(&class.ID, &class.Name, &class.JadwalKelas, &class.CreatedAt, &class.Teacher, &class.ClassCode); err != nil {
			return nil, fmt.Errorf("failed to scan class row: %v", err)
		}
		classes = append(classes, class)
	}

	return classes, nil
}

func (s *ClassService) GetClassByID(classID string) (*dto.ClassResponse, error) {
	query := `SELECT id, name, jadwal_kelas, teacher, class_code, created_at FROM classes WHERE id = $1`

	var class dto.ClassResponse
	err := s.DB.QueryRow(query, classID).Scan(
		&class.ID,
		&class.Name,
		&class.JadwalKelas,
		&class.Teacher,
		&class.ClassCode,
		&class.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &class, nil
}

func (s *ClassService) UpdateClass(classID string, req dto.UpdateClassRequest) (*dto.ClassResponse, error) {
	query := `UPDATE classes SET name = COALESCE($1, name), jadwal_kelas = COALESCE($2, jadwal_kelas), teacher = COALESCE($3, teacher), class_code = COALESCE($4, class_code) WHERE id = $5 RETURNING id, name, jadwal_kelas, teacher, class_code, created_at`

	var updatedClass dto.ClassResponse
	err := s.DB.QueryRow(
		query,
		req.Name,
		req.JadwalKelas,
		req.Teacher,
		req.ClassCode,
		classID,
	).Scan(
		&updatedClass.ID,
		&updatedClass.Name,
		&updatedClass.JadwalKelas,
		&updatedClass.Teacher,
		&updatedClass.ClassCode,
		&updatedClass.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("class not found")
		}
		return nil, fmt.Errorf("error updating class: %w", err)
	}

	return &updatedClass, nil
}

func (s *ClassService) JoinClass(userID int, classCode string) error {
	var classID int

	query := `SELECT id FROM classes WHERE class_code = $1`
	err := s.DB.QueryRow(query, classCode).Scan(&classID)
	if err != nil {
		log.Printf("Failed to find class by class_code: %v\n", err)
		return fmt.Errorf("class not found")
	}

	var userExists bool
	checkUserQuery := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`
	err = s.DB.QueryRow(checkUserQuery, userID).Scan(&userExists)
	if err != nil || !userExists {
		log.Printf("UserID: %d not found in users table\n", userID)
		return fmt.Errorf("user not found")
	}

	log.Printf("UserID: %d, is joining class ID: %d\n", userID, classID)

	insertQuery := `INSERT INTO class_members (user_id, class_id, role) VALUES ($1, $2, $3)`
	_, err = s.DB.Exec(insertQuery, userID, classID, "siswa")
	if err != nil {
		log.Printf("Failed to insert user into class_members: %v\n", err)
		return fmt.Errorf("failed to join class")
	}

	return nil
}

func (s *ClassService) AddMember(classID, userID int) error {
	query := `INSERT INTO class_members (class_id, user_id, role) VALUES ($1, $2, 'siswa')`
	_, err := s.DB.Exec(query, classID, userID)
	if err != nil {
		return fmt.Errorf("failed to add member: %v", err)
	}
	return nil
}

func (s *ClassService) RemoveMember(classID, userID int) error {
	query := `DELETE FROM class_members WHERE class_id = $1 AND user_id = $2`
	_, err := s.DB.Exec(query, classID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove member: %v", err)
	}
	return nil
}

func (s *ClassService) GetMembers(classID int) ([]model.User, error) {
	query := `SELECT u.id, u.username, u.role FROM users u
              JOIN class_members cm ON u.id = cm.user_id
              WHERE cm.class_id = $1`
	rows, err := s.DB.Query(query, classID)
	if err != nil {
		return nil, fmt.Errorf("failed to get members: %v", err)
	}
	defer rows.Close()

	var members []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Role); err != nil {
			return nil, fmt.Errorf("failed to scan member: %v", err)
		}
		members = append(members, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating members: %v", err)
	}

	return members, nil
}

func (s *ClassService) GetClassesByStudentID(studentID int) ([]model.Class, error) {
    query := `SELECT c.id, c.name, c.jadwal_kelas, c.created_at, c.teacher, c.class_code
              FROM classes c
              JOIN class_members cm ON c.id = cm.class_id
              WHERE cm.user_id = $1`
    rows, err := s.DB.Query(query, studentID)
    if err != nil {
        return nil, fmt.Errorf("failed to get classes: %v", err)
    }
    defer rows.Close()

    var classes []model.Class
    for rows.Next() {
        var class model.Class
        if err := rows.Scan(&class.ID, &class.Name, &class.JadwalKelas, &class.CreatedAt, &class.Teacher, &class.ClassCode); err != nil {
            return nil, fmt.Errorf("failed to scan class: %v", err)
        }
        classes = append(classes, class)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating classes: %v", err)
    }

    return classes, nil
}

func (s *ClassService) CountClassesByUserID(userID int) (int, error) {
    query := `
        SELECT COUNT(*)
        FROM classes c
        JOIN class_members cm ON c.id = cm.class_id
        WHERE cm.user_id = $1
    `
    var count int
    err := s.DB.QueryRow(query, userID).Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("failed to count classes: %w", err)
    }

    return count, nil
}