<img src="gowhatsbot.png" width="30%"/>

# GoWhatsBot : Apa itu ?
GoWhatsBot adalah Bot WhatsApp yang dibangun dengan Go-lang bebasiskan library [` whatsmeow `](https://github.com/tulir/whatsmeow).


# Setup

## Konfigurasi
Untuk menjalankan Bot, kita perlu untuk mengatur konfigurasi database pada berkas ` gowhatsbot.json `. Jika tidak terdapat pada direktori repo, maka kita bisa membuatnya dengan contoh isi konfigurasi sebagai berikut :
``` json
{
    "driver": "sqlite3", // nama driver database yang digunakan
    "sqlite3": "file:whatsapp.db?_foreign_keys=on", // alamat database
    "pgx": "postgres://user:pass@localhost:5432/wadb" // alamat database
}
```
Pada contoh diatas, driver yang akan di gunakan adalah ` sqlite3 ` dengan alamat ` file:whatsapp.db?__foreign_keys=on `.

Secara default ada 2 library driver database yang tesedia yaitu ` pgx ` dan ` go-sqlite3 `, perlu untuk menambahkan baris kode jika ingin menambahkan dukungan layanan database lainnya.

## Autoload
GoWhatsBot secara otomatis akan membuat berkas ` autoload.go ` untuk mendaftarkan packages pada `./plugins`. Daftar autoload tidak akan dibuat ulang jika tidak ada penambahan/penghapusan package pada direktori `./plugins`. Atau bisa di tambahkan secara manual dengan format standar :
```go
package main

import (
    _ "main/plugins/<nama_package>"
)

```

## Menjalankan & Kompilasi

### Menjalankan
Untuk menjalankan bot tanpa kompilasi cukup untuk menjalankan perintah :
```sh
go run .
```

### Kompilasi
Untuk kompilasi agar dapat mendukung driver database maupun library yang menggunakan sumber program dari bahasa C :
```sh
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -v -ldflags "-s -w" -o ./linux-x64

```

Selebihnya dapat dilihat pada berkas ` build `.

# Kontribusi ?
Silahkan jika ingin melakukan kontribusi dengan membuka issue, pull request maupun diskusi.

# Donasi ?
Silahkan untuk melakukan donasi jika berkenan.
- Saweria : [Ma'mur Rezeki](https://saweria.co/mamurrezeki)



# Library ?
- [whatsmeow](https://go.mau.fi/whatsmeow)
- [qrterminal](https://github.com/mdp/qrterminal)
- [pgx](https://github.com/jackc/pgx)
- [go-sqlite3](https://github.com/mattn/go-sqlite3)
- [barcode](https://github.com/boombuler/barcode)
- [go-qrcode](https://github.com/skip2/go-qrcode)
- [imageorient](https://github.com/disintegration/imageorient)
- [imaging](https://github.com/disintegration/imaging)
- [webp](https://github.com/chai2010/webp)
- [gozxing](https://github.com/makiuchi-d/gozxing)
