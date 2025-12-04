document.addEventListener('DOMContentLoaded', () => {
    const terminalOutput = document.getElementById('terminal-output');
    const terminalInput = document.getElementById('terminal-input');
    const terminalBody = document.querySelector('.card-terminal-retro');

    if (!terminalInput) return;

    // Simulated File System
    const fileSystem = {
        'contact.txt': 'Email: jmorgan3142001@gmail.com\nLinkedIn: linkedin.com/in/jmorgan3142001',
        'projects.txt': '- Auto-Caption Network\n- Lawless Lowcountry Living\n- CRM Pipeline Opt.',
        'skills.txt': '- Go (Golang)\n- Python (Django)\n- Distributed Systems\n- C++',
        'sys_log.log': '[INFO] System initialized successfully.\n[WARN] Caffeine levels critical.\n[INFO] Portfolio v2.1 loaded.'
    };

    // Command Logic
    const commands = {
        'help': () => 'Available commands: help, ls, cat [file], tail [file], grep [term] [file], whoami, uptime, clear, reboot',
        'ls': () => Object.keys(fileSystem).join('  '),
        'whoami': () => 'guest_user@portfolio',
        'uptime': () => 'up 42 days, 7:14, 1 user, load average: 0.08, 0.03, 0.01',
        'date': () => new Date().toString(),
        'clear': () => {
            terminalOutput.innerHTML = '';
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
                // Return last line simulation
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
                return matches.length > 0 ? matches.join('\n') : ''; // Silent return if no match
            }
            return `grep: ${file}: No such file or directory`;
        },
        'sudo': () => 'Permission denied: User is not in the sudoers file. This incident will be reported.',
        'reboot': () => {
            setTimeout(() => location.reload(), 1000);
            return 'System rebooting...';
        }
    };

    // Handle Enter Key
    terminalInput.addEventListener('keydown', (e) => {
        if (e.key === 'Enter') {
            const fullCommand = terminalInput.value.trim();
            if (fullCommand === '') return;

            // 1. Create history line
            const historyLine = document.createElement('div');
            historyLine.className = 'text-mono small mb-1';
            historyLine.innerHTML = `<span class="terminal-prompt">user@portfolio:~/library$</span> ${fullCommand}`;
            terminalOutput.appendChild(historyLine);

            // 2. Process Command
            const parts = fullCommand.split(' ');
            const cmd = parts[0].toLowerCase();
            const args = parts.slice(1);

            let response = '';
            if (commands[cmd]) {
                response = commands[cmd](args);
            } else {
                response = `bash: ${cmd}: command not found`;
            }

            // 3. Create response line
            if (response) {
                const responseLine = document.createElement('div');
                responseLine.className = 'text-mono small mb-3 text-secondary';
                responseLine.style.whiteSpace = 'pre-wrap';
                responseLine.textContent = response;
                terminalOutput.appendChild(responseLine);
            }

            // 4. Reset Input and Scroll
            terminalInput.value = '';
            terminalBody.scrollTop = terminalBody.scrollHeight;
        }
    });

    terminalBody.addEventListener('click', () => {
        terminalInput.focus();
    });
});