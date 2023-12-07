# be_healhero

This repository contains the backend (BE) codebase for the HealHero application.

## Recent Updates

### Model
- `type.go`
1. go.mongodb.org/mongo-driver/bson/primitive : berfungsi untuk menyediakan tipe data primitif yang sesuai dengan format BSON (Binary JSON). BSON adalah format serialisasi data biner yang digunakan oleh MongoDB untuk menyimpan dan mentransfer data.
2. Beberapa tipe data primitif yang paling umum digunakan dalam BSON adalah:
- ObjectID: Digunakan sebagai pengenal unik untuk dokumen dalam koleksi MongoDB. ObjectID umumnya digunakan sebagai nilai _id dalam dokumen.
- Timestamp: Merepresentasikan waktu dalam detik sejak waktu tertentu (epoch). Timestamp sering digunakan dalam ObjectID untuk mencatat waktu pembuatan dokumen.
- Binary: Menyimpan data biner, seperti gambar atau dokumen biner lainnya.
- Time : bahasa pemrograman Go (Golang) menyediakan fungsionalitas untuk bekerja dengan waktu dan durasi.

### Module
- `controller.go`
  1. "context"menyediakan mekanisme untuk mengatur dan membatasi waktu, pembatalan, dan nilai-nilai konteks yang dikirimkan antar goroutine.
  2. "crypto/rand"menyediakan fungsi-fungsi untuk menghasilkan angka acak kriptografis. Ini umumnya digunakan untuk menghasilkan kunci atau nilai acak yang digunakan dalam keamanan kriptografi.
  3. "encoding/hex"menyediakan fungsi-fungsi untuk mengonversi data biner menjadi representasi heksadesimal dan sebaliknya.
  4. "errors"menyediakan tipe dan fungsi-fungsi untuk menangani kesalahan atau error dalam program.
  5. "fmt"digunakan untuk formatting dan mencetak output
  6. "os"menyediakan fungsi-fungsi untuk berinteraksi dengan sistem operasi, seperti membaca atau menulis ke file, mengakses variabel lingkungan, dan lainnya.
  7. "strings"menyediakan fungsi-fungsi untuk manipulasi string, seperti penggabungan, pemisahan, dan pencarian substring.
  8. "github.com/HealHeroo/be_healhero/model"paket khusus ( model) dari github.com/HealHeroo/be_healhero repositori. berisi struct atau tipe Go yang mewakili model data yang digunakan dalam aplikasi HealHero.
  9. "github.com/badoux/checkmail"Paket ini menyediakan cara sederhana untuk memeriksa apakah alamat email valid. Ini dapat digunakan untuk melakukan validasi dasar alamat email.
 10. "go.mongodb.org/mongo-driver/bson"Paket ini adalah bagian dari driver resmi MongoDB Go dan menyediakan fungsionalitas untuk menyandikan dan mendekode data BSON (Binary JSON).
 11. "go.mongodb.org/mongo-driver/bson/primitive"Paket ini menyediakan tipe BSON primitif, seperti ObjectID, Timestamp, dan Binary.
 12. "go.mongodb.org/mongo-driver/mongo"Paket ini adalah bagian dari driver resmi MongoDB Go dan menyediakan API tingkat tinggi untuk berinteraksi dengan MongoDB.
 13. "go.mongodb.org/mongo-driver/mongo/options"Paket ini menyediakan opsi untuk mengonfigurasi berbagai aspek driver MongoDB Go, seperti pengaturan koneksi, opsi kueri, dan banyak lagi.
 14. "golang.org/x/crypto/argon2"Paket ini menyediakan implementasi algoritma hashing kata sandi Argon2. Argon2 dirancang untuk melakukan hashing kata sandi dengan aman dan tahan terhadap jenis serangan tertentu.

