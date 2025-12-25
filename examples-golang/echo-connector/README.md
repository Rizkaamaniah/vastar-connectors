# Echo Connector

Echo Connector adalah contoh **HTTP generic connector** yang dibuat menggunakan  
**Vastar Workflow Connector SDK (Go)**.

Connector ini mengirim payload JSON ke public Echo API
(`https://postman-echo.com/post`) dan menerima kembali payload yang sama sebagai response.

---

## 🎯 Tujuan
Connector ini dibuat untuk:
- Menunjukkan penggunaan Vastar WF Connector SDK selain LLM
- Memberikan contoh connector REST API yang sederhana
- Membuktikan bahwa komunikasi eksternal dilakukan melalui Connector Runtime

---

## 🛠️ Cara Kerja Singkat
1. Connector membangun request HTTP menggunakan SDK Vastar
2. Request dikirim melalui **Vastar Connector Runtime**
3. Runtime meneruskan request ke Echo API
4. Response JSON dikembalikan ke connector

Connector **tidak melakukan HTTP request langsung** menggunakan `net/http`.

---

## ▶️ Cara Menjalankan

```bash
go run main.go
