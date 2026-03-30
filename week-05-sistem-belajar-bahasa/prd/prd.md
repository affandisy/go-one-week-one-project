# Foreign Language Learning System – Complete System Design

**Version:** 1.0  
**Date:** 2026‑03‑29  
**Author:** Product Manager & Software Architect

---

## 1. Product Requirements Document (PRD)

### 1.1 Executive Summary

The Foreign Language Learning System is a web‑based application that enables beginner users to master basic language fundamentals through a linear, locked progression. The system ensures strong language foundations by guiding users through three progressive levels: Alphabet, Vocabulary, and Sentences. Each level must be completed (≥80% correct) before the next level unlocks. The application focuses on a simple, user‑friendly experience with interactive learning cards, visual quizzes, and sentence‑building exercises.

**Primary Goals:**
- Provide a structured, non‑skippable learning path.
- Reinforce memory using visual and audio cues.
- Track user progress persistently (local storage + optional server sync).
- Limit session fatigue (max 10 questions per session).
- Deliver a seamless, mobile‑responsive interface.

### 1.2 Target Users

- Absolute beginners with no prior language knowledge.
- Casual learners who want short, focused study sessions.
- Users who prefer gamified, quiz‑based learning.

### 1.3 Core Features

#### 1.3.1 User Management
- Registration & login (username + password).
- Persistent progress saved to database (or local storage fallback).
- Automatic session resumption.

#### 1.3.2 Module Structure (Locked Progression)

| Level | Focus                     | Output Goal                         |
|-------|---------------------------|-------------------------------------|
| 1     | Alphabet (characters & sounds) | Recognise letters and pronunciation |
| 2     | Basic Vocabulary (nouns, verbs) | Know 50–100 common words           |
| 3     | Simple Sentences (S-V-O structure) | Construct “I eat an apple”         |

- Level 1 is initially unlocked; levels 2 and 3 are locked.
- Unlock condition: achieve ≥80% correct answers in the current level’s quiz.
- Once unlocked, a level remains accessible for review.

#### 1.3.3 Learning Session Flow
1. **Learn Mode** – Display study material (card with image, word, audio).
2. **Quiz Mode** – Present a set of questions (≤10 per session).
   - Multiple‑choice (choose correct letter/word/image).
   - Match word to image.
   - Unscramble words to form a sentence (level 3 only).
3. **Immediate Feedback** – Show correct/incorrect after each answer.
4. **Progress Validation** – After completing the quiz, calculate score.
   - If score ≥80% and current level is not yet completed → mark level as completed and unlock next level.
   - If score <80% → allow retry (no penalty).

#### 1.3.4 Progress Persistence
- User progress stored in database (server) and optionally cached in local storage.
- When user logs in, progress is restored.
- Offline fallback: if server unreachable, store progress locally and sync later.

#### 1.3.5 Additional Acceptance Criteria
- Each quiz session contains **maximum 10 questions**.
- The “Next” button is **disabled** until the user selects/answers the current question.
- Progress is saved after each completed quiz, even if the browser is closed.

---

## 2. Technical Design (Clean Architecture)

### 2.1 Architecture Overview

The system follows **Clean Architecture** principles, separating concerns into distinct layers. All business logic (progression rules, scoring, module unlocking) resides in the **service layer**. The database is used only for CRUD operations – no stored procedures or database‑side logic.

```
┌─────────────────────────────────────────────────────────────┐
│                     Presentation Layer                       │
│                    (Svelte + Tailwind CSS)                   │
└───────────────────────────┬─────────────────────────────────┘
                            │ HTTP/JSON
┌───────────────────────────▼─────────────────────────────────┐
│                     Delivery Layer (Handlers)                 │
│                   (Fiber, Middlewares, DTOs)                  │
├───────────────────────────────────────────────────────────────┤
│                     Usecase / Service Layer                    │
│     (Business logic: scoring, progression, session limits)    │
├───────────────────────────────────────────────────────────────┤
│                     Repository Layer                           │
│                (GORM – database CRUD operations)               │
└───────────────────────────┬─────────────────────────────────┘
                            │ SQL
┌───────────────────────────▼─────────────────────────────────┐
│                     PostgreSQL Database                       │
│                 (Only stores/fetches data)                    │
└───────────────────────────────────────────────────────────────┘
```

