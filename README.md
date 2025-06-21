# Game of Life

This project implements [Conway's Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life) as a web application using Go and the [go-app](https://github.com/maxence-charriere/go-app) library.

## Features

- Interactive 64x64 grid for toggling cell states (alive/dead)
- Start, pause, and resume the simulation
- State is encoded in the URL for sharing and persistence
- Responsive UI built with go-app
- Simple, idiomatic Go codebase

## Getting Started

### Prerequisites

- Go 1.18 or later

### Installation

Clone the repository:

```sh
git clone https://github.com/richardwooding/gameoflife.git
cd gameoflife
```

### Running the Application

```sh
make run
```

Then open [http://localhost:8000](http://localhost:8000) in your browser.

## Usage

- Click "Make Colony" to initialize the grid.
- Click on any cell to toggle its state (alive/dead).
- Use the play (▶️) and pause (⏸️) buttons to control the simulation.
- The current state is encoded in the URL, so you can bookmark or share it.

## Development

- Main logic is in `pkg/life/life.go`.
- UI styling is in `web/gameoflife.css`.
- Uses [go-app](https://github.com/maxence-charriere/go-app) for frontend logic.
- No external dependencies except go-app and emoji.

### Modifying the Grid Size

The grid size is set using the `GridWidth` and `GridHeight` constants in `pkg/life/life.go`.  
If you change these, also update the CSS grid in `web/gameoflife.css` for correct display.

## License

MIT License. See [LICENSE](LICENSE) for details.

