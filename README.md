# isdf

Single-function package which takes an input [image.Image](https://godoc.org/image#Image) and returns an \*[image.Gray16](https://godoc.org/image#Gray16) Signed Distance Field representation. See example code for simple usage in a command line app.

![img1](http://i.imgur.com/bNI5Ujc.png)

This code has been adpated from [image-sdf](https://github.com/mattdesl/image-sdf).

## Usage Example

```
package main

import (
  ...
  "github.com/vrav/isdf"
)

func main() {
  // import an image, decode it if necessary
  // see isdf-example.go for an example
  ...
  
  // convert to 1/4 size SDF image with spread of 32
  sdfImage := isdf.ImageToSDF(inputImage, 32.0, 4.0)
  
  // continue on to save the image or use as-is
  ...
}
```

## License

MIT, see [LICENSE.md](http://github.com/vrav/isdf/blob/master/LICENSE.md) for details.
