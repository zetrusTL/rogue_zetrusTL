# 🗡️ Rogue Terminal Game

Classic roguelike dungeon crawler built in Go — directly in your terminal.

<p align="center">
  <img src="docs/demo.gif" width="700"/>
</p>

---

## ✨ Features

- Procedural dungeon generation
- Fog of war
- Inventory system (weapons, food, elixirs, scrolls, keys)
- Color-coded doors and keys
- Special enemies:
  - 👻 Ghosts (invisible)
  - 📦 Mimics (disguised as items)
- Turn-based combat
- Save / Load
- Statistics & leaderboard
- Game Over / Victory screens
- Fully keyboard controlled

---

## 🎮 Controls

### Movement
- Arrow keys / WASD — move

### Inventory
- `I` — open inventory  
- `1–9` — use / select item  
- `H` — weapons  
- `J` — food  
- `K` — elixirs  
- `E` — scrolls  
- `0` — unequip weapon  

### Game
- `ESC` — menu / exit  
- `R` — restart after Game Over  

---

## 🛠 Tech Stack

- **Go**
- **tcell** (terminal UI)
- Clean Architecture
- JSON persistence

---

## 🧱 Architecture

- presentation/ → rendering + input
- domain/ → game logic & entities
- data/ → repositories & storage


Presentation communicates with domain only through interfaces.

---

## ▶ Run

```bash
go run main.go


