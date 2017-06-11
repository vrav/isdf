package main

import (
  "flag"
  "fmt"
  "os"
  
  "image"
  "image/png"
  
  "github.com/vrav/isdf"
)

var inFilename = flag.String("in", "in.png", "input image to process")
var outFilename = flag.String("out", "out.png", "output file to save to")
var spread = flag.Float64("spread", 32.0, "spread of distance field")
var downscale = flag.Float64("scale", 4.0, "image downscaling; 2 = half size")

func main() {
  flag.Parse()
  fmt.Println("Input file:", *inFilename)

  // try to load input file
  inFile, err := os.Open(*inFilename)
  if err != nil {
    panic(err)
  }
  defer inFile.Close()

  // decode input file data, import more to support more formats (ie, png)
  inImage, _, err := image.Decode(inFile)
  if err != nil {
    panic(err)
  }
  
  // convert input image to Signed Distance Field representation
  outImage := isdf.ImageToSDF(inImage, *spread, *downscale)
  
  // create output file
  out, err := os.Create(*outFilename)
  if err != nil {
    panic(err)
  }
  defer out.Close()
  
  // encode output image as png
  err = png.Encode(out, outImage)
  if err != nil {
    panic(err)
  }
  fmt.Println("Wrote", *outFilename)
}
