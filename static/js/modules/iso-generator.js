/**
 * ISO Generator Module
 * Handles ISO file generation, upload, build process, and status monitoring
 */

// Global variables for ISO generation
let currentBuildID = null;
let buildLogs = [];
let buildStatus = 'idle';

// API base URL
const API_BASE = '/api/v1';

/**
 * Initialize ISO generation form
 */
function initISOForm() {
    // Initialize source type toggle
    initSourceTypeToggle();
    
    // Initialize file upload handlers
    initFileUploadHandlers();
}

/**
 * Initialize source type toggle
 */
function initSourceTypeToggle() {
    // This function is called by the radio button onchange event
    // No need to add event listeners here
}

/**
 * Toggle source type between local file and download
 */
function toggleSourceType() {
    const sourceType = document.querySelector('input[name="sourceType"]:checked').value;
    const localIsoSection = document.getElementById('localIsoSection');
    const downloadIsoSection = document.getElementById('downloadIsoSection');
    
    if (sourceType === 'local') {
        localIsoSection.style.display = 'block';
        downloadIsoSection.style.display = 'none';
    } else {
        localIsoSection.style.display = 'none';
        downloadIsoSection.style.display = 'block';
    }
}

/**
 * Initialize file upload handlers
 */
function initFileUploadHandlers() {
    const fileInput = document.getElementById('isoFileInput');
    if (fileInput) {
        fileInput.addEventListener('change', function(e) {
            handleFileSelect(e.target);
        });
    }
}

/**
 * Handle file selection
 */
async function handleFileSelect(input) {
    const file = input.files[0];
    if (!file) return;
    
    // Validate file type
    if (!file.name.toLowerCase().endsWith('.iso')) {
        showStatus('userdataStatus', 'error', 'Please select a valid ISO file');
        return;
    }
    
    // Validate file size (max 10GB)
    const maxSize = 10 * 1024 * 1024 * 1024; // 10GB
    if (file.size > maxSize) {
        showStatus('userdataStatus', 'error', 'File size exceeds 10GB limit');
        return;
    }
    
    // Show file info
    showFileInfo(file);
    
    // Upload file
    try {
        const uploadResult = await uploadFile(file);
        if (uploadResult.success) {
            showStatus('userdataStatus', 'success', 'File uploaded successfully');
            // Store file info for later use
            window.uploadedFileInfo = {
                filename: file.name,
                size: file.size,
                uploadId: uploadResult.uploadId
            };
        } else {
            showStatus('userdataStatus', 'error', `Upload failed: ${uploadResult.error}`);
        }
    } catch (error) {
        console.error('File upload error:', error);
        showStatus('userdataStatus', 'error', 'File upload failed');
    }
}

/**
 * Upload file to server
 */
function uploadFile(file) {
    return new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest();
        const formData = new FormData();
        formData.append('file', file);
        
        // Show upload progress
        xhr.upload.onprogress = function (e) {
            if (e.lengthComputable) {
                const percentComplete = (e.loaded / e.total) * 100;
                updateUploadProgress(percentComplete);
            }
        };
        
        xhr.onload = function () {
            if (xhr.status === 200) {
                try {
                    const response = JSON.parse(xhr.responseText);
                    resolve(response);
                } catch (e) {
                    resolve({ success: false, error: 'Invalid response format' });
                }
            } else {
                resolve({ success: false, error: `HTTP ${xhr.status}: ${xhr.statusText}` });
            }
        };
        
        xhr.onerror = function () { 
            reject(new Error('Network error during upload')); 
        };
        
        xhr.open('POST', `${API_BASE}/upload`);
        xhr.send(formData);
    });
}

/**
 * Update upload progress
 */
function updateUploadProgress(percent) {
    const progressBar = document.getElementById('uploadProgressFill');
    const progressText = document.getElementById('uploadProgressText');
    const progressContainer = document.getElementById('uploadProgressContainer');
    
    if (progressBar) {
        progressBar.style.width = `${percent}%`;
    }
    if (progressText) {
        progressText.textContent = `${Math.round(percent)}%`;
    }
    if (progressContainer) {
        progressContainer.style.display = 'block';
    }
}

/**
 * Show file information
 */
function showFileInfo(file) {
    const fileInfo = document.getElementById('fileInfo');
    const fileName = document.getElementById('fileName');
    const fileSize = document.getElementById('fileSize');
    
    if (fileInfo && fileName && fileSize) {
        fileName.textContent = file.name;
        fileSize.textContent = formatFileSize(file.size);
        fileInfo.style.display = 'block';
    }
}

/**
 * Remove selected file
 */
