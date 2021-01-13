package main

import (
	"database/sql"
	"net/url"
)

type endpoint struct {
	ID      int    `json:"id"`
	URL     string `json:"url"`
	Content string `json:"content"`
}

func (e *endpoint) isValidURL() bool {
	// only URL safe chars are allowed in the endpoint so we
	// escape and unescape the endpoint url and verify that all are the same
	escapedURL := url.PathEscape(e.URL)
	unescapedURL, err := url.PathUnescape(e.URL)

	return err == nil && escapedURL == unescapedURL
}

func (e *endpoint) getEndpoint(db *sql.DB) error {
	return db.QueryRow("SELECT url, content FROM endpoints WHERE id=$1",
		e.ID).Scan(&e.URL, &e.Content)
}

func (e *endpoint) getEndpointByURL(db *sql.DB) error {
	return db.QueryRow("SELECT id, url, content FROM endpoints WHERE url=$1",
		e.URL).Scan(&e.ID, &e.URL, &e.Content)
}

func (e *endpoint) updateEndpoint(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE endpoints SET url=$1, content=$2 WHERE id=$3",
			e.URL, e.Content, e.ID)

	return err
}

func (e *endpoint) deleteEndpoint(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM endpoints WHERE id=$1", e.ID)

	return err
}

func (e *endpoint) createEndpoint(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO endpoints(url, content) VALUES($1, $2) RETURNING id",
		e.URL, e.Content).Scan(&e.ID)

	if err != nil {
		return err
	}

	return nil
}

func getEndpoints(db *sql.DB, start, count int) ([]endpoint, error) {
	rows, err := db.Query(
		"SELECT id, url, content FROM endpoints LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	endpoints := []endpoint{}

	for rows.Next() {
		var p endpoint
		if err := rows.Scan(&p.ID, &p.URL, &p.Content); err != nil {
			return nil, err
		}
		endpoints = append(endpoints, p)
	}

	return endpoints, nil
}
