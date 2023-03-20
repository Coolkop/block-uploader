package repository

type Entity struct {
	FileID    string `db:"file_id"`
	StorageID int    `db:"storage_id"`
	Order     int    `db:"part_order"`
}
