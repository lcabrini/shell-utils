package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

const (
	imageHeight  int    = 400
	ankiImageDir string = "%s/Immagini/Anki"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No images supplied.")
		os.Exit(1)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Cannot get user home directory")
		os.Exit(1)
	}

	destinationDir := fmt.Sprintf(ankiImageDir, home)
	if _, err := os.Stat(destinationDir); err != nil {
		fmt.Println("Anki image directory does not exist")
		os.Exit(1)
	}

	for _, fn := range os.Args[1:] {
		if _, err := os.Stat(fn); err != nil {
			fmt.Printf("skipping '%s': %s\n", fn, err)
			continue
		}

		mtype, err := mimetype.DetectFile(fn)
		if err != nil {
			fmt.Printf("skipping '%s': %s\n", fn, err.Error())
			continue
		}
		if !strings.HasPrefix(mtype.String(), "image") {
			fmt.Printf("skipping '%s': it does not seem to be an image\n", fn)
			continue
		}

		strippedName := strings.TrimSuffix(fn, filepath.Ext(fn))
		pngName := fmt.Sprintf("%s.png", strippedName)
		if pngName != fn {
			if _, err := os.Stat(pngName); err == nil {
				fmt.Printf("skipping '%s': PNG image already exists\n", fn)
				continue
			}
		}

		height, err := getImageHeight(fn)
		if err != nil {
			fmt.Println("unable to get image height: %s\n", err)
			continue
		}

		magickCmd := []string{fn}

		if height < imageHeight {
			fmt.Printf("warning: height of '%s' is less than %d, not resizing\n", fn, imageHeight)
		} else {
			extra := []string{
				"-resize",
				fmt.Sprintf("x%d", imageHeight),
			}
			magickCmd = append(magickCmd, extra...)
		}
		magickCmd = append(magickCmd, pngName)

		cmd := exec.Command("magick", magickCmd...)
		_, err = cmd.Output()
		if err != nil {
			fmt.Println(err)
			continue
		}

		dest := fmt.Sprintf("%s/%s", destinationDir, pngName)
		if _, err = os.Stat(dest); err == nil {
			fmt.Printf("skipping '%s': desination file already exists\n", pngName)
			continue
		}

		if err := os.Rename(pngName, dest); err != nil {
			fmt.Printf("Unable to move '%s' to Anki image directory: %s", pngName, err)
			continue
		}

		if err = os.Remove(fn); err != nil {
			fmt.Printf("Unable to remove '%s': %s\n", fn, err)
		}
	}
}

func getImageHeight(fn string) (int, error) {
	cmd := exec.Command("magick", "identify", "-format", "%h", fn)
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	height, err := strconv.Atoi(string(out))
	if err != nil {
		return 0, err
	}

	return height, nil
}
