# JTI/JTN - TUGAS Backend

1. Menggunakan Golang pada Framework FIBER ✓
2. Database menggunakan MYSQL
3. Penginputan Data Bisa Menggunakan Data PLAIN Menggunakan Metode ENCRYPT/DECRYPT ✓
4. Untuk Menampilkan Update Data Di Form Output bisa Menggunakan Metode REFRESH TIMER
Menggunakan WEB SOCKET ✓

# Tutorial Menjalankan Aplikasi Backend
1. Create Database dengan nama `jelajah_teknologi_negeri`
2. ketik perintah `cd application` untuk masuk ke directory aplikasi utama
3. Setelah itu jalankan migration tabel yang sudah dibuatkan schema pada folder `models` aplikasi golang dengan perintah `go run main.go -migrate=migrate`
4. Setelah migrate berhasil kemudian hentikan aplikasi tersebut dengan `ctrl + c` pada windows, `control + c` pada macbook
5. Jalankan seeder untuk create otomatis pada tabel user dan provider dengan perintah `go run main.go -seed=all`
6. Jika sudah berjalan seeder nya bisa dihentikan aplikasi tersebut dengan `ctrl + c` pada windows, `control + c` pada macbook
7. Setelah sudah dilakukan migrate dan seeder, langkah selanjutnya jalankan aplikasi tersebut dengan hot reload menggunakan perintah `air` jika tidak jalankan normal pada umumnya `go run main.go`