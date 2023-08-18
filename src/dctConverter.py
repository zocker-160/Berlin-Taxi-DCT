#! /usr/bin/env python3

import io
import os
import sys

from PIL import Image
import easystruct as es

class DCT:
    magic = b"DC2"
    scale: float = 1.0
    xRes: int
    yRes: int
    bpp: int
    unknown: int = 0
    numResolutions: int = 1

    data: bytes

    @staticmethod
    def parse(f: io.BufferedReader):
        assert f.read(3) == DCT.magic, "unsupported magic"
        assert es.read_float_buff(f) == DCT.scale, "unsupported scale factor"

        xRes = es.read_uint32_buff(f)
        yRes = es.read_uint32_buff(f)

        bpp = es.read_uint8_buff(f)
        assert bpp in (24, 32), "unsupported BPP"

        f.read(1) # we skip 1 unknown byte

        numRes = es.read_uint8_buff(f)
        assert numRes > 0, "number of resolutions is invalid"

        size = calcSize(xRes, yRes, bpp, numRes)
        return DCT(xRes, yRes, bpp, numRes, data=f.read(size))

    @staticmethod
    def parseNoCheck(f: io.BufferedReader):
        magic = f.read(3)
        scale = es.read_float_buff(f)

        xRes = es.read_uint32_buff(f)
        yRes = es.read_uint32_buff(f)

        bpp = es.read_uint8_buff(f)
        unknown = es.read_uint8_buff(f)
        numRes = es.read_uint8_buff(f)

        size = calcSize(xRes, yRes, bpp, numRes)

        dct = DCT(xRes, yRes, bpp, numRes, data=f.read(size))
        dct.magic = magic
        dct.scale = scale
        dct.unknown = unknown

        return dct

    def __init__(self, xRes: int, yRes: int, bpp: int, numRes: int, data: bytes = b''):
        self.xRes = xRes
        self.yRes = yRes
        self.bpp = bpp
        self.numResolutions = numRes
        self.data = data

    def toBytes(self) -> bytes:
        # we only return the first image and do not give a shit about the rest
        dataSize = self.xRes * self.yRes * (self.bpp // 8)
        return self.data[:dataSize]

    def assemble(self, f: io.BufferedWriter) -> bytes:
        f.write(self.magic)
        es.write_float_buff(f, self.scale)
        es.write_uint32_buff(f, self.xRes)
        es.write_uint32_buff(f, self.yRes)
        
        es.write_uint8_buff(f, self.bpp)
        es.write_uint8_buff(f, DCT.unknown)
        es.write_uint8_buff(f, DCT.numResolutions) # we only export one resolution
        f.write(self.toBytes())


def calcSize(xRes: int, yRes: int, bpp: int, numRes: int) -> int:
    size = 0

    for i in range(numRes):
        size += ((xRes * yRes) / 2**(2*i)) * (bpp / 8)

    return int(size)


def dct2png(file: str, name: str, dstFolder: str):
    print(f"converting {name} to png")

    with open(file, "rb") as f:
        dct = DCT.parse(f)

    mode = "RGBA" if dct.bpp == 32 else "RGB"
    image = Image.frombytes(mode, (dct.xRes, dct.yRes), dct.toBytes())

    if dct.bpp == 24:
        b, g, r = image.split()
        image = Image.merge(mode, (r, g, b))
    else:
        b, g, r, a = image.split()
        image = Image.merge(mode, (r, g, b, a))

    image.save(os.path.join(dstFolder, name+".png"))


def png2dct(file: str, name: str, dstFolder: str):
    print(f"converting {name} to dct")

    image = Image.open(file)

    print("IMAGE MODE", image.mode)
    assert image.mode in ("RGB", "RGBA"), "unsupported image mode"
    
    # we need to switch R and B channel when in 24 bit mode
    if image.mode == "RGB":
        bpp = 24

        r, g, b = image.split()
        image = Image.merge("RGB", (b, g, r))

    else:
        bpp = 32

        r, g, b, a = image.split()
        image = Image.merge("RGBA", (b, g, r, a))

    rawData = image.tobytes()

    dct = DCT(image.width, image.height, bpp, 1, rawData)

    loc = os.path.join(dstFolder, name+".dct")
    with open(loc, "wb") as f:
        dct.assemble(f)

def main():
    try:
        sourcefile = sys.argv[1]
        dstFolder = sys.argv[2]
    except IndexError:
        print("usage: dctconverter <imagefile> <destfolder>")
        exit()

    name, ext = os.path.basename(sourcefile).split(os.path.extsep)

    if ext.lower() == "dct":
        dct2png(sourcefile, name, dstFolder)
    elif ext.lower() == "png":
        png2dct(sourcefile, name, dstFolder)
    else:
        raise TypeError(f"unsupported file type ({ext})")

main()
