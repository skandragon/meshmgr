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
