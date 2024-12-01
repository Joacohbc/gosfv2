package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"gosfV2/src/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/argon2"
)

// Configuración de Argon2
type Argon2Config struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
	SaltLen uint32
}

var argon2Config = Argon2Config{
	Time:    2,
	Memory:  64 * 1024,
	Threads: 4,
	KeyLen:  32,
	SaltLen: 16,
}

// Maneja los errores de los archivos
func HandleUserError(err error) error {
	if err == models.ErrUserNotFound {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}

// Genera una contraseña con Argon2id
func GeneratePassword(password *string) error {
	salt := make([]byte, argon2Config.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("error generando salt: %w", err)
	}

	hash := argon2.IDKey([]byte(*password), salt, argon2Config.Time, argon2Config.Memory, argon2Config.Threads, argon2Config.KeyLen)
	*password = base64.StdEncoding.EncodeToString(append(salt, hash...))

	return nil
}

// Compara el password con el hash de la base de datos
func CheckPassword(password, dbHash string) (bool, error) {
	decodedHash, err := base64.StdEncoding.DecodeString(dbHash)
	if err != nil {
		return false, fmt.Errorf("error decodificando hash: %w", err)
	}

	salt := decodedHash[:argon2Config.SaltLen]
	hashPart := decodedHash[argon2Config.SaltLen:]

	testHash := argon2.IDKey([]byte(password), salt, argon2Config.Time, argon2Config.Memory, argon2Config.Threads, argon2Config.KeyLen)

	return subtle.ConstantTimeCompare(testHash, hashPart) == 1, nil
}
