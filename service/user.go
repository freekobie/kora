package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/freekobie/kora/session"
	"github.com/freekobie/kora/mail"
	"github.com/freekobie/kora/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	VERIFICATION   = "verification"
	AUTHENTICATION = "authentication"
)

type UserService struct {
	store model.UserStore
	mail  *mail.Mailer
}

func NewUserService(us model.UserStore, m *mail.Mailer) *UserService {
	return &UserService{
		store: us,
		mail:  m,
	}
}

// CreateUser creates a new user with the given details
func (s *UserService) CreateUser(ctx context.Context, name, email, password string) (*model.User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return nil, ErrFailedOperation
	}
	now := time.Now().UTC()
	user := &model.User{
		Id:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    now,
		LastModifed:  now,
		Verified:     false,
	}

	if err := s.store.InsertUser(ctx, user); err != nil {
		return nil, err
	}

	otpString := generateOTP()
	slog.Debug("OTP verificatio code", "code", otpString) //TODO: delete this line later
	otpHash := hashString(otpString)

	userAddr := mail.Address{Name: user.Name, Email: user.Email}
	data := mail.Data{
		Address: userAddr,
		Code:    otpString,
	}

	token := model.UserToken{
		Hash:      otpHash,
		UserId:    user.Id,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		Scope:     VERIFICATION,
	}

	_ = s.store.InsertToken(ctx, &token)

	s.sendEmail([]mail.Address{userAddr}, "verify_email.html", data)

	return user, nil
}

func (us *UserService) VerifyUser(ctx context.Context, code string, email string) (*model.User, error) {

	hash := hashString(code)
	user, err := us.store.GetUserForToken(ctx, hash, VERIFICATION, email)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	user.Verified = true

	err = us.store.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// Delete otp after successful verification
	_ = us.store.DeleteToken(ctx, hash, VERIFICATION)

	address := mail.Address{Name: user.Name, Email: user.Email}
	us.sendEmail([]mail.Address{address}, "welcome_email.html", mail.Data{Address: address})

	return user, nil
}

func (us *UserService) ResendVerificationEmail(ctx context.Context, email string) error {
	user, err := us.store.GetUserByMail(ctx, email)
	if err != nil {
		return err
	} else if user.Verified {
		return errors.New("user already verified")
	}

	otpString := generateOTP()
	otpHash := hashString(otpString)

	slog.Debug("OTP verification code", "code", otpString) //TODO: delete this line later

	userAddr := mail.Address{Email: email, Name: user.Name}
	data := mail.Data{
		Address: userAddr,
		Code:    otpString,
	}

	token := model.UserToken{
		Hash:      otpHash,
		UserId:    user.Id,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		Scope:     VERIFICATION,
	}

	err = us.store.InsertToken(ctx, &token)
	if err != nil {
		return err
	}

	us.sendEmail([]mail.Address{userAddr}, "verify_email.html", data)

	return nil
}

func (us *UserService) NewSession(ctx context.Context, email string, password string) (*session.UserSession, error) {
	user, err := us.store.GetUserByMail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !user.Verified {
		return nil, ErrUnverifiedUser
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidCredentials
		}
		slog.Error("failed to compare password and hash", "error", err.Error())
		return nil, ErrFailedOperation
	}

	ttl := 15 * (24 * time.Hour)
	refresh, err := session.GenerateToken(user.Id, user.Email, ttl, session.TokenTypeRefresh)
	if err != nil {
		return nil, ErrFailedOperation
	}

	expiresAt := time.Now().Add(ttl)
	tokenHash := hashString(refresh)
	token := model.UserToken{
		Hash:      tokenHash,
		UserId:    user.Id,
		ExpiresAt: expiresAt,
		Scope:     AUTHENTICATION,
	}

	err = us.store.InsertToken(ctx, &token)
	if err != nil {
		return nil, err
	}

	session := &session.UserSession{
		User:         user,
		RefreshToken: refresh,
		ExpiresAt:    expiresAt,
	}

	return session, nil
}

func (us *UserService) RefreshSession(ctx context.Context, refreshToken string) (*session.UserAccess, error) {
	claims, err := session.ValidateToken(refreshToken, session.TokenTypeRefresh)
	if err != nil {
		slog.Error("failed token validation", "error", err.Error())
		return nil, session.ErrInvalidToken
	}

	hash := hashString(refreshToken)

	user, err := us.store.GetUserForToken(ctx, hash, AUTHENTICATION, claims.Email)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	ttl := 2 * time.Hour // TODO: make time shorter
	accessToken, err := session.GenerateToken(user.Id, user.Email, ttl, session.TokenTypeAccess)
	if err != nil {
		return nil, err
	}
	// FIXME: obtain expiry time from GenerateToken function
	useracc := &session.UserAccess{
		AccessToken: accessToken,
		ExpiresAt:   time.Now().Add(ttl),
	}

	return useracc, nil
}

// UpdateUser updates an existing user's details
func (us *UserService) UpdateUser(ctx context.Context, userData map[string]any) (*model.User, error) {
	id, ok := userData["id"]
	if !ok {
		return nil, errors.New("user id not found")
	}

	user, err := us.store.GetUser(ctx, id.(uuid.UUID))
	if err != nil {
		return nil, err
	}

	name, ok := userData["name"]
	if ok {
		user.Name = name.(string)
	}

	profilePhoto, ok := userData["profilePhoto"]
	if ok {
		user.ProfilePhoto = profilePhoto.(string)
	}

	password, ok := userData["password"]
	if ok {
		if len(password.(string)) < 8 || len(password.(string)) > 20 {
			return nil, ErrInvalidPassword
		}
		err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password.(string)))
		if err != nil {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				hash, err := bcrypt.GenerateFromPassword([]byte(password.(string)), bcrypt.DefaultCost)
				if err != nil {
					return nil, ErrFailedOperation
				}

				user.PasswordHash = hash
			} else {
				slog.Error("failed to compare password and hash", "error", err.Error())
				return nil, ErrFailedOperation
			}
		}

	}

	user.LastModifed = time.Now().UTC()

	err = us.store.UpdateUser(ctx, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FetchUser retrieves a user by ID or email
func (s *UserService) FetchUser(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.store.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteUser removes a user from the system
func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.store.DeleteUser(ctx, id)
}
