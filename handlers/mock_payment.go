package handlers

import (
	"strings"

	"github.com/ehanz12/api-SneakHub/services"
	"github.com/gofiber/fiber/v2"
)

func MockPayPageHandler(c *fiber.Ctx) error {
	orderID := c.Params("order_id")
	amount := c.Query("amount", "0")
	name := c.Query("name", "Pelanggan")

	html := `<!DOCTYPE html>
<html lang="id">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Mock Payment - SneakHub</title>
<style>
body { font-family: system-ui, sans-serif; background: #f5f5f5; display: flex; justify-content: center; align-items: center; min-height: 100vh; margin: 0; }
.card { background: #fff; border-radius: 12px; padding: 32px; width: 100%%; max-width: 420px; box-shadow: 0 4px 20px rgba(0,0,0,.08); text-align: center; }
h2 { margin-top: 0; color: #111; }
.order-id { font-size: 12px; color: #888; word-break: break-all; }
.amount { font-size: 32px; font-weight: 700; color: #16a34a; margin: 16px 0; }
.badge { display: inline-block; background: #fff4e5; color: #b45309; font-size: 12px; font-weight: 600; padding: 4px 12px; border-radius: 999px; margin-bottom: 16px; }
button { width: 100%%; padding: 14px; border: 0; border-radius: 8px; background: #16a34a; color: #fff; font-size: 16px; font-weight: 600; cursor: pointer; }
button:disabled { opacity: .6; cursor: not-allowed; }
button.reset { background: #ef4444; margin-top: 8px; }
#result { margin-top: 16px; font-size: 14px; display: none; }
#result.ok { color: #16a34a; display: block; }
#result.err { color: #ef4444; display: block; }
</style>
</head>
<body>
<div class="card">
<h2>Halaman Bayar (Mock)</h2>
<span class="badge">MODE SIMULASI</span>
<div class="order-id">Order: ` + orderID + `</div>
<div class="amount">Rp ` + amount + `</div>
<p>Penerima: ` + name + `</p>
<button id="payBtn">Bayar Sekarang</button>
<button id="failBtn" class="reset">Simulasikan Gagal</button>
<div id="result"></div>
</div>
<script>
const orderID = "` + orderID + `";
async function settle(status) {
  const btn = document.getElementById("payBtn");
  const failBtn = document.getElementById("failBtn");
  const result = document.getElementById("result");
  btn.disabled = true; failBtn.disabled = true;
  result.className = ""; result.style.display = "none";
  try {
    const res = await fetch("/api/payments/mock-settle", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ order_id: orderID, status: status })
    });
    const data = await res.json();
    result.className = res.ok ? "ok" : "err";
    result.textContent = data.message || (res.ok ? "Pembayaran berhasil" : "Terjadi kesalahan");
  } catch (e) {
    result.className = "err";
    result.textContent = "Gagal menghubungi server";
  }
}
document.getElementById("payBtn").onclick = () => settle("PAID");
document.getElementById("failBtn").onclick = () => settle("FAILED");
</script>
</body>
</html>`

	c.Set(fiber.HeaderContentType, "text/html; charset=utf-8")
	return c.SendString(html)
}

func MockSettleHandler(c *fiber.Ctx) error {
	var body struct {
		OrderID string `json:"order_id"`
		Status  string `json:"status"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "request gagal"})
	}
	if strings.TrimSpace(body.OrderID) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "order_id wajib diisi"})
	}

	status := strings.ToUpper(strings.TrimSpace(body.Status))
	switch status {
	case "PAID", "EXPIRED", "FAILED", "REFUND":
	default:
		status = "PAID"
	}

	if err := services.HandlePaymentNotificationService(body.OrderID, status, "MOCK-TXN-"+body.OrderID, 0); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Pembayaran " + status + " diproses"})
}
