// Meshtastic client store using Svelte 5 runes
import { MeshDevice, Protobuf, Types } from '@meshtastic/core';
import { TransportWebSerial } from '@meshtastic/transport-web-serial';

export interface DeviceInfo {
	myNodeNum?: number;
	firmwareVersion?: string;
	hwModel?: string;
	hasWifi?: boolean;
	hasBluetooth?: boolean;
	region?: string;
	modemPreset?: string;
}

class MeshtasticStore {
	private device = $state<MeshDevice | null>(null);

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

			// Create transport (prompts user for device selection)
			const transport = await TransportWebSerial.create(115200);

			// Create device with transport
			this.device = new MeshDevice(transport);

			// Set up event listeners
			this.setupEventListeners();

			this.isConnected = true;
			console.log('[Meshtastic] Connected to device');
		} catch (err: any) {
			this.error = err.message || 'Failed to connect to Meshtastic device';
			this.isConnected = false;
			this.device = null;
			throw err;
		} finally {
			this.isConnecting = false;
		}
	}

	async disconnect(): Promise<void> {
		try {
			this.error = null;

			if (this.device) {
				// The device/transport will handle cleanup
				this.device = null;
			}

			this.isConnected = false;
			this.deviceInfo = null;
			this.lastMessage = null;

			console.log('[Meshtastic] Disconnected');
		} catch (err: any) {
			this.error = err.message || 'Failed to disconnect from Meshtastic device';
			throw err;
		}
	}

	private setupEventListeners(): void {
		if (!this.device) return;

		const events = this.device.events;

		// Listen for device metadata
		events.onDeviceMetadataPacket.subscribe((metadata) => {
			console.log('[Meshtastic] Device Metadata:', metadata.data);
			const data = metadata.data;
			this.deviceInfo = {
				...this.deviceInfo,
				firmwareVersion: data.firmwareVersion,
				hwModel: data.hwModel ? Protobuf.Mesh.HardwareModel[data.hwModel] : undefined,
				hasWifi: data.hasWifi,
				hasBluetooth: data.hasBluetooth
			};
		});

		// Listen for my node info
		events.onMyNodeInfo.subscribe((info) => {
			console.log('[Meshtastic] My Node Info:', info);
			this.deviceInfo = {
				...this.deviceInfo,
				myNodeNum: info.myNodeNum
			};
		});

		// Listen for config updates
		events.onConfigPacket.subscribe((config) => {
			console.log('[Meshtastic] Config:', config.data);
			const data = config.data;

			if (data.payloadVariant.case === 'lora' && data.payloadVariant.value) {
				const lora = data.payloadVariant.value;
				this.deviceInfo = {
					...this.deviceInfo,
					region: lora.region ? Protobuf.Config.Config_LoRaConfig_RegionCode[lora.region] : undefined,
					modemPreset: lora.modemPreset ? Protobuf.Config.Config_LoRaConfig_ModemPreset[lora.modemPreset] : undefined
				};
			}
		});

		// Listen for text messages
		events.onMessagePacket.subscribe((packet) => {
			console.log('[Meshtastic] Message Packet:', packet);
			if (packet.data.decoded?.portnum === Protobuf.Portnums.PortNum.TEXT_MESSAGE_APP) {
				const textDecoder = new TextDecoder();
				const text = textDecoder.decode(packet.data.decoded.payload);
				this.lastMessage = text;
				console.log('[Meshtastic] Text Message:', text);
			}
		});

		// Listen for connection status
		events.onDeviceStatus.subscribe((status) => {
			console.log('[Meshtastic] Device Status:', status);
		});
	}

	async sendMessage(text: string, channelIndex: number = 0): Promise<void> {
		if (!this.device || !this.isConnected) {
			throw new Error('Not connected to a Meshtastic device');
		}

		try {
			this.error = null;

			await this.device.sendText(text);
			console.log('[Meshtastic] Sent message:', text);
		} catch (err: any) {
			this.error = err.message || 'Failed to send message';
			throw err;
		}
	}
}

// Export singleton instance
export const meshtasticStore = new MeshtasticStore();
