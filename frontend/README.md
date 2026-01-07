# HRIS Frontend

A modern, responsive Human Resource Information System frontend built with React, TypeScript, and Vite.

## Tech Stack

- **React 19.2.0** - UI library
- **TypeScript 5.9.3** - Type-safe development
- **Vite 7.2.4** - Build tool and dev server
- **React Router DOM 7.11.0** - Client-side routing
- **TailwindCSS 3.4.17** - Utility-first CSS
- **Radix UI** - Accessible component primitives
- **TanStack Query 5.90.16** - Data fetching and caching
- **React Hook Form 7.70.0** - Form management
- **Zod 4.3.5** - Schema validation
- **Axios 1.13.2** - HTTP client

## Project Structure

The project follows **Feature-Based Architecture** with clean architecture principles:

```
frontend/
├── src/
│   ├── components/           # Reusable components
│   │   ├── ui/              # Base UI components (Button, Input, etc.)
│   │   └── layout/          # Layout components (DashboardLayout, Sidebar)
│   ├── features/            # Feature-based modules (Self-contained)
│   │   ├── auth/            # Authentication feature
│   │   │   ├── components/  # Auth-specific components
│   │   │   │   └── LoginForm.tsx
│   │   │   ├── hooks/       # Auth hooks
│   │   │   │   └── useAuth.tsx
│   │   │   └── types.ts     # Auth types (LoginPayload, LoginResponse)
│   │   └── user/            # User management feature
│   │       ├── components/  # User-specific components
│   │       │   ├── GeneralForm.tsx
│   │       │   └── PasswordForm.tsx
│   │       ├── hooks/       # User hooks
│   │       │   └── useProfile.tsx
│   │       └── types.ts     # User types (UserProfile, PasswordPayload)
│   ├── pages/               # Route components
│   │   ├── login/
│   │   └── profile/
│   ├── lib/                 # Utilities and API client
│   │   └── axios.ts        # Axios instance with interceptors
│   └── main.tsx            # App entry point
├── public/                  # Static assets
└── index.html               # HTML template
```

### Architecture Principles

1. **Feature-Based Organization**: Each feature is self-contained with its own:
   - `components/` - Feature-specific UI components
   - `hooks/` - Feature-specific custom hooks
   - `types.ts` - TypeScript interfaces and types

2. **Separation of Concerns**:
   - UI components in `components/ui/` are generic and reusable
   - Feature components are co-located with their business logic
   - Types are defined per-feature for better maintainability

3. **Clean Imports**: Use path aliases for clean imports:
   ```typescript
   import { useProfile } from "@/features/user/hooks/useProfile";
   import { UserProfile } from "@/features/user/types";
   ```

## Quick Start

### Using Docker Compose (Recommended)

From the root directory:

```bash
# Copy environment template
cp .env.example .env

# Start all services
docker compose up -d --build

# Check logs
docker compose logs -f frontend
```

The app will be available at `http://localhost:8080`

### Local Development

```bash
cd frontend

# Install dependencies
pnpm install

# Start dev server
pnpm dev
```

The app will be available at `http://localhost:5173`

## Available Scripts

```bash
pnpm dev          # Start development server
pnpm build        # Build for production
pnpm preview      # Preview production build locally
pnpm lint         # Run ESLint
```

## Key Features

- **Authentication**: Login, session management, JWT tokens
- **User Profile**: View and update profile, change password, upload avatar
- **Modern UI**: Radix UI components with TailwindCSS styling
- **Type Safety**: Full TypeScript with strict mode and type imports
- **Data Fetching**: TanStack Query for API state management and caching
- **Form Validation**: React Hook Form + Zod schemas
- **File Upload**: Multipart/form-data support for profile photos
- **Responsive**: Mobile-first design
- **Clean Architecture**: Feature-based modular structure

## Path Aliases

The project uses `@/*` as an alias for `./src/*`:

```typescript
import { Button } from "@/components/ui/button";
import { useProfile } from "@/features/user/hooks/useProfile";
import { GeneralForm } from "@/features/user/components/GeneralForm";
import type { UserProfile } from "@/features/user/types";
```

## Configuration

### Environment Variables

Key variables (see root `.env.example` for complete list):

| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_API_URL` | Backend API URL | http://localhost:8080 |
| `VITE_PORT` | Dev server port | 5173 |

### TypeScript Configuration

- Target: ES2022
- Strict mode: Enabled
- Path aliases: `@/*` → `./src/*`

### TailwindCSS

- Dark mode: Class-based
- Custom theme: CSS variables
- Design tokens: HSL colors

## Component Guidelines

### File Organization

1. **Reusable UI Components**: Place in `components/ui/` (Button, Input, Form, etc.)
2. **Layout Components**: Place in `components/layout/` (DashboardLayout, Sidebar)
3. **Feature-Specific Components**: Organize under `features/{feature}/components/`

### Creating a New Feature

When adding a new feature, follow this structure:

```
features/{feature-name}/
├── components/      # Feature-specific UI components
├── hooks/          # Feature-specific custom hooks
├── types.ts        # TypeScript interfaces
└── index.ts        # Optional exports barrel
```

### Best Practices

1. **Keep features self-contained** - Each feature should have its own components, hooks, and types
2. **Use type imports** - Use `import type { ... }` for type-only imports
3. **Co-locate related code** - Keep components close to where they're used
4. **Define types per-feature** - Create `types.ts` in each feature folder
5. **Export from hooks** - Export types from hooks for component use
6. **Use composition** - Build complex UIs from simple components
7. **Implement error boundaries** - Wrap features with error handling

### Example: Adding a New Feature

```typescript
// features/attendance/types.ts
export interface AttendanceRecord {
  id: number;
  check_in: string;
  check_out: string;
}

// features/attendance/hooks/useAttendance.ts
import { useQuery } from "@tanstack/react-query";
import type { AttendanceRecord } from "../types";

export const useAttendance = () => {
  return useQuery<AttendanceRecord[]>({
    queryKey: ["attendance"],
    queryFn: async () => {
      // ...
    },
  });
};

// features/attendance/components/AttendanceList.tsx
import { useAttendance } from "../hooks/useAttendance";

export const AttendanceList = () => {
  const { data } = useAttendance();
  // ...
};
```

## API Integration

The app uses:
- **Axios** for HTTP requests with interceptors
- **TanStack Query** for data fetching, caching, and state management
- **Zod** for runtime type validation

### Axios Configuration

The Axios instance (`lib/axios.ts`) includes:
- Automatic JWT token injection from localStorage
- FormData handling for file uploads
- Response interceptors for error handling (401 redirects)
- Automatic Content-Type header management

### File Uploads

For multipart/form-data requests (e.g., profile photo uploads):

```typescript
const formData = new FormData();
formData.append("phone_number", values.phone_number);
if (selectedFile) {
  formData.append("photo", selectedFile);
}

// Axios automatically detects FormData and sets correct headers
api.put("/users/profile", formData);
```

## Development Workflow

```bash
# Type check
tsc --noEmit

# Lint
pnpm lint

# Format (if using prettier)
pnpm format
```

## Troubleshooting

**Dependencies issues**
```bash
rm -rf node_modules pnpm-lock.yaml
pnpm install
```

**Port already in use**
```bash
# Change port in package.json or use:
pnpm dev --port 3000
```

**Build errors**
```bash
# Clear Vite cache
rm -rf dist
pnpm build
```

## Browser Support

Modern browsers with ES2022 support:
- Chrome/Edge 90+
- Firefox 88+
- Safari 14+

## License

See LICENSE file in the root directory.
