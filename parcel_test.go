package main

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEqualValues(t, 0, id)

	row, err := store.Get(id)
	require.NoError(t, err)

	assert.Equal(t, parcel.Address, row.Address, "Address is incorrect")
	assert.Equal(t, parcel.Client, row.Client, "Client is incorrect")
	assert.Equal(t, parcel.Status, row.Status, "Status is incorrect")
	assert.Equal(t, parcel.CreatedAt, row.CreatedAt, "CreatedAt is incorrect")

	err = store.Delete(parcel.Number)
	require.NoError(t, err)

	row, err = store.Get(parcel.Number)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)

	newAddress := "new test address"

	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	row, err := store.Get(id)
	require.NoError(t, err)

	assert.Equal(t, newAddress, row.Address, "Address is incorrect")

}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEqualValues(t, 0, id)

	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err)

	p, _ := store.Get(id)
	assert.Equal(t, ParcelStatusSent, p.Status, "Status is not correct")
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotEmpty(t, id)

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Len(t, storedParcels, len(parcels))

	// получите список посылок по идентификатору клиента, сохранённого в переменной client
	// убедитесь в отсутствии ошибки
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных

	// check
	for _, parcel := range storedParcels {
		require.NotEmpty(t, parcelMap[parcel.Number])
		assert.Equal(t, parcelMap[parcel.Number].Client, parcel.Client)
		assert.Equal(t, parcelMap[parcel.Number].Address, parcel.Address)
		assert.Equal(t, parcelMap[parcel.Number].CreatedAt, parcel.CreatedAt)
		assert.Equal(t, parcelMap[parcel.Number].Status, parcel.Status)

		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно
	}
}