- `handler.go`
  1. "encoding/json" Paket ini menyediakan fungsionalitas untuk menyandikan dan mendekode data JSON di Go. Ini biasanya digunakan untuk membuat serial struktur data Go ke dalam format JSON dan membatalkan serialisasi data JSON ke dalam struktur data Go.   
  2. "net/http" Paket ini memberikan landasan untuk membangun server dan klien HTTP di Go. Ini mencakup fungsionalitas untuk menangani permintaan HTTP, melayani respons HTTP, dan mengelola koneksi HTTP.
  3. "os" menyediakan cara untuk berinteraksi dengan sistem operasi, termasuk fungsi untuk membaca variabel lingkungan, bekerja dengan file, dan mengelola proses.
  4. "github.com/HealHeroo/be_healhero/model"paket khusus ( model) dari github.com/HealHeroo/be_healhero repositori. berisi struct atau tipe Go yang mewakili model data yang digunakan dalam aplikasi HealHero.
  5. "go.mongodb.org/mongo-driver/bson/primitive" berfungsi untuk menyediakan tipe data primitif yang sesuai dengan format BSON (Binary JSON). BSON adalah format serialisasi data biner yang digunakan oleh MongoDB untuk menyimpan dan mentransfer data.

- `paseto.go`
  1. "encoding/json" Paket ini menyediakan fungsionalitas untuk menyandikan dan mendekode data JSON di Go. Ini biasanya digunakan untuk membuat serial struktur data Go ke dalam format JSON dan membatalkan serialisasi data JSON ke dalam struktur data Go.   
  2. "fmt" digunakan untuk formatting dan mencetak output
  3. "time" bahasa pemrograman Go (Golang) menyediakan fungsionalitas untuk bekerja dengan waktu dan durasi.
  4. "aidanwoods.dev/go-paseto" merupakan alternatif aman untuk JWT (JSON Web Tokens). Ini memungkinkan Anda membuat dan memverifikasi token untuk komunikasi yang aman.
  5. "github.com/HealHeroo/be_healhero/model" paket khusus ( model) dari github.com/HealHeroo/be_healhero repositori. berisi struct atau tipe Go yang mewakili model data yang digunakan dalam aplikasi HealHero.
  6. "go.mongodb.org/mongo-driver/bson/primitive" berfungsi untuk menyediakan tipe data primitif yang sesuai dengan format BSON (Binary JSON). BSON adalah format serialisasi data biner yang digunakan oleh MongoDB untuk menyimpan dan mentransfer data.
  
### Dev Test
1. "crypto/rand" Paket ini menyediakan fungsionalitas untuk menyandikan dan mendekode data JSON di Go. Ini biasanya digunakan untuk membuat serial struktur data Go ke dalam format JSON dan membatalkan serialisasi data JSON ke dalam struktur data Go.   
2. "encoding/hex" menyediakan fungsi-fungsi untuk mengonversi data biner menjadi representasi heksadesimal dan sebaliknya.
3. "fmt" digunakan untuk formatting dan mencetak output
4. "testing" Paket ini merupakan bagian dari perpustakaan standar dan menyediakan dukungan untuk pengujian otomatis di Go
5. "github.com/HealHeroo/be_healhero/model" paket khusus ( model) dari github.com/HealHeroo/be_healhero repositori. berisi struct atau tipe Go yang mewakili model data yang digunakan dalam aplikasi HealHero.
6. "github.com/HealHeroo/be_healhero/module" merupakan paket khusus lain dari github.com/HealHeroo/be_healhero repositori, berisi fungsionalitas modular atau fungsi utilitas yang digunakan dalam aplikasi HealHero.
7. "go.mongodb.org/mongo-driver/bson" Paket ini adalah bagian dari driver resmi MongoDB Go dan menyediakan fungsionalitas untuk menyandikan dan mendekode data BSON (Binary JSON).
8. "go.mongodb.org/mongo-driver/bson/primitive"" berfungsi untuk menyediakan tipe data primitif yang sesuai dengan format BSON (Binary JSON). BSON adalah format serialisasi data biner yang digunakan oleh MongoDB untuk menyimpan dan mentransfer data.
9. "golang.org/x/crypto/argon2" Paket ini menyediakan implementasi algoritma hashing kata sandi Argon2. Argon2 dirancang untuk melakukan hashing kata sandi dengan aman dan tahan terhadap jenis serangan tertentu.
