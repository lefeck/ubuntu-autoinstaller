/**
 * User data processing module
 * Handles user-data preview and generation
 */

// Preview configuration
function previewUserData() {
    console.log('Previewing configuration...');

    try {
        // Get config object
        const config = window.ConfigManager.buildConfig();

        // Show configuration preview
        const configPreview = document.getElementById('configPreview');
        const configStatus = document.getElementById('configStatus');

        if (configPreview) {
            // Format JSON and display
            configPreview.innerHTML = `<pre>${JSON.stringify(config, null, 2)}</pre>`;
            configPreview.style.display = 'block';

            // Update status
            if (configStatus) {
                configStatus.innerHTML = '<div class="status-success">Configuration generated successfully!</div>';
            }
        }
    } catch (error) {
        console.error('Configuration preview failed:', error);
        const configStatus = document.getElementById('configStatus');
        if (configStatus) {
            configStatus.innerHTML = `<div class="status-error">Configuration preview failed: ${error.message}</div>`;
        }
    }
}

// Generate user-data
function generateUserData() {
    console.log('Generating user-data...');

    try {
        // Get config object
        const config = window.ConfigManager.buildConfig();

        // Prefer backend generation to ensure password hashing and schema
        fetch(`${API_BASE}/userdata/generate`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ config })
        })
        .then(res => res.json())
        .then(result => {
            const userdataResult = document.getElementById('userdataResult');
            const userdataStatus = document.getElementById('userdataStatus');
            const userDataContent = document.getElementById('userDataContent');

            if (result && result.success && (result.userData || result['user-data'])) {
                const yamlContent = result.userData || result['user-data'];
                if (userdataResult) {
                    userdataResult.innerHTML = `<pre>${yamlContent}</pre>`;
                    userdataResult.style.display = 'block';
                }
                if (userdataStatus) {
                    userdataStatus.innerHTML = '<div class="status-success">User-data generated successfully!</div>';
                }
                if (userDataContent) {
                    userDataContent.value = yamlContent;
                }
            } else {
                throw new Error(result && result.error ? result.error : 'Failed to generate user-data');
            }
        })
        .catch(err => {
            console.warn('Backend generation failed, fallback to local:', err.message);
            const yamlContent = convertToYAML(config);
            const userdataResult = document.getElementById('userdataResult');
            const userdataStatus = document.getElementById('userdataStatus');
            const userDataContent = document.getElementById('userDataContent');
            if (userdataResult) {
                userdataResult.innerHTML = `<pre>${yamlContent}</pre>`;
                userdataResult.style.display = 'block';
            }
            if (userdataStatus) {
                userdataStatus.innerHTML = `<div class="status-success">User-data generated locally</div>`;
            }
            if (userDataContent) {
                userDataContent.value = yamlContent;
            }
        });
    } catch (error) {
        console.error('User-data generation failed:', error);
        const userdataStatus = document.getElementById('userdataStatus');
        if (userdataStatus) {
            userdataStatus.innerHTML = `<div class="status-error">User-data generation failed: ${error.message}</div>`;
        }
    }
}

// Simple JSON to YAML conversion function (naive implementation)
function convertToYAML(json) {
    // Note: this is a simplified implementation; consider a real YAML library in production
    let yaml = '#cloud-config\n';

    // Extract main config from autoinstall structure
    if (json && json.autoinstall) {
        const autoinstall = json.autoinstall;

        yaml += yamlify(autoinstall, 0);
    }

    return yaml;
}

// Recursively convert JSON object to YAML
function yamlify(obj, depth = 0) {
    let result = '';
    const indent = '  '.repeat(depth);

    if (Array.isArray(obj)) {
        if (obj.length === 0) {
            result += '[]\n';
        } else {
            for (const item of obj) {
                if (typeof item === 'object' && item !== null) {
                    result += `${indent}- ${yamlify(item, depth + 1)}`;
                } else {
                    result += `${indent}- ${item}\n`;
                }
            }
        }
    } else if (typeof obj === 'object' && obj !== null) {
        for (const [key, value] of Object.entries(obj)) {
            if (value === undefined || value === null) {
                continue;
            }

            if (typeof value === 'object') {
                if (Array.isArray(value) && value.length === 0) {
                    result += `${indent}${key}: []\n`;
                } else if (Object.keys(value).length === 0 && !Array.isArray(value)) {
                    result += `${indent}${key}: {}\n`;
                } else {
                    result += `${indent}${key}:\n${yamlify(value, depth + 1)}`;
                }
            } else {
                result += `${indent}${key}: ${value}\n`;
            }
        }
    } else {
        result += `${obj}\n`;
    }

    return result;
}

// Expose functions to global scope
window.previewUserData = previewUserData;
window.generateUserData = generateUserData;

// Add CSS styles for preview status blocks
document.addEventListener('DOMContentLoaded', function() {
    // Success and error status styles
    const style = document.createElement('style');
    style.textContent = `
        .status-success {
            padding: 10px;
            background-color: #d4edda;
            color: #155724;
            border-radius: 4px;
            margin-bottom: 10px;
        }
        
        .status-error {
            padding: 10px;
            background-color: #f8d7da;
            color: #721c24;
            border-radius: 4px;
            margin-bottom: 10px;
        }
        
        pre {
            background-color: #f5f5f5;
            padding: 15px;
            border-radius: 4px;
            overflow-x: auto;
            white-space: pre-wrap;
            word-break: break-word;
        }
    `;
    document.head.appendChild(style);
});
