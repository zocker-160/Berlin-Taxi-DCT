# DCT image file converter for Berlin Taxi / Taxi Challenge Berlin

Script to convert DCT texture files used in Berlin Taxi / Taxi Challenge Berlin.

## Usage
```bash
dctConverter.py <image> <destination folder>
```
**NOTE:** only the conversion dct <-> png is supported currently

## DCT File Format

```c
char[3] magic = "DC2"
float scale // unsure but is always 1.0
u32 xResolution
u32 yResolution
u8 bpp // bit per pixel, should be 24 or 32
u8 unknown // no idea what this is
u8 numResolutions // [1]
u8[] data // [2]
```
[1]:
One DCT file can contain multiple resolutions of the same texture used for the different texture quality settings in the game.

[2]:
Texture data encoded RAW either in `BGR888` (24bit) or `BGRA8888` (32bit) format.
The length of the data can be calculated using the following formula:

$$ dataLenth = \sum_{k=0}^{numResolutions-1} { xRes * yRes * bpp \over 2^{2k} * 8 } $$
