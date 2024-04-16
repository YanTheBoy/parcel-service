package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

const (
	errWrongStatus = "Wrong status"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) "+
		"VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		log.Println(err)
		return 0, err
	}

	id, err := res.LastInsertId()
	return int(id), err

}

func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}

	row := s.db.QueryRow("SELECT * FROM parcel WHERE number=:number", sql.Named("number", number))
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		log.Println(err)
	}

	p = Parcel{
		Number:    p.Number,
		Client:    p.Client,
		Status:    p.Status,
		Address:   p.Address,
		CreatedAt: p.CreatedAt,
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	p := Parcel{}
	var res []Parcel
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client=:client",
		sql.Named("client", client))
	if err != nil {
		log.Println(err)
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		res = append(res, p)
		if err != nil {
			log.Println(err)
			return res, err
		}
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status=:status WHERE number=:number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {

	row, err := s.Get(number)
	if err != nil {
		log.Println(err)
		return err
	}
	if row.Status != ParcelStatusRegistered {
		fmt.Errorf("cannot set address for parcel: %w", errWrongStatus)
		return errors.New("Wrong status")
	}

	_, err = s.db.Exec("UPDATE parcel SET address=:address WHERE number=:number",
		sql.Named("address", address),
		sql.Named("number", number))
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	res, err := s.Get(number)
	if err != nil {
		log.Println(err)
		return err
	}
	if res.Status != ParcelStatusRegistered {
		fmt.Errorf("cannot delete parcel: %w", errWrongStatus)
		return errors.New("Wrong status")
	}

	_, err = s.db.Exec("DELETE FROM parcel WHERE number=:number", sql.Named("number", number))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
