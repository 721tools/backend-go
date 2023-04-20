package files

type FileIO interface {
	Read(data interface{}) error
	Write(data interface{}) error
}
