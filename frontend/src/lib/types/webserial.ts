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

// Web Serial API TypeScript definitions
// Based on https://wicg.github.io/serial/

export interface SerialPortRequestOptions {
	filters?: SerialPortFilter[];
}

export interface SerialPortFilter {
	usbVendorId?: number;
	usbProductId?: number;
}

export interface SerialOptions {
	baudRate: number;
	dataBits?: 7 | 8;
	stopBits?: 1 | 2;
	parity?: 'none' | 'even' | 'odd';
	bufferSize?: number;
	flowControl?: 'none' | 'hardware';
}

export interface SerialPortInfo {
	usbVendorId?: number;
	usbProductId?: number;
}

export interface Serial extends EventTarget {
	onconnect: ((this: Serial, ev: Event) => any) | null;
	ondisconnect: ((this: Serial, ev: Event) => any) | null;
	getPorts(): Promise<SerialPort[]>;
	requestPort(options?: SerialPortRequestOptions): Promise<SerialPort>;
	addEventListener(
		type: 'connect' | 'disconnect',
		listener: (this: Serial, ev: Event) => any,
		options?: boolean | AddEventListenerOptions
	): void;
}

export interface SerialPort extends EventTarget {
	readonly readable: ReadableStream<Uint8Array> | null;
	readonly writable: WritableStream<Uint8Array> | null;
	open(options: SerialOptions): Promise<void>;
	close(): Promise<void>;
	getInfo(): SerialPortInfo;
	forget(): Promise<void>;
	setSignals(signals: SerialOutputSignals): Promise<void>;
	getSignals(): Promise<SerialInputSignals>;
}

export interface SerialOutputSignals {
	dataTerminalReady?: boolean;
	requestToSend?: boolean;
	break?: boolean;
}

export interface SerialInputSignals {
	dataCarrierDetect: boolean;
	clearToSend: boolean;
	ringIndicator: boolean;
	dataSetReady: boolean;
}

declare global {
	interface Navigator {
		serial?: Serial;
	}
}
