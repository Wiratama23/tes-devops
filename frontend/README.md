# Frontend - Next.js Admin & E-Commerce Platform

A modern, full-featured frontend application built with [Next.js](https://nextjs.org) for managing articles and products with admin functionality.

## 🚀 Features

- **Next.js 14+** - React framework with App Router
- **TypeScript** - Full type safety
- **Responsive UI** - Components and UI library
- **Admin Dashboard** - Administrative interface for content management
- **Article Management** - Create, read, update, delete articles
- **Product Management** - Manage e-commerce products
- **Authentication** - User authentication and authorization
- **File Uploads** - Handle image and file uploads
- **State Management** - Query caching and state handling
- **Unit Testing** - Vitest and MSW for testing
- **Linting & Formatting** - ESLint configuration
- **Bun Runtime** - Fast JavaScript runtime

## 📋 Project Structure

```
src/
├── app/                 # Next.js App Router pages
│   ├── (admin)/        # Admin dashboard routes
│   ├── (site)/         # Public site routes
│   ├── layout.tsx      # Root layout
│   ├── providers.tsx   # App providers (QueryClient, etc.)
│   └── globals.css     # Global styles
├── components/         # React components
│   ├── admin/          # Admin-specific components
│   ├── articles/       # Article-related components
│   ├── products/       # Product-related components
│   ├── site/           # Public site components
│   └── ui/             # Reusable UI components
├── services/           # API service calls
│   ├── articles.ts     # Article API integration
│   ├── auth.ts         # Authentication service
│   ├── products.ts     # Product API integration
│   ├── uploads.ts      # File upload service
│   └── users.ts        # User API integration
├── schemas/            # Data validation schemas
├── tools/              # Utilities and helpers
│   ├── api.ts          # API client setup
│   ├── auth.ts         # Auth utilities
│   ├── logger.ts       # Logging utility
│   ├── query-client.ts # React Query config
│   └── utils.ts        # General utilities
├── types/              # TypeScript type definitions
│   └── api.ts          # API response types
└── proxy.ts        # Admin route auth proxy

tests/                  # Test files
├── setup.ts            # Test configuration
├── admin/              # Admin component tests
├── articles/           # Article tests
├── auth/               # Auth tests
├── components/         # Component tests
├── lib/                # Library tests
├── msw/                # Mock Service Worker setup
├── products/           # Product tests
└── schemas/            # Schema validation tests
```

## 🛠️ Getting Started

### Prerequisites

- [Bun](https://bun.sh) (latest version)
- Node.js 18+ (for compatibility)

### Installation

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies using Bun:
```bash
bun install
```

### Development

Run the development server:

```bash
bun run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser to see the application.

The app will automatically reload as you make changes to the code.

## 📦 Available Scripts

```bash
# Development
bun run dev           # Start development server

# Building
bun run build         # Build for production
bun start            # Start production server

# Testing
bun run test         # Run unit tests
bun run test:cov # Run tests with coverage report

# Linting
bun run lint         # Run ESLint

# Type Checking
bun run tsgo --noEmit # Check TypeScript compilation

# Code Quality
bun run format       # Format code (if configured)
```

## 🧪 Testing

This project uses **Vitest** for unit testing and **MSW** (Mock Service Worker) for API mocking.

Run tests:
```bash
bun run test
```

Run tests with coverage:
```bash
bun run test:coverage
```

## 🔐 Authentication

The application includes built-in authentication features:
- User login/logout
- JWT token management
- Protected routes and API endpoints
- Role-based access control (Admin/User)

## 📱 API Integration

Services are configured to communicate with the backend API:
- Articles API
- Products API
- Users API
- Authentication endpoints
- File upload endpoints

Configure the API base URL in `src/tools/api-base.ts`

## 🎨 Styling

- **PostCSS** - CSS processing
- **Tailwind CSS** - Utility-first CSS framework (if configured)
- **CSS Modules** - Component-scoped styles

## 📚 Dependencies

Key dependencies:
- **next** - React framework
- **react** & **react-dom** - UI library
- **typescript** - Type safety
- **swr** - Server state management
- **vitest** - Unit testing
- **msw** - API mocking
- **eslint** - Code linting

See `package.json` for complete dependency list.

## 🚢 Deployment

### Build for Production
```bash
bun run build
```

### Run Production Server
```bash
bun start
```

### Deploy to Vercel (Recommended)

1. Push your code to GitHub
2. Connect your repository to [Vercel](https://vercel.com)
3. Select "Next.js" as the framework
4. Deploy

## 📖 Documentation

- [Next.js Documentation](https://nextjs.org/docs)
- [Bun Documentation](https://bun.sh/docs)
- [TypeScript Documentation](https://www.typescriptlang.org/docs/)
- [Vitest Documentation](https://vitest.dev/)

## 🤝 Contributing

1. Create a feature branch (`git checkout -b feature/amazing-feature`)
2. Commit your changes (`git commit -m 'Add amazing feature'`)
3. Push to the branch (`git push origin feature/amazing-feature`)
4. Open a Pull Request

## 📝 License

This project is part of the Mites DevOps platform. See the LICENSE file in the root directory for details.

## 🐛 Issues & Support

For issues or questions, please open an issue in the repository or contact the development team.
