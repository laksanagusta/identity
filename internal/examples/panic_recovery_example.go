package examples

import "github.com/gofiber/fiber/v2"

// Contoh handler yang bisa panic - JANGAN KHAWATIR!
// Recovery middleware akan menangkapnya dan aplikasi tetap berjalan

// ❌ BAD EXAMPLE - Handler yang bisa panic tanpa protection
func BadHandlerExample(c *fiber.Ctx) error {
	users := []string{"Alice", "Bob"}
	// Ini akan panic jika index out of range!
	name := users[5] // Index out of range!
	return c.JSON(fiber.Map{"name": name})
}

// ✅ GOOD EXAMPLE 1 - Handler dengan proper error checking
func GoodHandlerExample1(c *fiber.Ctx) error {
	users := []string{"Alice", "Bob"}
	index := 5

	if index >= len(users) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Index out of range",
		})
	}

	name := users[index]
	return c.JSON(fiber.Map{"name": name})
}

// ✅ GOOD EXAMPLE 2 - Handler dengan defer recover manual (opsional)
// Note: Ini opsional karena middleware global sudah menangani panic
func GoodHandlerExample2(c *fiber.Ctx) (err error) {
	// Defer recover ini opsional, middleware global sudah menangani
	// tapi bisa digunakan untuk custom error handling
	defer func() {
		if r := recover(); r != nil {
			err = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":  "Custom error handling",
				"detail": r,
			})
		}
	}()

	users := []string{"Alice", "Bob"}
	name := users[5] // Akan panic tapi di-catch oleh defer recover
	return c.JSON(fiber.Map{"name": name})
}

// ✅ BEST PRACTICE - Validasi input dan error handling yang jelas
func BestPracticeHandler(c *fiber.Ctx) error {
	// Parse request
	type Request struct {
		Index int `json:"index"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	users := []string{"Alice", "Bob", "Charlie"}

	// Validate index
	if req.Index < 0 || req.Index >= len(users) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":     "Index out of range",
			"max_index": len(users) - 1,
		})
	}

	// Safe access
	name := users[req.Index]
	return c.JSON(fiber.Map{
		"success": true,
		"name":    name,
	})
}

// 💡 CATATAN PENTING:
//
// 1. Recovery middleware SUDAH AKTIF secara global
//    - Semua panic akan tertangkap otomatis
//    - Aplikasi tidak akan crash
//    - Error akan di-log dengan detail
//
// 2. Namun, LEBIH BAIK menggunakan proper error handling:
//    - Validasi input
//    - Check boundary
//    - Return error secara eksplisit
//    - Panic hanya untuk kondisi truly exceptional
//
// 3. Recovery middleware adalah SAFETY NET, bukan primary error handling
//
// 4. Monitor logs untuk panic yang terjadi - ini indikasi bug!
