🚀 Vastar Enterprise API Gateway & Benchmark Tool
📋 Overview
Project ini adalah sebuah Developer-Grade Connector yang berfungsi sebagai "Pintu Masuk" (Gateway) untuk menguji performa Vastar Workflow Runtime . Berbeda dengan bot biasa, alat ini dirancang untuk mensimulasikan beban kerja nyata (real-world workload) guna membandingkan efisiensi jalur data Vastar melawan kompetitor seperti n8n atau Temporal .
+1.

Sesuai dengan instruksi keamanan sistem, connector ini bertindak sebagai Adapter Terkelola yang mengisolasi setiap permintaan di dalam lingkungan yang aman.

🎯 Tujuan Utama

Benchmarking: Mengukur latensi komunikasi IPC (Inter-Process Communication) menggunakan protokol FlatBuffers.
+1

Keamanan & Isolasi: Membuktikan bahwa setiap tugas dikelola dalam Tenant yang terisolasi sesuai standar arsitektur WASM perusahaan.


Observability: Menyediakan data angka (metrics) secara real-time yang bisa dibaca oleh sistem monitoring seperti Prometheus atau Grafana.
+1

⚙️ Cara Instalasi (Tutorial)
Siapkan Lingkungan: Pastikan kamu sudah berada di folder project di terminal WSL :

Bash

cd ~/vastar-wf-connector-sdk-bin/examples-golang/api-gateway-benchmark
Inisialisasi Modul: Jalankan perintah ini untuk menyambungkan SDK Vastar :

Bash

go mod init api-gateway-benchmark
go mod edit -replace github.com/fullstack-aidev/vastar-wf-connector-sdk-bin/sdk-golang=../../sdk-golang
go mod tidy
Nyalakan Vastar Runtime: Buka terminal baru, masuk ke folder utama, dan jalankan runtime :

Bash

./start_runtime.sh
🚀 Cara Penggunaan
Jalankan Gateway: Di terminal utama, ketik:

Bash

go run main.go
Status: Program sekarang menunggu kiriman data di Port 3000.

Kirim Data (Simulasi Tester): Buka terminal lain, dan gunakan perintah curl untuk mengirim data JSON :

Bash

curl -X POST http://localhost:3000/process \
-H "Content-Type: application/json" \
-d '{"nama": "Rizka", "tugas": "Benchmark Vastar vs n8n"}'
📊 Cara Melihat Hasil Benchmark
Untuk melihat performa sistem secara mendalam, buka browser kamu (Chrome/Edge) dan akses alamat berikut:

Papan Skor Digital: http://localhost:3000/metrics

Apa yang harus diperhatikan?


vastar_processed_tasks_total: Jumlah total barang/data yang berhasil lewat.

vastar_ipc_latency_ms: Kecepatan rata-rata. Jika angkanya di bawah 0.1ms, berarti Vastar sukses mengalahkan n8n dalam hal kecepatan jalur data.
+1

🛡️ Mengapa Ini Aman? (Pesan untuk Mentor)
No Direct Connection: Connector ini patuh pada aturan keamanan; tidak melakukan koneksi internet langsung, melainkan melalui vastar.ExecuteHTTP.

Tenant Isolation: Menggunakan TenantID unik untuk setiap koneksi guna mencegah kebocoran data antar pengguna.


Traceability: Setiap proses dilengkapi dengan Trace-ID untuk memudahkan pelacakan jika terjadi error pada sistem .

🛠️ Troubleshooting

Error: Connection Refused: Pastikan ./start_runtime.sh sudah dijalankan di terminal sebelah .


Error: Module Not Found: Pastikan kamu sudah menjalankan go mod tidy di folder yang benar .