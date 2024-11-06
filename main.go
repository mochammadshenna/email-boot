package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gopkg.in/gomail.v2"
)

// var db *sql.DB

func main() {
	// var err error
	// Connect to PostgreSQL
	// connStr := "user=yourusername dbname=email_boot sslmode=disable"
	// db, err = sql.Open("postgres", connStr)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	r := gin.Default()

	r.POST("/send-email", sendEmail)

	r.Run(":9000") // Run on port 8080
}

func sendEmail(c *gin.Context) {
	var json struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
		Name    string `json:"name"`
		UserID  string `json:"user_id"`
		Email   string `json:"email"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subject := "Penawaran Harga Cable Tray and Ladder, PT Arbrion Asia"

	// Create a template for the email body
	emailBody := `
		<!doctype html>
		<html amp4email>
		<head>
			<meta charset="utf-8">
			<script async src="https://cdn.ampproject.org/v0.js"></script>
			<script async custom-element="amp-form" src="https://cdn.ampproject.org/v0/amp-form-0.1.js"></script>
			<style amp-custom>
				body {
					color: #000;
					font-size: 13px;
					font-family: Helvetica, Arial, sans-serif;
					background: #fafafa;
				}
				a {
					color: #00205b;
					text-decoration: none;
				}
				.wrapper {
					max-width: 640px;
					margin: 0 auto;
					background: #fff;
					padding: 10px;
					margin-top: 20px;
					margin-bottom: 20px;
					border: solid 1px #d1d1d1;
				}
				.header-logo {
					text-align: center;
					border-bottom: solid 2px #ccccd1;
					padding: 10px;
				}
				.header-logo img {
					margin: 0 auto;
				}
				.container .guestname {
					text-align: left;
					margin: 15px 5px;
				}
				.promotion-text {
					text-align: left;
					margin: 5px 5px;
				}
				.promotion-text a {
					color: #04a9f5;
					font-family: Helvetica, Arial, sans-serif;
					font-size: 14px;
					text-decoration: underline;
					font-weight: bold;
				}
			</style>
		</head>
		<body>
			<div class="wrapper">
				<div class="header-logo">
					<img alt="Your Company Logo" width="307" height="50" src="https://static.pbahotels.com/Assets/images/Hotel/exterior/c3628711b086cca959673fbf01d201a2c661583f.png" />
				</div>
				<div class="container">
					<div class="guestname">
						Kepada<br />
						Yth. Bapak/Ibu<br />
						Di tempat ,
					</div>
					<div class="promotion-text">
						Kami adalah perusahaan yang bergerak di bidang yang menyediakan berbagai macam jenis cable tray dan cable ladder. Dengan berbagai pengalaman dan kepercayaan pelanggan sejak tahun 2007. PT Arbrion Asia akan membantu kebutuhan pembangunan infrastruktur untuk mendukung kebutuhan pengguna jasa.<br /><br />

						Untuk lebih jelasnya, silahkan klik link berikut :<br /><br />
						<a href="https://arbrion-asia.com/">Lihat Profil Perusahaan Kami</a><br /><br />

						Dengan ini kami sampaikan bahwa kami ingin mengajukan penawaran harga untuk produk Cable Tray and Ladder<br /><br />

						Hormat Kami,<br />
						<b>Indria Wigati</b><br />
						<b>PT Arbrion Asia</b><br /><br />
						<p><a href="https://wa.me/6289657664445" style="color: #04a9f5; text-decoration: underline;">+62 896 5766 4445</a></p>
						<img alt="Arbrion Asia" width="500" height="150" src="https://arbrion-asia.com/" />
					</div>
				</div>
			</div>
		</body>
		</html>
	`

	// Send email using Yahoo SMTP
	m := gomail.NewMessage()
	m.SetHeader("From", "indriarbrion@yahoo.com")
	m.SetHeader("To", json.To)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", emailBody)

	d := gomail.NewDialer("smtp.mail.yahoo.com", 587, "indriarbrion@yahoo.com", "ifulmzvmmcrmnriq")

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Email sent successfully"})
}
