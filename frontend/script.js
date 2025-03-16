// Unicode символы для фигур
const pieces = {
    white: {
        king: '♔', queen: '♕', rook: '♖', bishop: '♗', knight: '♘', pawn: '♙'
    },
    black: {
        king: '♚', queen: '♛', rook: '♜', bishop: '♝', knight: '♞', pawn: '♟'
    }
};

// Инициализация доски (8x8). Используем null для пустых клеток.
let board = [];
for (let i = 0; i < 8; i++) {
    board[i] = new Array(8).fill(null);
}

// Заполнение доски фигурами (стандартная расстановка)
function initBoard() {
    // Черные
    board[0][0] = {color: 'black', type: 'rook'};
    board[0][1] = {color: 'black', type: 'knight'};
    board[0][2] = {color: 'black', type: 'bishop'};
    board[0][3] = {color: 'black', type: 'queen'};
    board[0][4] = {color: 'black', type: 'king'};
    board[0][5] = {color: 'black', type: 'bishop'};
    board[0][6] = {color: 'black', type: 'knight'};
    board[0][7] = {color: 'black', type: 'rook'};
    for (let j = 0; j < 8; j++) {
        board[1][j] = {color: 'black', type: 'pawn'};
    }
    // Белые
    board[7][0] = {color: 'white', type: 'rook'};
    board[7][1] = {color: 'white', type: 'knight'};
    board[7][2] = {color: 'white', type: 'bishop'};
    board[7][3] = {color: 'white', type: 'queen'};
    board[7][4] = {color: 'white', type: 'king'};
    board[7][5] = {color: 'white', type: 'bishop'};
    board[7][6] = {color: 'white', type: 'knight'};
    board[7][7] = {color: 'white', type: 'rook'};
    for (let j = 0; j < 8; j++) {
        board[6][j] = {color: 'white', type: 'pawn'};
    }
}

// Отрисовка доски
let selectedCell = null;
let currentTurn = 'white';
const boardContainer = document.getElementById('board-container');

function renderBoard() {
    boardContainer.innerHTML = '';
    // Создаем таблицу с дополнительной колонкой для нумерации слева
    let table = document.createElement('table');
    for (let i = 0; i < 8; i++) {
        let tr = document.createElement('tr');
        // Первая ячейка строки — номер ряда (от 8 до 1)
        let tdNum = document.createElement('td');
        tdNum.textContent = 8 - i;
        tdNum.classList.add('side');
        tr.appendChild(tdNum);
        for (let j = 0; j < 8; j++) {
            let td = document.createElement('td');
            td.dataset.row = i;
            td.dataset.col = j;
            // Задаем класс клетки (светлая/темная)
            if ((i + j) % 2 === 0) {
                td.classList.add('light');
            } else {
                td.classList.add('dark');
            }
            // Если в данной клетке есть фигура – отображаем её
            if (board[i][j]) {
                let piece = board[i][j];
                td.textContent = pieces[piece.color][piece.type];
            }
            td.addEventListener('click', cellClick);
            tr.appendChild(td);
        }
        table.appendChild(tr);
    }
    boardContainer.appendChild(table);
}

// Обновление списка захваченных фигур
const whiteCaptures = [];
const blackCaptures = [];

function updateCaptured() {
    const whiteList = document.getElementById('white-captures-list');
    const blackList = document.getElementById('black-captures-list');
    whiteList.textContent = whiteCaptures.join(' ');
    blackList.textContent = blackCaptures.join(' ');
}

// Обработка клика по клетке
function cellClick(e) {
    let row = parseInt(e.currentTarget.dataset.row);
    let col = parseInt(e.currentTarget.dataset.col);
    let cell = e.currentTarget;

    // Если фигура не выбрана, выбираем фигуру текущего игрока
    if (!selectedCell) {
        if (board[row][col] && board[row][col].color === currentTurn) {
            selectedCell = {row, col};
            cell.classList.add('selected');
        }
    } else {
        // Клетка назначения выбрана, выполняем ход
        let from = selectedCell;
        // Если в клетке назначения находится фигура
        if (board[row][col]) {
            // Если фигура принадлежит противнику, фиксируем захват
            if (board[row][col].color !== currentTurn) {
                let capturedPiece = board[row][col];
                if (currentTurn === 'white') {
                    whiteCaptures.push(pieces[capturedPiece.color][capturedPiece.type]);
                } else {
                    blackCaptures.push(pieces[capturedPiece.color][capturedPiece.type]);
                }
            } else {
                // Если фигура своей стороны – отменяем ход
                clearSelection();
                return;
            }
        }
        // Перемещаем фигуру
        board[row][col] = board[from.row][from.col];
        board[from.row][from.col] = null;

        // Меняем ход
        currentTurn = currentTurn === 'white' ? 'black' : 'white';

        clearSelection();
        updateCaptured();
        renderBoard();
    }
}

function clearSelection() {
    selectedCell = null;
    document.querySelectorAll('td').forEach(td => td.classList.remove('selected'));
}

// Инициализация доски и отрисовка
initBoard();
renderBoard();
updateCaptured();
``