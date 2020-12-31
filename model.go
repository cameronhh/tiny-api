package main

import (
	"database/sql"
)

type endpoint struct {
	ID      int    `json:"id"`
	URL     string `json:"url"`
	Content string `json:"content"`
}

func (p *endpoint) getEndpoint(db *sql.DB) error {
	return db.QueryRow("SELECT url, content FROM endpoints WHERE id=$1",
		p.ID).Scan(&p.URL, &p.Content)
}

func (p *endpoint) getEndpointByURL(db *sql.DB) error {
	return db.QueryRow("SELECT id, url, content FROM endpoints WHERE url=$1",
		p.URL).Scan(&p.ID, &p.URL, &p.Content)
}

func (p *endpoint) updateEndpoint(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE endpoints SET url=$1, content=$2 WHERE id=$3",
			p.URL, p.Content, p.ID)

	return err
}

func (p *endpoint) deleteEndpoint(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM endpoints WHERE id=$1", p.ID)

	return err
}

func (p *endpoint) createEndpoint(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO endpoints(url, content) VALUES($1, $2) RETURNING id",
		p.URL, p.Content).Scan(&p.ID)

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
