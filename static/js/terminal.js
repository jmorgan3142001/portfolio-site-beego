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
    const WELCOME_TEXT = "Welcome to Portfolio Shell v1.0. Type 'help' to see available commands.<br />Warning: User is not privileged. Do not attempt 'sudo' commands.";

    // Simulated File System
    const fileSystem = {
        'contact.txt':
            'Email: jmorgan3142001@gmail.com\n' +
            'LinkedIn: https://www.linkedin.com/in/jake-morgan-\n' +
            'Resume (PDF): https://docs.google.com/document/d/1qCF9Oe2GXS9ayBZgxiA9xq1ZOHGdxF9deOKG7T9Cwvg/export?format=pdf\n' +
            'Note: Email is best for quick questions; LinkedIn is best for professional outreach.',

        'projects.txt':
            '- Auto-Caption Network: distributed, low-latency captioning system for synchronized multi-endpoint workflows.\n' +
            '- Lawless Lowcountry Living: production content site with a mobile-first, accessibility-minded front end.\n' +
            '- CRM Pipeline Optimization: refactored CI/CD and pipeline orchestration to shorten build and deploy times.\n' +
            '- Portfolio Website: this site; a small interactive CLI and static site demonstrating full-stack work.\n' +
            '- Open Source (NCI): Django5 Forms Fieldset and Django5 Scheduler contributions focused on accessibility and scheduling.',

        'skills.txt':
            'Languages: Python, Go, TypeScript, C++, C#, SQL\n' +
            'Backend: Django, .NET Core, REST APIs, gRPC\n' +
            'Frontend: TypeScript, Angular, Bootstrap, Vite\n' +
            'Databases and Infra: PostgreSQL, SQL Server, Azure, Docker\n' +
            'Areas: distributed systems, performance, testing, CI/CD, observability',

        'about.txt':
            "Hi, I'm Jake Morgan, a full-stack engineer focused on building pragmatic, reliable software that scales. " +
            "I prefer clean, testable code and practical solutions that help teams move faster. " +
            "Recent work has centered on captioning automation, modernizing legacy systems, and improving developer workflows.\n\n" +
            "Outside of work I stay active, tinker with tech, and live with two pocket pitties (see pets.txt).",

        'experience.txt':
            'NATIONAL CAPTIONING INSTITUTE - Software Engineer (Feb 2025 - Present)\n' +
            '  - Building automated captioning systems that meet accessibility and broadcast standards; redesigned testing and improved frontend performance.\n\n' +
            'UNCOMMON GIVING - Software Engineer (2023 - Present)\n' +
            '  - Full-stack development (TypeScript, Python, Flutter); optimized CI/CD to parallelize tasks and reduce deployment time.\n\n' +
            'MUSC - Systems Programmer II (2023 - 2025)\n' +
            '  - Led full-stack modernization to .NET Core and improved data workflows and maintainability.\n\n' +
            'DISTRICT 186 - Computer Programmer (2022 - 2023)\n' +
            '  - Built admin tools using PHP and Oracle SQL to streamline staff and student workflows.',

        'education.txt':
            'B.S. Computer Science\n' +
            'MSCS - in progress\n' +
            'Relevant coursework and project details available on LinkedIn and resume.',

        'resume.txt':
            'Full resume (latest): https://docs.google.com/document/d/1qCF9Oe2GXS9ayBZgxiA9xq1ZOHGdxF9deOKG7T9Cwvg/edit?usp=sharing\n' +
            'PDF export available via the Resume (PDF) link in contact.txt',

        'links.txt':
            'LinkedIn: https://www.linkedin.com/in/jake-morgan-\n' +
            'Portfolio: /\n' +
            'GitHub: https://github.com/jmorgan3142001\n' +
            'Email: jmorgan3142001@gmail.com',

        'pets.txt':
            'Ashe & Rhaenyra - the Staffies!\n' +
            '  - Ashe: blue-gray, muscular, affectionate, and convinced he is the main character.\n' +
            '  - Rhaenyra (Rainy): tan, lean, energetic, and loves sunbeams and snacks.\n' +
            'Both are affectionate and a big part of life outside work.',

        'hobbies.txt':
            'Hobbies and interests:\n' +
            '  - Tech and coding: always tinkering with small tools or side projects\n' +
            '  - Gaming: competitive, co-op, and strategy play\n' +
            '  - Hiking and biking: exploring trails and moving fast\n' +
            '  - Singing: I love to sing, especially in musicals, and always have\n' +
            '  - Chess: slow, strategic play that sharpens thinking\n' +
            '  - Building and crafting: miniatures, models, and hands-on projects\n' +
            '  - Working out: consistency and progression for mental clarity\n' +
            '  - Rock climbing: body-centered problem solving\n' +
            '  - Diving: calm, focused resets',

        'games.txt':
            'Gaming preferences:\n' +
            '  - Genres: sci-fi, fantasy, strategy, co-op, competitive\n' +
            '  - Favorite games: Hearts of Iron IV, Total War, ESO, Battlefield, Holdfast\n' +
            '  - Playstyle: I enjoy meaningful decision-making, teamwork, and games that reward planning and skill.',

        'movies.txt':
            'Movies and genres I like:\n' +
            '  - Science fiction\n' +
            '  - Fantasy\n' +
            '  - Westerns\n' +
            '  - War movies\n' +
            '  - Horror\n' +
            '  - Musicals\n' +
            '  - Favorites: Star Wars, Lord of the Rings, Tangled, Evil Dead Rising',

        'personal.txt':
            'Personal notes:\n' +
            '  - Based near Charleston\n' +
            '  - I value clean design and clean code\n' +
            '  - I balance tech, fitness, and life with two staffies\n' +
            '  - Always open to learning or building something new',

        'sys_log.log': `[INFO] System initialized at ${startTime.toISOString()}.\n` +
                        `[WARN] Caffeine levels critical.\n` +
                        `[INFO] Portfolio v2.1 loaded.\n`
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
        line.innerHTML = text;
        terminalOutput.appendChild(line);
        terminalBody.scrollTop = terminalBody.scrollHeight;
    };

    const asciiArt = {
        cat: `
            .: :=:             .-=..                                             
            :. ::=:::=====--:=+:.=                                              
            :=. . . :== =. ..::.+.                                              
                --.   .+:...:    .=.                                               
            .:.  :....=.:=-.    =                                               
            .:..:.=*-:.  :==- ::  .                                              
            :=:-. ..:=   :::..=+-.:                                              
            .-:.::=. .  .  ::..  .:                                              
            .:- .::.:** -===-  .:::                                             
                :-:--: :::.:==-::: :  .:::.                                        
            .::=-=-.      :=-=:.=:   ::=-==:                                     
            ::=:::  .::   :::+:  := + =*:::=:                                   
            .::*.:*+.    :=+==   :+ .: :+..-:--                                  
                *:.            .::   *..--  :=.+=                                 
                =.=:-=-     .:::   .=::+=-  .:.::-                                
                .:.    :::.       ::.:-:::   . =.:.                               
                =-:    :.:*.   ..  =:   :       ::                               
                ..:.... -===.           .   :. -.-.                              
                .:.-  .  .====-.  .-=.  :.  .:= .==.                             
                .-:         .=.      .-+::::::-::.=.                            
                    + :    ::     .:..  .:#:+  ++   :=                            
                    + .    =.           ::#:+.    .:::                           
                    ::-    :     .:.   .#-=      :-=-                           
                    :::        .   -.  ++:::.  :::::=.                          
                    .::::-=. : ...    ::    .*=-:   .*.                         
                        ::     :+           =       ..   *.                        
                        =.-   :: : .. .:     .. .==..:.:.:                        
                        +-..:#:=::. :.      .      .:: :=                        
                        : . .=+:   .=::         .:*.:: .+                        
                        =.  # .   .=:--... .   :+: :-...=.                       
                        - . .. .:. :=:=--:::...+:.:::.:===                        
                    .::.. -.      .. .:   .::..::-: :=.                         
                        :::-:*.     .:- =.:-=-:..:-  ==:                           
                            + :: -      .-:.:.   ==.                             
                                .         .=-::-=:                                 
        `, 
        dog: `
                    :=*#***+-.                         .-=*##**=:.                           
                -#-         :*=:.                  :%*.         :#*:                        
                *%*        :.     :%*             =%.                **                      
                .#=@+       :          +*        :%=.                    =%                    
            :+  .-              ..=  .-:=#=::+.                        .#                   
            =:-.*=-=+:          :+.     .=.:==  ::                        ::                 
            .%-.**%    .=       =%%@=    .:.+:+..                           #:                
            :=  :==...@* := :      ::     :  .-=   .                         :*                
        :*    #--:.*%**   +:      .    +   *+    .+.  .                    *                
        +:.    =%=+*#+     +    =:   .=.:.. *      ==  *.                  *                
        +      :.          =:        :%.%  :.       :=                    +=                
        :%#*::::*.    .:            .*.=    :        .           .       :+                 
            ==     :*.   .           .*-*                          -    + *=                  
            -+.     =*.       *   .+: #      -                   %.  -=:==                   
                -*:     :+*+++*= -==  .%:     #                   *. .:=.=*                    
                    :=*%=     =:     -:      :#                 .%*      *                     
                    .%:                  .%.                ** =+:::  =-                    
                        *.                :#                :%.    :=-:=:*                    
                        .#               :                :@:      .@:*:%:                    
                        =:                             ++.          :::                      
                        -#                          +=                                      
                            :*+                     .*                                        
                            :=:=:               .*:                                         
                            :=     %-*#@##%:   :%=                                          
                            *:     %.  :#.      #.                                          
                            *.      *.   :-     .%                                           
                            %*   *: *=    ::     :#                                           
                            :=-::::*.     *      --                                           
                            :=--        **.     *:                                           
                                        *  * :.-#                                            
                                        *-= .**-                                             
                                            ...                                                
        `  
    };

    // --- Command Logic ---

    const commands = {
        'help': () => 'Commands: help, ls, cd [page], cat [file], grep [term] [file], uptime, ping [host], theme [dark|light|matrix], clear, exit (some hidden)\nWarning: User <i>really</i> is not privileged. <u>Do not</u> attempt \'sudo\' commands.',
        
        'ls': () => Object.keys(fileSystem).join('  '),
                
        'uptime': () => `${new Date().toLocaleTimeString()} ${getUptime()}, 1 user, ${getSystemStats()}`,
        
        'date': () => new Date().toString(),
        
        'history': () => commandHistory.map((cmd, i) => ` ${i + 1}  ${cmd}`).join('\n'),

        'clear': () => {
            terminalOutput.innerHTML = '';
            // Redisplay welcome text after clear
            const welcomeLine = document.createElement('div');
            welcomeLine.className = 'text-mono small mb-3 text-secondary';
            welcomeLine.innerHTML = WELCOME_TEXT;
            terminalOutput.appendChild(welcomeLine);
            return ''; 
        },

        'cat': (args) => {
            if (!args[0]) {
                return `Usage: cat [filename]\n${asciiArt.cat}`;
            }
            if (fileSystem[args[0]]) return fileSystem[args[0]];
            return `cat: ${args[0]}: No such file or directory`;
        },

        'dog': () => {
            return `Usage: Woof and whine!\n${asciiArt.dog}`;
        },

        'pit': () => {
            return `Usage: Pitbulls are banned in 12 countries! (but not this terminal)\n${asciiArt.dog}`;
        },

        'staff': () => {
            return `Usage: This Staffy's tail is a registered weapon!\n${asciiArt.dog}`;
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
            if (args[0] === 'coffee') return 'make: *** No rule to make target `coffee`. Stop.\n(Hint: Make it yourself!)';
            return 'make: *** No targets specified and no makefile found. Stop.';
        },

        'sudo': () => {
            setTimeout(() => {
                window.open('https://www.youtube.com/watch?v=dQw4w9WgXcQ', '_blank'); 
            }, 1500); // Rick Rolled :D
            return 'Access Denied: User is not in the sudoers file. Should have listened! Goodbye :)';
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