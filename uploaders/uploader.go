package uploaders

type Uploader interface {
	Upload(src string) (url string, err error)
}
