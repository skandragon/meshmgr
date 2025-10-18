// Meshtastic client store using Svelte 5 runes
import { Client } from '@meshtastic/core';
import { ISerialConnection } from '@meshtastic/transport-web-serial';
import type { Protobuf } from '@meshtastic/protobufs';

export interface DeviceInfo {
	myNodeNum?: number;
	firmwareVersion?: string;
	hwModel?: string;
	hasWifi?: boolean;
	hasBluetooth?: boolean;
	macAddress?: string;
	region?: string;
	modemPreset?: string;
	numOnlineLocalNodes?: number;
}

export interface MeshtasticState {
	isSupported: boolean;
	client: Client | null;
	connection: ISerialConnection | null;
	isConnected: boolean;
	isConnecting: boolean;
	error: string | null;
	deviceInfo: DeviceInfo | null;
	lastMessage: string | null;
}

class MeshtasticStore {
	private client = $state<Client | null>(null);
	private connection = $state<ISerialConnection | null>(null);

	isSupported = $state(false);
	isConnected = $state(false);
	isConnecting = $state(false);
	error = $state<string | null>(null);
	deviceInfo = $state<DeviceInfo | null>(null);
	lastMessage = $state<string | null>(null);

	constructor() {
		// Check if Web Serial API is supported
		if (typeof window !== 'undefined') {
			this.isSupported = 'serial' in navigator;
		}
	}

	async connect(): Promise<void> {
		if (!this.isSupported) {
			this.error = 'Web Serial API is not supported in this browser';
			throw new Error(this.error);
		}

		if (this.isConnected || this.isConnecting) {
			return;
		}

		try {
			this.isConnecting = true;
			this.error = null;

			// Create a new client
			this.client = new Client();

			// Create a serial connection
			this.connection = new ISerialConnection();

			// Set up event listeners
			this.setupEventListeners();

			// Connect to the device
			await this.connection.connect({
				baudRate: 115200,
				concurrentLogOutput: false
			});

			// Connect the client to the connection
			await this.client.connect({
				connection: this.connection
			});

			this.isConnected = true;

			// Request device info
			await this.requestDeviceInfo();
		} catch (err: any) {
			this.error = err.message || 'Failed to connect to Meshtastic device';
			this.isConnected = false;
			throw err;
		} finally {
			this.isConnecting = false;
		}
	}

	async disconnect(): Promise<void> {
		try {
			this.error = null;

			if (this.client) {
				this.client.disconnect();
				this.client = null;
			}

			if (this.connection) {
				await this.connection.disconnect();
				this.connection = null;
			}

			this.isConnected = false;
			this.deviceInfo = null;
			this.lastMessage = null;
		} catch (err: any) {
			this.error = err.message || 'Failed to disconnect from Meshtastic device';
			throw err;
		}
	}

	private setupEventListeners(): void {
		if (!this.client) return;

		// Listen for device metadata updates (includes myNodeInfo)
		this.client.on('myNodeInfo', (myNodeInfo: Protobuf.Mesh.MyNodeInfo) => {
			console.log('[Meshtastic] My Node Info:', myNodeInfo);
			this.deviceInfo = {
				...this.deviceInfo,
				myNodeNum: myNodeInfo.myNodeNum
			};
		});

		// Listen for device metadata updates
		this.client.on('deviceMetadata', (metadata: Protobuf.Mesh.DeviceMetadata) => {
			console.log('[Meshtastic] Device Metadata:', metadata);
			this.deviceInfo = {
				...this.deviceInfo,
				firmwareVersion: metadata.firmwareVersion,
				hwModel: Protobuf.Mesh.HardwareModel[metadata.hwModel || 0],
				hasWifi: metadata.hasWifi,
				hasBluetooth: metadata.hasBluetooth
			};
		});

		// Listen for config updates
		this.client.on('config', (config: Protobuf.LocalConfig.LocalConfig) => {
			console.log('[Meshtastic] Config:', config);
			if (config.lora) {
				this.deviceInfo = {
					...this.deviceInfo,
					region: Protobuf.Config.Config_LoRaConfig_RegionCode[config.lora.region || 0],
					modemPreset: Protobuf.Config.Config_LoRaConfig_ModemPreset[config.lora.modemPreset || 0]
				};
			}
		});

		// Listen for node info updates
		this.client.on('nodeInfo', (nodeInfo: Protobuf.Mesh.NodeInfo) => {
			console.log('[Meshtastic] Node Info:', nodeInfo);
		});

		// Listen for text messages
		this.client.on('messagePacket', (packet: Protobuf.Mesh.MeshPacket) => {
			console.log('[Meshtastic] Message Packet:', packet);
			if (packet.decoded?.portnum === Protobuf.Portnums.PortNum.TEXT_MESSAGE_APP) {
				const textDecoder = new TextDecoder();
				const text = textDecoder.decode(packet.decoded.payload);
				this.lastMessage = text;
				console.log('[Meshtastic] Text Message:', text);
			}
		});

		// Listen for connection state changes
		this.connection?.on('connectionStatus', (status: any) => {
			console.log('[Meshtastic] Connection Status:', status);
		});
	}

	private async requestDeviceInfo(): Promise<void> {
		if (!this.client) return;

		try {
			// Request my node info
			await this.client.sendPacket(
				Protobuf.Mesh.MeshPacket.toBinary({
					decoded: {
						portnum: Protobuf.Portnums.PortNum.ADMIN_APP,
						payload: Protobuf.Admin.AdminMessage.toBinary({
							getOwnerRequest: true
						}),
						wantResponse: true
					}
				})
			);

			// Request device metadata
			await this.client.sendPacket(
				Protobuf.Mesh.MeshPacket.toBinary({
					decoded: {
						portnum: Protobuf.Portnums.PortNum.ADMIN_APP,
						payload: Protobuf.Admin.AdminMessage.toBinary({
							getDeviceMetadataRequest: true
						}),
						wantResponse: true
					}
				})
			);

			// Request config
			await this.client.sendPacket(
				Protobuf.Mesh.MeshPacket.toBinary({
					decoded: {
						portnum: Protobuf.Portnums.PortNum.ADMIN_APP,
						payload: Protobuf.Admin.AdminMessage.toBinary({
							getConfigRequest: Protobuf.Admin.AdminMessage_GetConfigRequest_ConfigType.LORA_CONFIG
						}),
						wantResponse: true
					}
				})
			);
		} catch (err: any) {
			console.error('[Meshtastic] Failed to request device info:', err);
			this.error = err.message || 'Failed to request device information';
		}
	}

	async sendMessage(text: string, channelIndex: number = 0): Promise<void> {
		if (!this.client || !this.isConnected) {
			throw new Error('Not connected to a Meshtastic device');
		}

		try {
			this.error = null;
			const encoder = new TextEncoder();
			const payload = encoder.encode(text);

			await this.client.sendPacket(
				Protobuf.Mesh.MeshPacket.toBinary({
					decoded: {
						portnum: Protobuf.Portnums.PortNum.TEXT_MESSAGE_APP,
						payload
					},
					channel: channelIndex
				})
			);

			console.log('[Meshtastic] Sent message:', text);
		} catch (err: any) {
			this.error = err.message || 'Failed to send message';
			throw err;
		}
	}
}

// Export singleton instance
export const meshtasticStore = new MeshtasticStore();
