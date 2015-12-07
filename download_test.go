package main

import (
	"fmt"
	"testing"
)

func TestDownload(t *testing.T) {
	str, err := DownloadPage()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(str)
}
