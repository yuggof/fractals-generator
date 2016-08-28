package main

import (
  "image"
  "image/color"
  "image/png"
  "os"
  "math/cmplx"
  "flag"
)

type flags struct {
  Width, Height int
  LowerLeft, UpperRight, Z complex128
  Output string
}

func parseFlags() *flags {
  w := flag.Int("w", 800, "")
  h := flag.Int("h", 800, "")
  llr := flag.Float64("llr", -1, "")
  lli := flag.Float64("lli", -1, "")
  urr := flag.Float64("urr", 1, "")
  uri := flag.Float64("uri", 1, "")
  zr := flag.Float64("zr", 0, "")
  zi := flag.Float64("zi", 0, "")
  o := flag.String("o", "output.png", "")
  flag.Parse()

  return &flags{
    *w, *h,
    complex(*llr, *lli), complex(*urr, *uri), complex(*zr, *zi),
    *o,
  }
}

func gradient(a, b, i float64) float64 {
  return a + (b - a) * i
}

func generateFractal(flags *flags) *image.RGBA {
  img := image.NewRGBA(image.Rect(0, 0, flags.Width, flags.Height))

  for x := 0; x < flags.Width; x++ {
    for y := 0; y < flags.Height; y++ {
      c := complex(
        gradient(real(flags.LowerLeft), real(flags.UpperRight), (float64(x) / float64(flags.Width))),
        gradient(imag(flags.LowerLeft), imag(flags.UpperRight), (float64(y) / float64(flags.Height))),
      )
      z := flags.Z

      for i := 0; i < 10; i++ {
        z = z * z + c
      }

      g := float64(cmplx.Abs(z)) / 1.0
      img.Set(x, y, color.RGBA{
        uint8(255.0 * g),
        uint8(255.0 * g),
        uint8(255.0 * g),
        255,
      })
    }
  }

  return img
}

func main() {
  fs := parseFlags()

  img := generateFractal(fs)

  f, _ := os.OpenFile(fs.Output, os.O_WRONLY | os.O_CREATE, 0600)
  defer f.Close()

  png.Encode(f, img)
}
