document.addEventListener('DOMContentLoaded', () => {
    const terminalOutput = document.getElementById('terminal-output');
    const terminalInput = document.getElementById('terminal-input');
    const terminalBody = document.querySelector('.card-terminal-retro');

    if (!terminalInput) return;

    // Track state
    const startTime = new Date();
    const commandHistory = [];
    let historyIndex = -1;

    // Help Text Constant
    const WELCOME_TEXT = "Welcome to Portfolio Shell v1.0. Type 'help' to see available commands.";

    // Simulated File System
    const fileSystem = {
        'contact.txt': 'Email: jmorgan3142001@gmail.com\nLinkedIn: linkedin.com/in/jmorgan3142001',
        'projects.txt': '- Auto-Caption Network\n- Lawless Lowcountry Living\n- CRM Pipeline Opt.',
        'skills.txt': '- Python (Django)\n- Go (Golang)\n- Typescript\n- SQL\n- Distributed Systems\n- C++',
        'sys_log.log': `[INFO] System initialized at ${startTime.toISOString()}.\n[WARN] Caffeine levels critical.\n[INFO] Portfolio v2.1 loaded.`
    };

    // --- Helper Functions ---

    const getUptime = () => {
        const now = new Date();
        const diff = Math.floor((now - startTime) / 1000); 
        const hours = Math.floor(diff / 3600);
        const minutes = Math.floor((diff % 3600) / 60);
        const seconds = diff % 60;
        return `up ${hours} hrs, ${minutes} min, ${seconds} sec`;
    };

    const getSystemStats = () => {
        // 1. Cores 
        const cores = navigator.hardwareConcurrency || 'unknown';
        
        // 2. Memory
        let memory = 'restricted';
        if (performance && performance.memory) {
            const used = Math.round(performance.memory.usedJSHeapSize / 1024 / 1024);
            const total = Math.round(performance.memory.totalJSHeapSize / 1024 / 1024);
            memory = `${used}MB/${total}MB heap`;
        }

        // 3. Connection
        let net = 'online';
        if (navigator.connection) {
            net = `${navigator.connection.effectiveType || 'unknown'} (${navigator.connection.rtt || '?'}ms rtt)`;
        }

        return `HW: ${cores} cores | MEM: ${memory} | NET: ${net}`;
    };

    const printLine = (text, isHtml = false) => {
        const line = document.createElement('div');
        line.className = 'text-mono small mb-1'; 
        line.style.whiteSpace = 'pre-wrap';
        if (isHtml) line.innerHTML = text;
        else line.textContent = text;
        terminalOutput.appendChild(line);
        terminalBody.scrollTop = terminalBody.scrollHeight;
    };

    // --- Command Logic ---

    const commands = {
        'help': () => 'Commands: help, ls, cd [page], cat [file], tail [file], grep [term] [file], uptime, ping [host], theme [dark|light|matrix], clear, exit\nNot root user. Do not attempt \'sudo\' commands.',
        
        'ls': () => Object.keys(fileSystem).join('  '),
                
        'uptime': () => `${new Date().toLocaleTimeString()} ${getUptime()}, 1 user, ${getSystemStats()}`,
        
        'date': () => new Date().toString(),
        
        'history': () => commandHistory.map((cmd, i) => ` ${i + 1}  ${cmd}`).join('\n'),

        'clear': () => {
            terminalOutput.innerHTML = '';
            // Redisplay welcome text after clear
            const welcomeLine = document.createElement('div');
            welcomeLine.className = 'text-mono small mb-3 text-secondary';
            welcomeLine.textContent = WELCOME_TEXT;
            terminalOutput.appendChild(welcomeLine);
            return ''; 
        },

        'cat': (args) => {
            if (!args[0]) return 'Usage: cat [filename]';
            if (fileSystem[args[0]]) return fileSystem[args[0]];
            return `cat: ${args[0]}: No such file or directory`;
        },

        'tail': (args) => {
            if (!args[0]) return 'Usage: tail [filename]';
            if (fileSystem[args[0]]) {
                const lines = fileSystem[args[0]].split('\n');
                return lines[lines.length - 1];
            }
            return `tail: ${args[0]}: No such file or directory`;
        },

        'grep': (args) => {
            if (args.length < 2) return 'Usage: grep [search_term] [filename]';
            const term = args[0];
            const file = args[1];
            if (fileSystem[file]) {
                const lines = fileSystem[file].split('\n');
                const matches = lines.filter(line => line.toLowerCase().includes(term.toLowerCase()));
                return matches.length > 0 ? matches.join('\n') : '';
            }
            return `grep: ${file}: No such file or directory`;
        },

        'cd': (args) => {
            if (!args[0]) return 'Usage: cd [directory]';
            const dir = args[0].toLowerCase();
            const routes = { 'home': '/', '..': '/', '~': '/', 'library': '/library', 'research': '/research', 'about': '/about' };
            if (routes[dir]) {
                window.location.href = routes[dir];
                return `Navigating to ${dir}...`;
            }
            return `bash: cd: ${dir}: No such directory`;
        },

        // --- Interactive / Real Data Commands ---

        'ping': async (args) => {
            const host = args[0] || window.location.hostname;
            let url = host;
            if (!url.startsWith('http')) {
                url = 'https://' + url;
            }

            printLine(`PING ${host} (${url}): 56 data bytes`);

            let seq = 0;
            let successCount = 0;
            const maxPings = 2;
            let totalTime = 0;

            const runPing = async () => {
                if (seq >= maxPings) {
                    const loss = ((maxPings - successCount) / maxPings) * 100;
                    const avg = successCount > 0 ? (totalTime / successCount).toFixed(2) : 0;
                    
                    printLine(`\n--- ${host} ping statistics ---`);
                    printLine(`${maxPings} packets transmitted, ${successCount} packets received, ${loss}% packet loss`);
                    if (successCount > 0) {
                        printLine(`round-trip min/avg/max = ?/${avg}/? ms`);
                    }
                    return; 
                }

                const start = performance.now();
                try {
                    await fetch(url, { method: 'HEAD', mode: 'no-cors', cache: 'no-cache' });
                    const end = performance.now();
                    const rtt = (end - start).toFixed(2);
                    totalTime += parseFloat(rtt);
                    successCount++;
                    printLine(`64 bytes from ${host}: icmp_seq=${seq} ttl=64 time=${rtt} ms`);
                } catch (e) {
                    printLine(`Request timeout for icmp_seq=${seq}`);
                }

                seq++;
                setTimeout(runPing, 1000); 
            };

            runPing();
            return ''; 
        },

        'theme': (args) => {
            const themes = {
                'dark': { bg: '#2C2C2C', color: '#f8f8f2', input: '#f8f8f2', prompt: '#50fa7b' },
                'light': { bg: '#EAE0D0', color: '#2C2C2C', input: '#2C2C2C', prompt: '#E0A045' },
                'matrix': { bg: '#0D0208', color: '#00FF41', input: '#00FF41', prompt: '#008F11' }
            };
            
            const selected = themes[args[0]];
            if (selected) {
                // Apply background and main text color
                terminalBody.style.backgroundColor = selected.bg;
                terminalBody.style.color = selected.color;
                
                // Force input colors (Text and Caret)
                terminalInput.classList.remove('text-dark');
                terminalInput.style.color = selected.input;
                terminalInput.style.caretColor = selected.input;

                // Update all previous output lines to match new theme
                document.querySelectorAll('#terminal-output div').forEach(el => {
                    // Reset specific color classes so they inherit from parent
                    el.style.color = selected.color === '#2C2C2C' ? '#595959' : selected.color; 
                });
                
                // Update prompts to the theme's specific prompt color
                document.querySelectorAll('.terminal-prompt').forEach(el => {
                    el.style.color = selected.prompt;
                });
                return `Theme changed to ${args[0]}.`;
            }
            return 'Usage: theme [light|dark|matrix]';
        },

        'make': (args) => {
            if (args[0] === 'coffee') return 'make: *** No rule to make target `coffee`. Stop.\n(Hint: Error 418 I\'m a teapot)';
            return 'make: *** No targets specified and no makefile found. Stop.';
        },

        'sudo': () => {
            setTimeout(() => {
                window.open('https://www.youtube.com/watch?v=dQw4w9WgXcQ', '_blank'); 
            }, 2000); // Rick Rolled :D
            return 'Access Denied: User is not in the sudoers file. Goodbye :D';
        },

        'exit': () => {
            setTimeout(() => window.location.href = '/', 1000);
            return 'Logging out...';
        }
    };

    // --- Input Handling ---

    terminalInput.addEventListener('keydown', (e) => {
        if (e.key === 'ArrowUp') {
            if (historyIndex > 0) {
                historyIndex--;
                terminalInput.value = commandHistory[historyIndex];
            }
            e.preventDefault();
        } else if (e.key === 'ArrowDown') {
            if (historyIndex < commandHistory.length - 1) {
                historyIndex++;
                terminalInput.value = commandHistory[historyIndex];
            } else {
                historyIndex = commandHistory.length;
                terminalInput.value = '';
            }
            e.preventDefault();
        }

        if (e.key === 'Enter') {
            const fullCommand = terminalInput.value.trim();
            if (fullCommand === '') return;

            commandHistory.push(fullCommand);
            historyIndex = commandHistory.length;

            const historyLine = document.createElement('div');
            historyLine.className = 'text-mono small mb-1';
            
            // Capture current prompt color for history consistency
            const currentPromptColor = document.querySelector('.terminal-prompt').style.color || '';
            historyLine.innerHTML = `<span class="terminal-prompt" style="user-select:none; color: ${currentPromptColor};">user@portfolio:~/library$</span> ${fullCommand}`;
            terminalOutput.appendChild(historyLine);

            const parts = fullCommand.split(' ');
            const cmd = parts[0].toLowerCase();
            const args = parts.slice(1);

            let response = '';
            if (commands[cmd]) {
                response = commands[cmd](args);
            } else {
                response = `bash: ${cmd}: command not found`;
            }

            if (response) {
                printLine(response);
            }

            terminalInput.value = '';
            terminalBody.scrollTop = terminalBody.scrollHeight;
        }
    });

    terminalBody.addEventListener('click', () => {
        terminalInput.focus();
    });
});