### 2.2 Technology Stack

| Component           | Technology                                                                 |
|---------------------|----------------------------------------------------------------------------|
| **Backend**         | Go 1.22+, Fiber v2, GORM, `github.com/golang-jwt/jwt`                     |
| **Database**        | PostgreSQL 15+ (with `uuid-ossp`)                                          |
| **Frontend**        | Svelte (SvelteKit) + Tailwind CSS, DaisyUI / Skeleton UI components       |
| **State Management**| Svelte stores (user progress, current module)                             |
| **Auth**            | JWT (access token) – optional but recommended for multi‑user support      |
| **File Storage**    | Local file system or S3‑compatible (for images/audio)                     |
| **Deployment**      | Docker, Railway (backend), Vercel (frontend)                               |
| **Monitoring**      | Structured logs, Sentry                                                    |

### 2.3 Database Schema

Only three main tables are required, plus support for progress tracking.

```sql
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR UNIQUE NOT NULL,
    password_hash VARCHAR NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Modules (levels)
CREATE TABLE modules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR NOT NULL,          -- "Alphabet", "Vocabulary", "Sentences"
    description TEXT,
    level_order INTEGER UNIQUE NOT NULL,  -- 1, 2, 3
    is_locked BOOLEAN DEFAULT true        -- initially only level 1 is false
);

-- Learning materials (cards / quiz items)
CREATE TABLE materials (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    module_id UUID REFERENCES modules(id) ON DELETE CASCADE,
    content_type VARCHAR NOT NULL,   -- "learn_card", "quiz_mcq", "quiz_match", "quiz_unscramble"
    question TEXT NOT NULL,          -- e.g., "What letter is this?" or "Unscramble: apple eat I"
    correct_answer TEXT NOT NULL,    -- e.g., "A" or "I eat apple"
    options JSONB,                   -- for MCQ: ["A", "B", "C"]
    image_url TEXT,
    audio_url TEXT,
    display_order INTEGER DEFAULT 0
);

-- User progress (tracking which modules completed and scores)
CREATE TABLE user_progress (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    module_id UUID REFERENCES modules(id),
    completed BOOLEAN DEFAULT false,
    best_score DECIMAL(5,2) DEFAULT 0,   -- percentage (0-100)
    attempts INTEGER DEFAULT 0,
    last_quiz_date TIMESTAMPTZ,
    UNIQUE(user_id, module_id)
);

-- (Optional) Session log for analytics
CREATE TABLE quiz_attempts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    module_id UUID REFERENCES modules(id),
    score DECIMAL(5,2),
    attempted_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_materials_module ON materials(module_id);
CREATE INDEX idx_user_progress_user ON user_progress(user_id);
```

**GORM Model Example:**
```go
type User struct {
    ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
    Username     string    `gorm:"unique;not null"`
    PasswordHash string    `gorm:"not null"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
    Progress     []UserProgress `gorm:"foreignKey:UserID"`
}

