package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
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
func (a AuthLayer) GenerateHashPassword(c *gin.Context, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		// change error into errorx format

		err = errorx.Decorate(err, "Can not hash password")
		util.Logger.Error("Can not hash password", zap.Error(err))
		c.Set(model.Error_type, model.INTERNAL_SERVER_ERROR)
		return "", err
	}
	return string(bytes), nil
}

// compare for login
func (a AuthLayer) CompareHashPassword(c *gin.Context, password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		// change error into errorx format
		err = errorx.Decorate(err, "Can not compare password")
		util.Logger.Error("Error while comparing password", zap.Error(err))
		c.Set(model.Error_type, model.INTERNAL_SERVER_ERROR)
		return err
	}
	return nil
}

// encryptToken encrypts a refresh token using AES encryption
func (a AuthLayer) EncryptToken(c *gin.Context, token string , iv []byte) (string, error) {
	block, err := aes.NewCipher(model.Encriptionkey)
	if err != nil {
		// change error into errorx format
		err = errorx.Decorate(err, "Can not create encription block")
		util.Logger.Error("Can not create encription block", zap.Error(err))
		c.Set(model.Error_type, model.INTERNAL_SERVER_ERROR)
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
func (a AuthLayer) DecryptToken(c *gin.Context, encryptedToken string) (string, error) {
	block, err := aes.NewCipher(model.Encriptionkey)
	if err != nil {
		// change error into errorx format
		err = errorx.Decorate(err, "Can not create decription block")
		util.Logger.Error("Can not create decription block", zap.Error(err))
		c.Set(model.Error_type, model.INTERNAL_SERVER_ERROR)
		return "", err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedToken) 
	if err != nil {
		err := errorx.Decorate(errorx.IllegalArgument.New("Error while decrypting"),"Error while decrypting")
		util.Logger.Error("Error while decrypting", zap.Error(err))
		c.Set(model.Error_type, model.BAD_REQUEST)
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
