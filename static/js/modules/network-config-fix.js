/**
 * Network Configuration Fix Module
 * Fix collection of device-type specific configuration in the network module
 */

/**
 * Collect device-type specific configuration
 * @param {HTMLElement} device - Network device DOM element
 * @param {Object} deviceConfig - Device configuration object
 * @param {string} deviceType - Device type
 */
function collectDeviceSpecificConfig(device, deviceConfig, deviceType) {
    console.log(`Collect device specific configuration: ${deviceType}`);

    switch (deviceType) {
        case 'wifis':
            // WiFi specific configuration
            collectWiFiConfig(device, deviceConfig);
            break;

        case 'bridges':
            // Bridge specific configuration
            collectBridgeConfig(device, deviceConfig);
            break;

        case 'bonds':
            // Bond specific configuration
            collectBondConfig(device, deviceConfig);
            break;

        case 'vlans':
            // VLAN specific configuration
            collectVLANConfig(device, deviceConfig);
            break;

        case 'tunnels':
            // Tunnel specific configuration
            collectTunnelConfig(device, deviceConfig);
            break;

        case 'vrfs':
            // VRF specific configuration
            collectVRFConfig(device, deviceConfig);
            break;
    }
}

/**
 * Collect WiFi device specific configuration
 */
function collectWiFiConfig(device, deviceConfig) {
    const accessPoints = {};
    const apItems = device.querySelectorAll('.access-point-item');
    let hasAnyAccessPoint = false;

    apItems.forEach(apItem => {
        const ssidInput = apItem.querySelector('.ap-ssid-input');
        const passwordInput = apItem.querySelector('.ap-password-input');
        const hiddenInput = apItem.querySelector('.ap-hidden-input');
        const bandInput = apItem.querySelector('.ap-band-input');

        if (ssidInput && ssidInput.value.trim()) {
            const ssid = ssidInput.value.trim();
            accessPoints[ssid] = {};
            hasAnyAccessPoint = true;

            if (passwordInput && passwordInput.value) {
                accessPoints[ssid].password = passwordInput.value;
            }

            if (hiddenInput && hiddenInput.value === 'true') {
                accessPoints[ssid].hidden = true;
            }

            if (bandInput && bandInput.value) {
                accessPoints[ssid].band = bandInput.value;
            }
        }
    });

    if (hasAnyAccessPoint) {
        // Use Netplan-compliant hyphenated key
        deviceConfig['access-points'] = accessPoints;
        console.log('Collected WiFi access-points:', accessPoints);
    }
}

/**
 * Collect Bridge device specific configuration
 */
function collectBridgeConfig(device, deviceConfig) {
    const bridgeInterfacesInput = device.querySelector('.bridge-interfaces-input');
    const bridgeStpInput = device.querySelector('.bridge-stp-input');
    const bridgeForwardDelayInput = device.querySelector('.bridge-forward-delay-input');
    const bridgeAgeingTimeInput = device.querySelector('.bridge-ageing-time-input');

    if (bridgeInterfacesInput && bridgeInterfacesInput.value.trim()) {
        const interfaces = bridgeInterfacesInput.value.split(',').map(iface => iface.trim()).filter(iface => iface);
        if (interfaces.length > 0) {
            deviceConfig.interfaces = interfaces;
            console.log('Collected bridge interfaces:', interfaces);
        }
    }

    const parameters = {};
    let hasParameters = false;

    if (bridgeStpInput) {
        parameters.stp = bridgeStpInput.value === 'true';
        hasParameters = true;
    }

    if (bridgeForwardDelayInput && bridgeForwardDelayInput.value) {
        parameters['forward-delay'] = parseInt(bridgeForwardDelayInput.value);
        hasParameters = true;
    }

    if (bridgeAgeingTimeInput && bridgeAgeingTimeInput.value) {
        parameters['ageing-time'] = parseInt(bridgeAgeingTimeInput.value);
        hasParameters = true;
    }

    if (hasParameters) {
        deviceConfig.parameters = parameters;
        console.log('Collected bridge parameters:', parameters);
    }
}

/**
 * Collect Bond device specific configuration
 */
function collectBondConfig(device, deviceConfig) {
    const bondInterfacesInput = device.querySelector('.bond-interfaces-input');
    const bondModeInput = device.querySelector('.bond-mode-input');
    const bondPrimaryInput = device.querySelector('.bond-primary-input');
    const bondMiiInput = device.querySelector('.bond-mii-input');
    const bondLacpRateInput = device.querySelector('.bond-lacp-rate-input');
    const bondMinLinksInput = device.querySelector('.bond-min-links-input');

    if (bondInterfacesInput && bondInterfacesInput.value.trim()) {
        const interfaces = bondInterfacesInput.value.split(',').map(iface => iface.trim()).filter(iface => iface);
        if (interfaces.length > 0) {
            deviceConfig.interfaces = interfaces;
            console.log('Collected bond interfaces:', interfaces);
        }
    }

    const parameters = {};
    let hasParameters = false;

    if (bondModeInput && bondModeInput.value) {
        parameters.mode = bondModeInput.value;
        hasParameters = true;
    }

    if (bondPrimaryInput && bondPrimaryInput.value.trim()) {
        parameters.primary = bondPrimaryInput.value.trim();
        hasParameters = true;
    }

    if (bondMiiInput && bondMiiInput.value) {
        parameters['mii-monitor-interval'] = parseInt(bondMiiInput.value);
        hasParameters = true;
    }

    if (bondLacpRateInput && bondLacpRateInput.value) {
        parameters['lacp-rate'] = bondLacpRateInput.value;
        hasParameters = true;
    }

    if (bondMinLinksInput && bondMinLinksInput.value) {
        parameters['min-links'] = parseInt(bondMinLinksInput.value);
        hasParameters = true;
    }

    if (hasParameters) {
        deviceConfig.parameters = parameters;
        console.log('Collected bond parameters:', parameters);
    }
}

