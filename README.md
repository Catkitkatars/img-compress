# ImgCompress
An application written in Go for compressing images, converting them to JPEG, and adding a watermark.

## Fast Start

```bash

git clone https://github.com/Catkitkatars/img-compress.git
docker compose up -d --build

POST:
  - http://localhost:8087/img - can send 1 or more images
GET:
  - http://localhost:8087/img/1
```
Rules of compression and watermarking in /config/config.yaml
