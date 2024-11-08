package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gopkg.in/gomail.v2"
)

var db *sql.DB

func main() {
	var err error
	// Connect to PostgreSQL
	connStr := "user=postgres dbname=db_emails sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()
	r.POST("/send-email", sendEmail)
	r.Run(":9000") // Run on port 9000
}

func sendEmail(c *gin.Context) {
	var json struct {
		Batch   int64  `json:"batch"` // This will be used as a limit
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

	subject := "Penawaran Harga Cable Tray dan Ladder, PT Arbrion Asia"

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
					<img alt="Your Company Logo" width="307" height="50" src="https://static.pbahotels.com/Assets/images/Hotel/exterior/a54135947cedd3fd5597ccfe82ab3c3ab0be1c4a.png" />
				</div>
				<div class="container">
					<div class="guestname">
						Kepada<br />
						Yth. Bapak/Ibu<br />
						Di tempat ,
					</div>
					<div class="promotion-text">
						Kami adalah perusahaan yang bergerak di bidang yang menyediakan berbagai macam jenis cable tray dan cable ladder. Dengan berbagai pengalaman dan kepercayaan pelanggan sejak tahun 2007. PT Arbrion Asia akan membantu kebutuhan jasa pembangunan infrastruktur anda.<br /><br />

						Untuk lebih jelasnya, silahkan klik link berikut :<br /><br />
						<a href="https://arbrion-asia.com/">Lihat Profil Perusahaan Kami</a><br /><br />

						Dengan ini kami ingin mengajukan penawaran harga untuk produk Cable Tray dan Ladder. Berikut terlampir katalog dan harga produk kami. Semoga penawaran ini dapat memenuhi kebutuhan anda.<br /><br />

						Apabila ada kelanjutan dari penawaran ini, jangan sungkan untuk menghubungi kami serta kami akan dengan senang hati melakukan penyesuaian harga dan spesifikasi produk. Kami juga siap melakukan kunjungan langsung ke lokasi anda untuk mendiskusikan lebih lanjut mengenai kebutuhan anda.<br /><br />

						Hormat Kami,<br />
						<b>Indria Wigati</b><br />
						<b>PT Arbrion Asia</b><br /><br />
						<a href="https://wa.me/6289657664445" style="color: #04a9f5; text-decoration: underline;">+62 896 5766 4445</a><br /><br />
						<img alt="Arbrion Asia" width="500" height="150" src="https://arbrion-asia.com/" />
					</div>
				</div>
			</div>
		</body>
		</html>
	`

	// Query to get data based on batch limit where has_sent is false
	rows, err := db.Query("SELECT email, id FROM emails WHERE has_sent = false LIMIT $1", json.Batch)
	if err != nil {
		log.Printf("Failed to query database: %v", err) // Log the actual error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
		return
	}
	defer rows.Close() // Ensure rows are closed after use

	var wg sync.WaitGroup
	count := 0
	successfulEmails := []string{} // Slice to hold successfully sent emails

	for rows.Next() {
		var email string
		var id int64 // Assuming you have an ID column to identify the email record
		if err := rows.Scan(&email, &id); err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}

		wg.Add(1)
		go func(email string, id int64) {
			defer wg.Done()

			m := gomail.NewMessage()
			m.SetHeader("From", "indriarbrion@yahoo.com")
			m.SetHeader("To", email)
			m.SetHeader("Subject", subject)
			m.SetBody("text/html", emailBody)

			filePath1 := "Katalog PT.Arbrion Asia.pdf" // First attachment
			// filePath2 := "Price List.pdf"              // Second attachment

			// Check if the first file exists
			if _, err := os.Stat(filePath1); os.IsNotExist(err) {
				log.Printf("Failed to attach file: %v", err)
				return
			}
			m.Attach(filePath1) // Attach the first PDF file

			// Check if the second file exists
			// if _, err := os.Stat(filePath2); os.IsNotExist(err) {
			// 	log.Printf("Failed to attach file: %v", err)
			// 	return
			// }
			// m.Attach(filePath2) // Attach the second PDF file

			d := gomail.NewDialer("smtp.mail.yahoo.com", 587, "indriarbrion@yahoo.com", "ifulmzvmmcrmnriq")

			// Send the email
			if err := d.DialAndSend(m); err != nil {
				log.Printf("Failed to send email to %s: %v", email, err)
				return // Do not update the database if sending fails
			}

			// Update has_sent to true and set updated_at to the current timestamp in the database
			_, err := db.Exec("UPDATE emails SET has_sent = true, updated_at = NOW() WHERE id = $1", id)
			if err != nil {
				log.Printf("Failed to update has_sent for email %s: %v", email, err)
				return
			}

			// Add the successfully sent email to the slice
			successfulEmails = append(successfulEmails, email)
			count++
		}(email, id)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error occurred during rows iteration: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred during rows iteration"})
		return
	}

	wg.Wait() // Wait for all goroutines to finish

	// Return the status and the list of successfully sent emails
	c.JSON(http.StatusOK, gin.H{"status": "Emails sent successfully", "count": count, "sentEmails": successfulEmails})
}

// <img alt="Your Company Logo" width="307" height="50" src="https://static.pbahotels.com/Assets/images/Hotel/exterior/c3628711b086cca959673fbf01d201a2c661583f.png" />