/**
 * Collect VLAN device specific configuration
 */
function collectVLANConfig(device, deviceConfig) {
    const vlanIdInput = device.querySelector('.vlan-id-input');
    const vlanLinkInput = device.querySelector('.vlan-link-input');

    if (vlanIdInput && vlanIdInput.value) {
        deviceConfig.id = parseInt(vlanIdInput.value);
        console.log('Collected VLAN ID:', deviceConfig.id);
    }

    if (vlanLinkInput && vlanLinkInput.value.trim()) {
        deviceConfig.link = vlanLinkInput.value.trim();
        console.log('Collected VLAN parent link:', deviceConfig.link);
    }
}

/**
 * Collect Tunnel device specific configuration
 */
function collectTunnelConfig(device, deviceConfig) {
    const tunnelModeInput = device.querySelector('.tunnel-mode-input');
    const tunnelLocalInput = device.querySelector('.tunnel-local-input');
    const tunnelRemoteInput = device.querySelector('.tunnel-remote-input');
    const tunnelKeyInput = device.querySelector('.tunnel-key-input');
    const tunnelTtlInput = device.querySelector('.tunnel-ttl-input');
    const tunnelPortInput = device.querySelector('.tunnel-port-input');

    if (tunnelModeInput && tunnelModeInput.value) {
        deviceConfig.mode = tunnelModeInput.value;
    }

    if (tunnelLocalInput && tunnelLocalInput.value.trim()) {
        deviceConfig.local = tunnelLocalInput.value.trim();
    }

    if (tunnelRemoteInput && tunnelRemoteInput.value.trim()) {
        deviceConfig.remote = tunnelRemoteInput.value.trim();
    }

    if (tunnelKeyInput && tunnelKeyInput.value) {
        deviceConfig.key = parseInt(tunnelKeyInput.value);
    }

    if (tunnelTtlInput && tunnelTtlInput.value) {
        deviceConfig.ttl = parseInt(tunnelTtlInput.value);
    }

    if (tunnelPortInput && tunnelPortInput.value) {
        deviceConfig.port = parseInt(tunnelPortInput.value);
    }

    console.log('Collected tunnel configuration:', {
        mode: deviceConfig.mode,
        local: deviceConfig.local,
        remote: deviceConfig.remote
    });
}

/**
 * Collect VRF device specific configuration
 */
function collectVRFConfig(device, deviceConfig) {
    const vrfTableInput = device.querySelector('.vrf-table-input');
    const vrfRoutesInput = device.querySelector('.vrf-routes-input');

    if (vrfTableInput && vrfTableInput.value) {
        deviceConfig.table = parseInt(vrfTableInput.value);
        console.log('Collected VRF table ID:', deviceConfig.table);
    }

    if (vrfRoutesInput && vrfRoutesInput.value.trim()) {
        const routes = vrfRoutesInput.value.split(',').map(route => route.trim()).filter(route => route);
        if (routes.length > 0) {
            if (!deviceConfig.routes) {
                deviceConfig.routes = [];
            }
            routes.forEach(route => {
                deviceConfig.routes.push({ to: route });
            });
            console.log('Collected VRF routes:', deviceConfig.routes);
        }
    }
}

/**
 * Extend the original buildConfig function to collect device-specific configuration
 */
function fixNetworkConfig() {
    console.log('Applying network configuration fixes...');

    // Save the original buildConfig function
    const originalBuildConfig = window.ConfigManager.buildConfig;

    // Create an enhanced buildConfig function
    window.ConfigManager.buildConfig = function() {
        // Call the original buildConfig to get the base configuration
        const config = originalBuildConfig.call(this);

        // Re-collect network device configuration
        const networkDevices = document.querySelectorAll('.network-device');
        if (networkDevices.length > 0) {
            networkDevices.forEach(device => {
                const deviceType = device.querySelector('.device-type-input')?.value || 'ethernets';
                const deviceName = device.querySelector('.device-name-input')?.value || 'eth0';

                if (deviceType && deviceName && config.autoinstall && config.autoinstall.network) {
                    // Check if the device already exists in the configuration
                    if (config.autoinstall.network[deviceType] &&
                        config.autoinstall.network[deviceType][deviceName]) {

                        // Collect device-specific configuration and merge it into the existing config
                        collectDeviceSpecificConfig(
                            device,
                            config.autoinstall.network[deviceType][deviceName],
                            deviceType
                        );
                    }
                }
            });
        }

        return config;
    };

    console.log('Network configuration fixes applied');
}

// Initialize fixes
document.addEventListener('DOMContentLoaded', function() {
    console.log('Initializing network configuration fix module...');
    fixNetworkConfig();
});

// Export functions for use by other modules
window.NetworkConfigFix = {
    collectDeviceSpecificConfig,
    fixNetworkConfig
};
