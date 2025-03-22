package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Serve static files
	http.HandleFunc("/", serveHTML)
	
	log.Printf("Starting Tetris server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go Tetris</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
            background-color: #f0f0f0;
        }
        .container {
            display: flex;
            flex-direction: column;
            align-items: center;
        }
        #game-board {
            width: 300px;
            height: 600px;
            border: 2px solid #333;
            background-color: #efefef;
            position: relative;
            overflow: hidden;
        }
        .game-info {
            margin-top: 20px;
            display: flex;
            justify-content: space-between;
            width: 300px;
        }
        .score-container, .next-piece {
            border: 2px solid #333;
            padding: 10px;
            width: 120px;
            text-align: center;
        }
        .block {
            position: absolute;
            width: 30px;
            height: 30px;
            box-sizing: border-box;
            border: 1px solid rgba(0, 0, 0, 0.2);
        }
        button {
            margin-top: 10px;
            padding: 10px 20px;
            background-color: #4CAF50;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
        }
        button:hover {
            background-color: #45a049;
        }
        h1 {
            color: #333;
        }
        .controls {
            margin-top: 15px;
            text-align: center;
        }
        .controls p {
            margin: 5px 0;
        }
        .game-over {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background-color: rgba(0, 0, 0, 0.7);
            color: white;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            font-size: 24px;
            display: none;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Go Tetris</h1>
        
        <div id="game-board">
            <div class="game-over">
                <p>Game Over!</p>
                <button id="restart-btn">Play Again</button>
            </div>
        </div>
        
        <div class="game-info">
            <div class="score-container">
                <p>Score: <span id="score">0</span></p>
                <p>Lines: <span id="lines">0</span></p>
                <p>Level: <span id="level">1</span></p>
            </div>
            <div class="next-piece">
                <p>Next:</p>
                <div id="next-piece-preview" style="height: 90px; position: relative;"></div>
            </div>
        </div>
        
        <button id="start-btn">Start Game</button>
        
        <div class="controls">
            <p>Left/Right Arrow: Move</p>
            <p>Up Arrow: Rotate</p>
            <p>Down Arrow: Soft Drop</p>
            <p>Space: Hard Drop</p>
            <p>P: Pause Game</p>
        </div>
    </div>

    <script>
        document.addEventListener('DOMContentLoaded', () => {
            // Game constants
            const ROWS = 20;
            const COLS = 10;
            const BLOCK_SIZE = 30;
            const COLORS = [
                '#FF0D72', // I
                '#0DC2FF', // J
                '#0DFF72', // L
                '#F538FF', // O
                '#FF8E0D', // S
                '#FFE138', // T
                '#3877FF'  // Z
            ];
            
            // Tetromino shapes
            const SHAPES = [
                [[0, 0, 0, 0], [1, 1, 1, 1], [0, 0, 0, 0], [0, 0, 0, 0]], // I
                [[1, 0, 0], [1, 1, 1], [0, 0, 0]],                         // J
                [[0, 0, 1], [1, 1, 1], [0, 0, 0]],                         // L
                [[1, 1], [1, 1]],                                          // O
                [[0, 1, 1], [1, 1, 0], [0, 0, 0]],                         // S
                [[0, 1, 0], [1, 1, 1], [0, 0, 0]],                         // T
                [[1, 1, 0], [0, 1, 1], [0, 0, 0]]                          // Z
            ];
            
            // Game variables
            let board = createBoard();
            let currentPiece = null;
            let nextPiece = null;
            let gameInterval = null;
            let isPaused = false;
            let isGameOver = false;
            let score = 0;
            let lines = 0;
            let level = 1;
            let speed = 1000;
            
            // DOM elements
            const gameBoard = document.getElementById('game-board');
            const scoreElement = document.getElementById('score');
            const linesElement = document.getElementById('lines');
            const levelElement = document.getElementById('level');
            const startButton = document.getElementById('start-btn');
            const restartButton = document.getElementById('restart-btn');
            const nextPiecePreview = document.getElementById('next-piece-preview');
            const gameOverElement = document.querySelector('.game-over');
            
            // Functions
            function createBoard() {
                const board = [];
                for (let row = 0; row < ROWS; row++) {
                    board[row] = [];
                    for (let col = 0; col < COLS; col++) {
                        board[row][col] = 0;
                    }
                }
                return board;
            }
            
            function createPiece() {
                const shapeIndex = Math.floor(Math.random() * SHAPES.length);
                const shape = SHAPES[shapeIndex];
                const color = COLORS[shapeIndex];
                const piece = {
                    shape,
                    color,
                    row: 0,
                    col: Math.floor((COLS - shape[0].length) / 2),
                    shapeIndex
                };
                return piece;
            }
            
            function drawBlock(row, col, color) {
                const block = document.createElement('div');
                block.className = 'block';
                block.style.backgroundColor = color;
                block.style.top = row * BLOCK_SIZE + 'px';
                block.style.left = col * BLOCK_SIZE + 'px';
                gameBoard.appendChild(block);
            }
            
            function drawNextPiece() {
                // Clear previous next piece preview
                while (nextPiecePreview.firstChild) {
                    nextPiecePreview.removeChild(nextPiecePreview.firstChild);
                }
                
                // Draw next piece
                const shape = nextPiece.shape;
                const color = nextPiece.color;
                const offsetX = (4 - shape[0].length) / 2 * BLOCK_SIZE;
                const offsetY = (3 - shape.length) / 2 * BLOCK_SIZE;
                
                for (let row = 0; row < shape.length; row++) {
                    for (let col = 0; col < shape[row].length; col++) {
                        if (shape[row][col]) {
                            const block = document.createElement('div');
                            block.className = 'block';
                            block.style.backgroundColor = color;
                            block.style.top = row * BLOCK_SIZE + offsetY + 'px';
                            block.style.left = col * BLOCK_SIZE + offsetX + 'px';
                            block.style.width = BLOCK_SIZE + 'px';
                            block.style.height = BLOCK_SIZE + 'px';
                            nextPiecePreview.appendChild(block);
                        }
                    }
                }
            }
            
            function draw() {
                // Clear board
                while (gameBoard.childElementCount > 1) { // Keep game over element
                    gameBoard.removeChild(gameBoard.lastChild);
                }
                
                // Draw board
                for (let row = 0; row < ROWS; row++) {
                    for (let col = 0; col < COLS; col++) {
                        if (board[row][col]) {
                            drawBlock(row, col, board[row][col]);
                        }
                    }
                }
                
                // Draw current piece
                if (currentPiece) {
                    for (let row = 0; row < currentPiece.shape.length; row++) {
                        for (let col = 0; col < currentPiece.shape[row].length; col++) {
                            if (currentPiece.shape[row][col]) {
                                drawBlock(
                                    currentPiece.row + row,
                                    currentPiece.col + col,
                                    currentPiece.color
                                );
                            }
                        }
                    }
                }
            }
            
            function moveDown() {
                if (isPaused || isGameOver) return;
                
                currentPiece.row++;
                if (isCollision()) {
                    currentPiece.row--;
                    placePiece();
                    checkLines();
                    if (isGameOver) {
                        gameOverElement.style.display = 'flex';
                        clearInterval(gameInterval);
                    } else {
                        currentPiece = nextPiece;
                        nextPiece = createPiece();
                        drawNextPiece();
                    }
                }
                draw();
            }
            
            function moveLeft() {
                if (isPaused || isGameOver) return;
                
                currentPiece.col--;
                if (isCollision()) {
                    currentPiece.col++;
                }
                draw();
            }
            
            function moveRight() {
                if (isPaused || isGameOver) return;
                
                currentPiece.col++;
                if (isCollision()) {
                    currentPiece.col--;
                }
                draw();
            }
            
            function rotate() {
                if (isPaused || isGameOver) return;
                
                const oldShape = currentPiece.shape;
                currentPiece.shape = rotateShape(currentPiece.shape);
                
                if (isCollision()) {
                    currentPiece.shape = oldShape;
                }
                draw();
            }
            
            function rotateShape(shape) {
                const N = shape.length;
                const rotated = Array.from({ length: N }, () => Array(N).fill(0));
                
                for (let row = 0; row < N; row++) {
                    for (let col = 0; col < N; col++) {
                        rotated[col][N - 1 - row] = shape[row][col];
                    }
                }
                
                return rotated;
            }
            
            function hardDrop() {
                if (isPaused || isGameOver) return;
                
                while (!isCollision()) {
                    currentPiece.row++;
                }
                currentPiece.row--;
                placePiece();
                checkLines();
                if (isGameOver) {
                    gameOverElement.style.display = 'flex';
                    clearInterval(gameInterval);
                } else {
                    currentPiece = nextPiece;
                    nextPiece = createPiece();
                    drawNextPiece();
                }
                draw();
            }
            
            function isCollision() {
                for (let row = 0; row < currentPiece.shape.length; row++) {
                    for (let col = 0; col < currentPiece.shape[row].length; col++) {
                        if (currentPiece.shape[row][col]) {
                            const boardRow = currentPiece.row + row;
                            const boardCol = currentPiece.col + col;
                            
                            if (
                                boardRow < 0 ||
                                boardRow >= ROWS ||
                                boardCol < 0 ||
                                boardCol >= COLS ||
                                board[boardRow][boardCol]
                            ) {
                                return true;
                            }
                        }
                    }
                }
                return false;
            }
            
            function placePiece() {
                for (let row = 0; row < currentPiece.shape.length; row++) {
                    for (let col = 0; col < currentPiece.shape[row].length; col++) {
                        if (currentPiece.shape[row][col]) {
                            const boardRow = currentPiece.row + row;
                            const boardCol = currentPiece.col + col;
                            
                            if (boardRow < 0) {
                                isGameOver = true;
                                return;
                            }
                            
                            board[boardRow][boardCol] = currentPiece.color;
                        }
                    }
                }
            }
            
            function checkLines() {
                let linesCleared = 0;
                
                for (let row = ROWS - 1; row >= 0; row--) {
                    let isLineFull = true;
                    
                    for (let col = 0; col < COLS; col++) {
                        if (!board[row][col]) {
                            isLineFull = false;
                            break;
                        }
                    }
                    
                    if (isLineFull) {
                        // Clear the line
                        for (let r = row; r > 0; r--) {
                            for (let col = 0; col < COLS; col++) {
                                board[r][col] = board[r - 1][col];
                            }
                        }
                        
                        // Clear the top line
                        for (let col = 0; col < COLS; col++) {
                            board[0][col] = 0;
                        }
                        
                        // Check the same row again
                        row++;
                        linesCleared++;
                    }
                }
                
                if (linesCleared > 0) {
                    // Update score based on number of lines cleared
                    const points = [40, 100, 300, 1200];
                    score += points[linesCleared - 1] * level;
                    lines += linesCleared;
                    
                    // Update level
                    level = Math.floor(lines / 10) + 1;
                    
                    // Update speed
                    speed = Math.max(100, 1000 - (level - 1) * 100);
                    clearInterval(gameInterval);
                    gameInterval = setInterval(moveDown, speed);
                    
                    // Update UI
                    scoreElement.textContent = score;
                    linesElement.textContent = lines;
                    levelElement.textContent = level;
                }
            }
            
            function startGame() {
                // Reset game state
                board = createBoard();
                currentPiece = createPiece();
                nextPiece = createPiece();
                isPaused = false;
                isGameOver = false;
                score = 0;
                lines = 0;
                level = 1;
                speed = 1000;
                
                // Update UI
                scoreElement.textContent = score;
                linesElement.textContent = lines;
                levelElement.textContent = level;
                gameOverElement.style.display = 'none';
                
                // Draw initial state
                drawNextPiece();
                draw();
                
                // Start game loop
                clearInterval(gameInterval);
                gameInterval = setInterval(moveDown, speed);
                
                // Hide start button
                startButton.style.display = 'none';
            }
            
            function pauseGame() {
                if (isGameOver) return;
                
                isPaused = !isPaused;
                if (isPaused) {
                    clearInterval(gameInterval);
                } else {
                    gameInterval = setInterval(moveDown, speed);
                }
            }
            
            // Event listeners
            startButton.addEventListener('click', startGame);
            restartButton.addEventListener('click', startGame);
            
            document.addEventListener('keydown', (e) => {
                switch (e.key) {
                    case 'ArrowLeft':
                        moveLeft();
                        break;
                    case 'ArrowRight':
                        moveRight();
                        break;
                    case 'ArrowUp':
                        rotate();
                        break;
                    case 'ArrowDown':
                        moveDown();
                        break;
                    case ' ':
                        hardDrop();
                        break;
                    case 'p':
                    case 'P':
                        pauseGame();
                        break;
                }
            });
        });
    </script>
</body>
</html>
`))
}