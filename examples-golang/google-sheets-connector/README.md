

📊 Vastar Google Sheets Automation Connector
Konektor performa tinggi berbasis Go yang memungkinkan aplikasi untuk melakukan operasi baca/tulis data secara otomatis ke Google Sheets melalui infrastruktur Vastar. Proyek ini menggunakan arsitektur Enterprise IPC untuk menjamin keamanan dan kecepatan transmisi data.

📑 Ringkasan Proyek
Versi: 1.0.0-Stable

Teknologi: Google Sheets API v4

Transport: Unix Domain Sockets (UDS)

Protokol: FlatBuffers (Biner)

Metode: Append Row (Otomatis)

🏗️ Arsitektur Integrasi
Konektor ini tidak terhubung langsung ke internet, melainkan melalui Vastar Runtime Gateway untuk memastikan setiap paket data melewati proses audit internal:

[App Go] <--- IPC (Unix Socket) ---> [Vastar Runtime] <--- HTTPS (OAuth2) ---> [Google Sheets API]

🚀 Panduan Instalasi
1. Persiapan Google Cloud
Aktifkan Google Sheets API di Google Cloud Console.

Buat Service Account dan unduh file kredensial JSON.

Share file Google Sheets target kepada email Service Account sebagai Editor.

2. Konfigurasi Environment
Atur variabel lingkungan di terminal WSL Anda untuk mengaktifkan koneksi:

Bash

export GOOGLE_API_KEY='AIzaSyCc6hXgSeJI0Ou-UlZnlg5wMzBPdjUB-aQ'
export SHEETS_ID='your-spreadsheet-id-here'
3. Inisialisasi Modul
Jalankan perintah berikut untuk mengintegrasikan Vastar SDK secara lokal:

Bash

go mod init google-sheets-connector
go mod edit -replace github.com/fullstack-aidev/vastar-wf-connector-sdk-bin/sdk-golang=../../sdk-golang
go mod tidy
4. Menjalankan Konektor
Bash

go run main.go
🛠️ Fitur Enterprise
Low-Latency Logging: Menulis data ribuan karakter dengan overhead milidetik berkat jalur biner FlatBuffers.

Automated Handshake: Sistem secara otomatis memverifikasi koneksi ke /tmp/vastar-connector-runtime.sock sebelum memulai operasi penulisan.

Unexported Logic: Implementasi kode menggunakan standar enkapsulasi Go (huruf kecil) untuk menjaga integritas logika internal.

📊 Contoh Output Audit
Plaintext

================================================================================
              VASTAR SDK - GOOGLE SHEETS AUTOMATION ENGINE
================================================================================
🐧 Connected via Unix Socket: /tmp/vastar-connector-runtime.sock
✅ Data berhasil ditulis ke Sheets (Latensi: 450ms)
--------------------------------------------------------------------------------
✅ Workflow Task Completed.
================================================================================

Project: Enterprise AI & Data Workflow Integration with Vastar SDK