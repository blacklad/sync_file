package storage

import (
	"blacklad.com/sync_file/utils"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type DbFile struct {
	fileList []*DbFileStat
	path     string
	db       *sql.DB
}

type DbFileStat struct {
	Id int
	utils.FileStat
}

func NewDbFile(path string) (*DbFile, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	d := &DbFile{
		fileList: nil,
		path:     path,
		db:       db,
	}
	err = d.createTable()
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *DbFile) List() error {
	rows, err := d.db.Query("SELECT * FROM db_file")
	if err != nil {
		return err
	}

	var fileListTemp = make([]*DbFileStat, 0)
	for rows.Next() {
		var id int
		var path string
		var md5 string
		var fileType utils.FileType
		var lastModified int64
		var version int64
		err = rows.Scan(&id, &path, &md5, &fileType, &lastModified, &version)
		if err != nil {
			return err
		}

		dbFileStat := &DbFileStat{
			Id: id,
			FileStat: utils.FileStat{
				Path:         path,
				MD5:          md5,
				FileType:     fileType,
				LastModified: lastModified,
				Version:      version,
			},
		}

		fileListTemp = append(fileListTemp, dbFileStat)
	}

	d.fileList = fileListTemp
	return nil
}

func (d *DbFile) GetByPath(path string) (*DbFileStat, error) {
	rows, err := d.db.Query("SELECT * FROM db_file WHERE path = ?", path)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		var id int
		var path string
		var md5 string
		var fileType utils.FileType
		var lastModified int64
		var version int64
		err = rows.Scan(&id, &path, &md5, &fileType, &lastModified, &version)
		if err != nil {
			return nil, err
		}

		dbFileStat := &DbFileStat{
			Id: id,
			FileStat: utils.FileStat{
				Path:         path,
				MD5:          md5,
				FileType:     fileType,
				LastModified: lastModified,
				Version:      version,
			},
		}
		defer rows.Close()
		return dbFileStat, nil
	}
	return nil, nil
}

func (d *DbFile) GetMaxVersion() (int64, error) {
	rows, err := d.db.Query("SELECT max(version) FROM db_file")
	if err != nil {
		return 0, err
	}

	if rows.Next() {
		var version int64
		err = rows.Scan(&version)
		if err != nil {
			return 0, err
		}

		defer rows.Close()
		return version, nil
	}
	return 0, nil
}

func (d *DbFile) GetByHistoryVersion(version int64) ([]*DbFileStat, error) {
	rows, err := d.db.Query("SELECT * FROM db_file WHERE version < ?", version)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fileListTemp = make([]*DbFileStat, 0)
	if rows.Next() {
		var id int
		var path string
		var md5 string
		var fileType utils.FileType
		var lastModified int64
		var version int64
		err = rows.Scan(&id, &path, &md5, &fileType, &lastModified, &version)
		if err != nil {
			return nil, err
		}

		dbFileStat := &DbFileStat{
			Id: id,
			FileStat: utils.FileStat{
				Path:         path,
				MD5:          md5,
				FileType:     fileType,
				LastModified: lastModified,
				Version:      version,
			},
		}
		fileListTemp = append(fileListTemp, dbFileStat)
	}

	return fileListTemp, nil
}

func (d *DbFile) Add(fileStat *DbFileStat) (int64, error) {
	//插入数据
	stmt, err := d.db.Prepare("INSERT INTO db_file(path, md5, file_type, last_modified, version) values(?,?,?,?,?)")
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(fileStat.Path, fileStat.MD5, fileStat.FileType, fileStat.LastModified, fileStat.Version)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (d *DbFile) DeleteById(id int) (int64, error) {
	//删除数据
	stmt, err := d.db.Prepare("delete from db_file where id=?")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (d *DbFile) UpdateById(id int, fileStat *DbFileStat) (int64, error) {
	//删除数据
	stmt, err := d.db.Prepare("update db_file set path = ?, md5 = ?, file_type = ?, last_modified = ?, version = ? where id=?")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(fileStat.Path, fileStat.MD5, fileStat.FileType, fileStat.LastModified, fileStat.Version, id)
	if err != nil {
		return 0, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (d *DbFile) createTable() error {
	sql_table := `
	CREATE TABLE IF NOT EXISTS "db_file" (
	"id" INTEGER PRIMARY KEY AUTOINCREMENT,
	"path" VARCHAR(64)  NOT NULL,
	"md5" VARCHAR(64)  NOT NULL,
	"file_type" VARCHAR(64)  NOT NULL,
	"last_modified" INTEGER  NOT NULL,
	"version" INTEGER NOT NULL
	);`

	_, err := d.db.Exec(sql_table)
	return err
}
