package storage

type File struct {
	Id   int
	Name string
	Path string
}

type FileStorage interface {
	CreateMany(files []*File) error
	Create(file *File) error
	Delete(id int) error
	Read(id int) (File, error)
	ReadMany(page int, pageSize int) ([]File, error)
}
