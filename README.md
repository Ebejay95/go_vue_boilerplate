# Frontend Testing Strategy

## Überblick

Dieses Projekt verwendet eine umfassende Frontend-Test-Strategie mit mehreren Ebenen:

### Test-Pyramid

```
    E2E Tests (Cypress)
         /\
        /  \
       /    \
    Integration Tests
       /      \
      /        \
   Unit Tests (Vitest)
```

## Test-Technologien

### Unit Tests - Vitest

- **Zweck**: Testen einzelner Funktionen, Komponenten und Store-Module
- **Framework**: Vitest mit Vue Test Utils
- **Ausführung**: `npm run test:unit`
- **Coverage**: `npm run test:coverage`

### Component Tests - Cypress

- **Zweck**: Isoliertes Testen von Vue-Komponenten
- **Framework**: Cypress Component Testing
- **Ausführung**: `npm run test:component`

### E2E Tests - Cypress

- **Zweck**: End-to-End Testing der kompletten Anwendung
- **Framework**: Cypress E2E
- **Ausführung**: `npm run test:e2e`

### Visual Regression Tests - Percy

- **Zweck**: Erkennung visueller Änderungen
- **Framework**: Percy + Cypress
- **Ausführung**: `npm run test:visual`

## Test-Struktur

```
frontend/
├── src/
│   ├── components/
│   │   └── **/*.test.js        # Unit Tests
│   │   └── **/*.cy.js          # Component Tests
│   ├── store/
│   │   └── **/*.test.js        # Store Unit Tests
│   └── test/
│       ├── setup.js            # Test Setup
│       └── integration/        # Integration Tests
├── cypress/
│   ├── e2e/                    # E2E Tests
│   ├── fixtures/               # Test Data
│   ├── support/                # Helper Functions
│   └── component/              # Component Test Config
├── vitest.config.js            # Vitest Configuration
├── cypress.config.js           # Cypress Configuration
└── package.json                # Test Scripts
```

## Test-Commands

### Makefile Commands

```bash
# Frontend Tests
make test-frontend-unit          # Unit Tests
make test-frontend-coverage      # Coverage Report
make test-frontend-component     # Component Tests
make test-frontend-e2e          # E2E Tests (benötigt Backend)
make test-frontend-visual       # Visual Regression Tests
make test-frontend-lint         # Code Linting
make test-frontend-format       # Code Formatting Check
make test-frontend-all          # Alle Frontend Tests
make test-frontend-ci           # CI Pipeline

# Combined Tests
make test-all                   # Backend + Frontend Tests
make test-quick                 # Schnelle Test-Suite
```

### NPM Scripts

```bash
# Unit Tests
npm run test                    # Run unit tests
npm run test:watch             # Watch mode
npm run test:ui                # UI mode
npm run test:coverage          # With coverage

# E2E Tests
npm run test:e2e               # Run E2E tests
npm run test:e2e:open          # Interactive mode
npm run test:e2e:headless      # Headless mode

# Component Tests
npm run test:component         # Run component tests
npm run test:component:open    # Interactive mode

# Visual Tests
npm run test:visual            # Visual regression tests

# Code Quality
npm run lint                   # ESLint
npm run lint:fix              # Fix linting issues
npm run format                # Prettier formatting
```

## Test-Konfiguration

### Environment Variables

```bash
# Test URLs
VUE_APP_GRPC_WEB_URL=http://localhost:8081
VUE_APP_WS_URL=ws://localhost:8082

# Percy (Visual Testing)
PERCY_TOKEN=your_percy_token

# Lighthouse CI
LHCI_GITHUB_APP_TOKEN=your_lighthouse_token
```

### Coverage Thresholds

```javascript
// vitest.config.js
coverage: {
  thresholds: {
    global: {
      branches: 70,
      functions: 70,
      lines: 75,
      statements: 75
    }
  }
}
```

## Test-Patterns

### Unit Test Beispiel

```javascript
// src/components/base/Button.test.js
import { describe, it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import Button from "./Button.vue";

describe("Button Component", () => {
  it("renders correctly", () => {
    const wrapper = mount(Button, {
      slots: { default: "Click Me" }
    });

    expect(wrapper.text()).toBe("Click Me");
    expect(wrapper.classes()).toContain("btn");
  });
});
```

### E2E Test Beispiel

```javascript
// cypress/e2e/user-management.cy.js
describe("User Management", () => {
  beforeEach(() => {
    cy.intercept("POST", "**/user.UserService/ListUsers", {
      fixture: "users-list.json"
    }).as("listUsers");

    cy.visit("/");
  });

  it("creates a new user", () => {
    cy.get('input[type="text"]').type("John Doe");
    cy.get('input[type="email"]').type("john@example.com");
    cy.get('button[type="submit"]').click();

    cy.wait("@listUsers");
    cy.get(".toast").should("contain", "erfolgreich erstellt");
  });
});
```

