let editor;
let currentChallengeId = null;

window.onload = function() {
    editor = CodeMirror.fromTextArea(document.getElementById("code-editor"), {
        lineNumbers: true,
        mode: "python",
        theme: "default",
        viewportMargin: Infinity,
        indentUnit: 4
    });
};

function setTerminalTheme(themeName) {
    editor.setOption("theme", themeName);
}

function toggleContextPane() {
    const leftPane = document.getElementById('left-pane');
    const icon = document.getElementById('toggle-icon');
    const workspace = document.getElementById('workspace');
    const mobileShowBtn = document.getElementById('mobile-show-context');
    
    leftPane.classList.toggle('collapsed');
    
    if (leftPane.classList.contains('collapsed')) {
        icon.classList.remove('bi-chevron-left');
        icon.classList.add('bi-chevron-right');
        workspace.classList.add('left-collapsed');
        if(mobileShowBtn) mobileShowBtn.style.display = 'block';
    } else {
        icon.classList.remove('bi-chevron-right');
        icon.classList.add('bi-chevron-left');
        workspace.classList.remove('left-collapsed');
        if(mobileShowBtn) mobileShowBtn.style.display = 'none';
    }

    setTimeout(() => {
        editor.refresh();
    }, 310);
}

function loadChallenge(id, title, desc, starter, lang, diff, btnElement) {
    currentChallengeId = id;
    
    document.getElementById('welcome-state').style.display = 'none';
    document.getElementById('active-state').style.display = 'block';
    
    document.getElementById('challenge-title').innerText = title;
    document.getElementById('challenge-desc').innerText = desc;
    const diffEl = document.getElementById('challenge-diff');
    diffEl.innerText = diff.toUpperCase();

    diffEl.className = 'badge text-mono'; 

    if (diff === 'Easy') {
        diffEl.classList.add('bg-success', 'bg-opacity-10', 'text-success');
    } else if (diff === 'Medium') {
        diffEl.classList.add('bg-warning', 'bg-opacity-10', 'text-warning');
    } else if (diff == 'Hard') { 
        diffEl.classList.add('bg-danger', 'bg-opacity-10', 'text-danger');
    }
    
    editor.setValue(starter);
    document.getElementById('run-btn').disabled = false;
    
    document.getElementById('console-output').innerHTML = '<span class="text-secondary">> Module Loaded: ' + id + '</span>';

    const bsCollapse = new bootstrap.Collapse(document.getElementById('challengeList'), {
        toggle: false
    });
    bsCollapse.hide();

    document.querySelectorAll('.challenge-btn').forEach(b => b.classList.remove('active'));
    if(btnElement) btnElement.classList.add('active');
}

async function runCode() {
    if (!currentChallengeId) return;

    const outputDiv = document.getElementById('console-output');
    const runBtn = document.getElementById('run-btn');
    
    runBtn.disabled = true;
    runBtn.innerHTML = '<span class="spinner-border spinner-border-sm me-2"></span>';
    
    outputDiv.innerHTML = '<span class="text-accent">> Compiling on remote container...</span>';

    const userCode = editor.getValue();

    try {
        const response = await fetch('/challenges/run', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ 
                challenge_id: currentChallengeId, 
                user_code: userCode 
            })
        });

        const data = await response.json();
        outputDiv.innerHTML = ''; 

        if (data.passed) {
            outputDiv.innerHTML = `<div class="mb-2"><span class="badge bg-success">PASSED</span></div>` + formatOutput(data.output);
        } else {
            outputDiv.innerHTML = `<div class="mb-2"><span class="badge bg-danger">FAILED</span></div>` + formatOutput(data.output);
        }
    } catch (e) {
        outputDiv.innerHTML = '<span class="text-danger fw-bold">>> CONNECTION_ERROR</span>';
    } finally {
        runBtn.disabled = false;
        runBtn.innerHTML = '<i class="bi bi-play-fill me-1"></i> RUN';
    }
}

function formatOutput(rawText) {
    if (!rawText) return '';
    return rawText.replace(/\n/g, '<br/>')
                    .replace(/✓/g, '<span class="text-success fw-bold">✓</span>')
                    .replace(/✗/g, '<span class="text-danger fw-bold">✗</span>');
}