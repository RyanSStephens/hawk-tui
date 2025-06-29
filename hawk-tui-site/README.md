# Hawk TUI Website

Official website for [Hawk TUI](https://github.com/hawk-tui/hawk) - the universal TUI framework.

🌐 **Live Site**: [hawktui.dev](https://hawktui.dev)

## Overview

This is the marketing website for Hawk TUI, built with Next.js and deployed on Cloudflare Pages. It features:

- **Interactive Landing Page** with live TUI demos
- **Real-time Browser Demos** showing Hawk TUI capabilities
- **Code Examples** for multiple programming languages
- **Use Case Showcases** for different industries
- **Installation Instructions** and getting started guide

## Tech Stack

- **Framework**: Next.js 14 with App Router
- **Styling**: Tailwind CSS with custom terminal themes
- **Components**: React with TypeScript
- **Deployment**: Cloudflare Pages with Wrangler
- **Icons**: Lucide React
- **Animations**: Framer Motion

## Features

### Interactive TUI Demos
- **Real-time Logging Demo**: Simulated log streaming with filtering
- **Live Metrics Dashboard**: Dynamic gauges and charts
- **Database Migration**: Progress tracking visualization
- **Combined Dashboard**: Multi-widget interface

### Responsive Design
- Mobile-first responsive layout
- Terminal-inspired color scheme
- Professional dark theme
- Smooth animations and transitions

### Performance
- Static site generation for fast loading
- Optimized images and assets
- CDN delivery via Cloudflare
- Minimal JavaScript bundle

## Development

### Prerequisites
- Node.js 18 or higher
- npm or yarn

### Setup
```bash
# Clone the repository
git clone https://github.com/RyanSStephens/hawk-tui-site.git
cd hawk-tui-site

# Install dependencies
npm install

# Start development server
npm run dev
```

### Available Scripts
- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run start` - Start production server
- `npm run lint` - Run ESLint
- `npm run deploy` - Deploy to Cloudflare Pages

## Deployment

The site is automatically deployed to Cloudflare Pages when changes are pushed to the main branch.

### Manual Deployment
```bash
# Build the site
npm run build

# Deploy with Wrangler
npm run deploy
```

### Environment Setup
1. Install Wrangler CLI: `npm install -g wrangler`
2. Login to Cloudflare: `wrangler login`
3. Configure domain in `wrangler.toml`

## Project Structure

```
src/
├── app/
│   ├── layout.tsx          # Root layout with metadata
│   ├── page.tsx            # Home page
│   └── globals.css         # Global styles
├── components/
│   ├── Hero.tsx            # Landing page hero section
│   ├── Features.tsx        # Feature highlights
│   ├── LiveDemo.tsx        # Interactive TUI demos
│   ├── UseCases.tsx        # Industry use cases
│   ├── CodeExamples.tsx    # Multi-language examples
│   ├── Installation.tsx    # Setup instructions
│   ├── Navigation.tsx      # Site navigation
│   └── Footer.tsx          # Site footer
├── lib/
│   └── utils.ts            # Utility functions
└── hooks/                  # Custom React hooks
```

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-feature`
3. Make your changes
4. Test locally: `npm run dev`
5. Commit changes: `git commit -m "Add new feature"`
6. Push to branch: `git push origin feature/new-feature`
7. Create a Pull Request

### Guidelines
- Follow the existing code style
- Test all interactive demos
- Ensure mobile responsiveness
- Optimize for performance
- Update documentation as needed

## License

This website is part of the Hawk TUI project. See the main [LICENSE](https://github.com/hawk-tui/hawk/blob/main/LICENSE) for details.

## Links

- **Main Project**: [github.com/hawk-tui/hawk](https://github.com/hawk-tui/hawk)
- **Documentation**: [hawktui.dev/docs](https://hawktui.dev/docs)
- **Issues**: [github.com/hawk-tui/hawk/issues](https://github.com/hawk-tui/hawk/issues)
- **Discussions**: [github.com/hawk-tui/hawk/discussions](https://github.com/hawk-tui/hawk/discussions)