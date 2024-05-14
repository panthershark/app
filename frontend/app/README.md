# Getting Started

- Install the latest LTS of NodeJS
- Run `npm install` from the `./frontend` folder.

# Development Server

- Runs the app in the development mode. `npm run dev -w app`
- Open http://localhost:8081 to view it in the browser.
- The page will reload if you make edits; typescript, graphql, sass are included in the watcher.
- There are more operations in each `package.json`. Packages are self-explanatory so inspect those to find options.

# Testing

Every component, function, or public interface (e.g. exported) should have a test. Code without tests is instant tech debt.

**What should be tested**

- The html output of a rendering function.
- Events - ensure your component is correctly calling the correct handler or updating the value of a passed signal.
- Hooks - ensure they update state correctly.

**What should not be tested**

- The rendered output or screenshots of a browser.
- Other libraries. Just test code you write!
- Small wrappers that bind components that have are already tested.

## State management with Signals

Everything should be a functional component. It renders based on the call parameters.

As the user interacts with the UI, state changes will be necessary. There are two ways we do this:

When the **state is internal** to the component, use a local hook.

```tsx
const Counter = () => {
  const [count, setCount] = useSignal(0);

  return (
    <div>
      <p>Count: {count}</p>
      <button onClick={() => setCount(count + 1)}>Increment</button>
    </div>
  );
};
```

When the **state is external** to the component and a parent might need to change it, then use a signal.

```tsx
const Counter: FunctionalComponent<{ count: Signal<number> }> = ({ count }) => {
  return (
    <div>
      <p>Count: {count.value}</p>
      <button onClick={() => (count.value += 1)}>Increment</button>
    </div>
  );
};

const Page = () => {
  const count = useSignal<number>(0);

  return (
    <div>
      <Counter count={count} />;
      <button type="button" onClick={() => (count.value = 0)}>
        Reset
      </button>
    </div>
  );
};
```

## General Principals

**Make impossible states impossible**

Modeling & state management play a large role in the complexity cost. Increased complexity cost means more initial work (need to write tests), higher maintenance cost, and more chances for bugs.

**Code is a liability, not an asset**
More code and more complexity means higher changes for bugs.

1. Avoid being clever!
2. Avoid YAGNI. Don't build something you don't need yet. Be ruthless!

## Be Ruthless!

If the burden of writing a test seems too high, then the design of the interfaces probably needs more thought to make it easier to test; also, consider discussing with the team.

## Running Tests

- Use `npm run test-watch` for interactive watch mode.
- Use `npm test` to run tests one time; this is used in CI/CD.

# Build

- Build the app for production using `npm run build`. The bundled, minified, code split output will be written to `dist`.
- Pay attention to the size of the output. It might need code splitting.
- CI/CD works exactly like the local build. There is some added file copying to deploy multiple FE apps to static hosting.

# Tech Stack

- NodeJS
- ViteJS
- SASS
- PostCSS
- Preact
- Urql
- Typescript
- Graphql Codegen
- Vitest

# VS Code settings

Do not check in `.vscode` files. Here are some settings from my environment to get you setup.

**settings.json**

```json
{
  "typescript.format.semicolons": "insert",
  "[typescriptreact]": {
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },
  "[typescript]": {
    "editor.formatOnSave": true,
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },
  "editor.codeActionsOnSave": {
    "source.fixAll": true,
    "source.organizeImports": true,
    "source.sortMembers": true
  }
}
```
