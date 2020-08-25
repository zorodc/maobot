package imgfetch

import "encoding/base64"
import "bytes"
import "strings"
import "io/ioutil"
import "net/http"
import "errors"

func encode(data []byte) (encoded []byte) {
	buffer  := bytes.NewBuffer(nil)
	encoder := base64.NewEncoder(base64.StdEncoding, buffer)
	encoder.Write(data)
	encoder.Close() // flush out
	return buffer.Bytes()
}

// The set of supported mimetypes.
var supported = map[string]struct{} {
	"image/png":struct{}{},
	"image/jpeg":struct{}{},
	"image/gif":struct{}{},
	//	"image/bmp":struct{}{}, // not supported by the client
	// "image/webp":struct{}{},
	//	"image/x-icon":struct{}{},
	//	"image/svg+xml":struct{}{},
}

func FetchImage(url string) (data []byte, mimetype string, ok error) {
	if !(strings.HasPrefix(url, "http://")  ||
		  strings.HasPrefix(url, "https://") ||
		  strings.HasPrefix(url, "www.")) {
		ok = errors.New("Argument not a URL.")
		return
	}

	response, ok := http.Get(url)
	if ok == nil {
		defer response.Body.Close()
	} else { return; }

	bytes, ok := ioutil.ReadAll(response.Body)
	if ok != nil {
		return
	}

	if _, in := supported[http.DetectContentType(bytes)]; !in {
		ok = errors.New("Url does not point to an acceptible mimetype.")
		return
	}

	data = encode(bytes)
	return
}
