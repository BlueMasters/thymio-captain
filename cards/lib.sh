#!/usr/bin/env bash

page() {
    pageno=$1
    for i in $(seq 8); do
        x=$(../genid/genid -short -key "$SECRET")
        echo "$URL/$x" | qrencode -s 10 -o $temp/q.png
        composite -geometry +700+346  $temp/q.png $SOURCE $temp/res$i.png
    done
    outfile=$(printf 'page-%03i.png' $pageno)
    echo $outfile
    montage $temp/res[12345678].png -geometry +0+0 -tile 2x4 $temp/$outfile
    rm $temp/res[12345678].png
}

doc() {
    npages=$1
    home=$(pwd)
    temp=$(mktemp -d)

    for i in $(seq $npages); do
        page $i
    done

    convert $temp/page-???.png -page a4 $DEST
    rm $temp/page*
    rm $temp/q.png
    rmdir $temp
}