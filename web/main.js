'use strict';

const go = new Go();

let maxGuesses = 7;
let currentGuess = 0;

async function initWasm() {
  const result = await WebAssembly.instantiateStreaming(
    fetch('game.wasm'),
    go.importObject
  );
  // Do not await: go.run() blocks forever inside Go's select{}.
  // By the time JS regains control, window.newGame / window.makeGuess are set.
  go.run(result.instance);

  document.getElementById('loading').classList.add('hidden');
  document.getElementById('terminal').classList.remove('hidden');
  startNewGame();
}

function startNewGame() {
  const resp = JSON.parse(window.newGame());
  maxGuesses = resp.maxGuesses;
  currentGuess = 0;

  document.getElementById('history').innerHTML = '';
  document.getElementById('status').textContent = '';
  document.getElementById('instructions').textContent =
    `Guess a complex number a+bi where both a and b are between 1 and 10.\n` +
    `After each guess you receive modulus (distance from origin) and angle feedback.\n` +
    `You have ${maxGuesses} guesses. Format: 5+3i`;

  setInputEnabled(true);
  updatePrompt();
  const input = document.getElementById('guess-input');
  input.value = '';
  input.focus();
}

function updatePrompt() {
  document.getElementById('prompt').textContent =
    `Guess ${currentGuess + 1}/${maxGuesses}: `;
}

function setInputEnabled(enabled) {
  document.getElementById('guess-input').disabled = !enabled;
  document.getElementById('submit-btn').disabled = !enabled;
}

function appendLine(html) {
  const div = document.createElement('div');
  div.innerHTML = html;
  document.getElementById('history').appendChild(div);
}

function escapeHtml(str) {
  return String(str)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;');
}

function submitGuess() {
  const inputEl = document.getElementById('guess-input');
  const input = inputEl.value.trim();
  if (!input) return;

  const resp = JSON.parse(window.makeGuess(input));

  if (resp.error) {
    appendLine(`<span class="error-line">&#10060; ${escapeHtml(resp.error)}</span>`);
    inputEl.value = '';
    inputEl.focus();
    return;
  }

  currentGuess = resp.guessNumber;

  const modResult = resp.modulus.result;
  const modIcon = modResult === 'too_low' ? '&#128201;' :
                  modResult === 'too_high' ? '&#128200;' : '&#10003;';
  const modLabel = modResult === 'too_low' ? 'Modulus too low' :
                   modResult === 'too_high' ? 'Modulus too high' : 'Modulus correct';
  const modHtml =
    `<span class="modulus-${modResult}">${modIcon} ${modLabel} (${resp.modulus.value.toFixed(2)})</span>`;

  const angResult = resp.angle.result;
  const angIcon = angResult === 'left' ? '&#8630;' :
                  angResult === 'right' ? '&#8631;' : '&#10003;';
  const angLabel = angResult === 'left' ? 'Turn left' :
                   angResult === 'right' ? 'Turn right' : 'Angle correct';
  const angHtml =
    `<span class="angle-${angResult}">${angIcon} ${angLabel} (${resp.angle.degrees.toFixed(1)}&deg;)</span>`;

  const remainHtml = (!resp.gameOver && resp.guessesRemaining > 0)
    ? ` &mdash; ${resp.guessesRemaining} guess(es) remaining`
    : '';

  appendLine(
    `<div class="guess-row">` +
    `<span class="guess-num">#${resp.guessNumber}</span>` +
    modHtml + angHtml + remainHtml +
    `</div>`
  );

  if (resp.gameOver) {
    setInputEnabled(false);
    const statusEl = document.getElementById('status');
    if (resp.won) {
      statusEl.innerHTML =
        `<span class="win">&#9989; Correct! The number was ${escapeHtml(resp.target)}. ` +
        `You got it in ${resp.guessNumber} guess(es)! &#127881;</span>`;
    } else {
      statusEl.innerHTML =
        `<span class="lose">&#9760; Out of guesses! The number was ${escapeHtml(resp.target)}.</span>`;
    }
  } else {
    updatePrompt();
    inputEl.value = '';
    inputEl.focus();
  }
}

document.getElementById('submit-btn').addEventListener('click', submitGuess);
document.getElementById('guess-input').addEventListener('keydown', e => {
  if (e.key === 'Enter') submitGuess();
});
document.getElementById('new-game-btn').addEventListener('click', startNewGame);

initWasm().catch(err => {
  const loadingEl = document.getElementById('loading');
  loadingEl.textContent = '❌ Failed to load game: ' + err.message;
  console.error(err);
});
