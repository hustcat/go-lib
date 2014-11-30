package main

import (
    "fmt"
    "github.com/hustcat/go-lib/bitmap"
)

func main(){
    bitmap := bitmap.NewNumaBitmap()

    //node 0
    bitmap.SetBit(0, 1)
    bitmap.SetBit(5, 1)

    //node 1
    bitmap.SetBit(6, 1)
    bitmap.SetBit(11, 1)

    //node 0
    bitmap.SetBit(12, 1)
    bitmap.SetBit(17, 1)

    //node 1
    bitmap.SetBit(18, 1)
    bitmap.SetBit(23, 1)

    actual := bitmap.Get1BitOffsNuma(2)
    expected := [][]uint{
        []uint{0, 5, 12, 17},
        []uint{6, 11, 18, 23},
    }
    a := fmt.Sprintf("%v", expected)
    b := fmt.Sprintf("%v", actual)
    if a != b {
        fmt.Println("expected:%v, actual:%v", a, b)
    }else{
        fmt.Println("pass")
    }
}
