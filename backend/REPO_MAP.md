# Backend Repo Map

> Auto-generated. Do not edit manually.

## Entry Points

- backend/check_doc.go
- backend/cmd/api/main.go

## Packages

### backend/internal/config

- type Config
- func Load

### backend/internal/domain

- type NewsletterSubscription
- type Plan
- type User
- type UserSettings

### backend/internal/firebase

- type Client
  - func (Client) CreateSessionCookie
  - func (Client) VerifyIDToken
  - func (Client) VerifySessionCookie
- type FirebaseToken
- type TokenVerifier
- func NewClient

### backend/internal/handler

- type AuthHandler
  - func (AuthHandler) ServeLogin
  - func (AuthHandler) ServeLogout
  - func (AuthHandler) ServeMe
- type HealthHandler
  - func (HealthHandler) ServeHTTP
- type NewsletterHandler
  - func (NewsletterHandler) ServeGetSubscription
  - func (NewsletterHandler) ServeSubscribe
  - func (NewsletterHandler) ServeUnsubscribe
  - func (NewsletterHandler) ServeUpdateEmail
- type SubscriptionHandler
  - func (SubscriptionHandler) ServeCheckout
  - func (SubscriptionHandler) ServeGeneratePortalURL
  - func (SubscriptionHandler) ServeGetMySubscription
  - func (SubscriptionHandler) ServeGetPlans
  - func (SubscriptionHandler) ServeWebhook
- type UserHandler
  - func (UserHandler) ServeDeleteAvatar
  - func (UserHandler) ServeUpdateProfile
  - func (UserHandler) ServeUploadAvatar
- func NewAuthHandler
- func NewNewsletterHandler
- func NewSubscriptionHandler
- func NewUserHandler

### backend/internal/infrastructure/storage

- type GCSStorage
  - func (GCSStorage) Close
  - func (GCSStorage) Delete
  - func (GCSStorage) Save
- func NewGCSStorage

### backend/internal/infrastructure/univapay

- type Client
  - func (Client) CreateSubscription
  - func (Client) GetSubscription
- func NewClient

### backend/internal/middleware

- func CORS
- func ContextWithUser
- func Logger
- func Recovery
- func RequireAuth
- func RequireRole
- func RequireSubscription
- func UserFromContext

### backend/internal/repository

- type NewsletterRepository
  - func (NewsletterRepository) FindByUserID
  - func (NewsletterRepository) UpdateEmail
  - func (NewsletterRepository) UpdateIsActive
  - func (NewsletterRepository) Upsert
- type SubscriptionRepository
  - func (SubscriptionRepository) FindActivePlans
  - func (SubscriptionRepository) FindByUnivaPaySubscriptionID
  - func (SubscriptionRepository) FindPlanByID
  - func (SubscriptionRepository) UpdateSubscriptionStatus
  - func (SubscriptionRepository) UpdateUnivaPaySubscriptionID
- type UserRepository
  - func (UserRepository) ClearAvatarURL
  - func (UserRepository) Create
  - func (UserRepository) FindByID
  - func (UserRepository) FindByOAuth
  - func (UserRepository) UpdateAvatarURL
  - func (UserRepository) UpdateDisplayName
  - func (UserRepository) UpdateEmail
- func NewDB
- func NewNewsletterRepository
- func NewSubscriptionRepository
- func NewUserRepository

### backend/internal/router

- type Builder
  - func (Builder) Build
- type Handlers
- type Middlewares
- func NewBuilder

### backend/internal/usecase

- type AuthUsecase
  - func (AuthUsecase) CreateSessionCookie
  - func (AuthUsecase) GetCurrentUser
  - func (AuthUsecase) Login
- type NewsletterRepository
- type NewsletterUsecase
  - func (NewsletterUsecase) GetSubscription
  - func (NewsletterUsecase) Subscribe
  - func (NewsletterUsecase) Unsubscribe
  - func (NewsletterUsecase) UpdateEmail
- type ObjectStorage
- type SubscriptionInfo
- type SubscriptionRepository
- type SubscriptionUsecase
  - func (SubscriptionUsecase) Checkout
  - func (SubscriptionUsecase) GeneratePortalURL
  - func (SubscriptionUsecase) GetMySubscription
  - func (SubscriptionUsecase) GetPlans
  - func (SubscriptionUsecase) HandleWebhook
- type UnivaPayClient
- type UserProfileRepository
- type UserRepository
- type UserUsecase
  - func (UserUsecase) DeleteAvatar
  - func (UserUsecase) UpdateDisplayName
  - func (UserUsecase) UploadAvatar
- type WebhookData
- type WebhookPayload
- func NewAuthUsecase
- func NewNewsletterUsecase
- func NewSubscriptionUsecase
- func NewUserUsecase
- const EventSubscriptionCanceled
- const EventSubscriptionFailure
- const EventSubscriptionPayment
- const SubscriptionStatusActive
- const SubscriptionStatusCanceled
- const SubscriptionStatusPastDue
- const SubscriptionStatusTrialing
- const UnivaPayStatusCurrent
- const UnivaPayStatusUnconfirmed
- const UnivaPayStatusUnpaid
- var ErrAlreadySubscribed
- var ErrDisplayNameBlankOnly
- var ErrDisplayNameEmpty
- var ErrDisplayNameTooLong
- var ErrInvalidEmail
- var ErrInvalidToken
- var ErrNotFound
- var ErrStorageNotConfigured

