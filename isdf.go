// Package isdf provides a single function to convert an input image into its Signed Distance Field representation.
// See https://github.com/mattdesl/image-sdf for the reference implementation used.
package isdf

import (
  "math"
  
  "image"
  "image/color"
)

// ImageToSDF converts an image using the image.Image interface to its Signed Distance Field representation.
// Spread affects the width of the distance fields while downscale decreases the resolution of the output image.
// A downscale value of 2.0 will turn a 512px image to 256px; 4.0 will result in a 128px output image.
// Converting from a high resolution image to a low resolution image results in better quality.
func ImageToSDF(src image.Image, spread, downscale float64) *image.Gray16 {
  sb := src.Bounds()
  sw, sh := sb.Max.X, sb.Max.Y
  
  bitMask := image.NewGray(sb)
  for y := 0; y < sh; y++ {
    for x := 0; x < sw; x++ {
      grayValue, _, _, _ := color.GrayModel.Convert(src.At(x, y)).RGBA()
      floatValue := float32(grayValue) / 65535.0
      if floatValue > 0.5 {
        bitMask.SetGray(x, y, color.Gray{255})
      } else {
        bitMask.SetGray(x, y, color.Gray{0})
      }
    }
  }
  
  sdfImage := sdfCompute(bitMask, spread, downscale)
  return sdfImage
}

func sdfCompute(bitMask *image.Gray, spread, downscale float64) *image.Gray16 {
  bb := bitMask.Bounds()
  bw, bh := bb.Max.X, bb.Max.Y
  
  ow := math.Floor(float64(bw) / downscale)
  oh := math.Floor(float64(bh) / downscale)
  out := image.NewGray16(image.Rect(0, 0, int(ow), int(oh)))
  
  for y := 0.0; y < oh; y++ {
    for x := 0.0; x < ow; x++ {
      centerX := math.Floor(x * downscale + downscale * 2.0)
      centerY := math.Floor(y * downscale + downscale * 2.0)
      
      signedDistance := findSignedDistance(bitMask, bw, bh, centerX, centerY, spread)
      
      alpha := 0.5 + 0.5 * (signedDistance / spread)
      alpha = math.Floor(math.Min(math.Max(0.0, alpha), 1.0) * 65535.0)
      
      out.SetGray16(int(x), int(y), color.Gray16{uint16(alpha)})
    }
  }
  
  return out
}

func findSignedDistance(bitMask *image.Gray, bw, bh int, centerX, centerY, spread float64) float64 {
  base := bitMask.GrayAt(int(centerX), int(centerY))
  delta := math.Ceil(spread)
  startX := math.Max(0.0, centerX - delta)
  endX := math.Min(float64(bw - 1), centerX + delta)
  startY := math.Max(0.0, centerY - delta)
  endY := math.Min(float64(bh - 1), centerY + delta )
  
  closestSquareDist := delta * delta
  
  for y := startY; y <= endY; y++ {
    for x := startX; x <= endX; x++ {
      if base != bitMask.GrayAt(int(x), int(y)) {
        sqDist := squareDist(centerX, centerY, x, y)
        if sqDist < closestSquareDist {
          closestSquareDist = sqDist
        }
      }
    }
  }
  
  closestDist := math.Min(math.Sqrt(closestSquareDist), spread)
  if base.Y == 255 {
    return 1.0 * closestDist
  } else {
    return -1.0 * closestDist
  }
}

func squareDist(x1, y1, x2, y2 float64) float64 {
  dx := x1 - x2
  dy := y1 - y2
  return dx*dx + dy*dy
}
