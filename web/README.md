# SaaS POS Frontend

Enterprise-grade SaaS POS frontend built with Vue 3 and TypeScript, designed for scalability and maintainability.

## Tech Stack

- **Framework**: Vue 3 (Composition API)
- **Build Tool**: Vite
- **State Management**: Pinia
- **Styling**: Tailwind CSS (v4)
- **UI Components**: Shadcn Vue (via Reka UI)
- **HTTP Client**: Axios (wrapped in custom http module)
- **Icons**: Lucide Icons
- **Type Safety**: TypeScript
- **Validation**: Integrated form validation composables

## Project Structure

The codebase follows a modular structure within the `src/` directory:

- `api/`: Pure HTTP functions for backend integration.
- `assets/`: Global styles, including `main.css` (color tokens) and `base.css`.
- `components/`:
    - `ui/`: Shadcn UI base components.
    - `common/`: Shared reusable application components.
- `composables/`: Shared business logic and stateful utilities (e.g., `useFormErrors`, `usePagination`).
- `layouts/`: Application layout wrappers (e.g., `AppLayout.vue`).
- `router/`: Vue Router configuration and navigation guards.
- `stores/`: Pinia setup stores for global state management.
- `types/`: TypeScript interfaces and type definitions.
- `utils/`: Common utility functions (e.g., currency formatting).
- `views/`: Page-level components organized by feature.

## AI Agentic System

This codebase is developed using an **AI Agentic Architecture** orchestrated by the **Gemini CLI**. The majority of the logic and components are generated and maintained by specialized AI sub-agents to ensure consistency and speed.

### Orchestration Model
The system uses an **Orchestrator** (defined in `gemini.md`) that routes tasks to specialized sub-agents located in `.gemini/agents/`. This multi-agent approach ensures that each aspect of the application—from styling to state management—is handled by a focused expert.

### Sub-Agents Registry

| Agent | Responsibility |
|---|---|
| `api-integration-agent` | API types and endpoint integration. |
| `state-management-agent` | Pinia stores and business logic composables. |
| `ui-component-generator` | Creating UI components and feature views. |
| `styling-css-agent` | Layout, colors, themes, and Tailwind utility management. |
| `form-validator-agent` | Form logic, validation rules, and error messaging. |
| `error-handler-agent` | Error states, try/catch logic, and user notifications (toasts). |
| `code-reviewer-refactorer` | Code cleanup, optimization, and adherence to standards. |
| `performance-optimizer-agent` | Runtime performance and bundle optimization. |
| `responsive-design-checker-agent` | Mobile-first responsiveness and layout integrity. |
| `documentation-writer-agent` | Maintaining project documentation and READMEs. |

### How to Prompt
When interacting with the Gemini Orchestrator, use specific task-oriented prompts:
- "Create a new view for Category management using the category store."
- "Refactor the DiscountForm to use the new form validation rules."
- "Fix the responsive layout for the mobile sidebar."
- "Integrate the new Payroll API endpoints and update the store."

## Naming Conventions

To maintain consistency across the agentic workflow, the following naming conventions are strictly enforced:

| Type | Convention | Example |
|---|---|---|
| Component / View | PascalCase | `CategoryView.vue` |
| Store | camelCase + `.stores.ts` | `category.stores.ts` |
| API | camelCase + `.api.ts` | `category.api.ts` |
| Types | camelCase + `.types.ts` | `category.types.ts` |
| Composable | camelCase + `use` prefix | `useFormErrors.ts` |

## Getting Started

### Prerequisites
- Node.js (v20.19.0 or >=22.12.0)
- npm

### Installation
```bash
npm install
```

### Development
```bash
npm run dev
```

### Build for Production
```bash
npm run build
```

### Linting and Formatting
```bash
# Run all linters
npm run lint

# Run specific linters
npm run lint:oxlint
npm run lint:eslint

# Format code
npm run format
```
