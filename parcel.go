package main

import (
	"database/sql"
	"log"
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
		return 0, nil
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

	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	// заполните срез Parcel данными из таблицы

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

	row, _ := s.Get(number)
	if row.Status != ParcelStatusRegistered {
		log.Printf("Status is not %s\n", ParcelStatusRegistered)
		return nil
	}

	_, err := s.db.Exec("UPDATE parcel SET address=:address WHERE number=:number",
		sql.Named("address", address),
		sql.Named("number", number))
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	res, _ := s.Get(number)
	if res.Status != ParcelStatusRegistered {
		log.Printf("Status is not %s", ParcelStatusRegistered)
		return nil
	}

	_, err := s.db.Exec("DELETE FROM parcel WHERE number=:number", sql.Named("number", number))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
