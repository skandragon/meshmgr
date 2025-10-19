/*
 * Copyright (C) 2025 Michael Graff
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, version 3.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

// Web Serial connection store using Svelte 5 runes
import type { SerialPort } from '$lib/types/webserial';

export interface SerialConnectionState {
	isSupported: boolean;
	port: SerialPort | null;
	isConnected: boolean;
	isReading: boolean;
	error: string | null;
	lastMessage: string | null;
	deviceInfo: {
		vendorId?: number;
		productId?: number;
	} | null;
}

class SerialStore {
	private port = $state<SerialPort | null>(null);
	private reader: ReadableStreamDefaultReader<Uint8Array> | null = null;
	private writer: WritableStreamDefaultWriter<Uint8Array> | null = null;
	private decoder = new TextDecoder();
	private encoder = new TextEncoder();

	isSupported = $state(false);
	isConnected = $state(false);
	isReading = $state(false);
	error = $state<string | null>(null);
	lastMessage = $state<string | null>(null);
	deviceInfo = $state<{ vendorId?: number; productId?: number } | null>(null);

	constructor() {
		// Check if Web Serial API is supported
		if (typeof window !== 'undefined') {
			this.isSupported = 'serial' in navigator;
		}
	}

	async connect(baudRate: number = 115200): Promise<void> {
		if (!this.isSupported) {
			this.error = 'Web Serial API is not supported in this browser';
			throw new Error(this.error);
		}

		try {
			this.error = null;

			// Request port from user
			this.port = await navigator.serial!.requestPort({
				// Optional: filter for Meshtastic devices
				// filters: [{ usbVendorId: 0x303A }] // ESP32 vendor ID
			});

			// Open the port
			await this.port.open({
				baudRate,
				dataBits: 8,
				stopBits: 1,
				parity: 'none',
				flowControl: 'none'
			});

			// Get device info
			const info = this.port.getInfo();
			this.deviceInfo = {
				vendorId: info.usbVendorId,
				productId: info.usbProductId
			};

			// Setup writer
			if (this.port.writable) {
				this.writer = this.port.writable.getWriter();
			}

			this.isConnected = true;

			// Start reading
			this.startReading();
		} catch (err: any) {
			this.error = err.message || 'Failed to connect to serial port';
			this.isConnected = false;
			throw err;
		}
	}

	async disconnect(): Promise<void> {
		try {
			this.error = null;
			this.isReading = false;

			// Release reader
			if (this.reader) {
				await this.reader.cancel();
				this.reader.releaseLock();
				this.reader = null;
			}

			// Release writer
			if (this.writer) {
				this.writer.releaseLock();
				this.writer = null;
			}

			// Close port
			if (this.port) {
				await this.port.close();
				this.port = null;
			}

			this.isConnected = false;
			this.deviceInfo = null;
			this.lastMessage = null;
		} catch (err: any) {
			this.error = err.message || 'Failed to disconnect from serial port';
			throw err;
		}
	}

	private async startReading(): Promise<void> {
		if (!this.port || !this.port.readable) {
			return;
		}

		this.isReading = true;
		this.reader = this.port.readable.getReader();

		try {
			while (this.isReading) {
				const { value, done } = await this.reader.read();

				if (done) {
					// Reader has been canceled
					break;
				}

				if (value) {
					// Decode the received data
					const text = this.decoder.decode(value);
					this.lastMessage = text;

					// Log to console for now
					console.log('[Serial RX]:', text);
				}
			}
		} catch (err: any) {
			if (this.isReading) {
				// Only set error if we're still supposed to be reading
				this.error = err.message || 'Error reading from serial port';
				console.error('[Serial Error]:', err);
			}
		} finally {
			if (this.reader) {
				this.reader.releaseLock();
				this.reader = null;
			}
			this.isReading = false;
		}
	}

	async write(data: string): Promise<void> {
		if (!this.isConnected || !this.writer) {
			throw new Error('Not connected to a serial port');
		}

		try {
			this.error = null;
			const encoded = this.encoder.encode(data);
			await this.writer.write(encoded);
			console.log('[Serial TX]:', data);
		} catch (err: any) {
			this.error = err.message || 'Failed to write to serial port';
			throw err;
		}
	}

	async writeBytes(bytes: Uint8Array): Promise<void> {
		if (!this.isConnected || !this.writer) {
			throw new Error('Not connected to a serial port');
		}

		try {
			this.error = null;
			await this.writer.write(bytes);
			console.log('[Serial TX]:', bytes.length, 'bytes');
		} catch (err: any) {
			this.error = err.message || 'Failed to write to serial port';
			throw err;
		}
	}
}

// Export singleton instance
export const serialStore = new SerialStore();
