package persist

import (
	"archive/zip"
	"fmt"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/date"
	"io"
	"os"
	"path/filepath"
)

type ArchivableRecordContainer interface {
	GetArchivableRecords() []ArchivableRecord
	GetName() string
}

type ArchivableRecord interface {
	GetFileName() string
	Marshall(writer io.Writer) error
}

// ArchiveName returns the ArchivableRecordContainer file name with today's date
func ArchiveName(records ArchivableRecordContainer) string {
	return fmt.Sprintf("%s %s.zip", records.GetName(), date.Today().String())
}

// CreateRecordsArchive writes a zipped archive of all contained record files to the given writer
func CreateRecordsArchive(writer io.Writer, sites ArchivableRecordContainer) error {

	zw := zip.NewWriter(writer)

	for _, sr := range sites.GetArchivableRecords() {
		wfw, err := zw.Create(sr.GetFileName())
		if err != nil {
			return fmt.Errorf("could not create zip entry for site %s", sr.GetFileName())
		}
		err = sr.Marshall(wfw)
		if err != nil {
			return fmt.Errorf("could not write zip entry for site %s", sr.GetFileName())
		}
	}

	return zw.Close()
}

// CreateRecordsArchive writes a zipped archive of all contained record files to the given writer
func SaveRecordsArchive(path string, sites ArchivableRecordContainer) error {

	archiveFile, err := os.Create(filepath.Join(path, ArchiveName(sites)))
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	return CreateRecordsArchive(archiveFile, sites)
}
