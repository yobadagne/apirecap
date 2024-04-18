package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	

	//"github.com/joomcode/errorx"
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
func (a AuthLayer) GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {		
		util.Logger.Error("Can not hash password using bcrypt", zap.Error(err))
		// change error into errorx format
		err = model.ErrInternalServerErr.New("Can not hash password")
		model.Error_type = model.INTERNAL_SERVER_ERROR
		return "", err
	}

	return string(bytes), err
}

// compare for login
func (a AuthLayer) CompareHashPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		// change error into errorx format
		util.Logger.Error("password does not match",zap.Error(err))
		err = model.ErrBadRequest.New("Password not correct")
		model.Error_type = model.UNAUTHORIZED
		return err
	}
	return nil
}

// encryptToken encrypts a refresh token using AES encryption
func (a AuthLayer) EncryptToken(token string , iv []byte) (string, error) {
	block, err := aes.NewCipher(model.Encriptionkey)
	if err != nil {
		// change error into errorx format
		util.Logger.Error("Can not create encription block using AES encryption", zap.Error(err))
		err = model.ErrInternalServerErr.New("Can not create encryption block")
		model.Error_type = model.INTERNAL_SERVER_ERROR
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
func (a AuthLayer) DecryptToken(encryptedToken string) (string, error) {
	block, err := aes.NewCipher(model.Encriptionkey)
	if err != nil {
		// change error into errorx format
		util.Logger.Error("Can not create decryption block using AES decryption", zap.Error(err))
		err = model.ErrInternalServerErr.New("Can not create decryption block using AES decryption")
		model.Error_type = model.INTERNAL_SERVER_ERROR
		return "", err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedToken) 
	if err != nil {
		util.Logger.Error("Error while decrypting", zap.Error(err))
		err = model.ErrBadRequest.New("Error while decrypting")
		model.Error_type = model.BAD_REQUEST
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