function removeSelectedFile() {
    const fileInput = document.getElementById('isoFileInput');
    const fileInfo = document.getElementById('fileInfo');
    const progressContainer = document.getElementById('uploadProgressContainer');
    
    if (fileInput) fileInput.value = '';
    if (fileInfo) fileInfo.style.display = 'none';
    if (progressContainer) progressContainer.style.display = 'none';
    
    // Clear uploaded file info
    window.uploadedFileInfo = null;
    
    showStatus('userdataStatus', 'info', 'File selection cleared');
}

/**
 * Format file size in human readable format
 */
function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

/**
 * Generate ISO function - main entry point
 */
async function generateISO() {
    try {
        // Validate configuration first
        if (!await validateAllConfigs()) {
            return;
        }
        
        // Build configuration
        const config = buildCompleteConfig();
        if (!config) {
            showStatus('userdataStatus', 'error', 'Failed to build configuration');
            return;
        }
        
        // Show build progress
        isoShowBuildProgress();

        // Start build process
        const buildResult = await startBuildProcess(config);
        if (buildResult.success) {
            currentBuildID = buildResult.buildID;
            showStatus('userdataStatus', 'success', 'Build started successfully');
            
            // Start polling for status
            pollBuildStatus(buildResult.buildID);
        } else {
            showStatus('userdataStatus', 'error', `Build failed: ${buildResult.error}`);
            isoHideBuildProgress();
        }
        
    } catch (error) {
        console.error('ISO generation error:', error);
        showStatus('userdataStatus', 'error', 'ISO generation failed');
        isoHideBuildProgress();
    }
}

/**
 * Validate all configurations
 */
async function validateAllConfigs() {
    const validations = [];
    
    // Validate basic config
    if (window.ConfigManager && window.ConfigManager.validateConfig) {
        validations.push(window.ConfigManager.validateConfig());
    }
    
    // Validate storage config
    if (window.StorageManager && window.StorageManager.validateStorageConfig) {
        validations.push(window.StorageManager.validateStorageConfig());
    }
    
    // Validate network config
    if (window.NetworkManager && window.NetworkManager.validateNetworkConfig) {
        validations.push(window.NetworkManager.validateNetworkConfig());
    }
    
    // Wait for all validations
    const results = await Promise.all(validations);
    
    // Check if all validations passed
    const allValid = results.every(result => result.valid);
    if (!allValid) {
        const allErrors = results.flatMap(result => result.errors || []);
        showStatus('userdataStatus', 'error', `Configuration validation failed: ${allErrors.join(', ')}`);
        return false;
    }
    
    return true;
}

/**
 * Build complete configuration
 */
function buildCompleteConfig() {
    try {
        const config = {
            basic: {},
            apt: {},
            network: {},
            storage: {},
            ssh: {},
            advanced: {}
        };
        
        // Get basic config
        if (window.ConfigManager && window.ConfigManager.getCurrentConfig) {
            const basicConfig = window.ConfigManager.getCurrentConfig();
            config.basic = basicConfig.basic || {};
            config.apt = basicConfig.apt || {};
        }
        
        // Get storage config
        if (window.StorageManager && window.StorageManager.getStorageConfig) {
            config.storage = window.StorageManager.getStorageConfig();
        }
        
        // Get network config
        if (window.NetworkManager && window.NetworkManager.getNetworkConfig) {
            const networkConfig = window.NetworkManager.getNetworkConfig();
            config.network = networkConfig;
        }
        
        // Add source information
        const sourceType = document.querySelector('input[name="sourceType"]:checked').value;
        config.source = {
            type: sourceType,
            file: window.uploadedFileInfo || null,
            codename: document.getElementById('codename')?.value || 'jammy'
        };
        
        // Add other form data
        config.destinationISO = document.getElementById('destinationISO')?.value || '';
        config.packageList = document.getElementById('packageListTextarea')?.value || '';
        config.userDataContent = document.getElementById('userDataContent')?.value || '';
        config.useHWEKernel = document.getElementById('useHWEKernelCheckbox')?.checked || false;
        config.md5Checksum = document.getElementById('md5ChecksumCheckbox')?.checked || true;
        config.gpgVerify = document.getElementById('gpgVerifyCheckbox')?.checked || true;
        
        return config;
        
    } catch (error) {
        console.error('Configuration build error:', error);
        return null;
    }
}

/**
 * Start build process
 */
async function startBuildProcess(config) {
    try {
        const response = await fetch(`${API_BASE}/build/start`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(config)
        });
        
        if (response.ok) {
            const result = await response.json();
            return { success: true, buildID: result.buildID };
        } else {
            const error = await response.text();
            return { success: false, error: error };
        }
    } catch (error) {
        console.error('Build start error:', error);
        return { success: false, error: error.message };
    }
}

/**
 * Poll build status
 */
