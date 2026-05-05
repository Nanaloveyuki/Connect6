# Connect6 Architecture

## Responsibility Split

- Frontend only renders API data, handles user interaction, and manages view state.
- Frontend does not compute graph paths, merge GitHub data, maintain seed data, or own business rules.
- Backend owns GitHub data access, caching, graph expansion, shortest-path calculation, and future background refresh jobs.
- Any rule that changes the meaning of a relationship belongs in backend `service` or `domain`, never in frontend code.

## Backend Folder Rules

- `backend/cmd/api`: process entrypoint only. No business logic.
- `backend/configs`: local YAML defaults and environment-specific config files.
- `backend/internal/bootstrap`: application wiring, middleware registration, dependency assembly.
- `backend/internal/config`: config structs and config loading.
- `backend/internal/constant`: stable hardcoded constants such as default port, default timeouts, and default limits.
- `backend/internal/domain`: core business models and domain errors. Split by subject, for example `github` and `graph`.
- `backend/internal/dto`: transport-facing request and response structs. These are for HTTP payloads, not business logic.
- `backend/internal/handler/http`: Gin handlers only. Parse requests, call services, map errors to HTTP responses.
- `backend/internal/repository`: data access interfaces only.
- `backend/internal/repository/memory`: temporary in-memory implementations, cache adapters, and seed/demo data.
- `backend/internal/service`: orchestration and business logic. BFS, validation, dedupe, merge rules, refresh coordination.
- `backend/internal/transport/githubapi`: external GitHub API client and raw GitHub response structs.
- `backend/internal/router`: route registration only.

## Where Structs Go

- Config structs go in `backend/internal/config`.
- Domain structs go in `backend/internal/domain/*`.
- HTTP request and response structs go in `backend/internal/dto`.
- External API response structs go in `backend/internal/transport/githubapi`.
- Repository persistence structs stay beside their concrete implementation, for example under `backend/internal/repository/memory` or a future `postgres` package.

## Where Hardcoding Goes

- Stable application defaults belong in `backend/internal/constant`.
- Local dev seed data belongs in repository implementation folders, currently `backend/internal/repository/memory`.
- Secrets never go into constants or code. They come from env vars or config files.
- Endpoint paths are registered in `backend/internal/router`, not spread across handlers.

## Layer Boundaries

- `handler` can depend on `service` and `dto`.
- `service` can depend on `domain` and `repository` interfaces.
- `repository` implementations can depend on `domain` and external clients.
- `transport` adapters talk to third-party APIs and return domain-friendly data.
- `frontend` only consumes backend DTO-shaped JSON and should not import backend packages or reimplement backend rules.

## GitHub Relationship Modeling

- GitHub does not expose a direct "contact" API. Backend relationship edges must be derived from supported public API slices.
- The current online relationship source for path expansion is `following`, because it is a direct user-to-user edge.
- An empty online relationship result is still a valid result. Do not fall back to local seed data unless the relationship provider is explicitly unavailable, rate-limited, or unauthorized.
- Public organizations are useful context, but `org login` is not a user node. Do not inject organization logins into the current user-only BFS as if they were users.
- If organization membership should influence shortest paths later, add explicit organization nodes or a projected relation service first.
- Repository contributors and other collaboration signals should be added behind bounded expansion limits and caching, not directly inside HTTP handlers.

## Frontend Folder Rules

- `frontend/src/lib/api`: backend API calls only.
- `frontend/src/lib/components`: presentational and container components.
- `frontend/src/lib/config`: frontend runtime configuration such as API base URL.
- `frontend/src/lib/types`: frontend DTO mirrors for API responses.
- `frontend/src/routes`: page composition and route-level data flow.

## Frontend TypeScript Boundaries

- `routes` can orchestrate page state, loading state, error state, and component composition.
- `lib/api` can contain fetch wrappers, request helpers, and response parsing for backend DTOs.
- `lib/i18n` should hold locale dictionaries and typed message access for future internationalization.
- `lib/types` contains only DTO mirror types and UI-local view types. It must not contain business logic.
- Frontend TypeScript must not implement graph traversal, shortest-path search, relation scoring, GitHub data merging, or cache policy logic.
- Frontend constants may include UI defaults such as pagination size, animation timing, and API base URL. Product rules belong to the backend.
- Search defaults, control ranges, and graph runtime defaults should come from backend bootstrap data whenever the UI depends on product behavior.
- If a frontend helper starts interpreting relationship semantics, it is in the wrong layer and must move to backend `service` or `domain`.

## Frontend CSS Boundaries

- Global tokens, reset-like rules, and app-wide background treatment belong in `frontend/src/app.css`.
- Component-scoped styling belongs in the component `.svelte` file unless it is intentionally shared.
- Shared visual tokens should be promoted to CSS variables before repeating hardcoded values across multiple components.
- CSS must not encode product logic. State styling can react to `loading`, `error`, `active`, and layout state, but not infer graph meaning.
- Avoid text-heavy decorative panels. Visual hierarchy should prioritize the query form, graph area, node details, and actual product output.

## Frontend Copy Rules

- Do not place prompt-like helper text, development explanations, architecture notes, or unrelated filler text in the UI.
- Frontend copy must serve product use only: labels, placeholders, concise empty states, actionable errors, and graph-related metadata.
- Explanatory notes for developers belong in docs, code comments, or pull request text, not in the rendered interface.
- Temporary placeholder UI should stay visually neutral and minimal so it does not distort the eventual product experience.
- Frontend must not embed seed graph data or demo relationship datasets in component code unless there is an explicit offline requirement.

## Root Workspace Commands

- Run `pnpm run dev` at repo root to start both backend and frontend.
- Run `pnpm run dev backend` or `pnpm run dev frontend` to start only one side.
- Run `pnpm run check` to execute all checks sequentially.
- Run `pnpm run check backend` or `pnpm run check frontend` to validate a single side.
- Run `pnpm run build` to build all supported targets, or pass `backend` or `frontend` to scope it.
- Backend `check` must prefer `go test ./...` over plain compile-only checks whenever practical.

## Current Intentional Temporary Pieces

- `backend/internal/repository/memory/relationship_repository.go` contains hardcoded seed relationships for local development.
- The shortest-path algorithm already lives in `backend/internal/service/path_service.go`, so replacing seed data later does not require frontend changes.
