package helper

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// random string kode
func GenerateRandomString(length int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var sb strings.Builder

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < length; i++ {
		randomIndex := rng.Intn(len(charset))
		sb.WriteByte(charset[randomIndex])
	}

	return sb.String()
}

// Fungsi untuk mendapatkan salam acak
func GetRandomGreeting() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Daftar salam yang dapat dipilih secara acak
	greetings := []string{
		"Halo Wali Murid", "Halo Orang Tua Siswa", "Halo", "Hola",
		"Info Wali Murid", "Hai Wali Murid", "Hai", "Hi",
		"Halo Orang Tua Murid", "Hai Orang Tua Murid", "Hai Orang Tua Siswa",
		"Informasi Absen Murid", "Informasi Absen Siswa",
	}
	return greetings[rng.Intn(len(greetings))]
}

// genearate random text untuk masuk sekolah
func GetRandomMasukMessage(nama, kelas, datetime string) string {
	// Creating a new random generator instance
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	messages := []string{
		fmt.Sprintf("%s kelas %s telah memasuki lingkungan sekolah *(Masuk Sekolah)* pada tanggal %s pukul %s WIB", nama, kelas, FormatTanggalIndo(datetime), time.Now().Format("15:04")),
		fmt.Sprintf("%s kelas %s sudah sampai di sekolah *(Masuk Sekolah)* pada tanggal %s pukul %s WIB", nama, kelas, FormatTanggalIndo(datetime), time.Now().Format("15:04")),
	}
	return messages[rng.Intn(len(messages))]
}

// genearate random text untuk keluar sekolah
func GetRandomKeluarMessage(nama, kelas, datetime string) string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	messages := []string{
		fmt.Sprintf("%s kelas %s telah keluar dari area sekolah *(Pulang Sekolah)* pada tanggal %s pukul %s WIB", nama, kelas, FormatTanggalIndo(datetime), time.Now().Format("15:04")),
		fmt.Sprintf("%s kelas %s sudah keluar dari lingkungan sekolah *(Pulang Sekolah)* pada tanggal %s pukul %s WIB", nama, kelas, FormatTanggalIndo(datetime), time.Now().Format("15:04")),
	}
	return messages[rng.Intn(len(messages))]
}

// Fungsi untuk mendapatkan tanda tangan acak
func GetRandomSignature() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	signatures := []string{
		"Tim Absensi SMAN 5 KARIMUN", "Sekolah SMAN 5 KARIMUN", "SMAN 5 KARIMUN", "Petugas SMAN 5 KARIMUN",
	}

	return signatures[rng.Intn(len(signatures))]
}

// Fungsi untuk mendapatkan emoji acak
func GetRandomEmoji() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	emojis := []string{"😊", "🎉", "📚", "🚌", "🌟"}
	return emojis[rng.Intn(len(emojis))]
}

// Function to generate the "Masuk Sekolah" message
func GenerateMasukMessage(nama, kelas, datetime string) string {
	greeting := GetRandomGreeting()
	message := GetRandomMasukMessage(nama, kelas, datetime)
	emoji := GetRandomEmoji()
	signature := GetRandomSignature()

	return fmt.Sprintf("%s, %s %s\n\nDikirim oleh: %s", greeting, message, emoji, signature)
}

// Function to generate the "Pulang Sekolah" message
func GenerateKeluarMessage(nama, kelas, datetime string) string {
	greeting := GetRandomGreeting()
	message := GetRandomKeluarMessage(nama, kelas, datetime)
	emoji := GetRandomEmoji()
	signature := GetRandomSignature()

	return fmt.Sprintf("%s, %s %s\n\nDikirim oleh: %s", greeting, message, emoji, signature)
}

func FormatTanggalIndo(tanggal string) string {
	bulan := map[int]string{
		1: "Januari", 2: "Februari", 3: "Maret", 4: "April", 5: "Mei", 6: "Juni",
		7: "Juli", 8: "Agustus", 9: "September", 10: "Oktober", 11: "November", 12: "Desember",
	}

	split := strings.Split(tanggal, "-")
	if len(split) != 3 {
		return ""
	}

	day := split[2]
	month := split[1]
	year := split[0]

	monthInt := int(month[0]-'0')*10 + int(month[1]-'0')
	if monthInt < 1 || monthInt > 12 {
		return ""
	}
	return fmt.Sprintf("%s %s %s", day, bulan[monthInt], year)
}
