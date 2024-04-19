package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"net/http"

	"github.com/yobadagne/user_registration/model"
	"github.com/yobadagne/user_registration/util"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// adapter for auth layer
type AuthLayer struct {
}

func NewAuthLayer() model.AuthLayer {
	return &AuthLayer{}
}

func PKCS7Unpad(data []byte) []byte{
	length := len(data)
	unpadding := int(data[length-1]) 
	return  data[:(length -unpadding)]
}
func (a *AuthLayer) GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {		
		util.Logger.Error("Can not hash password using bcrypt, error while excuting auth.GenerateHashPassword()", zap.Error(err))
		// change error into errorx format
		err = model.MyError{
			Code: http.StatusInternalServerError,
			Message: "Can not hash password",
		}
		return "", err
	}

	return string(bytes), err
}

// compare for login
func (a *AuthLayer) CompareHashPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		// change error into errorx format
		util.Logger.Error("password does not match while comparing using bcrypt.CompareHashPassword method, error while excuting auth.CompareHashPassword()",zap.Error(err))
		err = model.MyError{
			Code: http.StatusUnauthorized,
			Message: "Password not correct",
		}
		return err
	}
	return nil
}

// encryptToken encrypts a refresh token using AES encryption
func (a *AuthLayer) EncryptToken(token string , iv []byte) (string, error) {
	block, err := aes.NewCipher(model.Encriptionkey)
	if err != nil {
		// change error into errorx format
		util.Logger.Error("Can not create encryption block using AES encryption, error while excuting auth.EncryptToken()", zap.Error(err))
		err = model.MyError{
			Code: http.StatusInternalServerError,
			Message: "Can not create encryption block",
		}
		return " ", err
	}
	padLength:= aes.BlockSize - len(token)%aes.BlockSize
	pad:= make([]byte, padLength)
	token += string(pad)
	ciphertext := make([]byte, aes.BlockSize + len(token)) 
	mode := cipher.NewCBCEncrypter(block, iv) 
	mode.CryptBlocks(ciphertext[aes.BlockSize:],[]byte(token))


	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decryptToken decrypts an encrypted refresh token using AES decryption
func (a *AuthLayer) DecryptToken(encryptedToken string) (string, error) {
	block, err := aes.NewCipher(model.Encriptionkey)
	if err != nil {
		// change error into errorx format
		util.Logger.Error("Can not create decryption block using AES decryption,error while excuting auth.DecryptToken()", zap.Error(err))
		err = model.MyError{
			Code: http.StatusInternalServerError,
			Message: "Can not create decryption block using AES decryption",
		}
		return "", err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedToken) 
	if err != nil {
		util.Logger.Error("Error while decrypting token using AES,error while excuting auth.DecryptToken()", zap.Error(err))
		err = model.MyError{
			Code: http.StatusInternalServerError,
			Message: "Error while decrypting",
		}
		return "", err
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(block, iv) 
	plaintext := make([]byte, len(ciphertext)) 
	mode.CryptBlocks(plaintext,ciphertext)

	//remove padding
	plaintext = PKCS7Unpad(plaintext)

	return string(plaintext), nil
}