type UserProgress struct {
    ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
    UserID       uuid.UUID
    ModuleID     uuid.UUID
    Completed    bool
    BestScore    float64
    Attempts     int
    LastQuizDate *time.Time
}
```

### 2.4 Service Layer – Business Logic

All core algorithms are implemented in Go service methods.

#### 2.4.1 Progression & Unlocking Logic

```go
// service/progress.go
func (s *ProgressService) EvaluateQuiz(userID uuid.UUID, moduleID uuid.UUID, answers []Answer) (*QuizResult, error) {
    // Fetch materials for this module
    materials, err := s.repo.GetMaterialsByModule(moduleID)
    if err != nil { ... }

    // Calculate score (percentage)
    correct := 0
    for i, ans := range answers {
        if ans.Selected == materials[i].CorrectAnswer {
            correct++
        }
    }
    score := (float64(correct) / float64(len(materials))) * 100

    // Get current user progress
    progress, err := s.repo.GetUserProgress(userID, moduleID)
    if err != nil { ... }

    // Update best score if higher
    if score > progress.BestScore {
        progress.BestScore = score
    }
    progress.Attempts++
    progress.LastQuizDate = time.Now()

    // Check completion condition (>=80%)
    if score >= 80 && !progress.Completed {
        progress.Completed = true
        // Unlock next module if exists
        nextModule, _ := s.repo.GetNextModule(moduleID)
        if nextModule != nil {
            s.repo.UnlockModule(userID, nextModule.ID)
        }
    }

    // Save progress
    s.repo.UpdateUserProgress(progress)

    // Log attempt
    s.repo.CreateQuizAttempt(userID, moduleID, score)

    return &QuizResult{Score: score, Completed: progress.Completed}, nil
}
```

#### 2.4.2 Session Limit Enforcement

```go
// service/quiz.go
const MaxQuestionsPerSession = 10

