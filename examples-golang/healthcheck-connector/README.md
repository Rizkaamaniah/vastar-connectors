# Healthcheck Connector

Healthcheck Connector adalah Workflow Connector utilitas
yang digunakan untuk mengecek status dan kesiapan sistem
sebelum atau selama eksekusi workflow.

Connector ini tidak melakukan komunikasi keluar dan
mengembalikan response status secara lokal.

---

## 🎯 Tujuan
- Menyediakan connector monitoring sederhana
- Memvalidasi kesiapan workflow
- Digunakan sebagai langkah awal (pre-check) dalam workflow
- Mengikuti arsitektur Connector Runtime Vastar

---

## 🧠 Konsep Arsitektur

- Workflow berjalan dalam lingkungan WASM (isolated)
- Workflow memanggil connector melalui Connector Runtime
- Connector mengembalikan status sistem secara lokal

---

## 🛠️ Response Contoh

```json
{
  "status": "ok",
  "service": "healthcheck-connector",
  "version": "v1.0.0",
  "timestamp": "2025-01-22T10:00:00Z"
}
