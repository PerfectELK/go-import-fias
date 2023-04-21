package download

import (
	"github.com/buger/goterm"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func File(
	link string,
	saveTo string,
) error {
	req, _ := http.NewRequest("GET", link, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	f, _ := os.OpenFile(saveTo, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	buf := make([]byte, 32*1024)
	var downloaded int64
	goterm.Clear()
	tBegin := time.Now()
	for {
		n, err := resp.Body.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if n > 0 {
			tNow := time.Now()
			if tNow.Unix()-tBegin.Unix() > 3 || downloaded == 0 {
				tBegin = tNow
				go func() {
					message := "Downloading... " + strconv.FormatFloat(float64(downloaded)/float64(resp.ContentLength)*100, 'f', 6, 64) + "%"
					goterm.MoveCursor(1, 1)
					goterm.Println(message)
					goterm.Flush()
				}()
			}

			f.Write(buf[:n])
			downloaded += int64(n)
		}
	}
	return nil
}
