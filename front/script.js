
// Couleurs possibles pour les boutons
const buttonColors = [
    '#ff5733', '#3357ff', '#ff33a1', '#a133ff'
];

// Rayons de bordure possibles pour les boutons
const buttonBorderRadius = [
    '5px', '10px', '15px', '20px', '25px'
];

// Couleurs possibles pour les barres de progression
const progressBarColors = [
    '#ff5733', '#3357ff', '#ff33a1', '#a133ff'
];

// Fonction pour créer les boutons dynamiques
function createButtons(mission) {

    const buttonsContainer = document.getElementById('buttons');
    const button = document.createElement('div');
    button.className = 'button';
    button.textContent = `${mission.verb} ${mission.action}`;
    button.id = `button_${mission.id}`;

    button.addEventListener('click', () => {
        temp = mission.id;
        console.log(temp);
        const message = {
            type: 'solve_mission',
            mission: {
                id: mission.id
            }
        };
        socket.send(JSON.stringify(message));
    });
    // Appliquer une couleur et un rayon de bordure aléatoires
    const color = buttonColors[Math.floor(Math.random() * buttonColors.length)];
    // const color = '#000'
    const borderRadius = buttonBorderRadius[Math.floor(Math.random() * buttonBorderRadius.length)];
    button.style.backgroundColor = color;
    button.style.borderRadius = borderRadius;

    buttonsContainer.appendChild(button);

}

function deleteTask(mission) {
    console.log(`to delete task_${mission.id}`)
    const taskElement = document.getElementById(`task_${mission.id}`);
    if (taskElement) {
        taskElement.parentNode.removeChild(taskElement);
    }
}


function deleteButton(mission) {
    const buttonElement = document.getElementById(`button_${mission.id}`);
    if (buttonElement) {
        buttonElement.classList.add('blink_success');
        setTimeout(() => {
            buttonElement.parentNode.removeChild(buttonElement);
        }, 1000);
    }
}

function spaceshipExplosion() {
    // Create explosion element
    document.body.innerHTML = '';
    const explosion = document.createElement('div');
    explosion.className = 'explosion';
    document.body.appendChild(explosion);

    // Remove all other elements after a short delay
    setTimeout(() => {
        document.body.innerHTML = '';

        // Create and show "GAME OVER" message
        const gameOverMessage = document.createElement('div');
        gameOverMessage.className = 'game-over';
        gameOverMessage.textContent = 'GAME OVER';
        document.body.appendChild(gameOverMessage);
    }, 3000); // Delay to allow explosion animation to play
}

function config(config) {
    document.getElementById('config-form-container').style.display = 'block';
    document.getElementById('config-form').style.display = 'block';
    document.getElementById('nb_players').value = config.nb_players;
    document.getElementById('nb_target_missions').value = config.nb_target_missions;
    document.getElementById('nb_solvable_missions').value = config.nb_solvable_missions;
    document.getElementById('fail_health').value = config.fail_health;
    document.getElementById('success_health').value = config.success_health;
    document.getElementById('default_timeout').value = config.default_timeout;
    document.getElementById('decrease_timeout').value = config.decrease_timeout;
}

function newPlayer(nbPLayer) {
    elem = document.getElementById('connected-players');
    if (elem) {
        elem.textContent = `Connected players: ${nbPLayer}`;
    }
}

function startGame() {


    const element = document.getElementById('config-form-container');
    if (element) {
        element.parentNode.removeChild(element);
    }

    const healthBar = document.getElementById('health-bar');
    if (healthBar) {
        healthBar.style.display = 'block';
    }

    // Create a countdown element
    const countdownElement = document.createElement('div');
    countdownElement.className = 'countdown';
    countdownElement.textContent = 'Game starts in 3';
    document.body.appendChild(countdownElement);

    let countdown = 3;
    const countdownInterval = setInterval(() => {
        countdown -= 1;
        if (countdown > 0) {
            countdownElement.textContent = `Game starts in ${countdown}`;
        } else {
            clearInterval(countdownInterval);
            countdownElement.parentNode.removeChild(countdownElement);
            // Start the game logic here
            console.log('Game started');
        }
    }, 1000);
}

function failButton(mission) {
    const buttonElement = document.getElementById(`button_${mission.id}`);
    if (buttonElement) {
        buttonElement.classList.add('blink_fail');
        setTimeout(() => {
            buttonElement.classList.remove('blink_fail');
        }, 1000);
    }
}

function failTask(mission) {
    const taskElement = document.getElementById(`task_${mission.id}`);
    if (taskElement) {
        const progressBar = taskElement.querySelector('.progress');
        if (progressBar) {
            progressBar.classList.add('blink_fail');
            setTimeout(() => {
                taskElement.parentNode.removeChild(taskElement);
            }, 1000);
        }
    }
}

function successButton(mission) {
    const buttonElement = document.getElementById(`button_${mission.id}`);
    if (buttonElement) {
        buttonElement.classList.add('blink_success');
        setTimeout(() => {
            buttonElement.classList.remove('blink_success');
        }, 1000);
    }
}
function updateHealthBar(health) {
    console.log(health)
    const healthBar = document.getElementById('health-bar');
    healthBar.style.width = `${health}%`;
}

// Fonction pour créer la liste des tâches avec barre de progression
function createTask(mission) {
    const tasksContainer = document.getElementById('tasks');
    const task = document.createElement('div');
    task.id = `task_${mission.id}`;
    task.className = 'task';
    task.innerHTML = `
            <div class="progress-bar">
                <div class="progress"></div>
                <span class="task-text">${mission.verb} ${mission.action}</span>
            </div>
            `;
    tasksContainer.appendChild(task);

    // Appliquer une couleur aléatoire à la barre de progression
    const progressBar = task.querySelector('.progress');
    const color = progressBarColors[Math.floor(Math.random() * progressBarColors.length)];
    progressBar.style.backgroundColor = color;
    progressBar.style.transition = `width ${mission.timeout_seconds}s linear`;
    // Animer la barre de progression
    setTimeout(() => {
        progressBar.style.width = '100%';
    }, 100); // Délai pour permettre l'application du style initial
}

function connectWebSocket(url) {
    socket = new WebSocket("wss://" + url);

    socket.addEventListener('open', (event) => {
        console.log('Connected to the WebSocket server');
    });

    socket.addEventListener('message', (event) => {
        const data = JSON.parse(event.data);
        console.log('Received message:', data);
        switch (data.type) {
            case 'new_target_mission':
                createTask(data.mission);
                break;
            case 'new_solvable_mission':
                createButtons(data.mission);
                break;

            case 'mission_target_solved':
                deleteTask(data.mission);
                break;
            case 'mission_solvable_solved':
                deleteButton(data.mission);
                break;

            case 'mission_target_failed':
                failTask(data.mission);
                break;
            case 'mission_solvable_failed':
                failButton(data.mission);
                break;

            case 'spaceship_explosion':
                spaceshipExplosion();
                break;

            case 'config':
                config(data.config);
                break;
            case 'start_game':
                startGame();
                break;
            case 'new_player':
                newPlayer(data.nb_players);
                break;

            case 'update_health':
                updateHealthBar(data.health);
                break;
            default:
                console.warn('Unknown message type:', data.type);
        }
    });

    socket.addEventListener('error', (event) => {
        console.error('WebSocket error:', event);
    });

    socket.addEventListener('close', (event) => {
        console.log('Disconnected from the WebSocket server');
    });
}
