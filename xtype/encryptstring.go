package xtype

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"github.com/gotidy/ptr"
	"io"
	"strings"
)

var EncryptStringSecret = []string{"20dd6fbb502a0465e070d3ed4b92a84e"}

type EncryptString struct {
	encrypt *string
	str     *string
}

func (es *EncryptString) GormDataType() string {
	return "varchar(256)"
}

func (es *EncryptString) String() string {
	if es.encrypt == nil && es.str == nil {
		return ""
	}
	if es.str == nil && es.encrypt != nil {
		str := dynamicDecrypt(*es.encrypt)
		es.str = &str
	}
	return *es.str
}

func (es *EncryptString) EncryptString() string {
	if es.encrypt == nil && es.str == nil {
		return ""
	}
	if es.encrypt == nil && es.str != nil {
		enc := dynamicEncrypt(*es.str)
		es.encrypt = &enc
	}
	return *es.encrypt
}

func (es *EncryptString) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", es.String())), nil
}

func (es *EncryptString) UnmarshalJSON(data []byte) error {
	raw := strings.TrimSuffix(strings.TrimPrefix(string(data), "\""), "\"")
	e := NewEncryptString(raw)
	if e == nil {
		return nil
	}
	*es = *e
	return nil
}

func encrypt(stringToEncrypt string, keyString string) (string, error) {
	sc, err := hex.DecodeString(keyString)
	if err != nil {
		return "", err
	}
	plaintext := []byte(stringToEncrypt)

	block, err := aes.NewCipher(sc)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext), nil
}

func decrypt(encryptedString string, keyString string) (string, error) {
	sc, err := hex.DecodeString(keyString)
	if err != nil {
		return "", err
	}
	enc, err := hex.DecodeString(encryptedString)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(sc)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", plaintext), nil
}

func dynamicDecrypt(str string) string {
	if EncryptStringSecret == nil || len(EncryptStringSecret) == 0 {
		return str
	}
	for _, sc := range EncryptStringSecret {
		decrypted, err := decrypt(str, sc)
		if err == nil {
			return decrypted
		}
	}
	return str
}

func dynamicEncrypt(str string) string {
	if EncryptStringSecret == nil || len(EncryptStringSecret) == 0 {
		return str
	}
	sc := EncryptStringSecret[len(EncryptStringSecret)-1]
	enc, err := encrypt(str, sc)
	if err == nil {
		return enc
	}
	return str
}

func (es *EncryptString) Scan(value any) error {
	if es == nil {
		return nil
	}
	*es = EncryptString{}
	var enc *string
	switch v := value.(type) {
	case []byte:
		enc = ptr.String(string(v))
	case string:
		enc = ptr.String(v)
	default:
		enc = nil
	}
	if enc != nil {
		str := dynamicDecrypt(*enc)
		if str == *enc {
			enc = ptr.String(dynamicEncrypt(str))
		}

		*es = EncryptString{
			str:     ptr.String(str),
			encrypt: enc,
		}
	}
	return nil
}

func (es *EncryptString) Value() (driver.Value, error) {
	s := es.EncryptString()
	if s == "" {
		return nil, nil
	}
	return es, nil
}

func NewEncryptString(str string) *EncryptString {
	if str == "" {
		return nil
	}
	enc := dynamicEncrypt(str)
	return &EncryptString{
		str:     &str,
		encrypt: &enc,
	}
}
