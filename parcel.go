package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	query := "INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)"
	result, err := s.db.Exec(query,
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	query := "SELECT * FROM parcel WHERE number = :number"
	res := s.db.QueryRow(query,
		sql.Named("number", number))

	p := Parcel{}
	err := res.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	query := "SELECT * FROM parcel WHERE client = :client"
	rows, err := s.db.Query(query,
		sql.Named("client", client))

	if err != nil {
		return nil, fmt.Errorf("GetByClient: no data with client id %d: %v", client, err)
	}
	defer rows.Close()

	var res []Parcel

	for rows.Next() {
		var p Parcel

		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	query := "UPDATE parcel SET status = :status WHERE number = number"
	_, err := s.db.Exec(query,
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		return fmt.Errorf("SetStatus error: %v", err)
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	query := "UPDATE parcel SET address = :address WHERE number = :number AND status = :status"

	_, err := s.db.Exec(query,
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return fmt.Errorf("SetAddress error: %v", err)
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	query := "DELETE FROM parcel WHERE number = :number and status = :status"

	_, err := s.db.Exec(query,
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return fmt.Errorf("Delete error: %s", err)
	}
	return nil
}
