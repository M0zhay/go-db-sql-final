package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
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

func ConnectToDb(driver, database string) (*sql.DB, error) {
	db, err := sql.Open(driver, database)
	if err != nil {
		return db, err
	}
	return db, err
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := ConnectToDb("sqlite", "tracker.db")
	require.NoError(t, err)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// get
	res, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, res.Client, parcel.Client)
	assert.Equal(t, res.Status, parcel.Status)
	assert.Equal(t, res.Address, parcel.Address)
	assert.Equal(t, res.CreatedAt, parcel.CreatedAt)

	// delete
	err = store.Delete(id)
	require.NoError(t, err)

	res, err = store.Get(parcel.Number)
	require.Error(t, err, sql.ErrNoRows)
	require.Empty(t, res)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := ConnectToDb("sqlite", "tracker.db")
	require.NoError(t, err)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// set address
	newAddress := "Night City"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	// check
	res, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, res.Address, "Night City")
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := ConnectToDb("sqlite", "tracker.db")
	require.NoError(t, err)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// set status
	err = store.SetStatus(id, ParcelStatusDelivered)
	require.NoError(t, err)

	// check
	res, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, res.Status, ParcelStatusDelivered)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := ConnectToDb("sqlite", "tracker.db")
	require.NoError(t, err)
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
		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.ElementsMatch(t, storedParcels, parcels)

	// check
	for _, parcel := range storedParcels {
		expectedParcel, ok := parcelMap[parcel.Number]
		require.True(t, ok)
		require.Equal(t, expectedParcel, parcel)
	}

}