func (s *QuizService) StartQuizSession(userID uuid.UUID, moduleID uuid.UUID) (*QuizSession, error) {
    materials, err := s.repo.GetMaterialsByModule(moduleID)
    if err != nil { ... }

    // Limit to 10 questions (random subset if more than 10)
    if len(materials) > MaxQuestionsPerSession {
        materials = randomSubset(materials, MaxQuestionsPerSession)
    }

    return &QuizSession{
        SessionID:  uuid.New(),
        Questions:  materials,
        StartedAt:  time.Now(),
    }, nil
}
```

#### 2.4.3 Answer Validation & "Next" Button Logic

The frontend enforces that an answer must be selected before moving to the next question. The backend also validates that each submitted answer is not empty and matches the expected question ID.

---

## 3. API Specifications (Fiber)

**Base URL:** `/api/v1`  
**Content-Type:** `application/json`  
**Auth Header:** `Authorization: Bearer <access_token>`

### 3.1 Authentication Endpoints

| Method | Endpoint          | Description                          | Access          |
|--------|-------------------|--------------------------------------|-----------------|
| POST   | `/auth/register`  | Register new user                    | Public          |
| POST   | `/auth/login`     | Login, returns JWT                   | Public          |
| GET    | `/me`             | Get current user profile + progress  | Authenticated   |

**Login Request:**
```json
{ "username": "learner1", "password": "secret123" }
```
**Response:**
```json
{
  "data": {
    "accessToken": "eyJhbG...",
    "user": {
      "id": "uuid",
      "username": "learner1"
    }
  }
}
```

### 3.2 Modules & Learning Materials

| Method | Endpoint                | Description                          | Access          |
|--------|-------------------------|--------------------------------------|-----------------|
| GET    | `/modules`              | List all modules (with lock status)  | Authenticated   |
| GET    | `/modules/{id}`         | Get module details                   | Authenticated   |
| GET    | `/modules/{id}/learn`   | Get learning cards for a module      | Authenticated   |
| GET    | `/modules/{id}/quiz`    | Start a quiz session (returns max 10 questions) | Authenticated |

**GET `/modules` Response:**
```json
{
  "data": [
    {
      "id": "uuid1",
      "title": "Alphabet",
      "levelOrder": 1,
      "isLocked": false,
      "userProgress": { "completed": false, "bestScore": 75.0 }
    },
    {
      "id": "uuid2",
      "title": "Vocabulary",
      "levelOrder": 2,
      "isLocked": true,
      "userProgress": null
    }
  ]
}
```

**GET `/modules/{id}/quiz` Response:**
```json
{
  "data": {
    "sessionId": "temp-session-uuid",
    "questions": [
      {
        "id": "mat1",
        "type": "mcq",
        "question": "What letter is this? 🍎",
        "options": ["A", "B", "C"],
        "imageUrl": "/images/apple.png"
      },
      {
        "id": "mat2",
        "type": "unscramble",
        "question": "Arrange: 'eat I apple'",
        "correctAnswer": "I eat apple"
      }
    ]
  }
}
```

### 3.3 Quiz Submission

| Method | Endpoint                    | Description                          | Access          |
|--------|-----------------------------|--------------------------------------|-----------------|
| POST   | `/modules/{id}/quiz/submit` | Submit answers for a module quiz     | Authenticated   |

**Request Body:**
```json
{
  "sessionId": "temp-session-uuid",
  "answers": [
    { "materialId": "mat1", "selected": "A" },
    { "materialId": "mat2", "selected": "I eat apple" }
  ]
}
```

**Response:**
```json
{
  "data": {
    "score": 85.0,
    "completed": true,
    "nextModuleUnlocked": {
      "id": "uuid2",
      "title": "Vocabulary"
    }
  }
}
```

### 3.4 Progress & Reset

| Method | Endpoint                    | Description                          | Access          |
|--------|-----------------------------|--------------------------------------|-----------------|
| GET    | `/progress`                 | Get full user progress (all modules) | Authenticated   |
| POST   | `/progress/reset`           | Reset progress (admin or user self)  | Authenticated   |

---

## 4. Frontend (Svelte) Implementation Notes

- **Pages:**
  - Login / Register
  - Dashboard (module list with lock icons and progress bars)
  - Learning Mode (cards with image, audio, text)
  - Quiz Mode (dynamic question components: MCQ, match, unscramble)
- **State Management:** Svelte stores for `user`, `currentModule`, `quizSession`.
- **Quiz Flow:**
  - Fetch quiz questions (max 10).
  - Render questions one by one.
  - Disable “Next” until an answer is selected.
  - On last question, “Submit” button appears.
  - After submission, show score and unlock notification.
- **Persistence:**
  - On each quiz submit, progress is saved to server.
  - Local storage fallback for offline scenarios (sync when back online).
- **Responsive Design:** Tailwind CSS ensures mobile‑friendly layout.

---

## 5. Implementation Roadmap (MVP)

| Phase | Duration | Features |
|-------|----------|----------|
| **Phase 1** | Week 1 | Project setup (Go + Fiber + GORM + PostgreSQL), user authentication, basic module CRUD. |
| **Phase 2** | Week 2 | Learning card display, image/audio upload endpoints, Svelte frontend structure. |
| **Phase 3** | Week 3 | Quiz engine (MCQ, match, unscramble), scoring logic, progression unlocking. |
| **Phase 4** | Week 4 | Progress tracking, dashboard, session limit enforcement, local storage sync. |
| **Phase 5** | Week 5 | Testing (unit + integration), deployment (Docker + Vercel), documentation. |

---

## 6. Non‑Functional Requirements

| Aspect          | Requirement                                                                 |
|-----------------|-----------------------------------------------------------------------------|
| **Performance** | Quiz load < 1s, submit response < 500ms.                                    |
| **Usability**   | Max 3 clicks to start a quiz; clear visual feedback on answer.              |
| **Security**    | Passwords hashed (bcrypt); JWT expiry 7 days.                               |
| **Reliability** | Progress saved on every quiz completion; database backups daily.            |
| **Offline Capability** | Local storage caches progress; sync when connection restored.        |

---

## 7. Success Metrics (KPIs)

- **Completion rate:** >60% of registered users finish Level 1.
- **Quiz retry rate:** <30% (i.e., most pass on first or second attempt).
- **Session length:** Average time per quiz <5 minutes.
- **User retention:** 40% return within 7 days.

---

## 8. Glossary

| Term               | Definition                                                                 |
|--------------------|-----------------------------------------------------------------------------|
| **Locked Progression** | Users cannot skip ahead; must complete current level to unlock next.     |
| **Learn Card**      | A flashcard showing a letter, word, or image with optional audio.          |
| **Unscramble**      | A question type where the user arranges jumbled words into a correct sentence. |
| **Best Score**      | Highest percentage achieved on a module’s quiz.                            |

---

**Document Version:** 1.0  
**Last Updated:** 2026‑03‑29