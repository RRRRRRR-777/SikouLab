# Frontend Repo Map

> Auto-generated. Do not edit manually.

## Routes (app/)

- frontend/app/(main)/layout.tsx → MainLayout
- frontend/app/(main)/page.tsx → Home
- frontend/app/(main)/settings/page.tsx → SettingsPageRoute
- frontend/app/layout.tsx → RootLayout
- frontend/app/login/page.tsx → Login
- frontend/app/subscription/page.tsx → Subscription

## Components

### frontend/components

- frontend/components/GlobalErrorHandler.tsx
  - GlobalErrorHandler
### frontend/components/auth

- frontend/components/auth/LoginPage.tsx
  - LoginPage
### frontend/components/layout

- frontend/components/layout/Header.tsx
  - Header
- frontend/components/layout/Sidebar.tsx
  - Sidebar
### frontend/components/newsletter

- frontend/components/newsletter/NewsletterWidget.tsx
  - NewsletterWidget
### frontend/components/settings

- frontend/components/settings/EmailSection.tsx
  - EmailSection
- frontend/components/settings/FaqSection.tsx
  - FaqSection
- frontend/components/settings/ProfileSection.tsx
  - ProfileSection
- frontend/components/settings/SettingsPage.tsx
  - SettingsPage
- frontend/components/settings/SubscriptionSection.tsx
  - SubscriptionSection
### frontend/components/subscription

- frontend/components/subscription/SubscriptionPage.tsx
  - SubscriptionPage
### frontend/components/ui

- frontend/components/ui/button.tsx
  - Button
  - buttonVariants
- frontend/components/ui/skeleton.tsx
  - Skeleton
- frontend/components/ui/sonner.tsx
  - Toaster

## Utilities (lib/)

- frontend/lib/api.ts
  - apiClient
- frontend/lib/auth/auth-api.ts
  - authApi
- frontend/lib/auth/auth-context.tsx
  - AuthProvider
  - useAuth
- frontend/lib/auth/firebase.ts
  - initializeFirebaseApp
  - getFirebaseAuth
  - signInWithGoogle
  - signInWithApple
  - signInWithX
  - signOut
  - onAuthStateChangedHelper
  - getIdToken
- frontend/lib/auth/types.ts
  - User
  - AuthUser
  - LoginResponse
  - MeResponse
  - ApiError
  - OAuthProvider
  - toAuthUser
- frontend/lib/logger.ts
  - LogLevel
  - LogEntry
  - log
  - logError
  - logWarn
  - logInfo
  - logDebug
- frontend/lib/settings/settings-api.ts
  - UpdateProfileRequest
  - UserResponse
  - AvatarResponse
  - SubscriptionResponse
  - PortalResponse
  - NewsletterSubscriptionResponse
  - NewsletterSubscribeRequest
  - NewsletterUpdateRequest
  - settingsApi
- frontend/lib/subscription/subscription-api.ts
  - Plan
  - subscriptionApi
- frontend/lib/subscription/univapay.ts
  - CheckoutWidgetConfig
  - openCheckoutWidget
- frontend/lib/utils.ts
  - cn
- frontend/lib/validation.ts
  - EMAIL_REGEX

## Other

- frontend/app/api/log/route.ts
  - POST
- frontend/app/providers.tsx
  - Providers
- frontend/middleware.ts
  - middleware
  - config
  - __internal__

