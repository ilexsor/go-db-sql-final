package models

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

type ParcelStore struct {
	db *sql.DB
}

// NewParcelStore создает объект хранилища посылок (sqlite DB)
func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

// Add метод структуры ParcelStore, добавляет в БД новую посылку
func (s ParcelStore) Add(p Parcel) (int, error) {

	resultSet, err := s.db.Exec("INSERT INTO parcel(client, status, address, created_at) VALUES(:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))

	if err != nil {
		return 0, err
	}
	log.WithFields(log.Fields{
		"Message": "Execute INSERT statment",
	}).Info()

	parcelId, err := resultSet.LastInsertId()
	if err != nil {
		return 0, err
	}
	log.WithFields(log.Fields{
		"Message": fmt.Sprintf("Fetch record id: %v", parcelId),
	}).Info()

	return int(parcelId), nil
}

// Get метод структуры ParcelStore, возвращает посылку по id из БД и ошибку при возникновении
func (s ParcelStore) Get(number int) (Parcel, error) {

	resultSet := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number", sql.Named("number", number))

	log.WithFields(log.Fields{
		"Message": "Execute SELECT statment get by id",
	}).Info()

	p := Parcel{}
	err := resultSet.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {

		log.WithFields(log.Fields{
			"error_message": fmt.Sprintf("scan results error: %v", err),
		}).Error()

		return p, err
	}

	log.WithFields(log.Fields{
		"Message": fmt.Sprintf("Parcel: %v", p),
	}).Info()

	return p, nil
}

// GetByClient метод структуры ParcelStore, возвращает посылку по id клиента из БД и ошибку при возникновении
func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {

	resultSet, err := s.db.Query("SELECT * FROM parcel WHERE client = :client", sql.Named("client", client))
	defer resultSet.Close()

	if err != nil {
		return []Parcel{}, err
	}

	log.WithFields(log.Fields{
		"Message": "Execute SELECT statment get by client",
	}).Info()

	var res []Parcel
	parcel := Parcel{}

	for resultSet.Next() {
		err := resultSet.Scan(&parcel.Number, &parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt)
		if err != nil {

			log.WithFields(log.Fields{
				"error_message": fmt.Sprintf("scan results error: %v", err),
			}).Error()

			return res, err
		}

		res = append(res, parcel)
	}

	return res, nil
}

// SetStatus метод структуры ParcelStore, обновляет статус посылки в БД
func (s ParcelStore) SetStatus(number int, status string) error {

	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))

	if err != nil {

		log.WithFields(log.Fields{
			"error_message": fmt.Sprintf("update status error: %v", err),
		}).Error()
		return err
	}

	return nil
}

// SetAddress метод структуры ParcelStore, обновляет адрес доставки посылки а БД
func (s ParcelStore) SetAddress(number int, address string) error {

	parcel, err := s.Get(number)
	if err != nil {
		return err
	}

	if parcel.Status == ParcelStatusRegistered {
		_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
			sql.Named("address", address),
			sql.Named("number", number))

		if err != nil {

			log.WithFields(log.Fields{
				"error_message": fmt.Sprintf("update address error: %v", err),
			}).Error()

			return err
		}
	}

	return nil
}

// Delete метод структуры ParcelStore, удаляет посылку из БД по id, возвращает ощибку при возникновении
func (s ParcelStore) Delete(number int) error {

	parcel, err := s.Get(number)
	if err != nil {
		return err
	}

	if parcel.Status == ParcelStatusRegistered {
		_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number",
			sql.Named("number", number))

		if err != nil {

			log.WithFields(log.Fields{
				"error_message": fmt.Sprintf("delete parcel error: %v", err),
			}).Error()

			return err
		}
	}

	return nil
}