### Component Test Beispiel

```javascript
// src/components/base/FormField.cy.js
import FormField from "./FormField.vue";

describe("FormField Component", () => {
  it("updates model value on input", () => {
    const onUpdate = cy.stub();

    cy.mount(FormField, {
      props: {
        field: { key: "name", label: "Name" },
        modelValue: "",
        "onUpdate:modelValue": onUpdate
      }
    });

    cy.get("input").type("John Doe");
    cy.then(() => {
      expect(onUpdate).to.have.been.calledWith("John Doe");
    });
  });
});
```

## Mock-Strategien

### gRPC-Web Mocks

```javascript
// src/test/setup.js
vi.mock("@/proto/user_grpc_web_pb", () => ({
  UserServiceClient: vi.fn(() => ({
    listUsers: vi.fn(),
    createUser: vi.fn(),
    deleteUser: vi.fn()
  }))
}));
```

### WebSocket Mocks

```javascript
// Test Setup
global.WebSocket = vi.fn(() => ({
  onopen: vi.fn(),
  onclose: vi.fn(),
  onmessage: vi.fn(),
  send: vi.fn(),
  close: vi.fn()
}));
```

### Store Mocks

```javascript
// Component Tests
const mockStore = createStore({
  modules: {
    users: {
      namespaced: true,
      state: { users: [] },
      actions: { fetchUsers: vi.fn() }
    }
  }
});
```

## CI/CD Integration

### GitHub Actions

```yaml
# .github/workflows/frontend-tests.yml
jobs:
  unit-tests:
    - name: Run unit tests
      run: npm run test:coverage

  e2e-tests:
    - name: Start backend services
      run: make dev-up-detached
    - name: Run E2E tests
      run: npm run test:e2e:headless
```

### Test Parallelisierung

```javascript
// cypress.config.js
e2e: {
  // Parallelisierung für CI
  retries: {
    runMode: 2,
    openMode: 0
  }
}
```

## Best Practices

### 1. Test-Isolation

- Jeder Test sollte unabhängig lauffähig sein
- Mocks nach jedem Test zurücksetzen
- Datenbankzustand vor Tests säubern

### 2. Aussagekräftige Tests

```javascript
// ✅ Gut: Beschreibt was getestet wird
it("shows error message when email is invalid", () => {});

// ❌ Schlecht: Zu generisch
it("works correctly", () => {});
```

### 3. Test-Daten

- Fixtures für konsistente Test-Daten verwenden
- Realistische aber anonymisierte Daten
- Separate Fixtures für verschiedene Szenarien

### 4. Accessibility Testing

```javascript
// cypress/e2e/accessibility.cy.js
it("has no accessibility violations", () => {
  cy.injectAxe();
  cy.checkA11y();
});
```

### 5. Performance Testing

```javascript
// Lighthouse CI Integration
it("loads within acceptable time", () => {
  const start = Date.now();
  cy.visit("/");
  cy.window().then(() => {
    expect(Date.now() - start).to.be.lessThan(3000);
  });
});
```

## Debugging

### Unit Tests

```bash
# Debug Mode
npm run test:ui                 # Visual debugging
npm run test:watch             # Watch mode
```

### E2E Tests

```bash
# Interactive Mode
npm run test:e2e:open          # Cypress GUI
```

### Screenshots & Videos

- Automatische Screenshots bei Fehlern
- Videos für alle E2E Tests
- Artifacts in CI/CD verfügbar

## Troubleshooting

### Häufige Probleme

1. **Tests hängen bei WebSocket Verbindung**

   ```javascript
   // Mock in beforeEach
   cy.window().then(win => {
     win.WebSocket = MockWebSocket;
   });
   ```

2. **gRPC-Web Calls schlagen fehl**

   ```javascript
   // Intercepts verwenden
   cy.intercept("POST", "**/UserService/**", { fixture: "users.json" });
   ```

3. **Component Tests finden Store nicht**

   ```javascript
   // Store in mount() übergeben
   cy.mount(Component, {
     global: { plugins: [store] }
   });
   ```

4. **Coverage zu niedrig**
   - Test-Abdeckung für kritische Pfade prüfen
   - Unused Code entfernen
   - Edge Cases testen

## Metriken & Reporting

### Coverage Reports

- HTML Report: `coverage/index.html`
- LCOV Format für CI/CD Integration
- Thresholds für Quality Gates

### Test Results

- JUnit XML für CI/CD
- Screenshot/Video Artifacts
- Performance Metriken von Lighthouse

### Visual Regression

- Percy Dashboard für visuelle Änderungen
- Baseline Screenshots
- Approval Workflow für Änderungen
