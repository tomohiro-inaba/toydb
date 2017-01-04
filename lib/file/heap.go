package file

import (
	"errors"
	"fmt"

	"github.com/inazo1115/toydb/lib/page"
	"github.com/inazo1115/toydb/lib/storage"
	"github.com/inazo1115/toydb/lib/table"
)

// HeapFile is the representation of the file and access methods. HeapFile is
// the linked list of pages.
type HeapFile struct {
	rootPid int64
	bm      *storage.BufferManager
	schema  *table.Schema
}

// NewHeapFile creates HeapFile struct and returns it's pointer.
func NewHeapFile(bm *storage.BufferManager, schema *table.Schema) *HeapFile {

	rootPage := page.NewDataPage(-1, -1, -1, schema.RecordSize())
	rootPid, err := bm.Create(rootPage)
	if err != nil {
		panic(err)
	}

	return &HeapFile{rootPid, bm, schema}
}

func (f *HeapFile) RootPid() int64 {
	return f.rootPid
}

// Scan scans all records. The traversing begins from given page id.
func (f *HeapFile) Scan(pid int64) ([]*table.Record, error) {

	// Prepare the result variable.
	ret := make([]*table.Record, 0)

	// Traverse the linked list of pages.
	p := &page.DataPage{}
	next := pid
	for next != -1 {

		// Read the page.
		err := f.bm.Read(next, p)
		if err != nil {
			return nil, err
		}

		// Read records.
		for i := 0; i < int(p.NumRecords()); i++ {
			rec, err := f.schema.DeserializeRecord(p.ReadRecord(int64(i)))
			if err != nil {
				return nil, err
			}
			ret = append(ret, rec)
		}

		next = p.Next()
	}

	return ret, nil
}

// Insert inserts a record into the page.
func (f *HeapFile) Insert(pid int64, record *table.Record) error {

	p := &page.DataPage{}

	err := f.bm.Read(pid, p)
	if err != nil {
		return err
	}

	fmt.Println("**********************")
	fmt.Println(record.Values())
	fmt.Println(record.Values()[0])
	fmt.Println(record.Values()[1])

	// Insert the record into this page.
	if p.HasFreeSpace() {
		b, err := f.schema.SerializeRecord(record)
		if err != nil {
			return err
		}
		p.AddRecord(b)
		f.bm.Update(p.Pid(), p)
		return nil
	}

	// Follow the link to the next page.
	if p.Next() != -1 {
		return f.Insert(p.Next(), record)
	}

	// Create the new page and insert the record into it.
	newPage := page.NewDataPage(-1, p.Pid(), -1, f.schema.RecordSize())
	b, err := f.schema.SerializeRecord(record)
	if err != nil {
		return err
	}
	newPage.AddRecord(b)
	newPid, err := f.bm.Create(newPage)
	if err != nil {
		return err
	}
	// Set the link to the next page.
	if err = f.bm.Read(p.Pid(), p); err != nil {
		return err
	}
	p.SetNext(newPid)
	f.bm.Update(p.Pid(), p)

	return nil
}

// Dump is a debug function.
func (f *HeapFile) Dump(pid int64) {

	p := &page.DataPage{}

	var err error
	err = f.bm.Read(pid, p)
	if err != nil {
		panic(err)
	}

	fmt.Println("dump")
	fmt.Println(pid)
	fmt.Println(p)
}

func SearchEq(rootPid int64) error {
	return errors.New("not implemented")
}

func SearchRange(rootPid int64) error {
	return errors.New("not implemented")
}

func Delete(rootPid int64) error {
	return errors.New("not implemented")
}

func (f *HeapFile) WriteBackAll() error {
	err := f.bm.WriteBackAll()
	if err != nil {
		return err
	}
	return nil
}
