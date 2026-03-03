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
- type SubscriptionHandler
  - func (SubscriptionHandler) ServeCheckout
  - func (SubscriptionHandler) ServeGetPlans
  - func (SubscriptionHandler) ServeWebhook
- func NewAuthHandler
- func NewSubscriptionHandler

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

- type SubscriptionRepository
  - func (SubscriptionRepository) FindActivePlans
  - func (SubscriptionRepository) FindByUnivaPaySubscriptionID
  - func (SubscriptionRepository) FindPlanByID
  - func (SubscriptionRepository) UpdateSubscriptionStatus
  - func (SubscriptionRepository) UpdateUnivaPaySubscriptionID
- type UserRepository
  - func (UserRepository) Create
  - func (UserRepository) FindByID
  - func (UserRepository) FindByOAuth
- func NewDB
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
- type SubscriptionRepository
- type SubscriptionUsecase
  - func (SubscriptionUsecase) Checkout
  - func (SubscriptionUsecase) GetPlans
  - func (SubscriptionUsecase) HandleWebhook
- type UnivaPayClient
- type UserRepository
- type WebhookData
- type WebhookPayload
- func NewAuthUsecase
- func NewSubscriptionUsecase
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
- var ErrInvalidToken

