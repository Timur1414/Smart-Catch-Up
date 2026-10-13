package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/domain"
	"github.com/Timur1414/Smart-Catch-Up/internal/app/auth/repository"
	"github.com/Timur1414/Smart-Catch-Up/pkg/crypto"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	Create(ctx context.Context, user domain.User) (int, error)
	GetById(ctx context.Context, id int) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Update(ctx context.Context, user domain.User) error
	IsExists(ctx context.Context, user domain.User) (domain.User, error)
	Setup2FA(ctx context.Context, userId int) (string, string, error)
	Enable2FA(ctx context.Context, userId int, code string, password string) ([]string, error)
	Disable2FA(ctx context.Context, userId int, code string) error
	Verify2FA(ctx context.Context, userId int, code string) (bool, error)
}

type User struct {
	repository repository.UserRepository
	crypto     *crypto.Encryptor
}

func NewUser(repo repository.UserRepository, crypto *crypto.Encryptor) *User {
	return &User{
		repository: repo,
		crypto:     crypto,
	}
}

func (obj *User) Create(ctx context.Context, user domain.User) (int, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	userToCreate := domain.User{
		Email: user.Email,
	}
	origPassword := []byte(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(origPassword, bcrypt.DefaultCost)
	if err != nil {
		log.Warn("failed to hash password", zap.Error(err))
		return -1, domain.ErrFailedToHashPassword
	}
	userToCreate.Password = string(hashedPassword)
	return obj.repository.Create(ctx, userToCreate)
}

func (obj *User) GetById(ctx context.Context, id int) (domain.User, error) {
	return obj.repository.GetById(ctx, id)
}

func (obj *User) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return obj.repository.GetByEmail(ctx, email)
}

func (obj *User) Update(ctx context.Context, user domain.User) error {
	return obj.repository.Update(ctx, user)
}

func (obj *User) IsExists(ctx context.Context, user domain.User) (domain.User, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	userFromDb, err := obj.GetByEmail(ctx, user.Email)
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userFromDb.Password), []byte(user.Password))
	if err != nil {
		log.Warn("failed to compare password", zap.Error(err))
		return domain.User{}, err
	}
	return userFromDb, nil
}

func (obj *User) Setup2FA(ctx context.Context, userId int) (string, string, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	user, err := obj.repository.GetById(ctx, userId)
	if err != nil {
		return "", "", err
	}
	if user.TotpEnabled {
		log.Error("TOTP is already enabled", zap.Int("userId", userId))
		return "", "", domain.ErrTotpAlreadyExists
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "SmartCatchUp",
		AccountName: user.Email,
	})
	if err != nil {
		log.Error("failed to generate TOTP", zap.Error(err))
		return "", "", err
	}
	secret := key.Secret()
	otpauthUrl := key.URL()
	encryptedSecret, err := obj.crypto.Encrypt(secret)
	if err != nil {
		return "", "", err
	}
	err = obj.repository.SaveTotpSecret(ctx, userId, encryptedSecret)
	if err != nil {
		return "", "", err
	}
	log.Info("2FA setup initialized", zap.Int("userId", userId))
	return secret, otpauthUrl, nil
}

func (obj *User) Enable2FA(ctx context.Context, userId int, code string, password string) ([]string, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	user, err := obj.repository.GetById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if user.TotpEnabled {
		log.Error("TOTP is already enabled", zap.Int("userId", userId))
		return []string{}, domain.ErrTotpAlreadyExists
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Warn("failed to compare password", zap.Error(err))
		return []string{}, err
	}
	if user.TotpSecret == "" {
		return nil, domain.ErrNotFound
	}
	cleanCode := strings.TrimSpace(code)
	decryptedSecret, err := obj.crypto.Decrypt(user.TotpSecret)
	if err != nil {
		return []string{}, err
	}
	if !totp.Validate(cleanCode, decryptedSecret) {
		log.Warn("invalid code", zap.Int("userId", userId))
		return nil, domain.ErrInvalidTotpCode
	}
	backupCodes, hashedCodes, err := generateBackupCodes(8)
	if err != nil {
		log.Error("failed to generate backup codes", zap.Error(err))
		return nil, err
	}
	err = obj.repository.EnableTotp(ctx, userId, hashedCodes)
	if err != nil {
		return nil, err
	}
	log.Info("2fa successfully enabled", zap.Int("user_id", userId))
	return backupCodes, nil
}

func (obj *User) Verify2FA(ctx context.Context, userId int, code string) (bool, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	user, err := obj.repository.GetById(ctx, userId)
	if err != nil {
		return false, err
	}
	if !user.TotpEnabled || user.TotpSecret == "" {
		return false, domain.ErrInvalidTotpCode
	}
	cleanCode := strings.TrimSpace(code)
	decryptedSecret, err := obj.crypto.Decrypt(user.TotpSecret)
	if err != nil {
		return false, err
	}
	if totp.Validate(cleanCode, decryptedSecret) {
		return true, nil
	}
	if len(cleanCode) != 8 {
		return false, domain.ErrInvalidTotpCode
	}
	for _, backupCode := range user.TotpBackupCodes {
		err = bcrypt.CompareHashAndPassword([]byte(backupCode), []byte(cleanCode))
		if err == nil {
			err = obj.repository.RemoveBackupCode(ctx, userId, backupCode)
			if err != nil {
				return false, err
			}
			log.Info("backup 2fa successfully used", zap.Int("user_id", userId))
			return true, nil
		}
	}
	return false, nil
}

func (obj *User) Disable2FA(ctx context.Context, userId int, code string) error {
	log := logger.GetLoggerWithRequestId(ctx)
	user, err := obj.repository.GetById(ctx, userId)
	if err != nil {
		return err
	}
	if !user.TotpEnabled || user.TotpSecret == "" {
		return nil
	}
	valid, err := obj.Verify2FA(ctx, userId, code)
	if err != nil {
		return err
	}
	if !valid {
		return domain.ErrInvalidTotpCode
	}
	err = obj.repository.DisableTotp(ctx, userId)
	if err != nil {
		return err
	}
	log.Info("2fa successfully disabled", zap.Int("user_id", userId))
	return nil
}

func generateBackupCodes(count int) ([]string, []string, error) {
	codes := make([]string, count)
	hashedCodes := make([]string, count)
	for i := range count {
		bytes := make([]byte, 4)
		if _, err := rand.Read(bytes); err != nil {
			return []string{}, []string{}, err
		}
		codes[i] = hex.EncodeToString(bytes)
		hashedCode, err := bcrypt.GenerateFromPassword([]byte(codes[i]), bcrypt.DefaultCost)
		if err != nil {
			return []string{}, []string{}, domain.ErrFailedToHashPassword
		}
		hashedCodes[i] = string(hashedCode)
	}
	return codes, hashedCodes, nil
}
