package db

func (s *storage) IncreaseLikeCount(userID, contentID int64, contentType string) error {
	query := `
		INSERT INTO likes (user_id, content_id, content_type)
		VALUES ($1, $2, $3)
	`

	_, err := s.pg.Exec(query, userID, contentID, contentType)

	if err != nil && IsDuplicationError(err) {
		return ErrAlreadyExists
	} else if err != nil {
		return err
	}

	return nil
}

func (s *storage) DecreaseLikeCount(userID, contentID int64, contentType string) error {
	query := `
		DELETE FROM likes
		WHERE user_id = $1 AND content_id = $2 AND content_type = $3;
	`

	res, err := s.pg.Exec(query, userID, contentID, contentType)

	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrNotFound
	}

	return nil
}
