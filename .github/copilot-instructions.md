# Copilot Instructions

- This repository implements Conway's Game of Life as a web application using Go and the go-app library.
- The main logic is in `pkg/life/life.go`, which manages the grid, cell state, and UI rendering.
- The UI is a 64x64 grid where each cell can be toggled between alive and dead.
- The application supports starting, pausing, and resuming the simulation.
- State is encoded in the URL for sharing and persistence.
- Use idiomatic Go and follow the structure of the existing code.
- UI changes should use go-app's declarative style.
- CSS is in `web/gameoflife.css` and controls the appearance of the grid and cells.
- Do not introduce external dependencies unless necessary.
- Keep code concise and readable.
