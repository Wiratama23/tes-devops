# Frontend - Next.js Admin & E-Commerce PlatformThis is a [Next.js](https://nextjs.org) project bootstrapped with [`create-next-app`](https://nextjs.org/docs/app/api-reference/cli/create-next-app).



A modern, full-featured frontend application built with [Next.js](https://nextjs.org) for managing articles and products with admin functionality.## Getting Started



## 🚀 FeaturesFirst, run the development server:



- **Next.js 14+** - React framework with App Router```bash

- **TypeScript** - Full type safetynpm run dev

- **Responsive UI** - Components and UI library# or

- **Admin Dashboard** - Administrative interface for content managementyarn dev

- **Article Management** - Create, read, update, delete articles# or

- **Product Management** - Manage e-commerce productspnpm dev

- **Authentication** - User authentication and authorization# or

- **File Uploads** - Handle image and file uploadsbun dev

- **State Management** - Query caching and state handling```

- **Unit Testing** - Vitest and MSW for testing

- **Linting & Formatting** - ESLint configurationOpen [http://localhost:3000](http://localhost:3000) with your browser to see the result.

- **Bun Runtime** - Fast JavaScript runtime

You can start editing the page by modifying `app/page.tsx`. The page auto-updates as you edit the file.

## 📋 Project Structure

This project uses [`next/font`](https://nextjs.org/docs/app/building-your-application/optimizing/fonts) to automatically optimize and load [Geist](https://vercel.com/font), a new font family for Vercel.

```

src/## Learn More

├── app/                 # Next.js App Router pages

│   ├── (admin)/        # Admin dashboard routesTo learn more about Next.js, take a look at the following resources:

│   ├── (site)/         # Public site routes

│   ├── layout.tsx      # Root layout- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.

│   ├── providers.tsx   # App providers (QueryClient, etc.)- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

│   └── globals.css     # Global styles

├── components/         # React componentsYou can check out [the Next.js GitHub repository](https://github.com/vercel/next.js) - your feedback and contributions are welcome!

│   ├── admin/          # Admin-specific components

│   ├── articles/       # Article-related components## Deploy on Vercel

│   ├── products/       # Product-related components

│   ├── site/           # Public site componentsThe easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

│   └── ui/             # Reusable UI components

├── services/           # API service callsCheck out our [Next.js deployment documentation](https://nextjs.org/docs/app/building-your-application/deploying) for more details.

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
└── proxy.ts            # API proxy configuration

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

Configure the API base URL in `src/tools/api.ts`

## 🎨 Styling

- **PostCSS** - CSS processing
- **Tailwind CSS** - Utility-first CSS framework (if configured)
- **CSS Modules** - Component-scoped styles

## 📚 Dependencies

Key dependencies:
- **next** - React framework
- **react** & **react-dom** - UI library
- **typescript** - Type safety
- **@tanstack/react-query** - Server state management
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
