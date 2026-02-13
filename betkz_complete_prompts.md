# BetKZ Complete Development Prompts

This document contains all prompts needed to build the BetKZ sportsbook platform with dynamic odds system.

---

## Table of Contents

1. [Phase 1: Project Setup & Architecture](#phase-1-project-setup--architecture)
2. [Phase 2: Backend Foundation](#phase-2-backend-foundation)
3. [Phase 3: Frontend Foundation](#phase-3-frontend-foundation)
4. [Phase 4: User Features](#phase-4-user-features)
5. [Phase 5: Admin Panel](#phase-5-admin-panel)
6. [Phase 6: Polish & Testing](#phase-6-polish--testing)
7. [Phase 7: Deployment & Production](#phase-7-deployment--production)

---

## Phase 1: Project Setup & Architecture

### Part 1.1: Project Initialization

```
I'm building a full-stack sportsbook web application called BetKZ using React (frontend) and Go (backend) with a DYNAMIC ODDS SYSTEM.

Project structure requirements:
- Monorepo setup with separate frontend and backend directories
- Frontend: React with Vite, TypeScript, TailwindCSS
- Backend: Go with Gin framework, PostgreSQL database, Redis for caching
- WebSocket support for real-time odds updates
- Docker Compose for local development
- Environment configuration for dev/prod

Tech stack:
- Backend: Go, Gin, PostgreSQL, Redis, gorilla/websocket
- Frontend: React 18, TypeScript, Vite, TailwindCSS, Zustand, React Query
- Real-time: WebSocket for odds broadcasting
- Cache: Redis for current odds and market data

Please provide:
1. Complete project directory structure with websocket and redis support
2. Package.json and go.mod configurations (including websocket libraries)
3. Docker Compose setup with PostgreSQL, Redis, and WebSocket
4. Environment variable templates
5. Basic README with setup instructions
```

---

### Part 1.2: Database Schema Design with Dynamic Odds

```
Design a PostgreSQL database schema for a sportsbook platform with DYNAMIC ODDS SYSTEM:

Core tables needed:

-- Users table
Users (
  id UUID PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  balance DECIMAL(15,2) DEFAULT 0.00,
  role VARCHAR(20) DEFAULT 'user', -- 'user' or 'admin'
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
)

-- Sports categories
Sports (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  slug VARCHAR(100) UNIQUE NOT NULL,
  is_active BOOLEAN DEFAULT TRUE
)

-- Events (matches/games)
Events (
  id UUID PRIMARY KEY,
  sport_id INTEGER REFERENCES sports(id),
  home_team VARCHAR(255) NOT NULL,
  away_team VARCHAR(255) NOT NULL,
  start_time TIMESTAMP NOT NULL,
  status VARCHAR(50) DEFAULT 'upcoming', -- upcoming, live, finished, cancelled
  final_score_home INTEGER,
  final_score_away INTEGER,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
)

-- Markets (betting markets for events)
Markets (
  id UUID PRIMARY KEY,
  event_id UUID REFERENCES events(id) ON DELETE CASCADE,
  market_type VARCHAR(50) NOT NULL, -- '1x2', 'over_under', 'both_teams_score', etc.
  status VARCHAR(50) DEFAULT 'open', -- open, locked, settled, cancelled
  margin_percentage DECIMAL(5,2) DEFAULT 5.00, -- House margin (5%)
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
)

-- Odds with dynamic calculation support
Odds (
  id UUID PRIMARY KEY,
  market_id UUID REFERENCES markets(id) ON DELETE CASCADE,
  outcome VARCHAR(100) NOT NULL, -- 'home', 'draw', 'away', 'over', 'under', etc.
  initial_odds DECIMAL(10,2) NOT NULL, -- Admin sets this initially
  current_odds DECIMAL(10,2) NOT NULL, -- Dynamically calculated
  total_stake DECIMAL(15,2) DEFAULT 0.00, -- Total money bet on this outcome
  bet_count INTEGER DEFAULT 0, -- Number of bets on this outcome
  last_calculated_at TIMESTAMP,
  is_manual_override BOOLEAN DEFAULT FALSE, -- If admin manually set odds
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(market_id, outcome)
)

-- Market pool tracking
MarketPools (
  id UUID PRIMARY KEY,
  market_id UUID REFERENCES markets(id) ON DELETE CASCADE,
  total_pool DECIMAL(15,2) DEFAULT 0.00, -- Total money in this market
  house_margin DECIMAL(5,2) DEFAULT 5.00,
  liability DECIMAL(15,2) DEFAULT 0.00, -- Worst-case payout
  calculated_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(market_id)
)

-- Bets placed by users
Bets (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id),
  bet_type VARCHAR(20) NOT NULL, -- 'single' or 'accumulator'
  stake DECIMAL(15,2) NOT NULL,
  potential_return DECIMAL(15,2) NOT NULL,
  actual_return DECIMAL(15,2) DEFAULT 0.00,
  status VARCHAR(50) DEFAULT 'pending', -- pending, won, lost, cancelled
  placed_at TIMESTAMP DEFAULT NOW(),
  settled_at TIMESTAMP
)

-- Individual selections in a bet (bet legs)
BetLegs (
  id UUID PRIMARY KEY,
  bet_id UUID REFERENCES bets(id) ON DELETE CASCADE,
  market_id UUID REFERENCES markets(id),
  odd_id UUID REFERENCES odds(id),
  outcome VARCHAR(100) NOT NULL,
  locked_odd_value DECIMAL(10,2) NOT NULL, -- Odds locked at bet placement time
  result VARCHAR(50), -- won, lost, push, pending
  created_at TIMESTAMP DEFAULT NOW()
)

-- Transaction history
Transactions (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id),
  type VARCHAR(50) NOT NULL, -- deposit, withdraw, bet_placed, bet_won, bet_lost
  amount DECIMAL(15,2) NOT NULL,
  balance_before DECIMAL(15,2) NOT NULL,
  balance_after DECIMAL(15,2) NOT NULL,
  reference_id UUID, -- bet_id if related to bet
  status VARCHAR(50) DEFAULT 'completed',
  created_at TIMESTAMP DEFAULT NOW()
)

-- Odds history for analytics
OddsHistory (
  id UUID PRIMARY KEY,
  odd_id UUID REFERENCES odds(id) ON DELETE CASCADE,
  odds_value DECIMAL(10,2) NOT NULL,
  total_stake DECIMAL(15,2) NOT NULL,
  bet_count INTEGER NOT NULL,
  recorded_at TIMESTAMP DEFAULT NOW()
)

-- Admin activity log
AdminLogs (
  id UUID PRIMARY KEY,
  admin_id UUID REFERENCES users(id),
  action VARCHAR(100) NOT NULL,
  entity_type VARCHAR(50), -- event, market, odds
  entity_id UUID,
  details JSONB,
  created_at TIMESTAMP DEFAULT NOW()
)

Key features:
- Odds locked at bet placement (locked_odd_value in BetLegs)
- Track total stake per outcome for dynamic calculation
- Market pool tracking for risk management
- Odds history for analytics and charting
- Admin override capability

Please provide:
1. Complete SQL schema with all constraints and indexes
2. Indexes on: market_id, user_id, total_stake, current_odds, placed_at
3. Database triggers:
   - Auto-update updated_at timestamps
   - Auto-log odds changes to OddsHistory
   - Auto-update MarketPools when stakes change
4. Migration files structure (up/down migrations)
5. Sample seed data:
   - 2-3 sports (Football, Basketball, Tennis)
   - 5-10 upcoming events
   - Markets with initial odds set by admin
   - 2 test users (1 regular, 1 admin)
```

---

## Phase 2: Backend Foundation

### Part 2.1: Authentication System

```
Create a complete authentication system for a Go backend using Gin framework:

Requirements:
- User registration with email/password validation
- Login with JWT token generation (access + refresh tokens)
- Password hashing with bcrypt (cost 12)
- Token refresh mechanism
- Middleware for protected routes
- Role-based access control (user, admin)
- Rate limiting on auth endpoints

Database: PostgreSQL with users table

Please provide:
1. User model (struct) and repository layer (CRUD operations)
2. Authentication service with business logic:
   - Register(email, password) -> user, error
   - Login(email, password) -> tokens, error
   - RefreshToken(refreshToken) -> newTokens, error
   - VerifyToken(token) -> claims, error
3. JWT utilities:
   - GenerateTokenPair(userId, role) -> accessToken, refreshToken
   - ValidateToken(token) -> claims, error
   - RefreshAccessToken(refreshToken) -> newAccessToken
4. HTTP handlers:
   - POST /api/auth/register
   - POST /api/auth/login
   - POST /api/auth/refresh
   - POST /api/auth/logout
5. Middleware:
   - AuthMiddleware() - validates JWT
   - AdminMiddleware() - requires admin role
6. Input validation (email format, password min 8 chars)
7. Error handling with proper HTTP status codes
8. Password strength validation
```

---

### Part 2.2: Sports & Events Management

```
Build the sports and events management module for the Go backend:

Features needed:
- CRUD operations for sports categories
- CRUD operations for events (matches)
- Event status management (upcoming → live → finished → cancelled)
- Filtering events by:
  - Sport category
  - Date range
  - Status
  - Team name (search)
- Pagination for events list (page, limit, offset)
- Admin-only endpoints for creating/updating events
- Public endpoints for browsing events

Database tables: sports, events (with foreign key to sports)

Please provide:
1. Models (structs):
   - Sport{ID, Name, Slug, IsActive}
   - Event{ID, SportID, HomeTeam, AwayTeam, StartTime, Status, FinalScores}
2. Repository layer:
   - CreateEvent, GetEvent, UpdateEvent, DeleteEvent
   - ListEvents(filters, pagination)
   - GetEventsBySport(sportId)
3. Service layer with business logic:
   - ValidateEventData
   - CanUpdateEvent (status rules)
   - ScheduleEventStart (for auto status change)
4. HTTP handlers and routes:
   - GET /api/sports (public)
   - GET /api/events (public, with filters)
   - GET /api/events/:id (public)
   - POST /api/admin/events (admin only)
   - PUT /api/admin/events/:id (admin only)
   - DELETE /api/admin/events/:id (admin only)
   - PATCH /api/admin/events/:id/status (admin only)
5. Request/Response DTOs:
   - CreateEventRequest, UpdateEventRequest
   - EventResponse, EventListResponse
6. Validation rules:
   - StartTime must be in future
   - Teams cannot be empty
   - Valid status transitions
7. API documentation examples
```

---

### Part 2.3: Dynamic Odds System Core

```
Implement the DYNAMIC ODDS CALCULATION ENGINE - the core of the betting system:

System Overview:
- Admin sets initial odds for each market outcome
- Odds automatically adjust based on betting volume
- Users see live odds that change in real-time
- Each bet locks the current odds for that user

Core Components:

1. **Odds Calculation Service**:

Formula to implement:
```go
// For each outcome in a market
func CalculateOdds(marketID uuid.UUID) error {
    // Get all odds for this market
    odds := GetOddsByMarket(marketID)
    
    // Calculate total pool
    totalPool := 0.0
    for _, odd := range odds {
        totalPool += odd.TotalStake
    }
    
    // If no bets yet, use initial odds
    if totalPool == 0 {
        // current_odds = initial_odds
        return nil
    }
    
    // Get market margin (default 5%)
    margin := GetMarketMargin(marketID) // e.g., 0.05
    
    // For each outcome, calculate new odds
    for _, odd := range odds {
        // Skip if manual override
        if odd.IsManualOverride {
            continue
        }
        
        // Implied probability from stakes
        impliedProb := odd.TotalStake / totalPool
        
        // Adjust with house margin
        adjustedProb := impliedProb * (1 + margin)
        
        // Prevent division by zero
        if adjustedProb <= 0 {
            adjustedProb = 0.01
        }
        
        // Calculate new odds
        newOdds := 1 / adjustedProb
        
        // Apply boundaries (1.01 to 100.00)
        newOdds = math.Max(1.01, math.Min(newOdds, 100.00))
        
        // Update database
        UpdateOdds(odd.ID, newOdds)
        
        // Log to history
        LogOddsHistory(odd.ID, newOdds, odd.TotalStake, odd.BetCount)
    }
    
    return nil
}
```

2. **Markets & Odds Management**:

Database tables: markets, odds, market_pools

Features:
- Create markets for events (1X2, Over/Under, Both Teams Score, etc.)
- Set initial odds for each outcome
- Lock markets when events start
- Calculate potential returns based on current odds
- Track market liquidity and risk

3. **Stake Tracking**:

When bet is placed:
```go
func UpdateStakeOnBet(oddID uuid.UUID, stakeAmount float64) error {
    // Begin transaction
    tx := BeginTransaction()
    
    // Update odds.total_stake
    UpdateOddsTotalStake(oddID, stakeAmount) // += stakeAmount
    
    // Increment odds.bet_count
    IncrementBetCount(oddID)
    
    // Get market_id
    marketID := GetMarketIDFromOdd(oddID)
    
    // Update market pool
    UpdateMarketPool(marketID, stakeAmount)
    
    // Recalculate odds for entire market
    CalculateOdds(marketID)
    
    // Commit transaction
    tx.Commit()
    
    return nil
}
```

4. **API Endpoints**:

Public endpoints:
- GET /api/events/:id/markets - Get all markets for an event
- GET /api/markets/:id/odds - Get current odds for a market

Admin endpoints:
- POST /api/admin/markets - Create market for event
- POST /api/admin/markets/:id/odds - Set initial odds
- PUT /api/admin/markets/:id/odds/:oddId - Update odds (manual override)
- PATCH /api/admin/markets/:id/margin - Update house margin
- DELETE /api/admin/markets/:id - Delete market

5. **Odds History Logging**:
- Every odds change logged to odds_history table
- Used for analytics and charts
- Timestamp-indexed for efficient queries

Please provide:
1. OddsService struct with all calculation methods
2. MarketsRepository with CRUD operations
3. HTTP handlers for all endpoints listed
4. Request/Response DTOs for odds and markets
5. Validation:
   - Initial odds must be >= 1.01
   - Market type must be valid
   - Cannot modify locked markets
6. Database transactions for atomic updates
7. Error handling for edge cases:
   - Division by zero
   - Negative stakes
   - Locked market modifications
8. Unit tests for odds calculation with sample data
```

---

### Part 2.4: WebSocket Server for Real-Time Odds

```
Create a WebSocket server for broadcasting live odds updates to all connected clients:

Technology: gorilla/websocket (Go)

Features:

1. **WebSocket Connection Manager**:
```go
type Client struct {
    ID     string
    Conn   *websocket.Conn
    Send   chan []byte
    Rooms  map[string]bool  // Subscribe to specific events/markets
}

type Hub struct {
    Clients    map[*Client]bool
    Broadcast  chan Message
    Register   chan *Client
    Unregister chan *Client
    Rooms      map[string]map[*Client]bool
}
```

2. **Message Types**:
```go
type Message struct {
    Type    string      `json:"type"`    // "odds_update", "market_status", "bet_placed"
    EventID string      `json:"eventId"`
    Data    interface{} `json:"data"`
}

type OddsUpdateMessage struct {
    MarketID  string  `json:"marketId"`
    Outcome   string  `json:"outcome"`
    OldOdds   float64 `json:"oldOdds"`
    NewOdds   float64 `json:"newOdds"`
    Change    float64 `json:"change"`     // Percentage change
    Direction string  `json:"direction"`  // "up" or "down"
    Timestamp int64   `json:"timestamp"`
}
```

3. **Room-Based Broadcasting**:
- Clients subscribe to specific event rooms: `/ws?eventId=xxx`
- Only receive updates for subscribed events
- Reduces bandwidth usage

4. **Broadcast Triggers**:
```go
// After odds calculation
func BroadcastOddsUpdate(marketID, outcome string, oldOdds, newOdds float64) {
    eventID := GetEventIDFromMarket(marketID)
    
    message := OddsUpdateMessage{
        MarketID:  marketID,
        Outcome:   outcome,
        OldOdds:   oldOdds,
        NewOdds:   newOdds,
        Change:    ((newOdds - oldOdds) / oldOdds) * 100,
        Direction: if newOdds > oldOdds { "up" } else { "down" },
        Timestamp: time.Now().Unix(),
    }
    
    hub.BroadcastToRoom(eventID, "odds_update", message)
}
```

5. **Connection Handling**:
- Client connects: `ws://localhost:8080/ws?eventId={id}`
- Server validates event exists
- Add client to event room
- Send current odds snapshot
- Listen for odds updates
- Broadcast to all clients in room
- Handle disconnects gracefully

6. **Integration with Betting Flow**:
```go
// In bet placement service
func PlaceBet(...) {
    // ... bet placement logic ...
    
    // Update stakes
    UpdateStakeOnBet(oddID, stake)
    
    // Recalculate odds
    newOdds := CalculateOdds(marketID)
    
    // Broadcast to all connected clients
    BroadcastOddsUpdate(marketID, outcome, oldOdds, newOdds)
}
```

Please provide:
1. WebSocket hub implementation with room support
2. Client connection handler
3. Message broadcasting logic
4. Room subscription/unsubscription
5. Integration with odds calculation service
6. Heartbeat/ping-pong for connection health
7. Error handling and reconnection logic
8. WebSocket endpoint: `/ws`
9. CORS configuration for WebSocket
10. Example client test (can be simple HTML/JS)
```

---

### Part 2.5: Betting Engine with Odds Locking

```
Create the betting engine that handles bet placement with ODDS LOCKING:

Key Concept: When a user places a bet, they lock in the current odds. Even if odds change immediately after, their bet uses the locked odds.

Core Features:

1. **Bet Placement Flow**:
```go
func PlaceBet(userID uuid.UUID, selections []Selection, stake float64, betType string) (*Bet, error) {
    // Begin database transaction (SERIALIZABLE isolation level)
    tx := BeginTransaction()
    defer tx.Rollback()
    
    // Step 1: Validate user balance
    user := GetUser(userID)
    if user.Balance < stake {
        return nil, errors.New("insufficient balance")
    }
    
    // Step 2: Lock current odds for each selection
    lockedOdds := []LockedOdd{}
    for _, sel := range selections {
        odd := GetOdd(sel.OddID)
        market := GetMarket(odd.MarketID)
        
        // Check market is still open
        if market.Status != "open" {
            return nil, errors.New("market is locked")
        }
        
        // Lock the current odds value
        lockedOdds = append(lockedOdds, LockedOdd{
            OddID:     sel.OddID,
            MarketID:  odd.MarketID,
            Outcome:   odd.Outcome,
            LockedValue: odd.CurrentOdds,  // Lock current odds
        })
    }
    
    // Step 3: Calculate potential return
    potentialReturn := CalculatePotentialReturn(stake, lockedOdds, betType)
    
    // Step 4: Create bet record
    bet := CreateBet(userID, betType, stake, potentialReturn, "pending")
    
    // Step 5: Create bet legs with locked odds
    for _, locked := range lockedOdds {
        CreateBetLeg(bet.ID, locked.MarketID, locked.OddID, locked.Outcome, locked.LockedValue)
    }
    
    // Step 6: Deduct stake from user balance
    UpdateUserBalance(userID, -stake)
    
    // Step 7: Create transaction record
    CreateTransaction(userID, "bet_placed", stake, bet.ID)
    
    // Step 8: Update odds total_stake
    for _, locked := range lockedOdds {
        UpdateOddsTotalStake(locked.OddID, stake / len(lockedOdds))  // Divide stake among selections
    }
    
    // Commit transaction
    tx.Commit()
    
    // Step 9: AFTER commit, recalculate odds and broadcast
    go func() {
        for _, locked := range lockedOdds {
            marketID := locked.MarketID
            CalculateOdds(marketID)
            BroadcastOddsUpdate(marketID)  // Notify WebSocket clients
        }
    }()
    
    return bet, nil
}
```

2. **Potential Return Calculation**:
```go
func CalculatePotentialReturn(stake float64, lockedOdds []LockedOdd, betType string) float64 {
    if betType == "single" {
        // Single bet: stake × locked_odds
        return stake * lockedOdds[0].LockedValue
    } else {
        // Accumulator: stake × (odds1 × odds2 × ... × oddsN)
        combinedOdds := 1.0
        for _, odd := range lockedOdds {
            combinedOdds *= odd.LockedValue
        }
        return stake * combinedOdds
    }
}
```

3. **Bet Settlement**:
```go
func SettleBet(betID uuid.UUID, results map[string]string) error {
    bet := GetBet(betID)
    betLegs := GetBetLegs(betID)
    
    allWon := true
    for _, leg := range betLegs {
        // Use the LOCKED odds value, not current odds
        if results[leg.MarketID] == leg.Outcome {
            leg.Result = "won"
        } else {
            leg.Result = "lost"
            allWon = false
        }
        UpdateBetLeg(leg)
    }
    
    if allWon {
        // Calculate payout using LOCKED odds
        payout := bet.PotentialReturn
        
        // Update user balance
        UpdateUserBalance(bet.UserID, payout)
        
        // Create transaction
        CreateTransaction(bet.UserID, "bet_won", payout, bet.ID)
        
        // Update bet
        UpdateBet(bet.ID, "won", payout)
    } else {
        UpdateBet(bet.ID, "lost", 0)
    }
    
    return nil
}
```

4. **Bet History & Tracking**:
- GET /api/bets/history - User's bet history with filters
- GET /api/bets/:id - Detailed bet information with locked odds
- GET /api/bets/pending - Active bets

5. **Validation Rules**:
- Minimum stake: $0.50
- Maximum stake: configurable per user
- Cannot bet on locked/settled markets
- Cannot combine incompatible markets
- Accumulator requires 2+ selections

Database tables: bets, bet_legs, transactions

Please provide:
1. Betting service with PlaceBet and SettleBet functions
2. Potential return calculator
3. HTTP handlers:
   - POST /api/bets/place
   - GET /api/bets/history
   - GET /api/bets/:id
   - POST /api/admin/bets/:id/settle (admin only)
4. Request DTOs:
   - PlaceBetRequest{Selections, Stake, BetType}
   - SettleBetRequest{BetID, Results}
5. Response DTOs:
   - BetResponse with locked odds displayed
   - BetHistoryResponse with pagination
6. Transaction handling with SERIALIZABLE isolation
7. Race condition prevention (optimistic locking)
8. Error handling for all edge cases
9. Unit tests for bet placement and settlement
```

---

## Phase 3: Frontend Foundation

### Part 3.1: React Project Setup with WebSocket

```
Set up a React frontend with TypeScript and WebSocket support for real-time odds:

Requirements:
- React 18 with Vite
- TypeScript strict mode
- React Router v6 for navigation
- TailwindCSS for styling
- Axios for REST API calls
- React Query (TanStack Query) for data fetching and caching
- Zustand for global state management
- WebSocket client for real-time odds updates
- Sonner for toast notifications

Project structure:
```
frontend/
├── src/
│   ├── components/
│   │   ├── layout/
│   │   │   ├── Header.tsx
│   │   │   ├── Footer.tsx
│   │   │   ├── Sidebar.tsx
│   │   │   └── Layout.tsx
│   │   ├── auth/
│   │   ├── events/
│   │   ├── betting/
│   │   └── common/
│   ├── pages/
│   │   ├── HomePage.tsx
│   │   ├── LoginPage.tsx
│   │   ├── RegisterPage.tsx
│   │   ├── EventDetailsPage.tsx
│   │   ├── ProfilePage.tsx
│   │   ├── BetHistoryPage.tsx
│   │   └── admin/
│   ├── services/
│   │   ├── api.ts           // Axios instance
│   │   ├── auth.service.ts
│   │   ├── events.service.ts
│   │   ├── bets.service.ts
│   │   └── websocket.service.ts  // WebSocket manager
│   ├── stores/
│   │   ├── authStore.ts
│   │   ├── betslipStore.ts
│   │   └── oddsStore.ts     // Real-time odds state
│   ├── hooks/
│   │   ├── useAuth.ts
│   │   ├── useWebSocket.ts
│   │   └── useOdds.ts       // Real-time odds hook
│   ├── types/
│   │   ├── auth.types.ts
│   │   ├── events.types.ts
│   │   ├── odds.types.ts
│   │   └── bets.types.ts
│   ├── utils/
│   │   ├── formatters.ts
│   │   ├── validators.ts
│   │   └── constants.ts
│   ├── App.tsx
│   ├── main.tsx
│   └── routes.tsx
├── package.json
├── tsconfig.json
├── vite.config.ts
└── tailwind.config.js
```

Routes needed:
- / - Home (events list)
- /login - Login page
- /register - Registration page
- /events/:id - Event details with betting
- /profile - User profile and balance
- /bets/history - Bet history
- /admin - Admin dashboard (protected)
- /admin/events - Event management
- /admin/odds - Odds monitoring

Please provide:
1. Complete Vite + React + TypeScript setup
2. Router configuration with protected routes
3. API client (Axios) with:
   - Base URL from env
   - Auth token interceptor
   - Error handling interceptor
4. WebSocket service:
   - Connection manager
   - Room subscription
   - Message handlers
   - Auto-reconnect logic
5. Layout components (Header, Footer, Layout)
6. Auth store (Zustand) with login/logout/token management
7. Odds store (Zustand) for real-time odds state
8. Protected route wrapper (RequireAuth)
9. TailwindCSS configuration with custom colors
10. Package.json with all dependencies
```

---

### Part 3.2: Authentication UI

```
Build authentication pages with modern, clean design:

Pages needed:
1. Login page
2. Registration page

Features:
- Email/password input with validation
- Password visibility toggle
- Form error display
- Loading states (disabled button with spinner)
- Remember me checkbox (login)
- Password strength indicator (registration)
- Success/error toast notifications
- Redirect to home after successful auth
- Link to switch between login/register

Design requirements:
- Clean, minimal design with TailwindCSS
- Centered form card with shadow
- Mobile-responsive (full-width on mobile)
- Accessibility (ARIA labels, keyboard navigation, focus states)
- Form validation:
  - Email format validation
  - Password min 8 characters
  - Password confirmation match
  - Display errors below inputs

Components to create:
- LoginPage.tsx
- RegisterPage.tsx
- Input.tsx (reusable input component with error state)
- Button.tsx (reusable button with loading state)
- PasswordStrengthBar.tsx

State management:
- Use React Query for API calls
- Store JWT in localStorage
- Update Zustand auth store on success
- Axios interceptor adds token to requests

API integration:
- POST /api/auth/register
- POST /api/auth/login
- Store access_token and refresh_token

Please provide:
1. LoginPage component with form handling
2. RegisterPage component with validation
3. Reusable Input component with error display
4. Reusable Button component with loading spinner
5. Password strength indicator component
6. Form validation hooks (useLoginForm, useRegisterForm)
7. Auth service with API calls
8. Auth store (Zustand) with:
   - login(email, password)
   - register(email, password)
   - logout()
   - isAuthenticated getter
9. Protected route HOC (RequireAuth.tsx)
10. Toast notifications for success/error
```

---

### Part 3.3: Events Listing Page

```
Create the main events listing page for browsing upcoming sports events:

Features:
- Display list of upcoming events
- Group events by sport category
- Sport filter tabs (All, Football, Basketball, Tennis, etc.)
- Search bar to filter by team name
- Date range filter (Today, Tomorrow, This Week, Custom)
- Event cards showing:
  - Sport icon/badge
  - Team names (Home vs Away)
  - Start date and time
  - Number of available markets
  - "View Markets" button
- Click event card → navigate to /events/:id
- Loading skeletons while fetching
- Empty state if no events found
- Pagination or infinite scroll
- Real-time updates (new events appear, started events move to "Live")

Layout:
- Grid layout (1 column mobile, 2-3 columns desktop)
- Sport category tabs at top
- Search and filters in sidebar (desktop) or collapsible (mobile)
- Responsive design

Data fetching:
- Use React Query for caching
- Auto-refetch every 60 seconds
- Optimistic updates

Components to create:
- EventsListPage.tsx (main page)
- EventCard.tsx (individual event card)
- SportFilterTabs.tsx (filter by sport)
- SearchBar.tsx (search events)
- DateRangeFilter.tsx (filter by date)
- EventCardSkeleton.tsx (loading state)
- EmptyState.tsx (no events found)

API integration:
- GET /api/events?sport=&search=&date_from=&date_to=&page=&limit=

Please provide:
1. EventsListPage with filters and search
2. EventCard component with responsive design
3. Sport filter tabs component
4. Search and date filter components
5. Loading skeleton components
6. Empty state component
7. React Query hooks (useEvents, useEventsFiltered)
8. Event service API functions
9. Responsive grid layout with TailwindCSS
10. Pagination component (if not using infinite scroll)
```

---

### Part 3.4: Event Details with Real-Time Odds

```
Build the event details page with LIVE ODDS UPDATES and betting interface:

This is the most critical page - users see markets, live odds, and place bets.

Features:

1. **Event Header**:
   - Team names with logos/icons
   - Sport badge
   - Start date and time (or "LIVE" if in progress)
   - Event status badge (Upcoming, Live, Finished)

2. **Markets Display**:
   - Group markets by type:
     - Match Winner (1X2)
     - Over/Under
     - Both Teams to Score
     - etc.
   - Each market shows:
     - Market name
     - All outcomes with current odds
     - Odds as clickable buttons

3. **Real-Time Odds with WebSocket**:
   - Connect to WebSocket on page load: `ws://localhost:8080/ws?eventId={id}`
   - Subscribe to odds updates for this event
   - When odds change:
     - Update odds value immediately
     - Show movement indicator (↑ green or ↓ red)
     - Flash/pulse animation on change
     - Display percentage change (e.g., "2.50 ↑ +5%")
   - Connection status indicator (connected/disconnected)

4. **Odds Button Component**:
```tsx
<OddsButton 
  outcome="Home" 
  odds={2.50} 
  previousOdds={2.38}  // For movement calculation
  onChange={(change) => ...}
  selected={isInBetslip}
  onClick={() => addToBetslip(...)}
/>
```

Display:
- Large odds value
- Outcome name
- Movement indicator if changed recently
- Highlight if selected (in betslip)
- Disabled if market locked

5. **Betslip Widget** (Sticky Sidebar):

Desktop: Fixed right sidebar
Mobile: Bottom sheet (swipe up to expand)

Contents:
- List of selected bets
- Each selection shows:
  - Event name
  - Outcome
  - Current odds (live updating)
  - "Odds changed" badge if odds moved since adding
- Stake input for each selection (single bets)
- Or single stake input (accumulators)
- Total stake display
- Potential return calculation (live updating)
- "Place Bet" button
- Clear betslip button
- Toggle: Single bets / Accumulator

6. **Odds Change Handling**:
```tsx
// When odds change via WebSocket while in betslip
if (oddsIncreasedByMore than 5%) {
  showWarning("Odds improved! New potential return: $X")
} else if (oddsDecreasedByMoreThan 5%) {
  showWarning("Odds dropped. Accept new odds or remove bet.")
}
```

- Show accept/reject buttons for significant changes
- Auto-update potential return

7. **Bet Placement Flow**:
```
User clicks odds → Add to betslip
User enters stake → Calculate potential return
User clicks "Place Bet" → Show confirmation modal
Modal shows:
  - Current odds (locked at this moment)
  - If odds changed since adding, highlight it
  - Final potential return
  - "Confirm" button
User confirms → API call to place bet
Success → Show success toast, clear betslip
```

8. **Market Info Display**:
- Total pool amount (how much money in market)
- Distribution bar: visual showing stake % per outcome
- Example: [Home: 45% | Draw: 20% | Away: 35%]
- Live betting volume indicator

9. **Components to Create**:
- EventDetailsPage.tsx (main page)
- EventHeader.tsx (event info)
- MarketsList.tsx (all markets grouped)
- MarketCard.tsx (single market with outcomes)
- OddsButton.tsx (clickable odds with movement indicators)
- BetslipWidget.tsx (sidebar with selections)
- BetslipItem.tsx (single bet in betslip)
- StakeInput.tsx (input with validation)
- PotentialReturnDisplay.tsx (calculated return)
- BetConfirmationModal.tsx (final confirmation)
- StakeDistributionBar.tsx (visual stake distribution)
- WebSocketStatus.tsx (connection indicator)

10. **WebSocket Integration**:
```tsx
const { isConnected, subscribe, unsubscribe } = useWebSocket();

useEffect(() => {
  // Connect to event room
  subscribe(`event:${eventId}`, (message) => {
    if (message.type === 'odds_update') {
      updateOdds(message.data);  // Update Zustand store
      showOddsChangeAnimation(message.data.outcome);
    }
  });
  
  return () => unsubscribe(`event:${eventId}`);
}, [eventId]);
```

11. **State Management** (Zustand):
```tsx
// oddsStore.ts
interface OddsStore {
  odds: Record<string, Odd>;  // Key: oddId
  updateOdd: (oddId, newValue, change) => void;
  getOddsMovement: (oddId) => { direction, percentage };
}

// betslipStore.ts
interface BetslipStore {
  selections: Selection[];
  addSelection: (selection) => void;
  removeSelection: (oddId) => void;
  clearBetslip: () => void;
  totalStake: number;
  potentialReturn: number;
  betType: 'single' | 'accumulator';
}
```

12. **API Integration**:
- GET /api/events/:id (event details)
- GET /api/events/:id/markets (markets and odds)
- POST /api/bets/place (place bet)
- WebSocket: ws://localhost:8080/ws?eventId={id}

Please provide:
1. EventDetailsPage with all components
2. Real-time odds updates via WebSocket
3. OddsButton with movement indicators
4. Betslip widget (responsive: sidebar on desktop, bottom sheet on mobile)
5. Odds change detection and warnings
6. Bet confirmation modal
7. Stake distribution visualization
8. useWebSocket hook for managing connections
9. useOdds hook for odds state and updates
10. useBetslip hook for betslip management
11. Potential return calculator
12. Complete responsive design
13. Loading and error states
14. Animations for odds changes (pulse, fade)
```

---

## Phase 4: User Features

### Part 4.1: User Profile & Balance Management

```
Create user profile and balance management features:

Features:

1. **Profile Page**:
   - Display user info:
     - Email
     - Account created date
     - Current balance (large, prominent)
   - Deposit form
   - Withdraw form
   - Recent transactions (last 10)
   - Link to full transaction history

2. **Deposit Flow**:
   - Amount input with presets ($10, $25, $50, $100, $250, Custom)
   - Payment method selector (mock for MVP: "Credit Card", "PayPal")
   - Confirm button
   - Success toast
   - Balance updates immediately

3. **Withdraw Flow**:
   - Amount input
   - Available balance display
   - Minimum withdrawal validation ($10)
   - Cannot withdraw more than balance
   - Processing message (mock: instant for MVP)
   - Success toast
   - Balance updates immediately

4. **Transaction History**:
   - Table/list showing:
     - Type (Deposit, Withdraw, Bet Placed, Bet Won, Bet Lost)
     - Amount (+ green for gains, - red for losses)
     - Balance after transaction
     - Date/time
     - Reference (bet ID if applicable)
   - Filter by type
   - Filter by date range
   - Pagination
   - Export to CSV button (optional)

5. **Balance Display in Header**:
   - Show current balance in navbar
   - Live update when transactions occur
   - Click to navigate to profile

Components:
- ProfilePage.tsx (main page)
- BalanceCard.tsx (display balance prominently)
- DepositForm.tsx (deposit interface)
- WithdrawForm.tsx (withdraw interface)
- TransactionList.tsx (recent transactions)
- TransactionItem.tsx (single transaction row)
- AmountPresets.tsx (preset amount buttons)

API endpoints:
- GET /api/users/profile
- GET /api/users/balance
- POST /api/transactions/deposit
- POST /api/transactions/withdraw
- GET /api/transactions/history?type=&page=&limit=

Please provide:
1. Profile page with balance and forms
2. Deposit form with amount presets
3. Withdraw form with validation
4. Transaction list component
5. Balance display in header (live updating)
6. React Query hooks for profile and transactions
7. Form validation (min/max amounts)
8. Success/error handling
9. Optimistic updates for balance
10. Responsive design
```

---

### Part 4.2: Bet History Page

```
Build comprehensive bet history page showing user's past and pending bets:

Features:

1. **Bet Filters**:
   - Status filter: All, Pending, Won, Lost, Cancelled
   - Date range filter
   - Bet type filter: All, Single, Accumulator
   - Sort by: Date (newest/oldest), Stake (high/low)

2. **Bet List Display**:

Desktop: Table view
Mobile: Card view

Each bet shows:
- Bet ID (unique identifier)
- Date placed
- Bet type badge (Single / Accumulator)
- Stake amount
- Status badge (color-coded):
  - Pending: Yellow
  - Won: Green
  - Lost: Red
  - Cancelled: Gray
- Potential return (for pending/won)
- Actual return (for won bets)
- Expand/collapse button for details

3. **Expandable Bet Details**:

When expanded, show:
- All selections (bet legs):
  - Event name (Team A vs Team B)
  - Market type
  - Outcome selected
  - Locked odds (the odds when bet was placed)
  - Result (Won/Lost/Pending)
- Total odds (for accumulator)
- Stake
- Potential return
- Actual return (if settled)
- Placed at timestamp
- Settled at timestamp (if settled)

4. **Statistics Summary**:
- Total bets placed
- Win rate (%)
- Total wagered
- Total returns
- Net profit/loss
- Biggest win

5. **Empty States**:
- No bets yet: "Place your first bet!"
- No results for filter: "No bets match your filters"

6. **Pagination**:
- Load 20 bets per page
- "Load more" button or traditional pagination

Components:
- BetHistoryPage.tsx (main page)
- BetFilters.tsx (filter controls)
- BetTable.tsx (desktop table view)
- BetCard.tsx (mobile card view)
- BetDetailsModal.tsx (expanded bet details, alternative to inline expand)
- BetStatistics.tsx (summary stats)
- BetStatusBadge.tsx (status indicator)
- EmptyBetsState.tsx

API endpoints:
- GET /api/bets/history?status=&betType=&dateFrom=&dateTo=&page=&limit=
- GET /api/bets/:id (detailed bet info)
- GET /api/bets/statistics (summary stats)

Please provide:
1. Bet history page with filters
2. Table view (desktop) and card view (mobile)
3. Expandable bet details showing all legs
4. Status badges with colors
5. Statistics summary component
6. Filter and sort functionality
7. Pagination component
8. Empty states
9. React Query hooks for data fetching
10. Responsive design
11. Display locked odds (odds at time of bet)
```

---

## Phase 5: Admin Panel

### Part 5.1: Admin Dashboard Overview

```
Create an admin dashboard for monitoring and managing the sportsbook:

Layout:
- Separate admin layout with sidebar navigation
- Admin-only access (role-based auth)
- Breadcrumbs for navigation

Dashboard Features:

1. **Key Metrics** (Top Cards):
   - Total Users
   - Active Bets (pending)
   - Total Revenue Today
   - Total Liability (worst-case payout)

2. **Recent Activity Feed**:
   - Last 10 bets placed (user, event, stake, status)
   - Last 5 user registrations
   - Last 5 withdrawals/deposits

3. **Charts** (Optional for MVP, can use mock data):
   - Betting volume over time (line chart)
   - Revenue by sport (pie chart)
   - Win/loss ratio

4. **Quick Actions**:
   - "Create Event" button → navigate to event form
   - "Settle Bets" button → navigate to settlement page
   - "View Odds Monitor" → navigate to odds dashboard

5. **System Health**:
   - Database connection status
   - WebSocket connection status
   - Number of connected users
   - API response time (average)

Sidebar Navigation:
- Dashboard (home)
- Events Management
- Odds Monitor
- Bet Settlement
- User Management
- Reports

Components:
- AdminLayout.tsx (layout with sidebar)
- AdminSidebar.tsx (navigation)
- AdminDashboard.tsx (main dashboard)
- MetricCard.tsx (stat display card)
- ActivityFeed.tsx (recent activity list)
- QuickActionsPanel.tsx (action buttons)

API endpoints:
- GET /api/admin/stats (summary statistics)
- GET /api/admin/activity (recent activity feed)

Please provide:
1. Admin layout with sidebar
2. Dashboard page with metrics
3. Metric cards component
4. Activity feed component
5. Quick actions panel
6. Admin route protection (require admin role)
7. Responsive design
8. Mock data for charts (if implementing)
```

---

### Part 5.2: Event & Odds Management Panel

```
Build comprehensive event and odds management interface for admins:

Section 1: **Events Management**

Features:
- List all events (upcoming, live, finished)
- Filter by sport, status, date
- Search by team name
- Table view with columns:
  - ID
  - Sport
  - Home vs Away
  - Start Time
  - Status
  - Markets Count
  - Total Bets
  - Actions (Edit, Delete, Manage Odds)

CRUD Operations:
- Create Event:
  - Form with fields:
    - Sport (dropdown)
    - Home Team (text input)
    - Away Team (text input)
    - Start Date/Time (datetime picker)
  - Validation: start time must be future
  - Submit → API call
  - Success → refresh list

- Edit Event:
  - Same form, pre-filled with data
  - Can update teams, start time
  - Cannot change sport
  - Cannot edit if event started

- Delete Event:
  - Confirmation modal: "Delete event? This will cancel all associated bets."
  - Only if no bets placed yet
  - Soft delete (status = cancelled)

- Change Event Status:
  - Dropdown: Upcoming → Live → Finished
  - Auto-lock all markets when status = Live
  - Cannot place bets when status = Finished

Section 2: **Odds Management**

Features:
- Select event → show all markets
- For each market, show all odds
- Set initial odds (when market created)
- View current odds (calculated dynamically)
- Manual override option
- Monitor total stakes and bet counts

Create Market:
- Form:
  - Event (dropdown, pre-selected if coming from event page)
  - Market Type (dropdown: 1X2, Over/Under, Both Teams Score, etc.)
  - Initial Odds for each outcome:
    - For 1X2: Home odds, Draw odds, Away odds
    - For Over/Under: Over odds, Under odds, Line (e.g., 2.5 goals)
  - House Margin (default 5%)
- Validation: odds must be >= 1.01
- Submit → create market and odds

Odds Table View:
Columns:
- Market Type
- Outcome
- Initial Odds (set by admin)
- Current Odds (calculated)
- Total Stake
- Bet Count
- Status (Open/Locked)
- Actions (Edit, Lock/Unlock, Override)

Manual Override:
- Click "Override" on an odd
- Modal with input: new odds value
- Checkbox: "Lock automatic calculation"
- Submit → update odd, set is_manual_override = true
- Admin can resume automatic calculation later

Lock/Unlock Market:
- Lock: prevent new bets, keep odds frozen
- Unlock: allow bets, resume odds calculation
- Bulk action: lock all markets for an event

Odds Monitoring Dashboard:
- Real-time view of all active markets
- Show odds movement (compare current vs 1 hour ago)
- Alert if odds moved >10% (potential issue)
- Show liability per market (worst-case payout)
- Graph: odds movement over time

Components:
- EventsManagementPage.tsx
- EventsTable.tsx
- EventForm.tsx (create/edit modal)
- DeleteEventConfirmation.tsx
- OddsManagementPage.tsx
- MarketsTable.tsx
- CreateMarketForm.tsx
- OddsTable.tsx
- ManualOddsOverrideModal.tsx
- OddsMonitorDashboard.tsx
- OddsMovementChart.tsx (line chart)
- LiabilityCalculator.tsx

API endpoints:
- GET /api/admin/events
- POST /api/admin/events
- PUT /api/admin/events/:id
- DELETE /api/admin/events/:id
- PATCH /api/admin/events/:id/status

- GET /api/admin/markets?eventId=
- POST /api/admin/markets
- PUT /api/admin/markets/:id
- DELETE /api/admin/markets/:id

- GET /api/admin/markets/:id/odds
- POST /api/admin/odds (set initial odds)
- PUT /api/admin/odds/:id (manual override)
- PATCH /api/admin/odds/:id/lock
- PATCH /api/admin/odds/:id/unlock

- GET /api/admin/odds/monitor (real-time odds data)
- GET /api/admin/odds/:id/history (odds movement over time)

Please provide:
1. Events management page with CRUD
2. Event form component (create/edit)
3. Events table with filters
4. Delete confirmation modal
5. Odds management page
6. Create market form with initial odds inputs
7. Odds table showing initial vs current odds
8. Manual override modal
9. Lock/unlock market functionality
10. Odds monitoring dashboard with real-time data
11. Odds movement chart
12. Liability calculator
13. Form validation
14. Success/error notifications
15. Responsive design
```

---

### Part 5.3: Bet Settlement System

```
Create comprehensive bet settlement interface for admins:

Purpose: After an event finishes, admin enters final results and system automatically settles all related bets.

Features:

1. **Unsettled Events List**:
   - Show all finished events that have unsettled bets
   - Display:
     - Event name
     - Finished time
     - Number of unsettled bets
     - Total liability (potential payout)
   - Sort by finish time
   - Click event → go to settlement page

2. **Settlement Page for Event**:

Display:
- Event details (teams, sport, final score)
- All markets for this event
- For each market:
  - Market type
  - All outcomes
  - Number of bets on each outcome
  - Total stake on each outcome

Result Entry Form:
- For Match Winner (1X2):
  - Dropdown: Home Win / Draw / Away Win
- For Over/Under:
  - Input: Total goals/points
  - System calculates: Over or Under
- For Both Teams Score:
  - Radio: Yes / No

Submit Results:
- Button: "Settle Bets"
- Confirmation modal: "Settle X bets? This cannot be undone."
- On confirm:
  - API call with results
  - Backend:
    - Determines winning outcome for each market
    - Finds all bets with selections on this event
    - For each bet:
      - Check if all legs won (for accumulators)
      - If won: calculate payout using LOCKED odds
      - Update user balance
      - Update bet status (won/lost)
      - Create transaction records
  - Frontend:
    - Show success message
    - Display settlement summary:
      - Total bets settled
      - Winning bets count
      - Total payout amount
    - Remove event from unsettled list

3. **Settlement Summary**:

After settlement, show:
- Number of bets settled
- Winning bets: X
- Losing bets: Y
- Total paid out: $Z
- Net profit for house: $W

4. **Manual Settlement Override**:

Sometimes results are disputed or corrected. Admin can:
- Re-open settled bets
- Manually set bet result (won/lost)
- Reverse payouts
- Add notes explaining override

Use cases:
- VAR decision in football changes result
- Score correction
- Event cancelled after bets placed

5. **Settlement History**:

- List of all settled events
- Show settlement date, admin who settled, results
- View details of each settlement
- Filter by sport, date range

Components:
- UnsettledEventsPage.tsx (list of events to settle)
- UnsettledEventCard.tsx
- EventSettlementPage.tsx (main settlement page)
- ResultEntryForm.tsx (input final results)
- SettlementConfirmationModal.tsx
- SettlementSummary.tsx (results after settlement)
- ManualSettlementOverride.tsx (for corrections)
- SettlementHistory.tsx (past settlements)

Backend Logic (for reference):
```go
func SettleEventBets(eventID, results map[string]string) (SettlementSummary, error) {
    // Get all markets for event
    markets := GetMarketsByEvent(eventID)
    
    settlementSummary := SettlementSummary{}
    
    for _, market := range markets {
        // Determine winning outcome
        winningOutcome := DetermineWinner(market, results)
        
        // Get all bets on this market
        bets := GetBetsForMarket(market.ID)
        
        for _, bet := range bets {
            betLegs := GetBetLegs(bet.ID)
            
            // Check this market's leg
            for _, leg := range betLegs {
                if leg.MarketID == market.ID {
                    if leg.Outcome == winningOutcome {
                        leg.Result = "won"
                    } else {
                        leg.Result = "lost"
                    }
                    UpdateBetLeg(leg)
                }
            }
            
            // Check if entire bet won (all legs must win for accumulators)
            allLegsWon := true
            for _, leg := range betLegs {
                if leg.Result != "won" {
                    allLegsWon = false
                    break
                }
            }
            
            if allLegsWon {
                // Calculate payout using LOCKED odds from bet_legs
                payout := CalculatePayoutFromLockedOdds(bet.ID)
                
                // Update user balance
                UpdateUserBalance(bet.UserID, payout)
                
                // Update bet
                UpdateBet(bet.ID, "won", payout)
                
                // Create transaction
                CreateTransaction(bet.UserID, "bet_won", payout, bet.ID)
                
                settlementSummary.WonBetsCount++
                settlementSummary.TotalPayout += payout
            } else {
                UpdateBet(bet.ID, "lost", 0)
                settlementSummary.LostBetsCount++
            }
            
            settlementSummary.TotalBetsSettled++
        }
        
        // Mark market as settled
        UpdateMarket(market.ID, "settled")
    }
    
    // Mark event as settled
    UpdateEvent(eventID, "settled")
    
    return settlementSummary, nil
}
```

API endpoints:
- GET /api/admin/events/unsettled (list of events needing settlement)
- GET /api/admin/events/:id/settlement-preview (show bets before settlement)
- POST /api/admin/events/:id/settle (settle all bets, provide results)
- POST /api/admin/bets/:id/override (manual settlement override)
- GET /api/admin/settlements/history (past settlements)

Please provide:
1. Unsettled events list page
2. Event settlement page with result entry
3. Result entry form (different inputs for different market types)
4. Settlement confirmation modal
5. Settlement summary display
6. Manual override interface
7. Settlement history page
8. Backend settlement service (Go)
9. Payout calculation using locked odds
10. Transaction creation for payouts
11. Success/error handling
12. Audit logging for settlements
```

---

## Phase 6: Polish & Testing

### Part 6.1: Error Handling & UX Polish

```
Improve error handling and user experience across the entire application:

Tasks:

1. **Global Error Boundary** (React):
   - Catch React errors
   - Display friendly error page
   - "Go to Home" button
   - Log errors to console (or error tracking service)

2. **API Error Handling**:
   - Axios interceptor for errors:
     - 401 Unauthorized → logout user, redirect to login
     - 403 Forbidden → show "Access Denied" message
     - 404 Not Found → show "Resource not found"
     - 500 Server Error → show "Something went wrong, try again"
     - Network Error → show "Check your connection"
   - Display errors as toast notifications
   - Form-specific errors → show below inputs

3. **Toast Notification System**:
   - Success toasts:
     - "Bet placed successfully!"
     - "Deposit successful!"
     - "Event created!"
   - Error toasts:
     - "Insufficient balance"
     - "Market is locked"
     - "Failed to place bet"
   - Info toasts:
     - "Odds changed for your selection"
     - "Event starting soon"
   - Use Sonner library for toasts
   - Position: top-right (desktop), top-center (mobile)
   - Auto-dismiss after 4 seconds
   - Swipe to dismiss on mobile

4. **Loading States**:
   - Button loading states (spinner + disabled)
   - Page loading states (skeleton screens)
   - Infinite scroll loading indicator
   - Global loading bar (top of page)
   - WebSocket connection loading

5. **Empty States**:
   - No events: "No events available. Check back later!"
   - No bets: "You haven't placed any bets yet."
   - No transactions: "No transactions found."
   - No search results: "No events match your search."
   - Add illustration or icon to empty states

6. **Form Validation Improvements**:
   - Real-time validation (on blur)
   - Clear error messages
   - Disable submit until valid
   - Show validation icon (✓ or ✗)

7. **Confirmation Modals**:
   - Delete event: "Are you sure? This cannot be undone."
   - Settle bets: "Settle X bets? This is final."
   - Logout: "Are you sure you want to log out?"
   - Cancel bet: "Remove this selection from betslip?"

8. **Retry Logic**:
   - Failed API requests → show "Retry" button
   - WebSocket disconnected → auto-reconnect with exponential backoff
   - Failed bet placement → allow retry

9. **404 and Error Pages**:
   - 404 Page: "Page not found" with link to home
   - 500 Error Page: "Something went wrong" with retry button
   - 403 Access Denied: "You don't have permission" with link to home

10. **Accessibility Improvements**:
    - All buttons have aria-labels
    - Form inputs have labels
    - Error messages have aria-live regions
    - Keyboard navigation works everywhere
    - Focus indicators visible

Components to create:
- ErrorBoundary.tsx (React error boundary)
- ErrorPage.tsx (500 error page)
- NotFoundPage.tsx (404 page)
- AccessDeniedPage.tsx (403 page)
- ConfirmationModal.tsx (reusable confirmation dialog)
- EmptyState.tsx (reusable empty state component)
- LoadingSpinner.tsx (reusable spinner)
- RetryButton.tsx (retry failed action)

Please provide:
1. Global error boundary
2. API error interceptor with user-friendly messages
3. Toast notification setup (Sonner)
4. Loading spinner and skeleton components
5. Empty state components with icons
6. Confirmation modal component
7. 404, 500, and 403 pages
8. Retry logic for failed requests
9. Form validation with clear error display
10. Accessibility attributes (aria-labels, aria-live)
```

---

### Part 6.2: Mobile Responsiveness & Animations

```
Ensure the entire application is mobile-responsive with smooth animations:

Focus Areas:

1. **Navigation** (Mobile):
   - Hamburger menu on mobile (< 768px)
   - Slide-in sidebar menu
   - Overlay backdrop when menu open
   - Close on backdrop click
   - Smooth open/close animation

2. **Betslip** (Mobile):
   - Bottom sheet on mobile (fixed to bottom)
   - Swipe up to expand
   - Shows summary when collapsed: "X bets, $Y stake"
   - Fully expanded when tapped
   - Smooth expand/collapse animation
   - Close button when expanded

3. **Tables → Cards** (Mobile):
   - Event tables → event cards
   - Bet history table → bet cards
   - Transaction table → transaction list items
   - Each card shows key info prominently

4. **Forms** (Mobile):
   - Full-width inputs
   - Larger touch targets (min 44px height)
   - Bottom sheets for date pickers
   - Numeric keyboard for amount inputs
   - Sticky submit button at bottom

5. **Modals** (Mobile):
   - Full-screen modals on mobile
   - Slide up from bottom animation
   - "X" close button top-right
   - Scrollable content

6. **Animations**:

Odds Changes:
- Pulse animation when odds update
- Green flash for odds increase
- Red flash for odds decrease
- Fade in/out (300ms)

Bet Placement:
- Confetti animation on successful bet (optional)
- Success checkmark animation
- Shake animation on error

Page Transitions:
- Fade in/out between pages
- Slide in from right for navigation forward
- Slide out to right for navigation back

Loading:
- Skeleton screens fade in
- Smooth spinner rotation
- Progress bar for multi-step forms

Interactions:
- Button press animation (scale down slightly)
- Hover effects on desktop
- Ripple effect on clicks (mobile)
- Smooth scrolling

7. **Responsive Breakpoints**:
   - Mobile: < 768px
   - Tablet: 768px - 1024px
   - Desktop: > 1024px

8. **Touch Optimizations**:
   - Minimum touch target size: 44x44px
   - Adequate spacing between interactive elements
   - Swipe gestures for betslip, modals
   - Pull-to-refresh on event list (optional)
   - Fast tap response (no 300ms delay)

9. **Testing Checklist**:
   - Test on common screen sizes:
     - iPhone SE (375px)
     - iPhone 14 (390px)
     - Pixel 7 (412px)
     - iPad (768px)
     - Desktop (1440px)
   - Test in portrait and landscape
   - Test with touch (mobile) and mouse (desktop)
   - Test keyboard navigation
   - Test with slow 3G network

Components to create or update:
- MobileNav.tsx (hamburger menu)
- MobileBetslip.tsx (bottom sheet)
- BottomSheet.tsx (reusable bottom sheet)
- ResponsiveTable.tsx (table that converts to cards on mobile)
- TouchOptimizedButton.tsx (44px min height)
- SwipeableCard.tsx (swipe actions)
- PullToRefresh.tsx (pull to refresh component)

Animations library: Framer Motion

Please provide:
1. Responsive navigation (hamburger menu on mobile)
2. Bottom sheet betslip for mobile
3. Table-to-card responsive components
4. Mobile-optimized forms (full-width, large inputs)
5. Full-screen modals on mobile
6. Animation components:
   - Pulse for odds changes
   - Fade transitions between pages
   - Loading skeleton animations
   - Button press animations
7. Swipe gesture handling
8. Touch optimization (large targets)
9. Media queries for all breakpoints
10. Testing guide for different screen sizes
11. TailwindCSS responsive utilities
12. Framer Motion animations setup
```

---

## Phase 7: Deployment & Production

### Part 7.1: Production Configuration & Deployment

```
Prepare the BetKZ application for production deployment:

Backend Production Setup:

1. **Environment Configuration**:
```env
# Production .env
NODE_ENV=production
PORT=8080
DATABASE_URL=postgresql://user:pass@db:5432/betkz
REDIS_URL=redis://redis:6379
JWT_SECRET=<strong-random-secret>
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=7d
CORS_ORIGINS=https://betkz.com
WEBSOCKET_ORIGINS=https://betkz.com
```

2. **Security Headers**:
   - CORS configuration (allow only production domain)
   - Helmet middleware for security headers
   - Rate limiting (100 requests per minute per IP)
   - CSRF protection
   - XSS protection headers
   - HSTS (HTTP Strict Transport Security)

3. **Database**:
   - Production PostgreSQL instance
   - Connection pooling (max 20 connections)
   - Automated backups (daily)
   - Migration strategy:
     - Run migrations before deployment
     - Rollback plan if deployment fails

4. **Logging**:
   - Structured JSON logging
   - Log levels: ERROR, WARN, INFO, DEBUG
   - Log to stdout (captured by container orchestration)
   - Error tracking (Sentry integration optional)
   - Request/response logging with correlation IDs

5. **Health Checks**:
   - GET /health - basic health check
   - GET /health/db - database connection
   - GET /health/redis - Redis connection
   - GET /metrics - Prometheus metrics (optional)

6. **Performance**:
   - Enable GZIP compression
   - HTTP/2 support
   - Database query optimization (indexes)
   - Redis caching for frequently accessed data:
     - Current odds
     - Event lists
     - User sessions

Frontend Production Setup:

1. **Build Optimization**:
   - Vite production build
   - Code splitting by route
   - Tree shaking
   - Minification
   - Asset optimization (images, fonts)
   - Bundle size analysis

2. **Environment Variables**:
```env
VITE_API_URL=https://api.betkz.com
VITE_WS_URL=wss://api.betkz.com/ws
VITE_ENV=production
```

3. **Error Tracking**:
   - Sentry for frontend errors (optional)
   - Log errors to console in dev, Sentry in prod

4. **Analytics** (Optional):
   - Google Analytics or Plausible
   - Track page views, bet placements, conversions

Docker Setup:

1. **Backend Dockerfile** (production):
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o betkz-api ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/betkz-api .
EXPOSE 8080
CMD ["./betkz-api"]
```

2. **Frontend Dockerfile** (production):
```dockerfile
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

3. **Docker Compose** (production):
```yaml
version: '3.8'
services:
  db:
    image: postgres:15-alpine
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      POSTGRES_DB: betkz
      POSTGRES_USER: betkz
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
    restart: unless-stopped

  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgresql://betkz:${DB_PASSWORD}@db:5432/betkz
      REDIS_URL: redis://redis:6379
      JWT_SECRET: ${JWT_SECRET}
    depends_on:
      - db
      - redis
    restart: unless-stopped

  frontend:
    build: ./frontend
    ports:
      - "80:80"
    depends_on:
      - backend
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
```

4. **Nginx Configuration** (for frontend):
```nginx
server {
    listen 80;
    server_name betkz.com;
    root /usr/share/nginx/html;
    index index.html;

    # Gzip compression
    gzip on;
    gzip_types text/plain text/css application/json application/javascript;

    # SPA routing
    location / {
        try_files $uri $uri/ /index.html;
    }

    # API proxy
    location /api {
        proxy_pass http://backend:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    # WebSocket proxy
    location /ws {
        proxy_pass http://backend:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
    }

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
}
```

Deployment Steps:

1. **Pre-deployment**:
   - Run all tests
   - Check database migrations
   - Review environment variables
   - Test Docker build locally

2. **Deployment** (manual or CI/CD):
```bash
# Build images
docker-compose build

# Run migrations
docker-compose run backend ./betkz-api migrate up

# Start services
docker-compose up -d

# Check health
curl http://localhost:8080/health
```

3. **Post-deployment**:
   - Verify all services running
   - Check logs for errors
   - Test critical user flows:
     - Login
     - View events
     - Place bet
     - WebSocket connection
   - Monitor performance metrics

Monitoring (Optional):
- Prometheus for metrics
- Grafana for dashboards
- Uptime monitoring (UptimeRobot, Pingdom)

Please provide:
1. Production Dockerfiles (backend and frontend)
2. Docker Compose production setup
3. Nginx configuration with API/WS proxying
4. Environment variable templates
5. Database migration runner script
6. Health check endpoints
7. Logging configuration
8. Security middleware setup
9. CORS configuration
10. Deployment documentation (step-by-step)
11. Rollback procedure documentation
```

---

### Part 7.2: Testing & Documentation

```
Create comprehensive testing suite and documentation for BetKZ:

Backend Tests (Go):

1. **Unit Tests**:

Test core services:
- Odds calculation engine:
  - Test odds calculation formula
  - Test edge cases (zero bets, single bet, many bets)
  - Test odds boundaries (min 1.01, max 100)
  - Test house margin application
- Bet placement:
  - Test stake validation
  - Test balance checks
  - Test odds locking
  - Test accumulator calculation
- Settlement:
  - Test payout calculation with locked odds
  - Test accumulator settlement (all legs must win)
  - Test balance updates

Example test:
```go
func TestCalculateOdds(t *testing.T) {
    tests := []struct {
        name           string
        totalPool      float64
        outcomeStake   float64
        margin         float64
        expectedOdds   float64
    }{
        {
            name:         "Equal stakes",
            totalPool:    100.0,
            outcomeStake: 50.0,
            margin:       0.05,
            expectedOdds: 1.90,  // 1 / (0.5 * 1.05)
        },
        {
            name:         "Heavy favorite",
            totalPool:    100.0,
            outcomeStake: 80.0,
            margin:       0.05,
            expectedOdds: 1.19,  // 1 / (0.8 * 1.05)
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := CalculateOdds(tt.totalPool, tt.outcomeStake, tt.margin)
            if math.Abs(result-tt.expectedOdds) > 0.01 {
                t.Errorf("expected %.2f, got %.2f", tt.expectedOdds, result)
            }
        })
    }
}
```

2. **Integration Tests**:

Test API endpoints:
- POST /api/auth/register
- POST /api/auth/login
- POST /api/bets/place
- GET /api/events/:id/markets

Example:
```go
func TestPlaceBet(t *testing.T) {
    // Setup: create test user, event, market, odds
    user := createTestUser()
    event := createTestEvent()
    market := createTestMarket(event.ID)
    odd := createTestOdd(market.ID, 2.5)

    // Test: place bet
    req := PlaceBetRequest{
        UserID: user.ID,
        Selections: []Selection{{OddID: odd.ID}},
        Stake: 10.0,
        BetType: "single",
    }

    resp, err := PlaceBet(req)
    assert.NoError(t, err)
    assert.Equal(t, 25.0, resp.PotentialReturn)  // 10 * 2.5

    // Verify: user balance decreased
    updatedUser := getUser(user.ID)
    assert.Equal(t, user.Balance-10.0, updatedUser.Balance)

    // Verify: odds total stake increased
    updatedOdd := getOdd(odd.ID)
    assert.Equal(t, odd.TotalStake+10.0, updatedOdd.TotalStake)
}
```

3. **Database Tests**:
- Test transaction isolation
- Test concurrent bet placements (race conditions)
- Test constraint violations

Run tests:
```bash
go test ./... -v -cover
```

Frontend Tests (React):

1. **Component Tests** (Vitest + React Testing Library):

Test critical components:
- Login form validation
- Bet slip calculations
- Odds button behavior
- Bet confirmation modal

Example:
```tsx
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import BetslipWidget from './BetslipWidget';

describe('BetslipWidget', () => {
  it('calculates potential return correctly', () => {
    const selections = [
      { oddId: '1', outcome: 'Home', odds: 2.0 },
      { oddId: '2', outcome: 'Over', odds: 1.5 },
    ];

    render(<BetslipWidget selections={selections} betType="accumulator" />);

    const stakeInput = screen.getByLabelText('Stake');
    fireEvent.change(stakeInput, { target: { value: '10' } });

    const potentialReturn = screen.getByText(/Potential Return/);
    expect(potentialReturn).toHaveTextContent('$30.00'); // 10 * 2.0 * 1.5
  });
});
```

2. **E2E Tests** (Playwright):

Test complete user flows:
- User registration and login
- Browse events and view odds
- Add bet to betslip
- Place bet successfully
- View bet in history

Example:
```typescript
import { test, expect } from '@playwright/test';

test('complete betting flow', async ({ page }) => {
  // Login
  await page.goto('/login');
  await page.fill('input[name="email"]', 'test@example.com');
  await page.fill('input[name="password"]', 'password123');
  await page.click('button[type="submit"]');

  // Navigate to event
  await page.goto('/events/123');

  // Click odds to add to betslip
  await page.click('[data-testid="odds-button-home"]');

  // Enter stake
  await page.fill('[data-testid="stake-input"]', '10');

  // Verify potential return
  const returnText = await page.textContent('[data-testid="potential-return"]');
  expect(returnText).toContain('$25.00');

  // Place bet
  await page.click('[data-testid="place-bet-button"]');
  await page.click('[data-testid="confirm-bet-button"]');

  // Verify success message
  await expect(page.locator('text=Bet placed successfully')).toBeVisible();
});
```

Run tests:
```bash
# Component tests
npm run test

# E2E tests
npm run test:e2e
```

Documentation:

1. **API Documentation** (OpenAPI/Swagger):

Generate API docs:
```yaml
openapi: 3.0.0
info:
  title: BetKZ API
  version: 1.0.0

paths:
  /api/auth/register:
    post:
      summary: Register new user
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                email:
                  type: string
                password:
                  type: string
      responses:
        '201':
          description: User created
        '400':
          description: Invalid input

  /api/bets/place:
    post:
      summary: Place a bet
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                selections:
                  type: array
                  items:
                    type: object
                stake:
                  type: number
                betType:
                  type: string
                  enum: [single, accumulator]
      responses:
        '201':
          description: Bet placed successfully
```

2. **User Guide**:

Create `docs/USER_GUIDE.md`:
```markdown
# BetKZ User Guide

## Getting Started
1. Create an account
2. Deposit funds
3. Browse events
4. Place your first bet

## How to Place a Bet
1. Navigate to Events page
2. Click on an event to view markets
3. Click odds to add to betslip
4. Enter stake amount
5. Click "Place Bet"
6. Confirm your bet

## How Dynamic Odds Work
- Odds start at values set by admin
- As users place bets, odds adjust automatically
- More money on an outcome = lower odds
- Your bet locks in the current odds
- Even if odds change after, your bet uses locked odds

## Bet Types
- **Single**: Bet on one outcome
- **Accumulator**: Combine multiple outcomes, all must win
```

3. **Admin Guide**:

Create `docs/ADMIN_GUIDE.md`:
```markdown
# BetKZ Admin Guide

## Creating Events
1. Navigate to Admin > Events
2. Click "Create Event"
3. Fill in details (sport, teams, start time)
4. Submit

## Setting Initial Odds
1. Go to Admin > Odds Management
2. Select event
3. Create market (e.g., Match Winner)
4. Set initial odds for each outcome
5. Set house margin (default 5%)
6. Submit

## Monitoring Odds
- View real-time odds changes
- Check total stakes per outcome
- Monitor liability (worst-case payout)
- Override odds manually if needed

## Settling Bets
1. After event finishes, go to Admin > Settlement
2. Select event
3. Enter final results
4. Click "Settle Bets"
5. System automatically:
   - Determines winners
   - Calculates payouts using locked odds
   - Updates user balances
```

4. **Developer Setup Guide**:

Update `README.md`:
```markdown
# BetKZ - Online Sportsbook Platform

## Prerequisites
- Docker & Docker Compose
- Go 1.21+
- Node.js 20+
- PostgreSQL 15
- Redis 7

## Quick Start

1. Clone repository:
```bash
git clone https://github.com/yourusername/betkz.git
cd betkz
```

2. Start services:
```bash
docker-compose up -d
```

3. Run migrations:
```bash
docker-compose exec backend ./betkz-api migrate up
```

4. Seed database:
```bash
docker-compose exec backend ./betkz-api seed
```

5. Access application:
- Frontend: http://localhost:5173
- Backend: http://localhost:8080
- API Docs: http://localhost:8080/swagger

## Running Tests

Backend:
```bash
cd backend
go test ./... -v -cover
```

Frontend:
```bash
cd frontend
npm run test        # Component tests
npm run test:e2e    # E2E tests
```

## Project Structure
See `docs/ARCHITECTURE.md`
```

Please provide:
1. Unit tests for odds calculation, bet placement, settlement
2. Integration tests for API endpoints
3. Frontend component tests
4. E2E test suite (Playwright)
5. OpenAPI/Swagger API documentation
6. User guide (how to use the platform)
7. Admin guide (how to manage the platform)
8. Developer setup guide (README)
9. Architecture documentation
10. Testing documentation (how to run tests)
11. Deployment guide (how to deploy to production)
```

---

## Project Timeline Summary

| Week | Phase | Focus | Completion |
|------|-------|-------|------------|
| 1 | Setup | Project structure, Database schema, Docker setup | 10% |
| 2-3 | Backend Core | Auth, Events, Dynamic Odds Engine, WebSocket, Betting | 40% |
| 3-4 | Frontend Core | Auth UI, Events list, Real-time Odds Display, Betslip | 60% |
| 5 | User Features | Profile, Balance, Bet History | 70% |
| 6-7 | Admin Panel | Dashboard, Event/Odds Management, Odds Monitoring, Settlement | 85% |
| 8 | Polish | Error handling, Mobile responsive, Animations | 95% |
| 9 | Deployment | Production config, Testing, Documentation | 100% |

---

## Complete Tech Stack

| Layer | Technologies |
|-------|-------------|
| Frontend | React 18, TypeScript, Vite, TailwindCSS, Zustand, React Query, Framer Motion, Sonner |
| Backend | Go, Gin, PostgreSQL, Redis, gorilla/websocket, JWT |
| Real-time | WebSocket (gorilla/websocket) for odds broadcasting |
| Deployment | Docker, Docker Compose, Nginx |
| Testing | Vitest, React Testing Library, Playwright, Go testing |

---

## Key Features

✅ Admin sets initial odds → baseline for each market  
✅ Odds adjust automatically → based on betting volume  
✅ Real-time updates → WebSocket broadcasts changes to all users  
✅ Odds locking → users lock in odds when bet is confirmed  
✅ Visual indicators → show odds movement (↑↓) with animations  
✅ Stake distribution → users see how money is distributed across outcomes  
✅ Manual override → admin can manually set odds if needed  
✅ Odds history → track changes over time for analytics  
✅ Risk management → monitor liability and adjust margins

---

**End of Prompts Document**
