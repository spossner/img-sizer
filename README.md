# Img-Sizer

[![Version](https://img.shields.io/github/v/release/spossner/img-sizer?include_prereleases&sort=semver)](https://github.com/spossner/img-sizer/releases)
[![Go](https://github.com/spossner/img-sizer/actions/workflows/go.yml/badge.svg)](https://github.com/spossner/img-sizer/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/spossner/img-sizer)](https://goreportcard.com/report/github.com/spossner/img-sizer)
![Docker Build](https://github.com/spossner/img-sizer/actions/workflows/docker-build.yml/badge.svg)


A high-performance image resizing service that fetches images from URL or a specified S3 bucket and processes them according to specified parameters. This helps reducing the internet traffic especially on mobile clients when not the full, original size is needed.

Knowing that there are several tools available out there, this was developed to pull in the original images from S3 directly while matching the existing URL schema, which can not be changed easily because of clients out in the field. 
Therefore Img-Sizer aims to be a slim, laser sharp solution. 
See https://imgproxy.net/ et al if you want a much more sophisticated and feature complete version - and especially a cleaner and more modern API, which is not dictated by 8y old predecessor. 

## Features

- Resize images to specified dimensions
- Support for density scaling (e.g., 2x, 3x for retina displays)
- Configurable allowed dimensions
- Configurable allowed sources with optional S3 bucket mapping
- JPEG output with quality control
- High-quality Lanczos resampling
- S3 integration
- Rate limiting
- Parameter validation
- Modular architecture

## Configuration

The service is configured through json configuration files:

```json
{
     "allowed_sources": [
        {
            "pattern": "example.com",
            "bucket": "main-bucket"
        },
        {
            "pattern": "images.example.com",
            "bucket": "images-bucket"
        },
        {
            "pattern": "*.example.com"
        }
    ],
    "allowed_dimensions": [
        {
            "width": 100,
            "height": 100
        },
        {
            "width": 200,
            "height": 200
        },
        {
            "width": 300,
            "height": 300
        },
        {
            "width": 400,
            "height": 300
        },
        {
            "width": 800,
            "height": 600
        },
        {
            "width": 1024,
            "height": 768
        }
    ],
    "rate_limit": {
        "max_requests": 100,
        "window": "1m"
    }
}
```

The `allowed_sources` configuration maps URL patterns to their corresponding S3 bucket names. This allows you to use production URLs while the service automatically maps them to the correct S3 buckets. The rules are processed and evaluated in the order given in configuration.
If no bucket is specified, the service will fetch the image data from the source URL.

Multiple config files can be provided in ./config folder follwing the pattern `<app-env>.json`. The desired one is chosen by using the APP_ENV environment variable with fallback to local. The value from APP_ENV is used as `<app-env>`.

## Environment Variables

The following environment variables can be used to configure the service:

- `PORT` (optional): The port the server listens on (default: 8080)
- `APP_ENV` (optional): config file selector (e.g. local, test or prod) (default: local)
- `LOG_LEVEL` (optional): Logging level (debug, info, warn, error) (default: info)
- `AWS_REGION` (required): The AWS region for S3 operations (e.g., eu-central-1)
- `AWS_ACCESS_KEY_ID` (required): AWS access key for S3 access
- `AWS_SECRET_ACCESS_KEY` (required): AWS secret key for S3 access

Example `.env` file:
```env
PORT=8080
AWS_REGION=eu-central-1
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
LOG_LEVEL=info
APP_ENV=local
```

## Usage

### HTTP Endpoints

#### 1. Resize Image
Simple resize of a given image. Deprecated. Use /v2/resize.jpg instead.

```
GET /resize.jpg?src=<url>&width=<width>&height=<height>&density=<density>&quality=<quality>
```

Parameters:
- `src`: URL of the source image (required). The URL must match one of the allowed patterns in the configuration.
- `width`: Target width
- `height`: Target height
- `density`: Scale factor for retina displays  (defaults to 1.0)
- `scale`: Alternative to density you can use scale - note that if scale is given, it overwrites any value specified in density.
- `quality`: Compression quality 1-100, passed directly to the JPEG compressor (defaults to 70)
- `background`: Color HEX to use as a background for flattening transparent images (PNG, GIF, etc.) (defaults to 000000)

Example:
```
/resize.jpg?src=https://images.example.com/photo.jpg&width=200
/resize.jpg?src=https://images.example.com/photo.jpg&width=300&height=200
/resize.jpg?src=https://images.example.com/photo.jpg&width=800&height=600&density=2&quality=90
```

#### 2. Crop Image
Simple cropping of a given image. Deprecated. User /v2/resize.jpg instead.

```
GET /crop.jpg?src=<url>&width=<width>&height=<height>&crop=<crop>&density=<density>&quality=<quality>
```

Parameters:
- `src`: URL of the source image (required)
- `width`: Crop width
- `height`: Crop height
- `x`: The left edge of the part to crop in the original image
- `y`: The upper edge of the part to crop in the original image
- `scale`: Scaling of original image
- `density`: Scale factor of the croppped image for retina displays (defaults to 1.0)
- `quality`: Compression quality 1-100, passed directly to the JPEG compressor (defaults to 70)
- `background`: Color HEX to use as a background for flattening transparent images (PNG, GIF, etc.) (defaults to 000000)

Example:
```
/crop.jpg?src=https://images.example.com/photo.jpg&height=100&scale=0.0626666667&width=100&x=10&y=4
/crop.jpg?src=https://images.example.com/photo.png&height=100&scale=0.0626666667&width=100&x=10&y=4&density=2.0&background=ffffff
```

#### 3. Combined cropping and resizing
```
GET /v2/resize.jpg?src=<url>&width=<width>&height=<height>&density=<density>&quality=<quality>
```

Parameters:
- `src`: URL of the source image (required)
- `width`: Target width 
- `height`: Target height
- `density`: Scale factor for retina displays (defaults to 1.0)
- `scale`: Alternative to density you can use scale - note that if scale is given, it overwrites any value specified in density.
- `quality`: Compression quality 1-100, passed directly to the JPEG compressor (defaults to 70)
- `background`: Color HEX to use as a background for flattening transparent images (PNG, GIF, etc.) (defaults to 000000)
- `crop[x]`: left offset of the crop zone (defaults to 0)
- `crop[y]`: top offset of the crop zone (defaults to 0)
- `crop[width]`: Width of the crop zone (*)
- `crop[height]`: Height of the crop zone (*)
- `crop[scale]`: Use crop zone at this scale (based on the original image size)

(*) If cropping should take place, both width and height must be specified and both must be positive - otherwise no croppping takes place

Example:
```
/v2/resize.jpg?width=570&height=320&density=1.2&src=https://images.example.com/photo.jpg
/v2/resize.jpg?crop[x]=0&crop[y]=163&crop[scale]=0.25&crop[width]=270&crop[height]=200&width=260&height=154&density=2&src=https://images.example.com/photo.jpg
```

### Docker

Build the image:
```bash
docker build -t img-sizer .
```

Run the container:
```bash
docker run -p 8080:8080 \
  -e AWS_ACCESS_KEY_ID=your_access_key \
  -e AWS_SECRET_ACCESS_KEY=your_secret_key \
  -e AWS_REGION=your_region \
  -e APP_ENV=test \
  img-sizer
```

### Make Commands

The project includes a Makefile with several useful commands:

- `make build`: Build the application
- `make dev`: Start development server with air
- `make docker-build`: Build Docker image
- `make docker-run`: Run Docker container locally with .env file
- `make docker-push`: Push Docker image to AWS ECR
- `make clean`: Clean up build artifacts
- `make create-ecr`: Create ECR repository if it doesn't exist
- `make deploy`: Deploy to AWS (build, push, create ECR)

Example usage:
```bash
# Build and run locally
make build
./bin/img-sizer

# Development with hot reloading
make dev

# Docker operations
make docker-build
make docker-run

# AWS deployment
make deploy
```

## Development

### Prerequisites

- Go 1.24.1 or later
- Docker (optional)
- Air (for hot reloading)

### Project Structure

```
.
├── cmd/
│   └── img-sizer/         # Main application entry point
├── config/                # Configuration files for different environments
│   ├── local.json         # Local development configuration
│   ├── test.json          # Test environment configuration
│   └── prod.json          # Production environment configuration
├── internal/
│   ├── config/            # Configuration package
│   ├── handlers/          # HTTP handlers
│   │   ├── handlers.go    # Main handlers
│   │   ├── helpers/       # Helper functions - especially S3 URL parsing
│   ├── processing/        # Image processing package
│   └── storage/           # Storage package to access S3
├── Dockerfile             # Docker build file
├── go.mod                 # Go module file
├── go.sum                 # Go module checksums
└── README.md              # This file
```

### Local Development

1. Clone the repository
2. Copy `.env.example` to `.env` and configure your AWS credentials
3. Run the service:
   ```bash
   make dev
   ```
   Or in docker container:
   ```bash
   make docker-run
   ```

### Testing

Run tests:
```bash
go test ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

The MIT License is a permissive license that is short and to the point. It lets people do anything they want with your code as long as they provide attribution back to you and don't hold you liable.