async function pollBuildStatus(buildID) {
    try {
        const response = await fetch(`${API_BASE}/build/status/${buildID}`);
        if (response.ok) {
            const status = await response.json();
            
            // Update build progress
            updateBuildProgress(status);
            
            // Check if build is complete
            if (status.status === 'completed') {
                showStatus('userdataStatus', 'success', 'Build completed successfully!');
                showDownloadSection(buildID);
                isoHideBuildProgress();
            } else if (status.status === 'failed') {
                showStatus('userdataStatus', 'error', `Build failed: ${status.error || 'Unknown error'}`);
                isoHideBuildProgress();
            } else {
                // Continue polling
                setTimeout(() => pollBuildStatus(buildID), 2000);
            }
        } else {
            throw new Error(`HTTP ${response.status}`);
        }
    } catch (error) {
        console.error('Status polling error:', error);
        showStatus('userdataStatus', 'error', 'Failed to get build status');
        isoHideBuildProgress();
    }
}

/**
 * Update build progress
 */
function updateBuildProgress(status) {
    // Update progress bar
    if (status.progress !== undefined) {
        const progressBar = document.getElementById('progressFill');
        if (progressBar) {
            progressBar.style.width = `${status.progress}%`;
        }
    }
    
    // Update steps
    if (status.steps) {
        updateStepsFromStatus(status.steps);
    }
    
    // Add log entries
    if (status.logs) {
        status.logs.forEach(log => {
            isoAddLogToUI(log.level, log.message);
        });
    }
}

/**
 * Update steps from status
 */
function updateStepsFromStatus(steps) {
    steps.forEach(step => {
        const stepElement = document.querySelector(`[data-step="${step.name}"]`);
        if (stepElement) {
            stepElement.classList.remove('active', 'completed');
            if (step.status === 'active') {
                stepElement.classList.add('active');
            } else if (step.status === 'completed') {
                stepElement.classList.add('completed');
            }
        }
    });
}

/**
 * Add log to UI
 */
function isoAddLogToUI(level, message) {
    const timestamp = new Date().toLocaleTimeString();
    const logEntry = document.createElement('div');
    logEntry.className = `log-entry ${level}`;
    logEntry.textContent = `[${timestamp}] ${message}`;
    
    const logContainer = document.getElementById('logContainer');
    if (logContainer) {
        logContainer.appendChild(logEntry);
        logContainer.scrollTop = logContainer.scrollHeight;
    }
    
    buildLogs.push({ level, message, timestamp });
}

/**
 * Show build progress
 */
function isoShowBuildProgress() {
    const buildProgress = document.getElementById('buildProgress');
    const buildLogsEl = document.getElementById('buildLogs');

    if (buildProgress) buildProgress.style.display = 'block';
    if (buildLogsEl) buildLogsEl.style.display = 'block';

    // Reset progress
    isoResetBuildProgress();
}

/**
 * Hide build progress
 */
function isoHideBuildProgress() {
    const buildProgress = document.getElementById('buildProgress');
    const buildLogsEl = document.getElementById('buildLogs');

    if (buildProgress) buildProgress.style.display = 'none';
    if (buildLogsEl) buildLogsEl.style.display = 'none';
}

/**
 * Reset build progress
 */
function isoResetBuildProgress() {
    // Reset steps
    document.querySelectorAll('.step').forEach(step => {
        step.classList.remove('active', 'completed');
    });
    
    // Reset progress bar
    const progressBar = document.getElementById('progressFill');
    if (progressBar) {
        progressBar.style.width = '0%';
    }
    
    // Clear logs
    const logContainer = document.getElementById('logContainer');
    if (logContainer) {
        logContainer.innerHTML = '';
    }
    
    buildLogs = [];
}

/**
 * Show download section
 */
function showDownloadSection(buildID) {
    currentBuildID = buildID;
    const downloadSection = document.getElementById('downloadSection');
    if (downloadSection) {
        downloadSection.style.display = 'block';
    }
}

/**
 * Download ISO function
 */
function downloadISO() {
    if (!currentBuildID) {
        showStatus('userdataStatus', 'error', 'No build ID available for download');
        return;
    }
    
    try {
        // Create download link
        const downloadUrl = `${API_BASE}/build/download/${currentBuildID}`;
        const link = document.createElement('a');
        link.href = downloadUrl;
        link.download = ''; // Let the server set the filename
        link.style.display = 'none';
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        
        showStatus('userdataStatus', 'success', 'Download started');
    } catch (error) {
        console.error('Download error:', error);
        showStatus('userdataStatus', 'error', `Download failed: ${error.message}`);
    }
}

// Export functions for use in other modules
window.ISOGenerator = {
    initISOForm,
    toggleSourceType,
    handleFileSelect,
    removeSelectedFile,
    generateISO,
    downloadISO,
    showBuildProgress: isoShowBuildProgress,
    hideBuildProgress: isoHideBuildProgress,
    resetBuildProgress: isoResetBuildProgress
};
