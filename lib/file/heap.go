package file

import (
	"errors"
	"fmt"

	"github.com/inazo1115/toydb/lib/page"
	"github.com/inazo1115/toydb/lib/storage"
)

type HeapFile struct {
	bm *storage.BufferManager
}

func NewHeapFile() *HeapFile {
	return &HeapFile{storage.NewBufferManager()}
}

// tmp
func deserializeRecord(b []byte) string {
	return string(b)
}

func serializeRecord(s string) []byte {
	return []byte(s)
}

//func Scan(rootPid int64) []*page.Record {
func (f *HeapFile) Scan(rootPid int64) (*[]string, error) {

	fmt.Println(rootPid)

	ret := make([]string, 0)

	p := &page.DataPage{}
	next := rootPid

	for next != -1 {

		_, err := f.bm.Read(int(next), p)
		if err != nil {
			return nil, err
		}

		for i := 0; i < int(p.NumRecords()); i++ {
			rec := deserializeRecord(p.ReadRecord(i))
			ret = append(ret, rec)
		}

		next = p.Next()
	}

	return &ret, nil
}

func (f *HeapFile) Insert(rootPid int64, record string) error {

	p := &page.DataPage{}

	var err error
	p, err = f.bm.Read(int(rootPid), p)
	if err != nil {
		return err
	}

	// Insert the record into this page.
	if p.HasFreeSpace() {
		p.AddRecord(serializeRecord(record))
		return nil
	}

	// Try to insert the record into the next page.
	if p.Next() != -1 {
		return f.Insert(p.Next(), record)
	}

	// Insert the record into the new page.
	newPage := page.NewDataPage(-1, p.Pid(), -1)
	newPage.AddRecord(serializeRecord(record))
	newPid, err := f.bm.Create(newPage)
	if err != nil {
		return err
	}
	p.SetNext(int64(newPid))

	return nil
}

func (f *HeapFile) Dump(pid int64) {

	p := &page.DataPage{}

	var err error
	p, err = f.bm.Read(int(pid), p)
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
