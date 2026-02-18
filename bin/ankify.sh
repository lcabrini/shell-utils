#! /bin/bash

# Prepare images for Anki.
#
# The following steps are performed:
#   1. Resize image so that the height is $image_height (default: 400 pixels)
#   2. Convert image to PNG
#   3. Move converted image to ~/Pictures/Anki
#
# Requirements:
#   ImageMagick
#
# Author: Lorenzo Cabrini <lorenzo.cabrini@gmail.com>

image_height=400
destdir=~/Immagini/Anki

if [[ $# -eq 0 ]]; then
    echo "No images supplied."
    exit 1
fi

for img in "$@"; do
    if [[ ! -f $img ]]; then
	echo "skipping '$img': file not found"
	continue
    fi

    mime=$(file -b --mime-type "$img")
    if [[ $mime != image/* ]]; then
	echo "skipping '$img': it does not seem to be an image"
	continue
    fi

    png=${img%.*}.png
    if [[ $png != $img ]]; then
	if [[ -f $png ]]; then
	    echo "skipping '$img': corresponding PNG image already exists"
	    continue
	fi
    fi

    height=$(magick identify -format "%h" "$img" 2> /dev/null)
    if [[ $height -lt $image_height ]]; then
	echo "warning: height of '$img' is less than 400px, not resizing"
	cp "$img" "$png" > /dev/null 2>&1
    else
	magick "$img" -resize x$image_height "$png" 2> /dev/null
    fi

    if [[ -f $destdir/$png ]]; then
	echo "warning: skipping '$png': destination file already exists"
	continue
    else
	mv "$png" $destdir
    fi

    rm -f "$img"
done
