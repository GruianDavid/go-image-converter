# Image Converter

A simple Go application to convert images in the `original/` directory to other formats.

## Features
- Converts images from JPG, JPEG, or PNG to WebP format
- Processes all images in the specified input directory
- Outputs three sizes for each image: large, medium, and small
- Multi-threaded processing for speed

## Usage

1. **Place images** you want to convert in the `original/` directory.
2. **Build the application:**
   ```sh
   go build -o image-converter main.go
   ```
3. **Run the application:**
   ```sh
   go run main.go -threads=6 -input=./original -output=./
   # or if built:
   ./image-converter -threads=6 -input=./original -output=./
   # on Windows:
   image-converter.exe -threads=6 -input=./original -output=./
   ```
   or on Windows:
   ```sh
   image-converter.exe
   ```
4. **Check the output** directory (if implemented) for converted images.

## Requirements
- Go 1.18 or newer
- [github.com/sunshineplan/imgconv](https://github.com/sunshineplan/imgconv) (install with `go get` if needed)

## Customization
- Modify `main.go` to change input/output formats or add processing logic.

## License
MIT License
